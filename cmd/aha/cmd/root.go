package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	aha "github.com/grokify/aha-go"
)

var (
	// Global flags
	apiKey       string
	subdomain    string
	verbose      bool
	outputFormat string

	// Shared client
	client *aha.Client
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "aha",
	Short: "CLI for Aha.io product management",
	Long: `aha is a command-line interface for interacting with Aha.io.

It provides commands for managing strategic canvases, features,
initiatives, and other Aha.io resources.

Environment variables:
  AHA_API_KEY    - Your Aha.io API key (required)
  AHA_SUBDOMAIN  - Your Aha.io subdomain (required)

Example:
  export AHA_API_KEY=your-api-key
  export AHA_SUBDOMAIN=your-company
  aha canvas create opportunity --product PROD --name "My Canvas"`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for help commands
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		// Get credentials from flags or environment
		if apiKey == "" {
			apiKey = os.Getenv("AHA_API_KEY")
		}
		if subdomain == "" {
			subdomain = os.Getenv("AHA_SUBDOMAIN")
		}

		// Build options slice
		var opts []aha.Option
		if subdomain != "" {
			opts = append(opts, aha.WithSubdomain(subdomain))
		}
		if apiKey != "" {
			opts = append(opts, aha.WithAPIKey(apiKey))
		}

		// Initialize client
		var err error
		client, err = aha.NewClient(opts...)
		if err != nil {
			return fmt.Errorf("failed to create Aha client: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Aha.io API key (or set AHA_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&subdomain, "subdomain", "", "Aha.io subdomain (or set AHA_SUBDOMAIN)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
}
