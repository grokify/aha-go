package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/canvas"
	"github.com/spf13/cobra"
)

// canvasCreateBMCCmd creates a Business Model Canvas.
var canvasCreateBMCCmd = &cobra.Command{
	Use:   "bmc",
	Short: "Create a Business Model Canvas",
	Long: `Create a new Business Model Canvas in Aha.io.

The Business Model Canvas follows Alexander Osterwalder's 9-block
structure for describing, designing, and analyzing business models.

Grid Layout:
  +---------------+---------------+---------------+---------------+
  |               | Key           |               | Customer      |
  | Key Partners  | Activities    | Value         | Relationships |
  |               +---------------+ Propositions  +---------------+
  |               | Key           |               | Channels      |
  |               | Resources     |               |               |
  +---------------+---------------+---------------+---------------+
  | Cost Structure                | Revenue Streams               |
  +-------------------------------+-------------------------------+

Example:
  aha canvas create bmc \
    --product PROD \
    --name "SaaS Business Model" \
    --description "Business model for our SaaS platform"`,
	RunE: runCreateBMCCanvas,
}

func init() {
	canvasCreateCmd.AddCommand(canvasCreateBMCCmd)
}

func runCreateBMCCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sm, err := canvas.Create(ctx, client, canvas.CreateOptions{
		ProductID:   productID,
		Name:        canvasName,
		Description: description,
		Kind:        canvas.KindBMC,
	})
	if err != nil {
		return fmt.Errorf("failed to create Business Model Canvas: %w", err)
	}

	printCanvasCreated("Business Model Canvas", sm)

	if verbose {
		fmt.Println()
		fmt.Println("Expected blocks for Business Model Canvas:")
		for _, block := range canvas.BMCBlocks {
			fmt.Printf("  - %s\n", block)
		}
	}

	return nil
}
