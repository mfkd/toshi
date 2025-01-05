package lib

import (
	"fmt"
	"net/url"
)

// pageURL returns the URL for the given search term and page number.
func pageURL(searchBaseURL, term string, page int) string {
	// Parse the base URL
	baseURL, err := url.Parse(searchBaseURL)
	if err != nil {
		panic(fmt.Sprintf("invalid base URL: %s", searchBaseURL))
	}

	// Add query parameters
	params := url.Values{}
	params.Add("req", term)
	params.Add("phrase", "1")
	params.Add("view", "simple")
	params.Add("column", "def")
	params.Add("sort", "def")
	params.Add("sortmode", "ASC")
	params.Add("page", fmt.Sprintf("%d", page)) // Add the page number

	// Encode the query parameters and attach them to the base URL
	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
