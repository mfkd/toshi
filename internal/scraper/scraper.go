package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

// Scrape sends a GET request to the given URL and returns the document.
func (s *Scraper) Scrape(url string) (*goquery.Document, error) {
	return s.ScrapeWithContext(context.Background(), url)
}

// ScrapeWithContext sends a GET request to the given URL and returns the document with context.
func (s *Scraper) ScrapeWithContext(ctx context.Context, url string) (*goquery.Document, error) {
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
		return nil, fmt.Errorf("HTTP request failed with status %d (%s) for URL: %s", resp.StatusCode, resp.Status, url)
	}

	// load html document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading document: %w", err)
	}

	return doc, nil
}

// CheckHead sends a HEAD request to the given URL and returns the status code.
func (s *Scraper) CheckHead(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(context.Background(), "HEAD", url, nil)
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

// DownloadFile downloads a file from the given URL and saves it to the download directory.
func (s *Scraper) DownloadFile(ctx context.Context, filename, downloadURL, downloadDir string) error {
	// Validate the URL
	if _, err := url.ParseRequestURI(downloadURL); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create new HTTP GET request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", s.UserAgent)

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create output file
	out, err := os.Create(filepath.Join(downloadDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
