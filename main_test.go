package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfkd/toshi/internal/catalog"
)

type localTransport struct{ target *url.URL }

func (transport localTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := r.Clone(r.Context())
	req.URL.Scheme, req.URL.Host = transport.target.Scheme, transport.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func TestRunDownloadsSelectedEPUB(t *testing.T) {
	epub, err := os.ReadFile("internal/catalog/testdata/book.epub")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/index.php":
			if got := r.URL.Query().Get("req"); got != "fooled by randomness" {
				t.Errorf("query = %q", got)
			}
			if !strings.HasPrefix(r.UserAgent(), "Mozilla/") {
				t.Errorf("user agent = %q", r.UserAgent())
			}
			fmt.Fprint(w, `<table id="tablelibgen"><tr>
				<td><a href="edition.php?id=456">Fooled by Randomness</a><a href="edition.php?id=456">ISBN</a></td>
				<td>Taleb, Nassim Nicholas</td><td>Texere</td><td>2004</td><td>English</td>
				<td></td><td><a href="/file.php?id=123">444 kB</a></td><td>EPUB</td>
				<td><a href="/missing">1</a> <a href="/redirect">2</a></td></tr></table>`)
		case "/redirect":
			http.Redirect(w, r, "/mirror/details", http.StatusFound)
		case "/mirror/details":
			fmt.Fprint(w, `<a href="get.php?id=123&amp;key=sample">GET</a>`)
		case "/mirror/get.php":
			if !strings.HasSuffix(r.Referer(), "/mirror/details") || r.URL.Query().Get("key") != "sample" {
				http.Error(w, "missing download context", http.StatusForbidden)
				return
			}
			w.Write(epub)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	var out, errOut bytes.Buffer
	cmd := command{
		in: strings.NewReader("1\n"), out: &out, errOut: &errOut, outputDir: dir,
		http: &http.Client{Transport: localTransport{target: target}},
	}
	if status := cmd.run(t.Context(), []string{"fooled", "by", "randomness", "-v"}); status != 0 {
		t.Fatalf("exit status = %d, stderr = %s", status, &errOut)
	}
	saved := filepath.Join(dir, "Taleb, Nassim Nicholas - Fooled by Randomness - Texere (2004).epub")
	got, err := os.ReadFile(saved)
	if err != nil || !bytes.Equal(got, epub) {
		t.Fatalf("saved file differs from fixture: %v", err)
	}
	if !strings.Contains(out.String(), "Book downloaded successfully") || !strings.Contains(errOut.String(), "mirror failed") {
		t.Fatalf("stdout = %s, stderr = %s", &out, &errOut)
	}
}

func TestSelectBookCancellation(t *testing.T) {
	in, writer := io.Pipe()
	defer in.Close()
	defer writer.Close()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var got error
	done := make(chan struct{})
	go func() {
		_, got = selectBook(ctx, in, io.Discard, []catalog.Book{{Title: "Test"}})
		close(done)
	}()
	cancel()
	<-done
	if got != context.Canceled {
		t.Fatalf("selectBook error = %v, want cancellation", got)
	}
}
