package lib

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/mfkd/toshi/internal/scraper"
)

func fetchBooks(ctx context.Context, s *scraper.Scraper, url string) ([]Book, error) {
	doc, err := s.ScrapeWithContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error scraping lib: %w", err)
	}

	var books []Book

	doc.Find("tr[valign=top]").Each(func(i int, s *goquery.Selection) {
		id := s.Find("td:nth-child(1)").Text()
		if _, err := strconv.Atoi(id); err != nil {
			// Skip rows where ID is not numeric (likely header)
			return
		}

		title, isbns := extractTitleAndISBN(s.Find("td:nth-child(3) a").Text())

		book := Book{
			ID:        id,
			Authors:   s.Find("td:nth-child(2)").Text(),
			Title:     title,
			ISBN:      isbns,
			Publisher: s.Find("td:nth-child(4)").Text(),
			Year:      s.Find("td:nth-child(5)").Text(),
			Pages:     s.Find("td:nth-child(6)").Text(),
			Language:  s.Find("td:nth-child(7)").Text(),
			Size:      s.Find("td:nth-child(8)").Text(),
			Extension: s.Find("td:nth-child(9)").Text(),
			Mirrors: []string{
				s.Find("td:nth-child(10) a:nth-child(1)").AttrOr("href", ""),
				s.Find("td:nth-child(11) a:nth-child(1)").AttrOr("href", ""),
			},
			Edit: s.Find("td:nth-child(11) a").AttrOr("href", ""),
		}
		books = append(books, book)
	})

	return books, nil
}

func fetchPagesURLs(ctx context.Context, s *scraper.Scraper, term string) ([]string, error) {
	var pages []string

	firstPage := pageURL(s.URL, term, 1)

	doc, err := s.ScrapeWithContext(ctx, firstPage)
	if err != nil {
		return nil, fmt.Errorf("error scraping lib: %w", err)
	}

	// Extract the <script> tag content
	var scriptContent string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if scriptContent == "" {
			scriptContent = s.Text()
		}
	})

	// Only one page found
	if scriptContent == "" {
		pages = append(pages, firstPage)
		return pages, nil
	}

	totalPages, err := totalPages(scriptContent)
	if err != nil {
		return nil, fmt.Errorf("error extracting total pages: %w", err)
	}

	return buildPageURLs(s.URL, term, totalPages), nil
}

func buildPageURLs(url, term string, totalPages int) []string {
	var urls []string
	for i := 1; i <= totalPages; i++ {
		urls = append(urls, pageURL(url, term, i))
	}
	return urls
}

func totalPages(content string) (int, error) {
	// Regex to match numbers followed by a comma
	re := regexp.MustCompile(`\b(\d+),`)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) != 3 {
		return 0, fmt.Errorf("unexpected number of page info matches: %d", len(matches))
	}

	totalPages, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0, fmt.Errorf("error converting total pages to integer: %w", err)
	}

	return totalPages, nil
}

func fetchAllBooks(ctx context.Context, s *scraper.Scraper, term string) ([]Book, error) {
	var books []Book

	pages, err := fetchPagesURLs(ctx, s, term)
	if err != nil {
		return nil, fmt.Errorf("error fetching page URLS: %w", err)
	}

	for _, page := range pages {
		booksOnPage, err := fetchBooks(ctx, s, page)
		if err != nil {
			return nil, fmt.Errorf("error fetching books from page: %w", err)
		}
		books = append(books, booksOnPage...)
		time.Sleep(s.RequestDelay)
	}

	return books, nil
}
