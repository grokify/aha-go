package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var ideaGetCmd = &cobra.Command{
	Use:   "get <idea-id>",
	Short: "Get an idea by ID or reference number",
	Long: `Get details about a specific idea.

Examples:
  aha idea get IDEA-123
  aha idea get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetIdea,
}

func init() {
	ideaCmd.AddCommand(ideaGetCmd)
}

func runGetIdea(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ideaID := args[0]

	idea, err := client.GetIdea(ctx, ideaID)
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
		return output.NewPrinter(format).Print(idea)
	}

	// Table output
	fmt.Printf("Idea: %s\n", idea.Name)
	fmt.Printf("  Reference:  %s\n", idea.ReferenceNum)
	fmt.Printf("  ID:         %s\n", idea.ID)
	fmt.Printf("  Votes:      %d\n", idea.Votes)

	if idea.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", idea.WorkflowStatus.Name)
	}

	if len(idea.Categories) > 0 {
		fmt.Printf("  Categories: ")
		for i, cat := range idea.Categories {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", cat.Name)
		}
		fmt.Println()
	}

	if idea.Feature != nil {
		fmt.Printf("  Feature:    %s (%s)\n", idea.Feature.Name, idea.Feature.ReferenceNum)
	}

	fmt.Printf("  Created:    %s\n", idea.CreatedAt.Format("2006-01-02"))
	fmt.Printf("  Updated:    %s\n", idea.UpdatedAt.Format("2006-01-02"))

	if idea.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", idea.Description)
	}

	return nil
}
