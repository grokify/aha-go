package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	page    int
	perPage int
)

var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List products (workspaces)",
	Long: `List all products (workspaces) you have access to.

Examples:
  aha product list
  aha product list --per-page 50
  aha product list --output json`,
	RunE: runListProducts,
}

func init() {
	productCmd.AddCommand(productListCmd)

	productListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	productListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListProducts(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	list, err := client.ListProducts(ctx)
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
	if len(list.Products) == 0 {
		fmt.Println("No products found.")
		return nil
	}

	fmt.Printf("Products (%d total):\n\n", list.Pagination.TotalRecords)
	for _, p := range list.Products {
		productLine := ""
		if p.ProductLine {
			productLine = " [product line]"
		}
		fmt.Printf("  %s  %s%s\n", p.ReferencePrefix, p.Name, productLine)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
