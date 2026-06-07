package cmd

import (
	"github.com/spf13/cobra"
)

// releaseCmd is the parent command for release operations.
var releaseCmd = &cobra.Command{
	Use:     "release",
	Aliases: []string{"releases"},
	Short:   "Manage Aha releases",
	Long: `Commands for managing Aha releases.

Releases are time-based containers for features in Aha.io.

Examples:
  aha release list --product PROD
  aha release get PROD-R-1`,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
