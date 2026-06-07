package cmd

import (
	"github.com/spf13/cobra"
)

// productCmd is the parent command for product operations.
var productCmd = &cobra.Command{
	Use:     "product",
	Aliases: []string{"products", "workspace", "workspaces"},
	Short:   "Manage Aha products (workspaces)",
	Long: `Commands for managing Aha products (workspaces).

Products are the top-level containers in Aha.io that hold features,
releases, ideas, and other resources.

Examples:
  aha product list
  aha product get PROD`,
}

func init() {
	rootCmd.AddCommand(productCmd)
}
