package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var requirementGetCmd = &cobra.Command{
	Use:   "get <requirement-id>",
	Short: "Get a requirement by ID or reference number",
	Long: `Get details about a specific requirement.

Examples:
  aha requirement get PROD-1-1
  aha requirement get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetRequirement,
}

func init() {
	requirementCmd.AddCommand(requirementGetCmd)
}

func runGetRequirement(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	requirementID := args[0]

	req, err := client.GetRequirement(ctx, requirementID)
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
		return output.NewPrinter(format).Print(req)
	}

	// Table output
	fmt.Printf("Requirement: %s\n", req.Name)
	fmt.Printf("  Reference:  %s\n", req.ReferenceNum)
	fmt.Printf("  ID:         %s\n", req.ID)

	if req.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", req.WorkflowStatus.Name)
	}
	if req.AssignedToUser != nil {
		fmt.Printf("  Assigned:   %s\n", req.AssignedToUser.Name())
	}

	if req.Feature != nil {
		fmt.Printf("  Feature:    %s (%s)\n", req.Feature.Name, req.Feature.ReferenceNum)
	}

	if req.OriginalEstimate > 0 {
		fmt.Printf("  Estimate:   %.1f\n", req.OriginalEstimate)
	}
	if req.RemainingEstimate > 0 {
		fmt.Printf("  Remaining:  %.1f\n", req.RemainingEstimate)
	}
	if req.WorkDone > 0 {
		fmt.Printf("  Work Done:  %.1f\n", req.WorkDone)
	}

	fmt.Printf("  URL:        %s\n", req.URL)
	fmt.Printf("  Created:    %s\n", req.CreatedAt.Format("2006-01-02"))

	if req.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", req.Description)
	}

	return nil
}
