package embed

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

const (
	scheme = "https://"
	path   = "/search.php"
)

//go:embed domains/domains.txt
var domains string

// validateDomain checks if the domain is valid
func validateDomain(domain string) bool {
	return len(domain) == 9
}

// buildURL constructs a URL from a domain
func buildURL(domain string) string {
	return fmt.Sprintf("%s%s%s", scheme, domain, path)
}

// GetUrls returns a list of URLs from domains.txt
func GetUrls() []string {
	if domains == "" {
		fmt.Println("No domains found in domains.txt")
		return []string{}
	}

	urlList := make([]string, 0)

	domainList := strings.Split(strings.TrimSpace(domains), "\n")
	for _, domain := range domainList {
		domain = strings.TrimSpace(domain)
		if !validateDomain(domain) {
			fmt.Printf("Invalid domain detected in domains.txt: %s", domain)
			os.Exit(1)
		}
		urlList = append(urlList, buildURL(domain))
	}

	return urlList
}
