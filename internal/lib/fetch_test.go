package lib

import "testing"

func TestBuildPageURLs(t *testing.T) {
    urls := buildPageURLs("https://books.xyz/search.php", "foo", 3)
    if len(urls) != 3 {
        t.Fatalf("expected 3 urls, got %d", len(urls))
    }
    want1 := "https://books.xyz/search.php?column=def&page=1&phrase=1&req=foo&sort=def&sortmode=ASC&view=simple"
    want3 := "https://books.xyz/search.php?column=def&page=3&phrase=1&req=foo&sort=def&sortmode=ASC&view=simple"
    if urls[0] != want1 || urls[2] != want3 {
        t.Fatalf("unexpected urls: %#v", urls)
    }
}

func TestTotalPages(t *testing.T) {
    // Three matches with trailing commas, first one is the total (implementation detail)
    content := "var a=50, b=10, c=5,"
    got, err := totalPages(content)
    if err != nil {
        t.Fatalf("totalPages error: %v", err)
    }
    if got != 50 {
        t.Fatalf("totalPages = %d, want 50", got)
    }
}

func TestTotalPages_Errors(t *testing.T) {
    cases := []struct{
        name    string
        content string
    }{
        {"no-matches", "var x=;"},
        {"one-match", "10,"},
        {"two-matches", "1, 2,"},
        {"four-matches", "1,2,3,4,"},
    }
    for _, tc := range cases {
        if _, err := totalPages(tc.content); err == nil {
            t.Fatalf("%s: expected error, got nil", tc.name)
        }
    }
}
