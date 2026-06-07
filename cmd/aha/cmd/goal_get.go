package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var goalGetCmd = &cobra.Command{
	Use:   "get <goal-id>",
	Short: "Get a goal by ID or reference number",
	Long: `Get details about a specific goal.

Examples:
  aha goal get PROD-G-1
  aha goal get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetGoal,
}

func init() {
	goalCmd.AddCommand(goalGetCmd)
}

func runGetGoal(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	goalID := args[0]

	goal, err := client.GetGoal(ctx, goalID)
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
		return output.NewPrinter(format).Print(goal)
	}

	// Table output
	fmt.Printf("Goal: %s\n", goal.Name)
	fmt.Printf("  Reference:  %s\n", goal.ReferenceNum)
	fmt.Printf("  ID:         %s\n", goal.ID)

	if goal.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", goal.WorkflowStatus.Name)
	}
	if goal.TimeFrame != nil {
		fmt.Printf("  Time Frame: %s\n", goal.TimeFrame.Name)
	}
	fmt.Printf("  Progress:   %.0f%%\n", goal.Progress*100)

	if goal.StartDate != nil {
		fmt.Printf("  Start:      %s\n", goal.StartDate.Format("2006-01-02"))
	}
	if goal.EndDate != nil {
		fmt.Printf("  End:        %s\n", goal.EndDate.Format("2006-01-02"))
	}

	fmt.Printf("  URL:        %s\n", goal.URL)
	fmt.Printf("  Created:    %s\n", goal.CreatedAt.Format("2006-01-02"))

	if goal.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", goal.Description)
	}

	return nil
}
