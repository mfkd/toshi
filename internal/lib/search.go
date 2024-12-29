package lib

import (
	"fmt"
	"net/url"
)

// Helper function to construct the search URL on ALL Column
func defaultSearchURL(term, searchBaseURL string) string {
	params := url.Values{}
	params.Add("req", term)
	params.Add("column", "def")
	return fmt.Sprintf("%s?%s", searchBaseURL, params.Encode())
}
