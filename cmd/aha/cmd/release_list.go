package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var releaseProductID string

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List releases for a product",
	Long: `List all releases for a specific product.

Examples:
  aha release list --product PROD
  aha release list --product PROD --per-page 50`,
	RunE: runListReleases,
}

func init() {
	releaseCmd.AddCommand(releaseListCmd)

	releaseListCmd.Flags().StringVarP(&releaseProductID, "product", "p", "", "Product ID (required)")
	releaseListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	releaseListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")

	_ = releaseListCmd.MarkFlagRequired("product")
}

func runListReleases(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	list, err := client.ListProductReleases(ctx, releaseProductID)
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
	if len(list.Releases) == 0 {
		fmt.Println("No releases found.")
		return nil
	}

	fmt.Printf("Releases (%d total):\n\n", list.Pagination.TotalRecords)
	for _, r := range list.Releases {
		status := ""
		if r.Released {
			status = " [released]"
		} else if r.ParkingLot {
			status = " [parking lot]"
		}

		dateInfo := ""
		if r.ReleaseDate != nil {
			dateInfo = fmt.Sprintf(" (%s)", r.ReleaseDate.Format("2006-01-02"))
		}

		fmt.Printf("  %s  %s%s%s\n", r.ReferenceNum, r.Name, dateInfo, status)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
