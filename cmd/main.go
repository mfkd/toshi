package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mfkd/toshi/internal/embed"
	"github.com/mfkd/toshi/internal/lib"
	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/scraper"
	"github.com/mfkd/toshi/internal/ui"
	"github.com/mfkd/toshi/internal/validate"
)

// parseArgs returns the search term and verbose flag
func parseArgs() (string, bool) {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsageAndExit()
	}

	var searchTerms []string
	verbose := false

	for _, arg := range args {
		if arg == "-v" {
			verbose = true
		} else if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Invalid flag detected: %s\n", arg)
			os.Exit(1)
		} else {
			searchTerms = append(searchTerms, arg)
		}
	}

	if len(searchTerms) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No search term provided.")
		os.Exit(1)
	}

	return strings.Join(searchTerms, " "), verbose
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, `Usage: toshi <searchterm> [options]
Example: toshi The Iliad Homer
Options:
  -v  Enable verbose output with debug logs
`)
	os.Exit(1)
}

// parseEnv returns the domain from the environment variable
func parseEnv() string {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return ""
	}

	if !validate.ValidateDomain(domain) {
		fmt.Printf("Invalid domain detected in environment variable: %s", domain)
		return ""
	}

	return validate.BuildURL(domain)
}

// selectURL returns the URL to use based on the environment variable or embedded URL.
// Environment variable takes precedence over embedded URLs
func selectURL(env string, embed []string) string {
	// TODO: Add support for multiple domains to improve reliability in the event of domain
	// resolution issues.
	if env != "" {
		return env
	}

	if len(embed) == 0 {
		return ""
	}

	return embed[0]
}

// Execute runs the CLI application
func Execute() {
	selected := selectURL(parseEnv(), embed.GetUrls())
	if selected == "" {
		fmt.Println("No valid domain found")
		fmt.Println("Please set the DOMAIN environment variable or add a valid domain to domains.txt")
		os.Exit(1)
	}

	s := scraper.NewScraper(selectURL(parseEnv(), embed.GetUrls()))

	searchTerm, verbose := parseArgs()

	if verbose {
		logger.Configure(logger.LevelDebug, nil)
		fmt.Println("DEBUG mode: Detailed logs are now enabled")
	}

	if err := lib.ProcessBooks(s, searchTerm, ui.CLI{}); err != nil {
		logger.Errorf("Error processing books: %v", err)
		os.Exit(1)
	}
}
