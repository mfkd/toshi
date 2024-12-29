package embed

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed urls/urls.txt
var urls string

func validateURL(url string) bool {
	return len(url) == 28
}

func GetUrls() []string {
	URLList := strings.Split(strings.TrimSpace(urls), "\n")
	for _, url := range URLList {
		if !validateURL(url) {
			fmt.Printf("Invalid URL detected in the urls.txt: %s", url)
			os.Exit(1)
		}
	}
	return URLList
}
