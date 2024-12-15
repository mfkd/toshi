package ui

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/mfkd/toshi/internal/libgen"
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
func displayBooksPaginated(books []libgen.Book, startIndex int) {
	terminalWidth := getTerminalWidth()
	endIndex := startIndex + booksPerPage
	if endIndex > len(books) {
		endIndex = len(books)
	}

	// Define color styles optimized for light and dark modes
	dividerColor := color.New(color.FgBlue)                   // Regular blue for dividers
	headerColor := color.New(color.FgHiWhite).Add(color.Bold) // Bold white for headers
	indexColor := color.New(color.FgYellow).Add(color.Bold)   // Regular yellow for index
	titleColor := color.New(color.FgCyan).Add(color.Bold)     // Bold cyan for titles
	labelColor := color.New(color.FgBlue).Add(color.Bold)     // Bold blue for labels
	valueColor := color.New(color.FgHiWhite).Add(color.Bold)  // Bold white for field values
	optionColor := color.New(color.FgRed).Add(color.Bold)     // Bold red for options

	// Print the centered header
	dividerColor.Println(strings.Repeat("=", terminalWidth))
	header := fmt.Sprintf("Books %d to %d of %d", startIndex+1, endIndex, len(books))
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", (terminalWidth-len(header))/2), headerColor.Sprint(header), strings.Repeat(" ", (terminalWidth-len(header))/2))
	dividerColor.Println(strings.Repeat("=", terminalWidth))

	// Print each book
	for i := startIndex; i < endIndex; i++ {
		book := books[i]

		// Print book index
		indexColor.Printf("#%d\n", i+1)

		// Print book details with refined alignment and hide empty fields
		if book.Title != "" {
			titleColor.Printf("  Title:       ")
			valueColor.Printf("%s\n", book.Title)
		}
		if book.Authors != "" {
			labelColor.Print("  Author(s):   ")
			valueColor.Println(book.Authors)
		}
		if book.Year != "" {
			labelColor.Print("  Year:        ")
			valueColor.Println(book.Year)
		}
		if book.Publisher != "" {
			labelColor.Print("  Publisher:   ")
			valueColor.Println(book.Publisher)
		}
		if book.Pages != "" {
			labelColor.Print("  Pages:       ")
			valueColor.Println(book.Pages)
		}
		if book.Language != "" {
			labelColor.Print("  Language:    ")
			valueColor.Println(book.Language)
		}
		if book.Size != "" {
			labelColor.Print("  Size:        ")
			valueColor.Println(book.Size)
		}
		if book.Extension != "" {
			labelColor.Print("  Format:      ")
			valueColor.Println(book.Extension)
		}
		if len(book.ISBN) > 0 {
			labelColor.Print("  ISBN(s):     ")
			valueColor.Println(strings.Join(book.ISBN, ", "))
		}

		// Add a dashed divider between books
		dividerColor.Println(strings.Repeat("-", terminalWidth))
		fmt.Println() // Add extra vertical space for better readability
	}

	// Final divider
	dividerColor.Println(strings.Repeat("=", terminalWidth))

	// Print options at the end
	optionColor.Println("Options:")
	optionColor.Println("Enter the number of the book to select it.")
	optionColor.Println("Enter 'q' to Quit.")
}
