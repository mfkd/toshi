package lib

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/scraper"
)

const downloadDir = "output"

func tryDownloadLinks(ctx context.Context, s *scraper.Scraper, downloadLinks []string, filename string) error {
	var err error
	for _, link := range downloadLinks {
		if err = s.DownloadFile(ctx, filename, link, downloadDir); err == nil {
			// TODO: Check if there is a way to handle this better.
			// Debug over Error as we want to try available links until we succeed.
			logger.Debugf("Successfully downloaded file from link: %s\n", link)
			break
		}
	}
	return err
}

func fetchDownloadLinks(ctx context.Context, s *scraper.Scraper, b Book) []string {
	var downloadLinks []string

	// TODO: Handle multiple mirrors
	doc, err := s.ScrapeWithContext(ctx, b.Mirrors[0])
	if err != nil {
		logger.Errorf("Error scraping download links: %v", err)
		return downloadLinks
	}

	doc.Find("div#download ul li a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// TODO: Make this more robust.
			if strings.Contains(href, "."+strings.TrimSpace(b.Extension)) {
				downloadLinks = append(downloadLinks, href)
			}
		}
	})

	if len(downloadLinks) == 0 {
		logger.Errorf("No download links found for book: %s", b.Title)
	}

	return downloadLinks
}
