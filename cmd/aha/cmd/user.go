package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd is the parent command for user operations.
var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage Aha users",
	Long: `Commands for managing Aha users.

Examples:
  aha user list
  aha user get user@example.com
  aha user me`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
