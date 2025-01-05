package lib

import (
	"fmt"

	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/scraper"
)

type UI interface {
	SelectBook(books []Book) *Book
}

// ProcessBooks handles the user selection, fetches download links, and attempts to download the selected book.
func ProcessBooks(c Collector, s *scraper.Scraper, searchTerm string, ui UI) error {
	books, err := fetchAllBooks(s, searchTerm)
	if err != nil {
		return fmt.Errorf("error fetching books from ages: %w", err)
	}

	// Allow the user to select a book from the filtered list (e.g., EPUB books)
	// TODO: Make user select extension type as an argument
	selectedBook := ui.SelectBook(filterEPUB(books))
	if selectedBook == nil {
		fmt.Println("No book selected.")
		return nil
	}

	fmt.Printf("Selected Book: %s\n", selectedBook.Title)

	// Fetch download links for the selected book
	downloadLinks := fetchDownloadLinks(c, *selectedBook)

	// Attempt to download the file
	fileName := fileName(*selectedBook)
	logger.Debugf("Attempting to download book to: %s\n", fileName)

	// Attempt to download the file
	if err := tryDownloadLinks(c, downloadLinks, fileName); err != nil {
		logger.Errorf("Failed to download file for book %s: %v", selectedBook.Title, err)
		return fmt.Errorf("failed to download book: %w", err)
	}

	fmt.Printf("Book downloaded successfully as %s\n", fileName)
	return nil
}
