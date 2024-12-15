package ui

import (
	"strings"

	"github.com/mfkd/toshi/internal/libgen"

	"github.com/fatih/color"
)

// Display a paginated list of books with single-line display
func DisplayBooksPaginated(books []libgen.Book, startIndex int) {
	endIndex := startIndex + booksPerPage
	if endIndex > len(books) {
		endIndex = len(books)
	}

	// Color styles
	dividerColor := color.New(color.FgHiBlack)
	headerColor := color.New(color.FgCyan).Add(color.Bold)
	bookColor := color.New(color.FgWhite)

	dividerColor.Println(strings.Repeat("=", 100))
	headerColor.Printf("Books %d to %d of %d\n", startIndex+1, endIndex, len(books))
	dividerColor.Println(strings.Repeat("=", 100))

	// Display each book on a single line
	for i := startIndex; i < endIndex; i++ {
		book := books[i]
		bookColor.Printf(
			"#%d | Title: %s | Author(s): %s | Year: %s | Publisher: %s | Pages: %s | Lang: %s | Size: %s | Format: %s | ISBN(s): %s\n",
			i+1,
			book.Title,
			book.Authors,
			book.Year,
			book.Publisher,
			book.Pages,
			book.Language,
			book.Size,
			book.Extension,
			strings.Join(book.ISBN, ", "),
		)
	}
	dividerColor.Println(strings.Repeat("=", 100))
}
