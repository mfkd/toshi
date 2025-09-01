package lib

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"

    "github.com/mfkd/toshi/internal/scraper"
)

func TestFetchDownloadLinks_FiltersByExtension(t *testing.T) {
    // Serve a mirror page with multiple links
    var serverURL string
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/mirror":
            fmt.Fprintf(w, `<!doctype html><div id="download"><ul>
                <li><a href="%s/file.pdf">PDF</a></li>
                <li><a href="%s/file.epub">EPUB</a></li>
                <li><a href="%s/file.mobi">MOBI</a></li>
            </ul></div>`, serverURL, serverURL, serverURL)
        default:
            http.NotFound(w, r)
        }
    }))
    serverURL = srv.URL
    t.Cleanup(srv.Close)

    s := scraper.NewScraper(serverURL)
    b := Book{Extension: "epub", Mirrors: []string{serverURL + "/mirror"}}
    links, err := fetchDownloadLinks(context.Background(), s, b)
    if err != nil {
        t.Fatalf("fetchDownloadLinks error = %v", err)
    }
    if len(links) != 1 || !strings.HasSuffix(links[0], "/file.epub") {
        t.Fatalf("unexpected links: %#v", links)
    }
}

func TestTryDownloadLinks_TriesUntilSuccess(t *testing.T) {
    // First URL fails, second succeeds
    hits := 0
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        hits++
        if strings.Contains(r.URL.Path, "fail") {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    }))
    t.Cleanup(srv.Close)

    s := scraper.NewScraper(srv.URL)

    // Isolate writes in a temp working directory so we don't touch the repo
    prevWD, err := os.Getwd()
    if err != nil { t.Fatalf("Getwd error: %v", err) }
    tmp := t.TempDir()
    if err := os.Chdir(tmp); err != nil { t.Fatalf("chdir temp error: %v", err) }
    t.Cleanup(func() { _ = os.Chdir(prevWD) })

    links := []string{srv.URL + "/fail", srv.URL + "/ok"}
    if err := tryDownloadLinks(context.Background(), s, links, "test.epub"); err != nil {
        t.Fatalf("tryDownloadLinks error = %v", err)
    }
    if hits < 2 {
        t.Fatalf("expected to try both links, hits=%d", hits)
    }
}
