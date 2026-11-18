package catalog

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseBooks(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<!doctype html><table>
		<tr><td>ID</td><td>Author</td><td>Title</td></tr>
		<tr>
			<td> 123 </td>
			<td>Doe; Smith</td>
			<td><a href="#">The Title 1234567890123</a></td>
			<td>Publisher Inc</td>
			<td>2024</td>
			<td>333</td>
			<td>EN</td>
			<td>1 MB</td>
			<td> EPUB </td>
			<td></td>
		</tr></table>`))
	if err != nil {
		t.Fatal(err)
	}
	books := parseBooks(doc)
	if len(books) != 1 {
		t.Fatalf("expected one book, got %d", len(books))
	}
	b := books[0]
	if b.Title != "The Title" || b.extension != "epub" || len(b.ISBN) != 1 {
		t.Fatalf("unexpected book parsed: %#v", b)
	}
}
