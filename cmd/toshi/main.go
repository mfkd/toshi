package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"

	"github.com/mfkd/toshi/internal/libgen"
	"github.com/mfkd/toshi/internal/ui"
)

func parseArgs() string {
	flag.Parse()
	args := flag.Args()

	// Ensure the positional argument "searchterm" is provided
	// TODO: Provide more informative output when user provides invalid input
	if len(args) < 1 {
		log.Println("Usage: toshi searchterm")
		log.Println("Example: toshi \"deep utopia\"")
		os.Exit(1)
	}

	return args[0]
}

func setupCollector() *colly.Collector {
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

// processBooks handles the user selection, fetches download links, and attempts to download the selected book.
func processBooks(c *colly.Collector, books []libgen.Book) {
	// Allow the user to select a book from the filtered list (e.g., EPUB books)
	// TODO: Make user select extension type as an argument
	selectedBook := ui.SelectBook(libgen.FilterEPUB(books))
	if selectedBook == nil {
		log.Println("No book selected.")
		return
	}

	fmt.Printf("Selected Book: %s\n", selectedBook.Title)

	// Fetch download links for the selected book
	downloadLinks := libgen.FetchDownloadLinks(c, *selectedBook)

	// Attempt to download the file
	fileName := libgen.FileName(*selectedBook)
	if err := libgen.TryDownloadLinks(c, downloadLinks, fileName); err != nil {
		log.Printf("Failed to download file for book %s: %v", selectedBook.Title, err)
	} else {
		fmt.Printf("Book downloaded successfully as %s\n", fileName)
	}
}

func main() {
	c := setupCollector()

	searchTerm := parseArgs()

	books, err := libgen.FetchAllBooks(c, searchTerm)
	if err != nil {
		log.Fatalf("Error fetching books from ages: %v", err)
	}

	processBooks(c, books)
}
