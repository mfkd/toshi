package validate

import "testing"

func TestValidateDomain(t *testing.T) {
    if !ValidateDomain("books.xyz") { // length 9
        t.Fatal("expected valid domain to pass")
    }
    if ValidateDomain("too-long.example") {
        t.Fatal("expected invalid domain to fail")
    }
}

func TestBuildURL(t *testing.T) {
    got := BuildURL("books.xyz")
    want := "https://books.xyz/search.php"
    if got != want {
        t.Fatalf("BuildURL = %q, want %q", got, want)
    }
}
