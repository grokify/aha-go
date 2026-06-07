package cmd

import (
	"github.com/spf13/cobra"
)

// goalCmd is the parent command for goal operations.
var goalCmd = &cobra.Command{
	Use:     "goal",
	Aliases: []string{"goals"},
	Short:   "Manage Aha goals",
	Long: `Commands for managing Aha goals.

Goals represent high-level objectives that guide product strategy.

Examples:
  aha goal list
  aha goal list --product PROD
  aha goal get PROD-G-1`,
}

func init() {
	rootCmd.AddCommand(goalCmd)
}
