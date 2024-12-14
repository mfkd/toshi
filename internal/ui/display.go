package ui

import (
	"strings"

	"github.com/mfkd/toshi/internal/libgen"

	"github.com/fatih/color"
)

// Display a paginated list of books
func DisplayBooksPaginated(books []libgen.Book, startIndex int) {
	endIndex := startIndex + booksPerPage
	if endIndex > len(books) {
		endIndex = len(books)
	}

	// Color styles
	dividerColor := color.New(color.FgHiBlack)
	headerColor := color.New(color.FgCyan).Add(color.Bold)
	fieldLabelColor := color.New(color.FgHiBlue).Add(color.Bold)
	fieldValueColor := color.New(color.FgWhite)

	dividerColor.Println(strings.Repeat("=", 100))
	headerColor.Printf("Books %d to %d of %d\n", startIndex+1, endIndex, len(books))
	dividerColor.Println(strings.Repeat("=", 100))

	// Display each book in detail
	for i := startIndex; i < endIndex; i++ {
		book := books[i]
		headerColor.Printf("Book #%d\n", i+1)
		fieldLabelColor.Printf("Title       : ")
		fieldValueColor.Println(book.Title)
		fieldLabelColor.Printf("Author(s)   : ")
		fieldValueColor.Println(book.Authors)
		fieldLabelColor.Printf("Year        : ")
		fieldValueColor.Println(book.Year)
		fieldLabelColor.Printf("Publisher   : ")
		fieldValueColor.Println(book.Publisher)
		fieldLabelColor.Printf("Pages       : ")
		fieldValueColor.Println(book.Pages)
		fieldLabelColor.Printf("Language    : ")
		fieldValueColor.Println(book.Language)
		fieldLabelColor.Printf("Size        : ")
		fieldValueColor.Println(book.Size)
		fieldLabelColor.Printf("Format      : ")
		fieldValueColor.Println(book.Extension)
		fieldLabelColor.Printf("ISBN(s)     : ")
		fieldValueColor.Println(strings.Join(book.ISBN, ", "))
		fieldLabelColor.Printf("Mirrors     : ")
		fieldValueColor.Println(strings.Join(book.Mirrors, ", "))
		fieldLabelColor.Printf("Edition     : ")
		fieldValueColor.Println(book.Edit)
		dividerColor.Println(strings.Repeat("-", 100))
	}
}
