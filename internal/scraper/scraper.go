package scraper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Scraper is a simple web scraper.
type Scraper struct {
	client       *http.Client
	UserAgent    string
	URL          string
	RequestDelay time.Duration
}

// NewScraper creates a new Scraper with the given URL.
func NewScraper(url string) *Scraper {
	return &Scraper{
		client:       &http.Client{},
		UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		URL:          url,
		RequestDelay: time.Second * 1,
	}
}

// Scrapes the given URL and returns the document.
func (s *Scraper) Scrape(url string) (*goquery.Document, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// load html document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading document: %w", err)
	}

	return doc, nil
}

// CheckHead sends a HEAD request to the given URL and returns the status code.
func (s *Scraper) CheckHead(url string) (int, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
