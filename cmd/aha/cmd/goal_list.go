//nolint:dupl // CLI list commands follow consistent pattern
package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var goalProductID string
var goalQuery string

var goalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List goals",
	Long: `List goals with optional filtering.

Examples:
  aha goal list
  aha goal list --product PROD
  aha goal list --query "revenue"
  aha goal list --output json`,
	RunE: runListGoals,
}

func init() {
	goalCmd.AddCommand(goalListCmd)

	goalListCmd.Flags().StringVarP(&goalProductID, "product", "p", "", "Filter by product ID")
	goalListCmd.Flags().StringVarP(&goalQuery, "query", "q", "", "Search query")
	goalListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	goalListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListGoals(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var list *aha.GoalList
	var err error

	if goalProductID != "" {
		list, err = client.ListProductGoals(ctx, goalProductID)
	} else {
		var opts []aha.ListGoalsOption
		if goalQuery != "" {
			opts = append(opts, aha.WithGoalQuery(goalQuery))
		}
		list, err = client.ListGoals(ctx, opts...)
	}

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
	if len(list.Goals) == 0 {
		fmt.Println("No goals found.")
		return nil
	}

	fmt.Printf("Goals (%d total):\n\n", list.Pagination.TotalRecords)
	for _, g := range list.Goals {
		fmt.Printf("  %s  %s\n", g.ReferenceNum, g.Name)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
