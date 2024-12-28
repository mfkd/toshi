package libgen

import "github.com/gocolly/colly/v2"

func SetupCollector() *colly.Collector {
	// Create a Colly collector
	c := colly.NewCollector(
		// Allow revisiting the base URL to fetch books from the current page and retrieve URLs of
		// other pages.
		// TODO: Enhance by reusing the base URL HTML for fetching of books and other pages
		colly.AllowURLRevisit(),
	)

	// Set headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	})

	return c
}
