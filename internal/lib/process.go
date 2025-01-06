package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/scraper"
)

const defaultTimeout = 30 * time.Second

type UI interface {
	SelectBook(books []Book) *Book
}

// ProcessBooks handles the user selection, fetches download links, and attempts to download the selected book.
func ProcessBooks(s *scraper.Scraper, searchTerm string, ui UI) error {
	// Create a context with a timeout for fetching books
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	books, err := fetchAllBooks(ctx, s, searchTerm)
	if err != nil {
		return fmt.Errorf("error fetching books from pages: %w", err)
	}

	// Allow the user to select a book from the filtered list (e.g., EPUB books)
	// TODO: Make user select extension type as an argument
	selectedBook := ui.SelectBook(filterEPUB(books))
	if selectedBook == nil {
		fmt.Println("No book selected.")
		return nil
	}

	fmt.Printf("Selected Book: %s\n", selectedBook.Title)

	// Create a new context for download link operations
	ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Fetch download links for the selected book
	downloadLinks := fetchDownloadLinks(ctx, s, *selectedBook)

	fileName := fileName(*selectedBook)
	logger.Debugf("Attempting to download book to: %s\n", fileName)

	// Attempt to download the file
	if err := tryDownloadLinks(ctx, s, downloadLinks, fileName); err != nil {
		logger.Errorf("Failed to download file for book %s: %v", selectedBook.Title, err)
		return fmt.Errorf("failed to download book: %w", err)
	}

	fmt.Printf("Book downloaded successfully as %s\n", fileName)
	return nil
}
