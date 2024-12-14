package libgen

import (
	"fmt"
	"regexp"
	"strings"
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

// Extract title and ISBN numbers from a string
func ExtractTitleAndISBN(input string) (string, []string) {

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
	return fmt.Sprintf("%s.%s", strings.ReplaceAll(b.Title, " ", "_"), strings.TrimSpace(b.Extension))
}

// Filter books to only include those with the "epub" extension
func FilterEPUB(books []Book) []Book {
	var filteredBooks []Book
	for _, b := range books {
		if b.Extension == "epub" {
			filteredBooks = append(filteredBooks, b)
		}
	}
	return filteredBooks
}
