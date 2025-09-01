package cmd

import (
	"os"
	"testing"
)

func TestSelectURL(t *testing.T) {
	// Env takes precedence
	if got := selectURL("https://env.example/search.php", []string{"https://embed1", "https://embed2"}); got != "https://env.example/search.php" {
		t.Fatalf("selectURL returned %q, want env URL", got)
	}
	// Embedded list fallback to first
	if got := selectURL("", []string{"https://embed1", "https://embed2"}); got != "https://embed1" {
		t.Fatalf("selectURL returned %q, want first embed", got)
	}
	// No sources
	if got := selectURL("", []string{}); got != "" {
		t.Fatalf("selectURL returned %q, want empty", got)
	}
}

func TestParseEnv(t *testing.T) {
	// Save and restore env
	old := os.Getenv("DOMAIN")
	t.Cleanup(func() { _ = os.Setenv("DOMAIN", old) })

	// Invalid domain -> empty string
	_ = os.Setenv("DOMAIN", "invalid-domain")
	if got := parseEnv(); got != "" {
		t.Fatalf("parseEnv with invalid domain = %q, want empty", got)
	}

	// Valid domain -> built URL
	_ = os.Setenv("DOMAIN", "books.xyz")
	if got := parseEnv(); got != "https://books.xyz/search.php" {
		t.Fatalf("parseEnv = %q, want https://books.xyz/search.php", got)
	}
}

func TestParseArgs_Success(t *testing.T) {
	// Save and restore os.Args
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"toshi", "The", "Iliad", "Homer", "-v"}
	term, verbose := parseArgs()
	if term != "The Iliad Homer" {
		t.Fatalf("term = %q, want %q", term, "The Iliad Homer")
	}
	if !verbose {
		t.Fatalf("verbose = false, want true")
	}
}
