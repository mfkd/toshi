package ui

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/mfkd/toshi/internal/lib"
	"golang.org/x/term"
)

// Function to get the terminal width dynamically
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(syscall.Stdout))
	if err != nil {
		return 80 // Fallback to a default width
	}
	return width
}

// Display a paginated list of books with dynamic dividers
func displayBooksPaginated(books []lib.Book, startIndex int) {
	terminalWidth := getTerminalWidth()
	endIndex := startIndex + booksPerPage
	if endIndex > len(books) {
		endIndex = len(books)
	}

	// Print the centered header
	fmt.Println(FgBlue + strings.Repeat("=", terminalWidth) + Reset)
	header := fmt.Sprintf("Books %d to %d of %d", startIndex+1, endIndex, len(books))
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", (terminalWidth-len(header))/2), Bold+FgBrightWhite+header+Reset, strings.Repeat(" ", (terminalWidth-len(header))/2))
	fmt.Println(FgBlue + strings.Repeat("=", terminalWidth) + Reset)

	// Print each book
	for i := startIndex; i < endIndex; i++ {
		book := books[i]

		// Print book index
		fmt.Printf("%s#%d%s\n", Bold+FgYellow, i+1, Reset)

		// Print book details with refined alignment and hide empty fields
		if book.Title != "" {
			fmt.Printf("  %sTitle:%s       %s%s\n", Bold+FgCyan, Reset, Bold+FgBrightWhite, book.Title)
		}
		if book.Authors != "" {
			fmt.Printf("  %sAuthor(s):%s   %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Authors)
		}
		if book.Year != "" {
			fmt.Printf("  %sYear:%s        %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Year)
		}
		if book.Publisher != "" {
			fmt.Printf("  %sPublisher:%s   %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Publisher)
		}
		if book.Pages != "" {
			fmt.Printf("  %sPages:%s       %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Pages)
		}
		if book.Language != "" {
			fmt.Printf("  %sLanguage:%s    %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Language)
		}
		if book.Size != "" {
			fmt.Printf("  %sSize:%s        %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Size)
		}
		if book.Extension != "" {
			fmt.Printf("  %sFormat:%s      %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, book.Extension)
		}
		if len(book.ISBN) > 0 {
			fmt.Printf("  %sISBN(s):%s     %s%s\n", FgBlue, Reset, Bold+FgBrightWhite, strings.Join(book.ISBN, ", "))
		}

		// Add a dashed divider between books
		fmt.Println(FgBlue + strings.Repeat("-", terminalWidth) + Reset)
		fmt.Println() // Add extra vertical space for better readability
	}

	// Final divider
	fmt.Println(FgBlue + strings.Repeat("=", terminalWidth) + Reset)

	// Print options at the end
	fmt.Printf("%sOptions:%s\n", Bold+FgRed, Reset)
	fmt.Printf("%sEnter the number of the book to select it.%s\n", Bold+FgRed, Reset)
	fmt.Printf("%sEnter 'q' to Quit.%s\n", Bold+FgRed, Reset)
}
