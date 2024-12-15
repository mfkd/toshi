package libgen

import (
	"fmt"
	"net/url"
)

// Base URL for the LibGen search
const libgenSearchBaseURL = "https://libgen.is/search.php"

// Helper function to construct the search URL on TITLE Column
func titleSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "title")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}

// Helper function to construct the search URL on ALL Column
func defaultSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "def")
	return fmt.Sprintf("%s?%s", libgenSearchBaseURL, params.Encode())
}
