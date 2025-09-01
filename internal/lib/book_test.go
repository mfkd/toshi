package lib

import "testing"

func TestExtractTitleAndISBN(t *testing.T) {
    title, isbns := extractTitleAndISBN("The Great Book, 1234567890123 and 9876543210")
    // Note: function does not collapse spaces, it only removes digits and trims.
    if title != "The Great Book,  and" {
        t.Fatalf("unexpected title: %q", title)
    }
    if len(isbns) != 2 || isbns[0] != "1234567890123" || isbns[1] != "9876543210" {
        t.Fatalf("unexpected isbns: %#v", isbns)
    }
}

func TestSanitizeComponent(t *testing.T) {
    in := "  bad:name/with*chars?  "
    got := sanitizeComponent(in)
    want := "bad_name_with_chars_"
    if got != want {
        t.Fatalf("sanitizeComponent = %q, want %q", got, want)
    }
}

func TestGetFirstItem(t *testing.T) {
    if got := getFirstItem("First; Second; Third"); got != "First" {
        t.Fatalf("getFirstItem = %q, want %q", got, "First")
    }
    if got := getFirstItem(""); got != "" {
        t.Fatalf("getFirstItem(empty) = %q, want empty", got)
    }
}

func TestFileName(t *testing.T) {
    b := Book{
        Authors:   "Doe; Someone Else",
        Title:     "My: Title?",
        Publisher: "Pub/Inc; Another",
        Year:      "2021",
        Extension: "epub",
    }
    got := fileName(b)
    want := "Doe - My_ Title_ - Pub_Inc (2021).epub"
    if got != want {
        t.Fatalf("fileName = %q, want %q", got, want)
    }
}

func TestFilterEPUB(t *testing.T) {
    books := []Book{
        {Title: "A", Extension: "pdf"},
        {Title: "B", Extension: "epub"},
        {Title: "C", Extension: "mobi"},
        {Title: "D", Extension: "epub"},
    }
    got := filterEPUB(books)
    if len(got) != 2 || got[0].Title != "B" || got[1].Title != "D" {
        t.Fatalf("filterEPUB unexpected result: %#v", got)
    }
}
