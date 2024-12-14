package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type Book struct {
	ID        string
	Authors   string
	Title     string
	Publisher string
	Year      string
	Pages     string
	Language  string
	Size      string
	Extension string
	Mirrors   []string
	Edit      string
}

// Base URL for the LibGen search
const libgenSearchBaseURL = "https://libgen.is/search.php"

// Helper function to construct the search URL on TITLE Column
func constructTitleSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "title")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}

// Helper function to construct the search URL on ALL Column
func constructDefaultSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "def")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}

func fetchBooks(term string) ([]Book, error) {
	var books []Book

	// Create a Colly collector
	c := colly.NewCollector(
		colly.AllowedDomains("libgen.is", "libgen.rs"),
	)

	// Set headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	})

	// Handle book rows
	c.OnHTML("tr[valign=top]", func(e *colly.HTMLElement) {
		book := Book{
			ID:        e.ChildText("td:nth-child(1)"),
			Authors:   e.ChildText("td:nth-child(2)"),
			Title:     e.ChildText("td:nth-child(3)"),
			Publisher: e.ChildText("td:nth-child(4)"),
			Year:      e.ChildText("td:nth-child(5)"),
			Pages:     e.ChildText("td:nth-child(6)"),
			Language:  e.ChildText("td:nth-child(7)"),
			Size:      e.ChildText("td:nth-child(8)"),
			Extension: e.ChildText("td:nth-child(9)"),
			Mirrors: []string{
				e.ChildAttr("td:nth-child(10) a:nth-child(1)", "href"),
				e.ChildAttr("td:nth-child(10) a:nth-child(2)", "href"),
			},
			Edit: e.ChildAttr("td:nth-child(11) a", "href"),
		}
		books = append(books, book)
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Construct the search URL using the helper function
	searchURL := constructDefaultSearchURL(term)

	// Visit the search page
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	return books, nil
}

func main() {
	term := "murakami city"

	books, err := fetchBooks(term)
	if err != nil {
		log.Fatalf("Error fetching books: %v", err)
	}

	fmt.Printf("Found %d books:\n", len(books))
	for _, book := range books {
		fmt.Printf("%+v\n", book)
	}
}
