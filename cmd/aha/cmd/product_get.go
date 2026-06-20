package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var productGetCmd = &cobra.Command{
	Use:   "get <product-id>",
	Short: "Get a product by ID or reference prefix",
	Long: `Get details about a specific product (workspace).

Examples:
  aha product get PROD
  aha product get 12345678
  aha product get PROD --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runGetProduct,
}

func init() {
	productCmd.AddCommand(productGetCmd)
}

func runGetProduct(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	productID := args[0]

	product, err := client.GetProduct(ctx, productID)
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
		return output.NewPrinter(format).Print(product)
	}

	// Table output
	fmt.Printf("Product: %s\n", product.Name)
	fmt.Printf("  ID:              %s\n", product.ID)
	fmt.Printf("  Reference:       %s\n", product.ReferencePrefix)
	if product.ProductLine {
		fmt.Printf("  Type:            Product Line\n")
	} else {
		fmt.Printf("  Type:            Product\n")
	}
	if product.WorkspaceType != "" {
		fmt.Printf("  Workspace Type:  %s\n", product.WorkspaceType)
	}
	if product.ParentID != "" {
		fmt.Printf("  Parent ID:       %s\n", product.ParentID)
	}
	fmt.Printf("  URL:             %s\n", product.URL)
	fmt.Printf("  Created:         %s\n", product.CreatedAt.Format("2006-01-02"))
	if product.UpdatedAt != nil {
		fmt.Printf("  Updated:         %s\n", product.UpdatedAt.Format("2006-01-02"))
	}
	if product.HasIdeas {
		fmt.Printf("  Ideas:           Enabled\n")
	}
	if product.HasMasterFeatures {
		fmt.Printf("  Master Features: Enabled\n")
	}
	if product.Description != "" {
		fmt.Printf("  Description:     %s\n", truncateString(product.Description, 80))
	}

	return nil
}

// truncateString truncates a string to maxLen characters, adding "..." if truncated.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
