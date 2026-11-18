package catalog

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const legacyResult = `<table><tr><td>1</td><td>Author</td><td><a>Available EPUB</a></td><td></td><td></td><td></td><td></td><td></td><td>epub</td><td><a href="/m">1</a></td></tr></table>`

func TestSearchKeepsResultsWhenLaterPageFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "1":
			fmt.Fprint(w, legacyResult, `<a href="?page=2">Next</a>`)
		case "2":
			http.Error(w, "Unavailable", http.StatusServiceUnavailable)
		default:
			t.Errorf("unexpected search page: %s", r.URL)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	source, err := url.Parse(server.URL + "/search.php")
	if err != nil {
		t.Fatal(err)
	}
	var logs bytes.Buffer
	client := NewClient(server.Client(), slog.New(slog.NewTextHandler(&logs, nil)))
	books, err := client.Search(t.Context(), []*url.URL{source}, "book")
	if err != nil || len(books) != 1 || books[0].Title != "Available EPUB" {
		t.Fatalf("Search = %v, %v", books, err)
	}
	if !strings.Contains(logs.String(), "showing partial search results") {
		t.Fatalf("missing partial results warning: %s", &logs)
	}
}

func TestSearchTriesNextSourceAfterEmptyOrInvalidPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/invalid":
			fmt.Fprint(w, "<title>Welcome to nginx!</title>")
		case "/empty":
			fmt.Fprint(w, `<input name="req">`)
		case "/search.php":
			fmt.Fprint(w, legacyResult)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	var sources []*url.URL
	for _, path := range []string{"/invalid", "/empty", "/search.php"} {
		source, err := url.Parse(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		sources = append(sources, source)
	}
	books, err := NewClient(server.Client(), slog.New(slog.DiscardHandler)).Search(t.Context(), sources, "book")
	if err != nil || len(books) != 1 || books[0].Title != "Available EPUB" {
		t.Fatalf("Search = %v, %v", books, err)
	}
}
