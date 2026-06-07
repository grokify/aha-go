package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var epicProductID string
var epicQuery string

var epicListCmd = &cobra.Command{
	Use:   "list",
	Short: "List epics",
	Long: `List epics with optional filtering.

Examples:
  aha epic list
  aha epic list --product PROD
  aha epic list --query "authentication"`,
	RunE: runListEpics,
}

func init() {
	epicCmd.AddCommand(epicListCmd)

	epicListCmd.Flags().StringVarP(&epicProductID, "product", "p", "", "Filter by product ID")
	epicListCmd.Flags().StringVarP(&epicQuery, "query", "q", "", "Search query")
	epicListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	epicListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListEpics(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var list *aha.EpicList
	var err error

	if epicProductID != "" {
		list, err = client.ListProductEpics(ctx, epicProductID)
	} else {
		var opts []aha.ListEpicsOption
		if epicQuery != "" {
			opts = append(opts, aha.WithEpicQuery(epicQuery))
		}
		list, err = client.ListEpics(ctx, opts...)
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
	if len(list.Epics) == 0 {
		fmt.Println("No epics found.")
		return nil
	}

	fmt.Printf("Epics (%d total):\n\n", list.Pagination.TotalRecords)
	for _, e := range list.Epics {
		fmt.Printf("  %s  %s\n", e.ReferenceNum, e.Name)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
