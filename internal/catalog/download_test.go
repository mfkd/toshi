package catalog

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDownloadLinksFiltersCandidates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `<div id="download"><a href="file.pdf">PDF</a><a href="file.epub">EPUB</a></div>`)
	}))
	defer server.Close()
	client := NewClient(server.Client(), slog.New(slog.DiscardHandler))
	links, err := client.downloadLinks(t.Context(), server.URL+"/mirror")
	if err != nil {
		t.Fatal(err)
	}
	want := downloadLink{url: server.URL + "/file.epub", referer: server.URL + "/mirror"}
	if len(links) != 1 || links[0] != want {
		t.Fatalf("downloadLinks = %v, want [%v]", links, want)
	}
}

func TestTryLinksRecoversAfterTimeoutAndPreservesExistingFile(t *testing.T) {
	epub, err := os.ReadFile("testdata/book.epub")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slow":
			<-r.Context().Done()
		case "/bad":
			fmt.Fprint(w, "<html>Unavailable</html>")
		case "/ok":
			w.Write(epub)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	first := filepath.Join(dir, "book.epub")
	if err := os.WriteFile(first, []byte("original contents"), 0600); err != nil {
		t.Fatal(err)
	}
	client := NewClient(server.Client(), slog.New(slog.DiscardHandler))
	links := []downloadLink{{url: server.URL + "/slow"}, {url: server.URL + "/bad"}, {url: server.URL + "/ok"}}
	saved, err := client.tryLinks(t.Context(), links, dir, "book.epub", 300*time.Millisecond)
	if err != nil || saved != filepath.Join(dir, "book (2).epub") {
		t.Fatalf("tryLinks = %q, %v", saved, err)
	}
	for name, want := range map[string][]byte{first: []byte("original contents"), saved: epub} {
		got, err := os.ReadFile(name)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("unexpected file %s: %v", name, err)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 2 {
		t.Fatalf("output entries = %v, error = %v", entries, err)
	}
}

func TestValidateEPUBRejectsMimetypeAlone(t *testing.T) {
	var buf bytes.Buffer
	archive := zip.NewWriter(&buf)
	entry, err := archive.CreateHeader(&zip.FileHeader{Name: "mimetype", Method: zip.Store})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprint(entry, "application/epub+zip"); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := validateEPUB(bytes.NewReader(buf.Bytes()), int64(buf.Len())); err == nil || !strings.Contains(err.Error(), "container") {
		t.Fatalf("validateEPUB error = %v, want missing container error", err)
	}
}
