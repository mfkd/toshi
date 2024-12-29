package lib

import (
	"fmt"
	"net/url"
)

// Base URL for the lib search
const libSearchBaseURL = "https://libgen.is/search.php"

// Helper function to construct the search URL on ALL Column
func defaultSearchURL(term string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "def")
	return fmt.Sprintf("%s?%s", libSearchBaseURL, params.Encode())
}
