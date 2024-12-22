package libgen

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
func FileName(b Book) string {
	parts := make([]string, 0)

	// Add author if available
	if authors := strings.TrimSpace(b.Authors); authors != "" {
		parts = append(parts, authors)
	}

	// Add title
	parts = append(parts, strings.TrimSpace(b.Title))

	// Add publisher and year in parentheses if either is available
	pubYear := make([]string, 0)
	if pub := strings.TrimSpace(b.Publisher); pub != "" {
		pubYear = append(pubYear, pub)
	}
	if year := strings.TrimSpace(b.Year); year != "" {
		pubYear = append(pubYear, year)
	}
	if len(pubYear) > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", strings.Join(pubYear, " ")))
	}

	// Clean the filename components and join with dashes
	filename := strings.Join(parts, " - ")

	// Replace problematic characters
	filename = regexp.MustCompile(`[<>:"/\\|?*]`).ReplaceAllString(filename, "_")

	// Add extension
	return fmt.Sprintf("%s.%s", filename, strings.TrimSpace(b.Extension))
}

// Filter books to only include those with the "epub" extension
// TODO: FilterEPUB could use a more generic filter function with a predicate, making it reusable for other extensions
func FilterEPUB(books []Book) []Book {
	var filteredBooks []Book
	for _, b := range books {
		if b.Extension == "epub" {
			filteredBooks = append(filteredBooks, b)
		}
	}
	return filteredBooks
}
