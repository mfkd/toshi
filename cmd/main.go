package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/mfkd/toshi/internal/libgen"
	"github.com/mfkd/toshi/internal/logger"
	"github.com/mfkd/toshi/internal/ui"
)

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

func Execute() {
	c := libgen.SetupCollector()
	searchTerm, verbose := parseArgs()

	if verbose {
		logger.Configure(logger.LevelDebug, nil)
		fmt.Println("DEBUG mode: Detailed logs are now enabled")
	}

	if err := libgen.ProcessBooks(c, searchTerm, ui.CLI{}); err != nil {
		logger.Errorf("Error processing books: %v", err)
		os.Exit(1)
	}
}
