package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"unicode"

	"github.com/mfkd/toshi/internal/catalog"
)

const usage = "Usage: toshi [-v] <search terms>\nSelect an EPUB by number or enter q to quit."

type command struct {
	in        io.Reader
	out       io.Writer
	errOut    io.Writer
	domain    string
	outputDir string
	http      *http.Client
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	cmd := command{
		in: os.Stdin, out: os.Stdout, errOut: os.Stderr,
		domain: os.Getenv("DOMAIN"), outputDir: "output", http: http.DefaultClient,
	}
	status := cmd.run(ctx, os.Args[1:])
	stop()
	os.Exit(status)
}

func (cmd command) run(ctx context.Context, args []string) int {
	query, verbose, err := parseArgs(args)
	if errors.Is(err, errHelp) {
		fmt.Fprintln(cmd.out, usage)
		return 0
	}
	if err != nil {
		fmt.Fprintln(cmd.errOut, err)
		fmt.Fprintln(cmd.errOut, usage)
		return 1
	}

	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(cmd.errOut, &slog.HandlerOptions{Level: level}))
	if err := cmd.searchAndDownload(ctx, query, logger); err != nil {
		fmt.Fprintln(cmd.errOut, err)
		if errors.Is(err, context.Canceled) {
			return 130
		}
		return 1
	}
	return 0
}

var errHelp = errors.New("help requested")

func parseArgs(args []string) (string, bool, error) {
	var words []string
	var verbose, positional bool
	for _, arg := range args {
		if positional {
			words = append(words, arg)
			continue
		}
		switch arg {
		case "--":
			positional = true
		case "-h", "--help":
			return "", false, errHelp
		case "-v":
			verbose = true
		default:
			if strings.HasPrefix(arg, "-") {
				return "", false, fmt.Errorf("unknown option %q", arg)
			}
			words = append(words, arg)
		}
	}
	query := strings.TrimSpace(strings.Join(words, " "))
	if query == "" {
		return "", false, errors.New("enter a search term")
	}
	return query, verbose, nil
}

func (cmd command) searchAndDownload(ctx context.Context, query string, logger *slog.Logger) error {
	sources, err := catalog.Sources(cmd.domain)
	if err != nil {
		return err
	}
	client := catalog.NewClient(cmd.http, logger)
	books, err := client.Search(ctx, sources, query)
	if err != nil {
		return err
	}
	if len(books) == 0 {
		fmt.Fprintf(cmd.out, "No EPUB results found for %q.\n", query)
		return nil
	}
	book, err := selectBook(ctx, cmd.in, cmd.out, books)
	if err != nil {
		return err
	}
	if book == nil {
		fmt.Fprintln(cmd.out, "No book selected.")
		return nil
	}

	fmt.Fprintf(cmd.out, "Selected book: %s\n", printable(book.Title))
	saved, err := client.Download(ctx, *book, cmd.outputDir)
	if err != nil {
		return fmt.Errorf("download %q: %w", book.Title, err)
	}
	fmt.Fprintf(cmd.out, "Book downloaded successfully to %s\n", printable(saved))
	return nil
}

func selectBook(ctx context.Context, in io.Reader, out io.Writer, books []catalog.Book) (*catalog.Book, error) {
	for i, book := range books {
		fmt.Fprintf(out, "%d. %s", i+1, printable(book.Title))
		for _, detail := range []string{book.Authors, book.Publisher, book.Year, book.Language, book.Pages, book.Size, strings.Join(book.ISBN, ", ")} {
			if detail = printable(detail); detail != "" {
				fmt.Fprintf(out, " | %s", detail)
			}
		}
		fmt.Fprintln(out)
	}
	scanner := bufio.NewScanner(in)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		fmt.Fprint(out, "Select an EPUB by number or enter q to quit: ")
		scanned := make(chan bool, 1)
		go func() { scanned <- scanner.Scan() }()
		var ok bool
		select {
		case <-ctx.Done():
			fmt.Fprintln(out)
			return nil, ctx.Err()
		case ok = <-scanned:
		}
		if !ok {
			fmt.Fprintln(out)
			return nil, scanner.Err()
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "q" {
			return nil, nil
		}
		n, err := strconv.Atoi(choice)
		if err == nil && n > 0 && n <= len(books) {
			return &books[n-1], nil
		}
		fmt.Fprintln(out, "Invalid selection.")
	}
}

func printable(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsControl(r) || unicode.IsSpace(r)
	}), " ")
}
