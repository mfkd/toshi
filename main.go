package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Book struct {
	ID        string
	Authors   string
	Title     string
	ISBN      []string
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
func titleSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "title")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}

// Helper function to construct the search URL on ALL Column
func defaultSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "def")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}

func extractTitleAndISBN(input string) (string, []string) {
	// Regular expression to match ISBN numbers
	isbnRegex := regexp.MustCompile(`\b\d{9,13}\b`)

	// Find all ISBN numbers
	isbns := isbnRegex.FindAllString(input, -1)

	// Remove ISBN numbers from the original string
	title := isbnRegex.ReplaceAllString(input, "")

	// Clean up the title (remove extra spaces and trailing commas)
	title = strings.TrimSpace(title)
	title = strings.TrimRight(title, ",")

	return title, isbns
}

func fetchBooks(c *colly.Collector, term string) ([]Book, error) {
	var books []Book

	// Handle book rows
	c.OnHTML("tr[valign=top]", func(e *colly.HTMLElement) {
		searchHandler(e, &books)
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Fetch Books Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Construct the search URL using the helper function
	searchURL := defaultSearchURL(term)

	// Visit the search page
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	return books, nil
}

func fetchDownloadLinks(c *colly.Collector, b Book) []string {
	var downloadLinks []string
	c.OnHTML("div#download ul li a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "gateway.ipfs.io") {
			downloadLinks = append(downloadLinks, href)
		}
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Download Links Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	err := c.Visit(b.Mirrors[0])
	if err != nil {
		log.Printf("Error visiting mirror link: %v", err)
	}
	return downloadLinks
}

func searchHandler(e *colly.HTMLElement, books *[]Book) {
	id := e.ChildText("td:nth-child(1)")
	if _, err := strconv.Atoi(id); err != nil {
		// Skip rows where ID is not numeric (likely header)
		return
	}

	title, isbns := extractTitleAndISBN(e.ChildText("td:nth-child(3) a"))

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

func parseArgs() string {
	flag.Parse()
	args := flag.Args()

	// Ensure the positional argument "searchterm" is provided
	if len(args) < 1 {
		fmt.Println("Usage: toshi searchterm")
		fmt.Println("Example: toshi \"deep utopia\"")
		os.Exit(1)
	}

	return args[0]
}

func downloadFile(c *colly.Collector, fileURL, savePath string) error {
	// Set up the OnResponse callback to save the file content
	c.OnResponse(func(r *colly.Response) {
		// Create the file in the current working directory
		savePath := filepath.Join(".", savePath)
		file, err := os.Create(savePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer file.Close()

		// Write the response body to the file
		reader := bytes.NewReader(r.Body)
		_, err = io.Copy(file, reader)
		if err != nil {
			log.Printf("Failed to write to file: %v", err)
			return
		}

		log.Printf("File successfully downloaded to: %s", savePath)
	})

	// Visit the file URL
	err := c.Visit(fileURL)
	if err != nil {
		log.Printf("Failed to visit file URL: %v", err)
		return err
	}

	return nil
}

func main() {

	searchTerm := parseArgs()

	// Create a Colly collector
	c := colly.NewCollector()

	// Set headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	})

	books, err := fetchBooks(c, searchTerm)
	if err != nil {
		log.Fatalf("Error fetching books: %v", err)
	}

	fmt.Printf("Found %d books:\n", len(books))
	for _, book := range books {
		fmt.Printf("%+v\n", book)
	}

	downloadLinks := fetchDownloadLinks(c, books[0])
	if err := downloadFile(c, downloadLinks[0], fileName(books[0])); err != nil {

	}
}

func fileName(b Book) string {
	return fmt.Sprintf("%s.%s", strings.ReplaceAll(b.Title, " ", "_"), strings.TrimSpace(b.Extension))
}
