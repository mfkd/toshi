package catalog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"

// Client searches catalogs and downloads EPUBs.
type Client struct {
	http *http.Client
	log  *slog.Logger
}

func NewClient(httpClient *http.Client, logger *slog.Logger) *Client {
	return &Client{http: httpClient, log: logger}
}

func (c *Client) get(ctx context.Context, target, referer string, html bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, errors.New("invalid request URL")
	}
	req.Header.Set("User-Agent", userAgent)
	if html {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	}
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	response, err := c.http.Do(req)
	if err != nil {
		var requestErr *url.Error
		if errors.As(err, &requestErr) {
			err = requestErr.Err
		}
		return nil, fmt.Errorf("GET %s: %w", req.URL.Hostname(), err)
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("GET %s: HTTP %s", response.Request.URL.Hostname(), response.Status)
	}
	return response, nil
}

func (c *Client) document(ctx context.Context, target string) (*goquery.Document, error) {
	response, err := c.get(ctx, target, "", true)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read HTML: %w", err)
	}
	doc.Url = response.Request.URL
	return doc, nil
}
