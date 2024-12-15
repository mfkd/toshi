package ui

import (
	"fmt"
	"strconv"

	"github.com/mfkd/toshi/internal/libgen"

	"github.com/fatih/color"
)

const booksPerPage = 5

func SelectBook(books []libgen.Book) *libgen.Book {
	startIndex := 0

	for {
		displayBooksPaginated(books, startIndex)

		color.New(color.FgYellow).Printf("\nOptions:\n")
		fmt.Println("Enter the number of the book to select it.")
		if startIndex > 0 {
			color.New(color.FgMagenta).Println("Enter 'p' for Previous page.")
		}
		if startIndex+booksPerPage < len(books) {
			color.New(color.FgMagenta).Println("Enter 'n' for Next page.")
		}
		color.New(color.FgRed).Println("Enter 'q' to Quit.")
		fmt.Print("Your choice: ")

		var input string
		fmt.Scanln(&input)

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
			color.New(color.FgRed).Println("Invalid input. Please try again.")
		}
	}
}
