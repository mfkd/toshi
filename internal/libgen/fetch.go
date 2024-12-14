package libgen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// Fetch a list of books based on the search term
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
