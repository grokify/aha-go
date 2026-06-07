package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/canvas"
	"github.com/spf13/cobra"
)

// canvasCreateLeanUXCmd creates a Lean UX Canvas.
var canvasCreateLeanUXCmd = &cobra.Command{
	Use:   "leanux",
	Short: "Create a Lean UX Canvas",
	Long: `Create a new Lean UX Canvas in Aha.io.

The Lean UX Canvas follows Jeff Gothelf's v2 8-block structure for
hypothesis-driven product development focused on outcomes.

Grid Layout:
  +-------------------+-------------------+
  | Business Problem  | Business Outcomes |
  +-------------------+--------+----------+
  | Users             | Benefits| Solutions|
  +-------------------+--------+----------+
  | Hypotheses (full width)               |
  +-------------------+-------------------+
  | Riskiest Assumption| Smallest Experiment|
  +-------------------+-------------------+

Example:
  aha canvas create leanux \
    --product PROD \
    --name "Checkout Flow Redesign" \
    --description "Reduce cart abandonment rate"`,
	RunE: runCreateLeanUXCanvas,
}

func init() {
	canvasCreateCmd.AddCommand(canvasCreateLeanUXCmd)
}

func runCreateLeanUXCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sm, err := canvas.Create(ctx, client, canvas.CreateOptions{
		ProductID:   productID,
		Name:        canvasName,
		Description: description,
		Kind:        canvas.KindLeanUX,
	})
	if err != nil {
		return fmt.Errorf("failed to create Lean UX Canvas: %w", err)
	}

	printCanvasCreated("Lean UX Canvas", sm)

	if verbose {
		fmt.Println()
		fmt.Println("Expected blocks for Lean UX Canvas:")
		for _, block := range canvas.LeanUXBlocks {
			fmt.Printf("  - %s\n", block)
		}
	}

	return nil
}
