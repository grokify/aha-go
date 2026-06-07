package cmd

import (
	"github.com/spf13/cobra"
)

// requirementCmd is the parent command for requirement operations.
var requirementCmd = &cobra.Command{
	Use:     "requirement",
	Aliases: []string{"requirements", "req"},
	Short:   "Manage Aha requirements",
	Long: `Commands for managing Aha requirements.

Requirements are detailed specifications that belong to features.

Examples:
  aha requirement list --feature PROD-1
  aha requirement get PROD-1-1
  aha requirement create --feature PROD-1 --name "Add validation"
  aha requirement delete PROD-1-1`,
}

func init() {
	rootCmd.AddCommand(requirementCmd)
}
