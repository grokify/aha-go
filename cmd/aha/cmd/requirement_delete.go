//nolint:dupl // CLI delete commands follow consistent pattern
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var requirementDeleteForce bool

var requirementDeleteCmd = &cobra.Command{
	Use:   "delete <requirement-id>",
	Short: "Delete a requirement",
	Long: `Delete a requirement by ID or reference number.

Examples:
  aha requirement delete PROD-1-1
  aha requirement delete PROD-1-1 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDeleteRequirement,
}

func init() {
	requirementCmd.AddCommand(requirementDeleteCmd)

	requirementDeleteCmd.Flags().BoolVar(&requirementDeleteForce, "force", false, "Skip confirmation prompt")
}

func runDeleteRequirement(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	requirementID := args[0]

	if !requirementDeleteForce {
		fmt.Printf("Are you sure you want to delete requirement %s? [y/N]: ", requirementID)
		var confirm string
		_, err := fmt.Scanln(&confirm)
		if err != nil || (confirm != "y" && confirm != "Y") {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := client.DeleteRequirement(ctx, requirementID); err != nil {
		return err
	}

	fmt.Printf("Deleted requirement: %s\n", requirementID)
	return nil
}
