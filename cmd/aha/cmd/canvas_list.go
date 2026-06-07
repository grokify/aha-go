package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var listKind string

// canvasListCmd lists strategic canvases.
var canvasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List strategic canvases",
	Long: `List strategic canvases (strategic models) in Aha.io.

You can filter by product and/or canvas kind.

Examples:
  # List all canvases
  aha canvas list

  # List canvases for a specific product
  aha canvas list --product PROD

  # List only Opportunity canvases
  aha canvas list --kind Opportunity`,
	RunE: runListCanvases,
}

func init() {
	canvasCmd.AddCommand(canvasListCmd)

	canvasListCmd.Flags().StringVarP(&productID, "product", "p", "", "Product ID or reference prefix")
	canvasListCmd.Flags().StringVarP(&listKind, "kind", "k", "", "Filter by canvas kind (Opportunity, Lean Canvas, etc.)")
}

func runListCanvases(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var opts []aha.ListStrategicModelsOption
	if listKind != "" {
		opts = append(opts, aha.WithStrategicModelKind(listKind))
	}

	var list *aha.StrategicModelList
	var err error

	if productID != "" {
		list, err = client.ListProductStrategicModels(ctx, productID, opts...)
	} else {
		list, err = client.ListStrategicModels(ctx, opts...)
	}
	if err != nil {
		return fmt.Errorf("failed to list canvases: %w", err)
	}

	if len(list.StrategicModels) == 0 {
		fmt.Println("No canvases found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "ID\tREFERENCE\tNAME\tKIND\n")
	_, _ = fmt.Fprintf(w, "--\t---------\t----\t----\n")
	for _, sm := range list.StrategicModels {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", sm.ID, sm.ReferenceNum, sm.Name, sm.Kind)
	}
	_ = w.Flush()

	if list.Pagination.TotalRecords > 0 {
		fmt.Printf("\nShowing %d of %d canvases\n", len(list.StrategicModels), list.Pagination.TotalRecords)
	}

	return nil
}
