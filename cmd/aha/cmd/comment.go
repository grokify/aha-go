package cmd

import (
	"github.com/spf13/cobra"
)

// commentCmd is the parent command for comment operations.
var commentCmd = &cobra.Command{
	Use:     "comment",
	Aliases: []string{"comments"},
	Short:   "Manage Aha comments",
	Long: `Commands for managing Aha comments.

Comments can be attached to features, ideas, releases, and other resources.

Examples:
  aha comment list --feature PROD-123
  aha comment list --idea IDEA-1
  aha comment get 12345678
  aha comment create --feature PROD-123 --body "Great progress!"`,
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
