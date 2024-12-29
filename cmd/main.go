package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/mfkd/toshi/internal/embed"
	"github.com/mfkd/toshi/internal/lib"
	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/ui"
	"github.com/mfkd/toshi/internal/validate"
)

// parseArgs returns the search term and verbose flag
func parseArgs() (string, bool) {
	// Define the verbose flag
	verbose := pflag.BoolP("verbose", "v", false, "Enable verbose output with debug logs")
	pflag.Parse()

	args := pflag.Args()

	// Ensure the positional argument <searchterm> is provided
	// TODO: Provide more informative output when user provides invalid input
	if len(args) < 1 {
		fmt.Println("Usage: toshi <searchterm>")
		fmt.Println("Example: toshi The Iliad Homer")
		fmt.Println("Flags:")
		fmt.Println("  -v  Enable verbose output with debug logs")
		os.Exit(1)
	}

	return strings.Join(args, " "), *verbose
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
	// TODO: Add support for multiple URLs to improve reliability in the event of server
	// outages.
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

	c := lib.SetupCollector(selectURL(parseEnv(), embed.GetUrls()))
	searchTerm, verbose := parseArgs()

	if verbose {
		logger.Configure(logger.LevelDebug, nil)
		fmt.Println("DEBUG mode: Detailed logs are now enabled")
	}

	if err := lib.ProcessBooks(c, searchTerm, ui.CLI{}); err != nil {
		logger.Errorf("Error processing books: %v", err)
		os.Exit(1)
	}
}
