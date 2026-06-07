//nolint:dupl // CLI list commands follow consistent pattern
package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var initiativeProductID string
var initiativeQuery string

var initiativeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List initiatives",
	Long: `List initiatives with optional filtering.

Examples:
  aha initiative list
  aha initiative list --product PROD
  aha initiative list --query "mobile"`,
	RunE: runListInitiatives,
}

func init() {
	initiativeCmd.AddCommand(initiativeListCmd)

	initiativeListCmd.Flags().StringVarP(&initiativeProductID, "product", "p", "", "Filter by product ID")
	initiativeListCmd.Flags().StringVarP(&initiativeQuery, "query", "q", "", "Search query")
	initiativeListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	initiativeListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListInitiatives(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var list *aha.InitiativeList
	var err error

	if initiativeProductID != "" {
		list, err = client.ListProductInitiatives(ctx, initiativeProductID)
	} else {
		var opts []aha.ListInitiativesOption
		if initiativeQuery != "" {
			opts = append(opts, aha.WithInitiativeQuery(initiativeQuery))
		}
		list, err = client.ListInitiatives(ctx, opts...)
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
	if len(list.Initiatives) == 0 {
		fmt.Println("No initiatives found.")
		return nil
	}

	fmt.Printf("Initiatives (%d total):\n\n", list.Pagination.TotalRecords)
	for _, i := range list.Initiatives {
		fmt.Printf("  %s  %s\n", i.ReferenceNum, i.Name)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
