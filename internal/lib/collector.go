package lib

import "github.com/gocolly/colly/v2"

type Collector struct {
	*colly.Collector
	url string
}

// SetupCollector initializes and returns a configured collector.
func SetupCollector(urls []string) Collector {
	// Create a custom collector
	c := Collector{
		colly.NewCollector(
			// Allow revisiting the base URL to fetch books from the current page and retrieve URLs
			// of other pages.
			// TODO: Enhance by reusing the base URL HTML for fetching of books and other pages
			colly.AllowURLRevisit(),
		),
		// TODO: Add support for multiple URLs to improve reliability in the event of server
		// outages.
		urls[0],
	}

	// Set headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	})

	return c
}
