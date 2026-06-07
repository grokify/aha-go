package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var epicGetCmd = &cobra.Command{
	Use:   "get <epic-id>",
	Short: "Get an epic by ID or reference number",
	Long: `Get details about a specific epic.

Examples:
  aha epic get PROD-E-1
  aha epic get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetEpic,
}

func init() {
	epicCmd.AddCommand(epicGetCmd)
}

func runGetEpic(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	epicID := args[0]

	epic, err := client.GetEpic(ctx, epicID)
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
		return output.NewPrinter(format).Print(epic)
	}

	// Table output
	fmt.Printf("Epic: %s\n", epic.Name)
	fmt.Printf("  Reference:  %s\n", epic.ReferenceNum)
	fmt.Printf("  ID:         %s\n", epic.ID)

	if epic.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", epic.WorkflowStatus.Name)
	}
	fmt.Printf("  Progress:   %.0f%%\n", epic.Progress*100)

	if epic.Release != nil {
		fmt.Printf("  Release:    %s\n", epic.Release.Name)
	}
	if epic.Initiative != nil {
		fmt.Printf("  Initiative: %s (%s)\n", epic.Initiative.Name, epic.Initiative.ReferenceNum)
	}

	if len(epic.Tags) > 0 {
		fmt.Printf("  Tags:       %s\n", strings.Join(epic.Tags, ", "))
	}

	if epic.StartDate != nil {
		fmt.Printf("  Start:      %s\n", epic.StartDate.Format("2006-01-02"))
	}
	if epic.DueDate != nil {
		fmt.Printf("  Due:        %s\n", epic.DueDate.Format("2006-01-02"))
	}

	fmt.Printf("  URL:        %s\n", epic.URL)
	fmt.Printf("  Created:    %s\n", epic.CreatedAt.Format("2006-01-02"))

	if epic.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", epic.Description)
	}

	return nil
}
