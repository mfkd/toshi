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
	if len(args) < 1 {
		log.Println("Usage: toshi searchterm")
		log.Println("Example: toshi \"deep utopia\"")
		os.Exit(1)
	}

	return args[0]
}

func main() {

	searchTerm := parseArgs()

	// Create a Colly collector
	c := colly.NewCollector()

	// Set headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	})

	books, err := libgen.FetchBooks(c, searchTerm)
	if err != nil {
		log.Fatalf("Error fetching books: %v", err)
	}

	selectedBook := ui.SelectBook(libgen.FilterEPUB(books))
	fmt.Println(selectedBook)

	downloadLinks := libgen.FetchDownloadLinks(c, *selectedBook)
	if err := libgen.TryDownloadLinks(c, downloadLinks, libgen.FileName(books[0])); err != nil {
		log.Printf("Failed to download file for book %s: %v", books[0].Title, err)
	}
}
