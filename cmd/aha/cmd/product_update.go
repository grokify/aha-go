package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	updateProductName          string
	updateProductPrefix        string
	updateProductDescription   string
	updateProductParentID      string
	updateProductWorkspaceType string
)

var productUpdateCmd = &cobra.Command{
	Use:   "update <product-id>",
	Short: "Update a product",
	Long: `Update an existing product (workspace).

Examples:
  # Update product name
  aha product update PROD --name "New Product Name"

  # Update product description
  aha product update PROD --description "Updated description"

  # Move product to a different product line
  aha product update PROD --parent PORT

  # Update multiple fields
  aha product update PROD --name "New Name" --description "New description"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdateProduct,
}

func init() {
	productCmd.AddCommand(productUpdateCmd)

	productUpdateCmd.Flags().StringVarP(&updateProductName, "name", "n", "", "New product name")
	productUpdateCmd.Flags().StringVarP(&updateProductPrefix, "prefix", "p", "", "New reference prefix")
	productUpdateCmd.Flags().StringVarP(&updateProductDescription, "description", "d", "", "New product description (HTML allowed)")
	productUpdateCmd.Flags().StringVar(&updateProductParentID, "parent", "", "New parent product line ID or prefix")
	productUpdateCmd.Flags().StringVar(&updateProductWorkspaceType, "workspace-type", "", "New workspace type")
}

func runUpdateProduct(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	productID := args[0]

	var opts []aha.UpdateProductOption

	if updateProductName != "" {
		opts = append(opts, aha.WithUpdateProductName(updateProductName))
	}
	if updateProductPrefix != "" {
		opts = append(opts, aha.WithUpdateProductReferencePrefix(updateProductPrefix))
	}
	if updateProductDescription != "" {
		opts = append(opts, aha.WithUpdateProductDescription(updateProductDescription))
	}
	if updateProductParentID != "" {
		opts = append(opts, aha.WithUpdateProductParentID(updateProductParentID))
	}
	if updateProductWorkspaceType != "" {
		opts = append(opts, aha.WithUpdateProductWorkspaceType(updateProductWorkspaceType))
	}

	if len(opts) == 0 {
		return fmt.Errorf("at least one field must be specified to update")
	}

	product, err := client.UpdateProduct(ctx, productID, opts...)
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

	// Human-readable output
	fmt.Printf("Updated product: %s (%s)\n", product.Name, product.ReferencePrefix)
	fmt.Printf("  ID:   %s\n", product.ID)
	fmt.Printf("  URL:  %s\n", product.URL)

	return nil
}
