package libgen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mfkd/toshi/internal/logger"
)

const downloadDir = "output"

func SetupDownloadCollector(c *colly.Collector, filename string) error {
	outputDir := downloadDir
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, filename)

	c.OnResponse(func(r *colly.Response) {
		file, err := os.Create(filePath)
		if err != nil {
			logger.Errorf("Failed to create file: %v", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, bytes.NewReader(r.Body))
		if err != nil {
			logger.Errorf("Failed to write to file: %v", err)
			return
		}

		fmt.Printf("File successfully downloaded to: %s", filePath)
	})

	return nil
}

func TryDownloadLinks(c *colly.Collector, downloadLinks []string, filename string) error {
	outputDir := downloadDir
	savePath := filepath.Join(outputDir, filename)

	if err := SetupDownloadCollector(c, filename); err != nil {
		return err
	}

	var err error
	for _, link := range downloadLinks {
		if err = downloadFile(c, link, savePath); err == nil {
			logger.Debugf("Successfully downloaded file from link: %s\n", link)
			break
		}
	}
	return err
}

func downloadFile(c *colly.Collector, fileURL, savePath string) error {
	err := c.Visit(fileURL)
	if err != nil {
		logger.Errorf("Failed to visit file URL: %v", err)
		return err
	}
	return nil
}

// FetchDownloadLinks extracts valid download links for a book from its mirror pages using a Colly collector.
func FetchDownloadLinks(c *colly.Collector, b Book) []string {
	// We need to try fetch download links from all Mirrors not just index 0
	var downloadLinks []string
	c.OnHTML("div#download ul li a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		// NOTE: gateway.ipfs.io seems to be the only working link
		if strings.Contains(href, "."+strings.TrimSpace(b.Extension)) {
			downloadLinks = append(downloadLinks, href)
		}
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		logger.Errorf("Download Links Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// TODO: Fetch all mirror links not just index 0
	err := c.Visit(b.Mirrors[0])
	if err != nil {
		logger.Errorf("Error visiting mirror link: %v", err)
	}
	return downloadLinks
}
