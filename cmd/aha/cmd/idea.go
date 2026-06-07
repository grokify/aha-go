package cmd

import (
	"github.com/spf13/cobra"
)

// ideaCmd is the parent command for idea operations.
var ideaCmd = &cobra.Command{
	Use:     "idea",
	Aliases: []string{"ideas"},
	Short:   "Manage Aha ideas",
	Long: `Commands for managing Aha ideas.

Ideas are user-submitted feature requests that can be tracked,
voted on, and converted to features.

Examples:
  aha idea list
  aha idea get IDEA-123`,
}

func init() {
	rootCmd.AddCommand(ideaCmd)
}
