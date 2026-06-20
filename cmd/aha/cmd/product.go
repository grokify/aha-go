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
  # List all products
  aha product list
  aha product list --with-idea-portals
  aha product list --updated-since 2024-01-01

  # Get a specific product
  aha product get PROD
  aha product get PROD --output json

  # Create a new product
  aha product create --name "My Product" --prefix PROD
  aha product create --name "Portfolio" --prefix PORT --product-line --product-line-type portfolio

  # Update a product
  aha product update PROD --name "New Name"
  aha product update PROD --description "Updated description"`,
}

func init() {
	rootCmd.AddCommand(productCmd)
}
