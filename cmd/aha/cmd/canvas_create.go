package cmd

import (
	"github.com/spf13/cobra"
)

// Shared flags for canvas creation
var (
	productID   string
	canvasName  string
	description string
)

// canvasCreateCmd represents the canvas create command.
var canvasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new strategic canvas",
	Long: `Create a new strategic canvas (strategic model) in Aha.io.

The canvas will be created in the specified product/workspace.
After creation, you can populate the canvas blocks via the Aha.io
web interface or by updating the canvas components via the API.

Examples:
  # Create an Opportunity Canvas
  aha canvas create opportunity --product PROD --name "Mobile App Opportunity"

  # Create a Lean UX Canvas with description
  aha canvas create leanux --product PROD --name "Checkout Redesign" \
    --description "Improve checkout conversion rates"`,
}

func init() {
	canvasCmd.AddCommand(canvasCreateCmd)

	// Common flags for all create subcommands
	canvasCreateCmd.PersistentFlags().StringVarP(&productID, "product", "p", "", "Product ID or reference prefix (required)")
	canvasCreateCmd.PersistentFlags().StringVarP(&canvasName, "name", "n", "", "Canvas name (required)")
	canvasCreateCmd.PersistentFlags().StringVarP(&description, "description", "d", "", "Canvas description")

	_ = canvasCreateCmd.MarkPersistentFlagRequired("product")
	_ = canvasCreateCmd.MarkPersistentFlagRequired("name")
}
