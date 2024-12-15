package libgen

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/gocolly/colly/v2"

	"github.com/mfkd/toshi/internal/logger"
)

// Fetch a list of books based on the search term
// TODO: Scrape all pages in response not just page 1.
func fetchBooks(c *colly.Collector, url string) ([]Book, error) {
	var books []Book

	// Handle book rows
	c.OnHTML("tr[valign=top]", func(e *colly.HTMLElement) {
		searchHandler(e, &books)
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		logger.Debugf("Fetch Books Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Visit the search page
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	return books, nil
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

func isValidPage(link string) bool {
	re := regexp.MustCompile(`search\.php.*&page=\d+`)
	return re.MatchString(link)
}

func fetchPagesURLs(c *colly.Collector, term string) ([]string, error) {
	var pages []string
	uniqueLinks := make(map[string]struct{}) // Map to store unique links
	var mu sync.Mutex

	// Capture pagination links
	// I wanted to make the GoQuery Selector more specific using "div#paginator_example_top ..."
	// but alas different selectors didn't work. So I opted for a more general approach for now.
	// We must recursively visit pages since only the next consecutive page is returned for each
	// page visited. They appear to be dynamic loading of some content using AJAX or JavaScript.
	c.OnHTML("tbody a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if isValidPage(href) {
			// Resolve relative URL to absolute URL
			fullURL := e.Request.AbsoluteURL(href)
			// Add only unique links
			mu.Lock()
			if _, exists := uniqueLinks[fullURL]; !exists {
				uniqueLinks[fullURL] = struct{}{}
				pages = append(pages, fullURL)
				mu.Unlock()
				// Recursively visit this page
				e.Request.Visit(fullURL)
			} else {
				mu.Unlock()
			}
		}
	})

	// Log errors with response details
	c.OnError(func(r *colly.Response, err error) {
		logger.Errorf("Fetch Pages Error: %v, Status Code: %d, Response: %s", err, r.StatusCode, string(r.Body))
	})

	// Build the search URL
	searchURL := defaultSearchURL(term)

	// Visit the search page
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting Libgen: %w", err)
	}

	c.Wait()
	// Return the collected unique pages
	return pages, nil
}

func fetchBooksFromURLs(c *colly.Collector, urls []string) ([]Book, error) {
	var books []Book
	for _, page := range urls {
		booksOnPage, err := fetchBooks(c, page)
		if err != nil {
			return nil, fmt.Errorf("error fetching books from URL %s: %w", page, err)
		}
		books = append(books, booksOnPage...)
	}
	return books, nil
}

func FetchAllBooks(c *colly.Collector, term string) ([]Book, error) {
	// Fetch the URLs of pages
	booksFromFirstPage, err := fetchBooks(c, defaultSearchURL(term))
	if err != nil {
		return nil, fmt.Errorf("error fetching books from first page 1 url: %w", err)
	}

	urls, err := fetchPagesURLs(c, term)
	if err != nil {
		return nil, fmt.Errorf("error fetching page URLs for term %q: %w", term, err)
	}

	// Fetch books from the URLs
	booksFromOtherPages, err := fetchBooksFromURLs(c, urls)
	if err != nil {
		return nil, fmt.Errorf("error fetching books from URLs: %w", err)
	}

	allBooks := append(booksFromFirstPage, booksFromOtherPages...)

	return allBooks, nil
}
