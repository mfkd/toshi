package catalog

import (
	"net/url"
	"strconv"
	"strings"
)

func resolveLink(base *url.URL, href string) *url.URL {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "#") {
		return nil
	}
	target, err := base.Parse(href)
	if err != nil || (target.Scheme != "https" && target.Scheme != "http") || target.Host == "" {
		return nil
	}
	target.Fragment = ""
	return target
}

func pageURL(source url.URL, term string, page int) string {
	params := source.Query()
	params.Set("req", term)
	if source.Path == "/index.php" {
		params.Set("res", "25")
	} else {
		params.Set("phrase", "1")
		params.Set("view", "simple")
		params.Set("column", "def")
		params.Set("sort", "def")
		params.Set("sortmode", "ASC")
	}
	if page > 1 || source.Path != "/index.php" {
		params.Set("page", strconv.Itoa(page))
	}
	source.RawQuery = params.Encode()
	return source.String()
}
