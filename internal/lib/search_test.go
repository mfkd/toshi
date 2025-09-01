package lib

import "testing"

func TestPageURL_PanicsOnInvalidBase(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Fatal("expected panic for invalid base URL")
        }
    }()
    _ = pageURL("://bad url", "term", 1)
}

func TestPageURL(t *testing.T) {
    base := "https://books.xyz/search.php"
    term := "The Iliad Homer"
    got := pageURL(base, term, 2)
    want := "https://books.xyz/search.php?column=def&page=2&phrase=1&req=The+Iliad+Homer&sort=def&sortmode=ASC&view=simple"
    if got != want {
        t.Fatalf("pageURL = %q, want %q", got, want)
    }
}
