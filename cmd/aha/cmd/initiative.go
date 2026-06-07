package cmd

import (
	"github.com/spf13/cobra"
)

// initiativeCmd is the parent command for initiative operations.
var initiativeCmd = &cobra.Command{
	Use:     "initiative",
	Aliases: []string{"initiatives"},
	Short:   "Manage Aha initiatives",
	Long: `Commands for managing Aha initiatives.

Initiatives are strategic themes that group related features together.

Examples:
  aha initiative list
  aha initiative list --product PROD
  aha initiative get PROD-I-1`,
}

func init() {
	rootCmd.AddCommand(initiativeCmd)
}
