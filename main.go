package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	// It uses testify. I don't like it
	"github.com/PuerkitoBio/goquery"

	colly "github.com/gocolly/colly/v2"
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

func fetchLibgen(term string) (string, error) {
	libgenTemplate := "https://libgen.is/search.php?req=%s&column=title"
	query := fmt.Sprintf(libgenTemplate, url.QueryEscape(term))

	// Create a channel to pass the response body
	responseChannel := make(chan string, 1)
	errorChannel := make(chan error, 1)

	// Create a new collector
	c := colly.NewCollector()

	// Set up a handler for the response
	c.OnResponse(func(r *colly.Response) {
		responseChannel <- string(r.Body) // Send the response body to the channel
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		errorChannel <- fmt.Errorf("request URL: %s failed with status: %d. Error: %w", r.Request.URL, r.StatusCode, err)
	})

	// Visit the URL
	go func() {
		err := c.Visit(query)
		if err != nil {
			errorChannel <- err
		}
	}()

	// Wait for either the response or an error
	select {
	case response := <-responseChannel:
		return response, nil
	case err := <-errorChannel:
		return "", err
	}
}

func main() {
	term := "deep utopia"
	body, err := fetchLibgen(term)
	if err != nil {
		log.Fatalf("Error fetching Libgen: %v", err)
	}

	// Use the response body
	fetchBooks(body)
}

func fetchBooks(html string) {

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// Slice to hold books
	var books []Book

	// Find table rows containing book details
	doc.Find("tr[valign=top]").Each(func(i int, row *goquery.Selection) {
		// Extract data
		book := Book{
			ID:        row.Find("td").Eq(0).Text(),
			Authors:   row.Find("td").Eq(1).Text(),
			Title:     row.Find("td").Eq(2).Text(),
			Publisher: row.Find("td").Eq(3).Text(),
			Year:      row.Find("td").Eq(4).Text(),
			Pages:     row.Find("td").Eq(5).Text(),
			Language:  row.Find("td").Eq(6).Text(),
			Size:      row.Find("td").Eq(7).Text(),
			Extension: row.Find("td").Eq(8).Text(),
			Mirrors: []string{
				row.Find("td").Eq(9).Find("a").Eq(0).AttrOr("href", ""),
				row.Find("td").Eq(9).Find("a").Eq(1).AttrOr("href", ""),
			},
			Edit: row.Find("td").Eq(10).Find("a").AttrOr("href", ""),
		}

		// Append to books slice
		books = append(books, book)
	})

	// Print the books
	for _, book := range books {
		fmt.Printf("%+v\n", book)
	}
}
