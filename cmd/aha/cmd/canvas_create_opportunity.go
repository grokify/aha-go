package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/canvas"
	"github.com/spf13/cobra"
)

// canvasCreateOpportunityCmd creates an Opportunity Canvas.
var canvasCreateOpportunityCmd = &cobra.Command{
	Use:   "opportunity",
	Short: "Create an Opportunity Canvas",
	Long: `Create a new Opportunity Canvas in Aha.io.

The Opportunity Canvas follows Jeff Patton's 10-block structure for
evaluating product opportunities before committing resources.

Grid Layout:
  +-------------------+-------------------+-------------------+
  | Users & Customers | Problems          | Solution Ideas    |
  +-------------------+-------------------+-------------------+
  | Solutions Today   | User Value        | Adoption Strategy |
  +-------------------+-------------------+-------------------+
  | User Metrics      | Business Problem  | Business Metrics  |
  +-------------------+-------------------+-------------------+
  | Budget (full width)                                       |
  +-----------------------------------------------------------+

Example:
  aha canvas create opportunity \
    --product PROD \
    --name "Mobile App Opportunity" \
    --description "Evaluate mobile app expansion"`,
	RunE: runCreateOpportunityCanvas,
}

func init() {
	canvasCreateCmd.AddCommand(canvasCreateOpportunityCmd)
}

func runCreateOpportunityCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sm, err := canvas.Create(ctx, client, canvas.CreateOptions{
		ProductID:   productID,
		Name:        canvasName,
		Description: description,
		Kind:        canvas.KindOpportunity,
	})
	if err != nil {
		return fmt.Errorf("failed to create Opportunity Canvas: %w", err)
	}

	printCanvasCreated("Opportunity Canvas", sm)

	if verbose {
		fmt.Println()
		fmt.Println("Expected blocks for Opportunity Canvas:")
		for _, block := range canvas.OpportunityBlocks {
			fmt.Printf("  - %s\n", block)
		}
	}

	return nil
}
