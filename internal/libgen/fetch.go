package libgen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// Fetch a list of books based on the search term
// TODO: Scrape all pages in response not just page 1.
func FetchBooks(c *colly.Collector, term string) ([]Book, error) {
	var books []Book

	// Handle book rows
	c.OnHTML("tr[valign=top]", func(e *colly.HTMLElement) {
		SearchHandler(e, &books)
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Fetch Books Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Construct the search URL using the helper function
	searchURL := DefaultSearchURL(term)

	// Visit the search page
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	return books, nil
}

func SearchHandler(e *colly.HTMLElement, books *[]Book) {
	id := e.ChildText("td:nth-child(1)")
	if _, err := strconv.Atoi(id); err != nil {
		// Skip rows where ID is not numeric (likely header)
		return
	}

	title, isbns := ExtractTitleAndISBN(e.ChildText("td:nth-child(3) a"))

	book := Book{
		ID:        id,
		Authors:   e.ChildText("td:nth-child(2)"),
		Title:     title,
		ISBN:      isbns,
		Publisher: e.ChildText("td:nth-child(4)"),
		Year:      e.ChildText("td:nth-child(5)"),
		Pages:     e.ChildText("td:nth-child(6)"),
		Language:  e.ChildText("td:nth-child(7)"),
		Size:      e.ChildText("td:nth-child(8)"),
		Extension: e.ChildText("td:nth-child(9)"),
		Mirrors: []string{
			e.ChildAttr("td:nth-child(10) a:nth-child(1)", "href"),
			e.ChildAttr("td:nth-child(11) a:nth-child(1)", "href"),
		},
		Edit: e.ChildAttr("td:nth-child(11) a", "href"),
	}
	*books = append(*books, book)
}

func isValidPage(link string) bool {
	re := regexp.MustCompile(`search\.php.*&page=\d+`)
	return re.MatchString(link)
}

func FetchPagesURLs(c *colly.Collector, term string) ([]string, error) {
	var pages []string
	uniqueLinks := make(map[string]struct{}) // Map to store unique links

	// Capture pagination links
	// I wanted to make the GoQuery Selector more specific using "div#paginator_example_top ..."
	// but alas different selectors didn't work. So I opted for a more general approach for now.
	// We must recursively visit pages since only the next consecutive page is returned for each
	// page visited. They appear to be dynamic loading of some content using AJAX or JavaScript.
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if isValidPage(href) {
			// Resolve relative URL to absolute URL
			fullURL := e.Request.AbsoluteURL(href)

			// Add only unique links
			if _, exists := uniqueLinks[fullURL]; !exists {
				uniqueLinks[fullURL] = struct{}{}
				pages = append(pages, fullURL)

				// Recursively visit this page
				e.Request.Visit(fullURL)
			}
		}
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Fetch Pages Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Build the search URL
	searchURL := DefaultSearchURL(term)

	// Visit the search page
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	// Return the collected unique pages
	return pages, nil
}
