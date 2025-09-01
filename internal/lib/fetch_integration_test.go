package lib

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/mfkd/toshi/internal/scraper"
)

func TestFetchPagesURLs_ParsesPagination(t *testing.T) {
    var base string
    // First page returns a <script> with three numbers ending in commas, first is total pages
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/search.php":
            fmt.Fprint(w, `<!doctype html><html><head></head><body><script>var total=3, other=2, x=1,</script></body></html>`)
        default:
            http.NotFound(w, r)
        }
    }))
    base = srv.URL
    t.Cleanup(srv.Close)

    s := scraper.NewScraper(base + "/search.php")
    urls, err := fetchPagesURLs(context.Background(), s, "foo bar")
    if err != nil {
        t.Fatalf("fetchPagesURLs error = %v", err)
    }
    if len(urls) != 3 {
        t.Fatalf("expected 3 urls, got %d", len(urls))
    }
    // spot check page=1 and page=3 exist
    if urls[0] == urls[2] || urls[0] == "" || urls[2] == "" {
        t.Fatalf("unexpected urls: %#v", urls)
    }
}

func TestFetchBooks_ParsesTable(t *testing.T) {
    html := `<!doctype html><table>
        <tr valign="top"><td>ID</td><td>Author</td><td>Title</td></tr>
        <tr valign="top">
          <td>123</td>
          <td>Doe; Smith</td>
          <td><a href="#">The Title 1234567890123</a></td>
          <td>Publisher Inc</td>
          <td>2024</td>
          <td>333</td>
          <td>EN</td>
          <td>1 MB</td>
          <td>epub</td>
          <td><a href="/m1">m1</a></td>
          <td><a href="/m2">m2</a></td>
        </tr>`

    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, html)
    }))
    t.Cleanup(srv.Close)

    s := scraper.NewScraper(srv.URL)
    books, err := fetchBooks(context.Background(), s, srv.URL)
    if err != nil {
        t.Fatalf("fetchBooks error = %v", err)
    }
    if len(books) != 1 {
        t.Fatalf("expected 1 book, got %d", len(books))
    }
    b := books[0]
    if b.ID != "123" || b.Extension != "epub" || len(b.ISBN) != 1 {
        t.Fatalf("unexpected book parsed: %#v", b)
    }
}

