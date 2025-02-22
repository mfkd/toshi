package embed

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/mfkd/toshi/internal/validate"
)

//go:embed domains/*
var embeddedFiles embed.FS

// GetUrls returns a list of URLs from domains.txt
func GetUrls() []string {
	data, err := embeddedFiles.ReadFile("domains/domains.txt")
	if err != nil {
		return []string{}
	}

	domains := string(data)

	if domains == "" {
		return []string{}
	}

	urlList := make([]string, 0)

	domainList := strings.Split(strings.TrimSpace(domains), "\n")
	for _, domain := range domainList {
		domain = strings.TrimSpace(domain)
		if !validate.ValidateDomain(domain) {
			fmt.Printf("Invalid domain detected in domains.txt: %s", domain)
			os.Exit(1)
		}
		urlList = append(urlList, validate.BuildURL(domain))
	}

	return urlList
}
