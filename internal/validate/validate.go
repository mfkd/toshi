package validate

import (
	"fmt"
)

const (
	scheme = "https://"
	path   = "/search.php"
)

// ValidateDomain checks if the domain is valid
func ValidateDomain(domain string) bool {
	return len(domain) == 9
}

// buildURL constructs a URL from a domain
func BuildURL(domain string) string {
	return fmt.Sprintf("%s%s%s", scheme, domain, path)
}
