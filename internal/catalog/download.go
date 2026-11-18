package catalog

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type downloadLink struct {
	url     string
	referer string
}

// Download saves the selected EPUB in dir and returns its final path.
func (c *Client) Download(ctx context.Context, book Book, dir string) (string, error) {
	if len(book.mirrors) == 0 {
		return "", errors.New("no download mirrors listed")
	}
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	filename := fileName(book)

	var failures []error
	for _, mirror := range book.mirrors {
		if ctx.Err() != nil {
			break
		}
		request, done := context.WithTimeout(ctx, 8*time.Second)
		links, err := c.downloadLinks(request, mirror)
		done()
		if err == nil {
			var saved string
			saved, err = c.tryLinks(ctx, links, dir, filename, 30*time.Second)
			if err == nil {
				return saved, nil
			}
		}

		label := "invalid URL"
		if target, parseErr := url.Parse(mirror); parseErr == nil {
			label = target.Hostname() + target.EscapedPath()
		}
		c.log.Debug("mirror failed", "mirror", label, "err", err)
		failures = append(failures, fmt.Errorf("mirror %s: %w", label, err))
	}
	return "", errors.Join(append(failures, ctx.Err())...)
}

func (c *Client) tryLinks(ctx context.Context, links []downloadLink, dir, filename string, timeout time.Duration) (string, error) {
	if len(links) == 0 {
		return "", errors.New("no EPUB download links found")
	}

	var failures []error
	for _, link := range links {
		if ctx.Err() != nil {
			break
		}
		attempt, cancel := context.WithTimeout(ctx, timeout)
		saved, err := c.saveEPUB(attempt, link, dir, filename)
		cancel()
		if err == nil {
			return saved, nil
		}
		failures = append(failures, err)
	}
	return "", errors.Join(append(failures, ctx.Err())...)
}

func (c *Client) downloadLinks(ctx context.Context, mirror string) ([]downloadLink, error) {
	if target, err := url.Parse(mirror); err == nil && strings.EqualFold(path.Ext(target.Path), ".epub") {
		return []downloadLink{{url: mirror}}, nil
	}
	doc, err := c.document(ctx, mirror)
	if err != nil {
		return nil, err
	}

	var links []downloadLink
	doc.Find(`#download a[href], a[href^="get.php?"], a[href^="/get.php?"]`).Each(func(_ int, link *goquery.Selection) {
		target := resolveLink(doc.Url, link.AttrOr("href", ""))
		if target == nil {
			return
		}
		ext := path.Ext(target.Path)
		switch {
		case strings.EqualFold(ext, ".epub"),
			path.Base(target.Path) == "get.php",
			strings.HasPrefix(target.Path, "/ipfs/") && ext == "",
			ext == "" && strings.EqualFold(strings.TrimSpace(link.Text()), "GET"):
			candidate := downloadLink{url: target.String()}
			if target.Scheme == doc.Url.Scheme && target.Host == doc.Url.Host {
				candidate.referer = doc.Url.String()
			}
			if !slices.Contains(links, candidate) {
				links = append(links, candidate)
			}
		}
	})
	return links, nil
}

func (c *Client) saveEPUB(ctx context.Context, link downloadLink, dir, filename string) (string, error) {
	response, err := c.get(ctx, link.url, link.referer, false)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}
	temp, err := os.CreateTemp(dir, ".toshi-*")
	if err != nil {
		return "", fmt.Errorf("create temporary download: %w", err)
	}
	defer os.Remove(temp.Name())

	size, saveErr := io.Copy(temp, response.Body)
	if saveErr == nil {
		saveErr = validateEPUB(temp, size)
	}
	closeErr := temp.Close()
	if saveErr != nil || closeErr != nil {
		return "", fmt.Errorf("save EPUB: %w", errors.Join(saveErr, closeErr))
	}
	return publishEPUB(ctx, temp.Name(), dir, filename)
}

func publishEPUB(ctx context.Context, temporary, dir, filename string) (string, error) {
	stem := strings.TrimSuffix(filename, ".epub")
	for i := 1; ; i++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if i > 1 {
			filename = fmt.Sprintf("%s (%d).epub", stem, i)
		}
		destination := filepath.Join(dir, filename)
		// Link atomically publishes the completed file and fails if the name exists.
		err := os.Link(temporary, destination)
		if err == nil {
			return destination, nil
		}
		if !errors.Is(err, fs.ErrExist) {
			return "", fmt.Errorf("publish EPUB: %w", err)
		}
	}
}

func validateEPUB(file io.ReaderAt, size int64) error {
	archive, err := zip.NewReader(file, size)
	if err != nil {
		return fmt.Errorf("read EPUB archive: %w", err)
	}
	const mediaType = "application/epub+zip"
	mimetype, err := archive.Open("mimetype")
	if err != nil {
		return fmt.Errorf("read EPUB mimetype: %w", err)
	}
	defer mimetype.Close()
	content, err := io.ReadAll(io.LimitReader(mimetype, int64(len(mediaType)+1)))
	if err != nil {
		return fmt.Errorf("read EPUB mimetype: %w", err)
	}
	if string(content) != mediaType {
		return errors.New("invalid EPUB mimetype")
	}
	var container struct {
		XMLName   xml.Name `xml:"container"`
		Rootfiles []struct {
			Path string `xml:"full-path,attr"`
		} `xml:"rootfiles>rootfile"`
	}
	if err := decodeEPUBXML(archive, "META-INF/container.xml", &container); err != nil {
		return fmt.Errorf("read EPUB container: %w", err)
	}
	if len(container.Rootfiles) == 0 || !fs.ValidPath(container.Rootfiles[0].Path) {
		return errors.New("EPUB container has no valid package path")
	}
	var manifest struct {
		XMLName xml.Name `xml:"package"`
	}
	if err := decodeEPUBXML(archive, container.Rootfiles[0].Path, &manifest); err != nil {
		return fmt.Errorf("read EPUB package: %w", err)
	}
	return nil
}

func decodeEPUBXML(archive *zip.Reader, name string, dest any) error {
	entry, err := archive.Open(name)
	if err != nil {
		return err
	}
	defer entry.Close()
	return xml.NewDecoder(entry).Decode(dest)
}
