package catalog

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Book is an EPUB search result.
type Book struct {
	Authors   string
	Title     string
	ISBN      []string
	Publisher string
	Year      string
	Pages     string
	Language  string
	Size      string
	extension string
	mirrors   []string
}

var isbnPattern = regexp.MustCompile(`\b\d{9,13}\b`)

func extractTitleAndISBN(input string) (string, []string) {
	title := isbnPattern.ReplaceAllString(input, "")
	return strings.TrimSpace(strings.TrimRight(strings.TrimSpace(title), ",")), isbnPattern.FindAllString(input, -1)
}

// Generate a filename for the book
var invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\p{Cc}]`)

// sanitizeComponent removes extra whitespace and replaces invalid characters.
func sanitizeComponent(input string) string {
	input = strings.TrimSpace(input)
	return invalidFilenameChars.ReplaceAllString(input, "_")
}

// getFirstItem extracts the first item from a semicolon-separated list and sanitizes it.
func getFirstItem(input string) string {
	first, _, _ := strings.Cut(input, ";")
	return sanitizeComponent(first)
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

	const maxStemBytes = 230
	if len(filename) > maxStemBytes {
		end := maxStemBytes
		for !utf8.RuneStart(filename[end]) {
			end--
		}
		filename = filename[:end]
	}
	filename = strings.TrimRight(filename, ". ")
	if filename == "" {
		filename = "book"
	}
	return filename + ".epub"
}

func filterEPUB(books []Book) []Book {
	var filteredBooks []Book
	for _, b := range books {
		if b.extension == "epub" {
			filteredBooks = append(filteredBooks, b)
		}
	}
	return filteredBooks
}
