package lib

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO: Enhance filtering by ordering books by most complete metadata

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

// Extract title and ISBN numbers from a string
// TODO: Consider a more robust ISBN regex and think about error handling
// isbnRegex := regexp.MustCompile(`\b(?:\d{9}[\dX]|\d{13})\b`)
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

// Generate a filename for the book
var invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*]`)

// sanitizeComponent removes extra whitespace and replaces invalid characters.
func sanitizeComponent(input string) string {
	input = strings.TrimSpace(input)
	return invalidFilenameChars.ReplaceAllString(input, "_")
}

// getFirstItem extracts the first item from a semicolon-separated list and sanitizes it.
func getFirstItem(input string) string {
	if input == "" {
		return ""
	}
	parts := strings.Split(input, ";")
	return sanitizeComponent(parts[0])
}

// fileName generates a filename from book details.
func fileName(b Book) string {
	var parts []string

	// Add author if available.
	if author := getFirstItem(b.Authors); author != "" {
		parts = append(parts, author)
	}

	// Add title.
	if title := sanitizeComponent(b.Title); title != "" {
		parts = append(parts, title)
	}

	// Add publisher and year.
	var pubYear string
	if publisher := getFirstItem(b.Publisher); publisher != "" {
		pubYear = publisher
	}
	if year := strings.TrimSpace(b.Year); year != "" {
		if pubYear != "" {
			pubYear = fmt.Sprintf("%s (%s)", pubYear, year)
		} else {
			pubYear = year
		}
	}
	if pubYear != "" {
		parts = append(parts, pubYear)
	}

	// Join parts with dashes.
	filename := strings.Join(parts, " - ")
	filename = sanitizeComponent(filename)

	return fmt.Sprintf("%s.%s", filename, strings.TrimSpace(b.Extension))
}

// Filter books to only include those with the "epub" extension
// TODO: filterEPUB could use a more generic filter function with a predicate, making it reusable for other extensions
func filterEPUB(books []Book) []Book {
	var filteredBooks []Book
	for _, b := range books {
		if b.Extension == "epub" {
			filteredBooks = append(filteredBooks, b)
		}
	}
	return filteredBooks
}
