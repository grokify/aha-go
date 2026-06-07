package cmd

import (
	"github.com/spf13/cobra"
)

// epicCmd is the parent command for epic operations.
var epicCmd = &cobra.Command{
	Use:     "epic",
	Aliases: []string{"epics"},
	Short:   "Manage Aha epics",
	Long: `Commands for managing Aha epics.

Epics are large bodies of work that can be broken down into features.

Examples:
  aha epic list
  aha epic list --product PROD
  aha epic get PROD-E-1`,
}

func init() {
	rootCmd.AddCommand(epicCmd)
}
