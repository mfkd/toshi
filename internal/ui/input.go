package ui

import (
	"fmt"
	"strconv"

	"github.com/mfkd/toshi/internal/lib"
)

const booksPerPage = 5

type CLI struct{}

func (CLI) SelectBook(books []lib.Book) *lib.Book {
	startIndex := 0

	for {
		displayBooksPaginated(books, startIndex)

		// Print options
		fmt.Printf("\n%sOptions:%s\n", FgYellow, Reset)
		fmt.Println("Enter the number of the book to select it.")
		if startIndex > 0 {
			fmt.Printf("%sEnter 'p' for Previous page.%s\n", FgMagenta, Reset)
		}
		if startIndex+booksPerPage < len(books) {
			fmt.Printf("%sEnter 'n' for Next page.%s\n", FgMagenta, Reset)
		}
		fmt.Printf("%sEnter 'q' to Quit.%s\n", FgRed, Reset)
		fmt.Print("Your choice: ")

		// Read user input
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Printf("%sError reading input. Please try again.%s\n", FgRed, Reset)
			continue
		}

		// Handle input
		if input == "n" && startIndex+booksPerPage < len(books) {
			startIndex += booksPerPage
		} else if input == "p" && startIndex > 0 {
			startIndex -= booksPerPage
		} else if input == "q" {
			return nil
		} else {
			selection, err := strconv.Atoi(input)
			if err == nil && selection > 0 && selection <= len(books) {
				return &books[selection-1]
			}
			fmt.Printf("%sInvalid input. Please try again.%s\n", FgRed, Reset)
		}
	}
}
