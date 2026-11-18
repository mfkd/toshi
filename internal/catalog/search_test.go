package catalog

import (
	"net/url"
	"testing"
)

func TestPageURL(t *testing.T) {
	tests := []struct {
		path string
		page int
		want string
	}{
		{"/search.php", 2, "https://books.xyz/search.php?column=def&page=2&phrase=1&req=The+Iliad+Homer&sort=def&sortmode=ASC&view=simple"},
		{"/index.php", 1, "https://books.xyz/index.php?req=The+Iliad+Homer&res=25"},
		{"/index.php", 2, "https://books.xyz/index.php?page=2&req=The+Iliad+Homer&res=25"},
	}
	for _, test := range tests {
		source := url.URL{Scheme: "https", Host: "books.xyz", Path: test.path}
		if got := pageURL(source, "The Iliad Homer", test.page); got != test.want {
			t.Errorf("pageURL(%s, %d) = %q, want %q", test.path, test.page, got, test.want)
		}
	}
}
