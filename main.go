package main

import (
	"fmt"
	"net/url"

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

func main() {

	term := "deep utopia"
	libgenTemplate := "https://libgen.is/search.php?req=%s&column=title"

	query := fmt.Sprintf(libgenTemplate, url.QueryEscape(term))
	// Create a new collector
	c := colly.NewCollector()

	// Set up a handler for the response
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response received:")
		fmt.Println(string(r.Body)) // Print the response body as a string
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with status: %d. Error: %s\n", r.Request.URL, r.StatusCode, err)
	})

	// Fetch the web page
	fmt.Println("Visiting ", query)
	c.Visit(query)
}
