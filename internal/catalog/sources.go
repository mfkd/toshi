package catalog

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

const defaultDomain = "libgen.bz"

var domainPattern = regexp.MustCompile(`(?i)^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

// Sources returns the default catalog or the ordered hostnames from DOMAIN.
func Sources(domain string) ([]*url.URL, error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		domain = defaultDomain
	}
	var sources []*url.URL
	var seen []string
	for _, name := range strings.Split(domain, ",") {
		host := strings.ToLower(strings.TrimSpace(name))
		if len(host) > 253 || !domainPattern.MatchString(host) {
			return nil, fmt.Errorf("invalid DOMAIN hostname %q", name)
		}
		if slices.Contains(seen, host) {
			continue
		}
		seen = append(seen, host)
		path := "/search.php"
		switch host {
		case "libgen.bz", "libgen.gl", "libgen.la", "libgen.li", "libgen.vg":
			path = "/index.php"
		}
		sources = append(sources, &url.URL{Scheme: "https", Host: host, Path: path})
	}
	return sources, nil
}
