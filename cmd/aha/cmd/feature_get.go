package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var featureGetCmd = &cobra.Command{
	Use:   "get <feature-id>",
	Short: "Get a feature by ID or reference number",
	Long: `Get details about a specific feature.

Examples:
  aha feature get PROD-123
  aha feature get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetFeature,
}

func init() {
	featureCmd.AddCommand(featureGetCmd)
}

func runGetFeature(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	featureID := args[0]

	feature, err := client.GetFeature(ctx, featureID)
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
		return output.NewPrinter(format).Print(feature)
	}

	// Table output
	fmt.Printf("Feature: %s\n", feature.Name)
	fmt.Printf("  Reference:  %s\n", feature.ReferenceNum)
	fmt.Printf("  ID:         %s\n", feature.ID)

	if feature.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", feature.WorkflowStatus.Name)
	}
	if feature.AssignedTo != nil {
		fmt.Printf("  Assigned:   %s\n", feature.AssignedTo.Name())
	}
	if feature.Release != nil {
		fmt.Printf("  Release:    %s\n", feature.Release.Name)
	}
	if len(feature.Tags) > 0 {
		fmt.Printf("  Tags:       %s\n", strings.Join(feature.Tags, ", "))
	}
	if feature.StartDate != nil {
		fmt.Printf("  Start:      %s\n", feature.StartDate.Format("2006-01-02"))
	}
	if feature.DueDate != nil {
		fmt.Printf("  Due:        %s\n", feature.DueDate.Format("2006-01-02"))
	}
	fmt.Printf("  URL:        %s\n", feature.URL)
	fmt.Printf("  Created:    %s\n", feature.CreatedAt.Format("2006-01-02"))

	if feature.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", feature.Description)
	}

	return nil
}
