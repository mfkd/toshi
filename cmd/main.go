package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/pflag"

	"github.com/mfkd/toshi/internal/libgen"
	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/ui"
)

var verbose *bool // Declare a global variable for verbose flag

func parseArgs() string {
	// Define the verbose flag
	verbose = pflag.BoolP("verbose", "v", false, "Enable verbose output with debug logs")

	pflag.Parse()

	args := pflag.Args()

	// Ensure the positional argument <searchterm> is provided
	// TODO: Provide more informative output when user provides invalid input
	if len(args) < 1 {
		fmt.Println("Usage: toshi <searchterm>")
		fmt.Println("Example: toshi The Iliad Homer")
		fmt.Println("Flags:")
		fmt.Println("  -v  Enable verbose output with debug logs")
		os.Exit(1)
	}

	return strings.Join(args, " ")
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
		fmt.Println("No book selected.")
		return
	}

	fmt.Printf("Selected Book: %s\n", selectedBook.Title)

	// Fetch download links for the selected book
	downloadLinks := libgen.FetchDownloadLinks(c, *selectedBook)

	// Attempt to download the file
	fileName := libgen.FileName(*selectedBook)

	if *verbose {
		fmt.Printf("Attempting to download book to: %s\n", fileName)
	}
	if err := libgen.TryDownloadLinks(c, downloadLinks, fileName); err != nil {
		logger.Errorf("Failed to download file for book %s: %v", selectedBook.Title, err)
	} else {
		fmt.Printf("Book downloaded successfully as %s\n", fileName)
	}
}

func main() {
	c := setupCollector()

	searchTerm := parseArgs()

	if *verbose {
		logger.Configure(logger.LevelDebug, nil)
		fmt.Println("DEBUG mode: Detailed logs are now enabled")
	}

	books, err := libgen.FetchAllBooks(c, searchTerm)
	if err != nil {
		logger.Fatalf("Error fetching books from ages: %v", err)
	}

	processBooks(c, books)
}
