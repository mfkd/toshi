package scraper

import (
    "context"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
)

func TestScrape_Success(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = io.WriteString(w, "<html><body><div id='ok'>hello</div></body></html>")
    }))
    t.Cleanup(srv.Close)

    s := NewScraper(srv.URL)
    doc, err := s.Scrape(srv.URL)
    if err != nil {
        t.Fatalf("Scrape() error = %v", err)
    }
    if got := doc.Find("#ok").Text(); got != "hello" {
        t.Fatalf("unexpected content: %q", got)
    }
}

func TestScrape_StatusNotOK(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusTeapot)
    }))
    t.Cleanup(srv.Close)

    s := NewScraper(srv.URL)
    if _, err := s.Scrape(srv.URL); err == nil {
        t.Fatal("expected error for non-200 response")
    }
}

func TestCheckHead(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodHead {
            t.Fatalf("expected HEAD, got %s", r.Method)
        }
        w.WriteHeader(http.StatusNoContent)
    }))
    t.Cleanup(srv.Close)

    s := NewScraper(srv.URL)
    code, err := s.CheckHead(context.Background(), srv.URL)
    if err != nil {
        t.Fatalf("CheckHead() error = %v", err)
    }
    if code != http.StatusNoContent {
        t.Fatalf("status = %d, want %d", code, http.StatusNoContent)
    }
}

func TestDownloadFile_Success(t *testing.T) {
    const content = "test file contents"
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = io.WriteString(w, content)
    }))
    t.Cleanup(srv.Close)

    s := NewScraper(srv.URL)
    dir := t.TempDir()
    filename := "file.txt"
    if err := s.DownloadFile(context.Background(), filename, srv.URL, dir); err != nil {
        t.Fatalf("DownloadFile() error = %v", err)
    }
    b, err := os.ReadFile(filepath.Join(dir, filename))
    if err != nil {
        t.Fatalf("reading downloaded file failed: %v", err)
    }
    if string(b) != content {
        t.Fatalf("downloaded contents = %q, want %q", string(b), content)
    }
}

func TestDownloadFile_BadStatus(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadGateway)
    }))
    t.Cleanup(srv.Close)

    s := NewScraper(srv.URL)
    dir := t.TempDir()
    if err := s.DownloadFile(context.Background(), "x.txt", srv.URL, dir); err == nil {
        t.Fatal("expected error on non-200 status")
    }
}

