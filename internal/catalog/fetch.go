package catalog

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseBooks(doc *goquery.Document) []Book {
	if table := doc.Find("#tablelibgen"); table.Length() > 0 {
		return parseIndexBooks(doc.Url, table)
	}

	var books []Book
	doc.Find("tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.ChildrenFiltered("td")
		if cells.Length() < 10 {
			return
		}
		cell := func(i int) string { return strings.TrimSpace(cells.Eq(i).Text()) }
		id := cell(0)
		if _, err := strconv.Atoi(id); err != nil {
			return
		}
		title, isbns := extractTitleAndISBN(cells.Eq(2).Find("a").Text())
		if title == "" {
			return
		}

		book := Book{
			Authors:   cell(1),
			Title:     title,
			ISBN:      isbns,
			Publisher: cell(3),
			Year:      cell(4),
			Pages:     cell(5),
			Language:  cell(6),
			Size:      cell(7),
			extension: strings.ToLower(cell(8)),
		}
		book.mirrors = mirrorLinks(doc.Url, cells.Eq(9).AddSelection(cells.Eq(10)))
		books = append(books, book)
	})
	return books
}

func parseIndexBooks(base *url.URL, table *goquery.Selection) []Book {
	var books []Book
	table.Find("tbody tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.ChildrenFiltered("td")
		if cells.Length() != 9 {
			return
		}
		title, isbns := extractTitleAndISBN(cells.Eq(0).ChildrenFiltered(`a[href*="edition.php"]`).First().Text())
		if title == "" || cells.Eq(6).Find(`a[href*="file.php"]`).Length() == 0 {
			return
		}
		cell := func(i int) string { return strings.TrimSpace(cells.Eq(i).Text()) }
		books = append(books, Book{
			Title:     title,
			ISBN:      isbns,
			Authors:   cell(1),
			Publisher: cell(2),
			Year:      cell(3),
			Language:  cell(4),
			Pages:     cell(5),
			Size:      cell(6),
			extension: strings.ToLower(cell(7)),
			mirrors:   mirrorLinks(base, cells.Eq(8)),
		})
	})
	return books
}

func mirrorLinks(base *url.URL, cells *goquery.Selection) []string {
	var mirrors []string
	cells.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		if strings.EqualFold(strings.TrimSpace(link.Text()), "Edit") {
			return
		}
		if target := resolveLink(base, link.AttrOr("href", "")); target != nil {
			if href := target.String(); !slices.Contains(mirrors, href) {
				mirrors = append(mirrors, href)
			}
		}
	})
	return mirrors
}

var paginatorPattern = regexp.MustCompile(`(?i)\bnew\s+Paginator\s*\(\s*["'][^"']+["']\s*,\s*(\d+)\s*,`)

func lastPage(doc *goquery.Document, term string) int {
	last := 1
	doc.Find("script").Each(func(_ int, script *goquery.Selection) {
		if match := paginatorPattern.FindStringSubmatch(script.Text()); match != nil {
			if count, err := strconv.Atoi(match[1]); err == nil {
				last = max(last, count)
			}
		}
	})
	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		target := resolveLink(doc.Url, link.AttrOr("href", ""))
		if target == nil || target.Host != doc.Url.Host || target.Path != doc.Url.Path {
			return
		}
		query := target.Query()
		if req := query.Get("req"); req != "" && req != term {
			return
		}
		if number, err := strconv.Atoi(query.Get("page")); err == nil {
			last = max(last, number)
		}
	})
	return last
}

// Search returns EPUB results from the first source with usable results.
func (c *Client) Search(ctx context.Context, sources []*url.URL, term string) ([]Book, error) {
	var failures []error
	var searched bool
	for _, source := range sources {
		books, err := c.searchSource(ctx, *source, term)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if err == nil {
			searched = true
		} else {
			failures = append(failures, fmt.Errorf("%s: %w", source.Host, err))
			c.log.Debug("search source failed", "source", source.Host, "err", err)
		}
		if len(books) > 0 {
			if err != nil {
				c.log.Warn("showing partial search results", "source", source.Host, "err", err)
			}
			return books, nil
		}
	}
	if searched {
		return nil, nil
	}
	if len(failures) == 0 {
		return nil, errors.New("no search sources configured")
	}
	return nil, fmt.Errorf("search for %q: %w", term, errors.Join(failures...))
}

const maxSearchPages = 3

func (c *Client) searchSource(ctx context.Context, source url.URL, term string) ([]Book, error) {
	var books []Book
	last := 1
	for page := 1; page <= min(last, maxSearchPages); page++ {
		if page > 1 {
			select {
			case <-ctx.Done():
				return books, ctx.Err()
			case <-time.After(time.Second):
			}
		}
		request, cancel := context.WithTimeout(ctx, 6*time.Second)
		doc, err := c.document(request, pageURL(source, term, page))
		cancel()
		if err != nil {
			return books, fmt.Errorf("search page %d: %w", page, err)
		}
		pageBooks := parseBooks(doc)
		if len(pageBooks) == 0 && doc.Find("input[name=req], #tablelibgen").Length() == 0 {
			title := strings.TrimSpace(doc.Find("title").First().Text())
			return books, fmt.Errorf("search page %d: expected a catalog response, got %q", page, title)
		}
		books = append(books, filterEPUB(pageBooks)...)
		last = max(last, lastPage(doc, term))
	}
	if last > maxSearchPages {
		c.log.Warn("search page limit reached", "limit", maxSearchPages, "source", source.Host)
	}
	return books, nil
}
