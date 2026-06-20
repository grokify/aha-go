package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	productName          string
	productPrefix        string
	productDescription   string
	productParentID      string
	productWorkspaceType string
	productIsLine        bool
	productLineType      string
)

var productCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new product",
	Long: `Create a new product (workspace) or product line.

Examples:
  # Create a product
  aha product create --name "My Product" --prefix PROD

  # Create a product with description
  aha product create --name "My Product" --prefix PROD --description "Product description"

  # Create a product under a product line
  aha product create --name "My Product" --prefix PROD --parent PORT

  # Create a product line
  aha product create --name "My Portfolio" --prefix PORT --product-line --product-line-type portfolio

  # Create with workspace type
  aha product create --name "IT Project" --prefix ITP --workspace-type it_workspace`,
	RunE: runCreateProduct,
}

func init() {
	productCmd.AddCommand(productCreateCmd)

	productCreateCmd.Flags().StringVarP(&productName, "name", "n", "", "Product name (required)")
	productCreateCmd.Flags().StringVarP(&productPrefix, "prefix", "p", "", "Reference prefix (required)")
	productCreateCmd.Flags().StringVarP(&productDescription, "description", "d", "", "Product description (HTML allowed)")
	productCreateCmd.Flags().StringVar(&productParentID, "parent", "", "Parent product line ID or prefix")
	productCreateCmd.Flags().StringVar(&productWorkspaceType, "workspace-type", "", "Workspace type (product_workspace, it_workspace, marketing_workspace, etc.)")
	productCreateCmd.Flags().BoolVar(&productIsLine, "product-line", false, "Create as a product line")
	productCreateCmd.Flags().StringVar(&productLineType, "product-line-type", "", "Product line type (required if --product-line is set)")

	_ = productCreateCmd.MarkFlagRequired("name")
	_ = productCreateCmd.MarkFlagRequired("prefix")
}

func runCreateProduct(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var opts []aha.CreateProductOption

	if productDescription != "" {
		opts = append(opts, aha.WithProductDescription(productDescription))
	}
	if productParentID != "" {
		opts = append(opts, aha.WithProductParentID(productParentID))
	}
	if productWorkspaceType != "" {
		opts = append(opts, aha.WithProductWorkspaceType(productWorkspaceType))
	}

	var product *aha.Product
	var err error

	if productIsLine {
		if productLineType == "" {
			return fmt.Errorf("--product-line-type is required when creating a product line")
		}
		product, err = client.CreateProductLine(ctx, productName, productPrefix, productLineType, opts...)
	} else {
		product, err = client.CreateProduct(ctx, productName, productPrefix, opts...)
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
		return output.NewPrinter(format).Print(product)
	}

	// Human-readable output
	fmt.Printf("Created product: %s (%s)\n", product.Name, product.ReferencePrefix)
	fmt.Printf("  ID:   %s\n", product.ID)
	fmt.Printf("  URL:  %s\n", product.URL)

	return nil
}
