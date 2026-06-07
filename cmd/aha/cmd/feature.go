package cmd

import (
	"github.com/spf13/cobra"
)

// featureCmd is the parent command for feature operations.
var featureCmd = &cobra.Command{
	Use:     "feature",
	Aliases: []string{"features"},
	Short:   "Manage Aha features",
	Long: `Commands for managing Aha features.

Features are the primary work items in Aha.io that represent
product capabilities to be built.

Examples:
  aha feature list
  aha feature get PROD-123
  aha feature create --release REL-1 --name "New Feature"`,
}

func init() {
	rootCmd.AddCommand(featureCmd)
}
