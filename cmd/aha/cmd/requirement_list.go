package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var requirementFeatureID string

var requirementListCmd = &cobra.Command{
	Use:   "list",
	Short: "List requirements for a feature",
	Long: `List requirements belonging to a feature.

Examples:
  aha requirement list --feature PROD-1
  aha requirement list -f PROD-1`,
	RunE: runListRequirements,
}

func init() {
	requirementCmd.AddCommand(requirementListCmd)

	requirementListCmd.Flags().StringVarP(&requirementFeatureID, "feature", "f", "", "Feature ID or reference number (required)")
	requirementListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	requirementListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
	_ = requirementListCmd.MarkFlagRequired("feature")
}

func runListRequirements(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	list, err := client.ListFeatureRequirements(ctx, requirementFeatureID)
	if err != nil {
		return err
	}

	// Parse output format
	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	// Handle structured output
	if format.IsStructured() {
		return output.NewPrinter(format).Print(list)
	}

	// Table output
	if len(list.Requirements) == 0 {
		fmt.Println("No requirements found.")
		return nil
	}

	fmt.Printf("Requirements for feature %s (%d total):\n\n", requirementFeatureID, list.Pagination.TotalRecords)
	for _, r := range list.Requirements {
		fmt.Printf("  %s  %s\n", r.ReferenceNum, r.Name)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
