package lib

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

func setupDownloadCollector(c Collector, filename string) error {
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

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		// TODO: Check if there is a way to handle this better.
		// Debug over Error as we want to try available links until we succeed.
		logger.Debugf("Download File Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	return nil
}

func tryDownloadLinks(c Collector, downloadLinks []string, filename string) error {
	if err := setupDownloadCollector(c, filename); err != nil {
		return err
	}

	var err error
	for _, link := range downloadLinks {
		if err = downloadFile(c, link); err == nil {
			// TODO: Check if there is a way to handle this better.
			// Debug over Error as we want to try available links until we succeed.
			logger.Debugf("Successfully downloaded file from link: %s\n", link)
			break
		}
	}

	return err
}

func downloadFile(c Collector, fileURL string) error {
	err := c.Visit(fileURL)
	if err != nil {
		// TODO: Check if there is a way to handle this better.
		// Debug over Error as we want to try available links until we succeed.
		logger.Debugf("Failed to visit file URL: %v", err)
		return err
	}
	c.Wait()

	return nil
}

// fetchDownloadLinks extracts valid download links for a book from its mirror pages using a Colly collector.
func fetchDownloadLinks(c Collector, b Book) []string {
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
		// TODO: Check if there is a way to handle this better.
		// Set logger to debug since we retry fetching pages until we succeed.
		logger.Debugf("Download Links Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// TODO: Fetch all mirror links not just index 0
	// Mirrors[0] contains the most reliable links so we will use it for now.
	err := c.Visit(b.Mirrors[0])
	if err != nil {
		logger.Errorf("Error visiting mirror link: %v", err)
	}
	c.Wait()

	return downloadLinks
}
