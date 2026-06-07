package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/canvas"
	"github.com/spf13/cobra"
)

var (
	contentFile      string
	componentName    string
	componentContent string
)

// canvasUpdateCmd updates a strategic canvas.
var canvasUpdateCmd = &cobra.Command{
	Use:   "update <canvas-id>",
	Short: "Update a strategic canvas",
	Long: `Update a strategic canvas (strategic model) in Aha.io.

You can update canvas components (blocks) in two ways:

1. From a JSON file mapping block names to content:
   aha canvas update PROD-SM-1 --file content.json

2. Update a single block directly:
   aha canvas update PROD-SM-1 --block "Problems" --content "<p>User problems...</p>"

JSON file format:
  {
    "Users & Customers": "<p>Target users...</p>",
    "Problems": "<p>Key problems...</p>",
    "Solution Ideas": "<ul><li>Idea 1</li><li>Idea 2</li></ul>"
  }

Note: Content should be HTML formatted.

Examples:
  # Update from JSON file
  aha canvas update PROD-SM-1 --file opportunity-content.json

  # Update single block
  aha canvas update PROD-SM-1 \
    --block "Business Problem" \
    --content "<p>High cart abandonment rate of 72%</p>"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdateCanvas,
}

func init() {
	canvasCmd.AddCommand(canvasUpdateCmd)

	canvasUpdateCmd.Flags().StringVarP(&contentFile, "file", "f", "", "JSON file with block name to content mapping")
	canvasUpdateCmd.Flags().StringVarP(&componentName, "block", "b", "", "Block name to update")
	canvasUpdateCmd.Flags().StringVarP(&componentContent, "content", "c", "", "HTML content for the block")
}

func runUpdateCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	canvasID := args[0]

	// Validate flags
	if contentFile == "" && componentName == "" {
		return fmt.Errorf("either --file or --block must be specified")
	}
	if componentName != "" && componentContent == "" {
		return fmt.Errorf("--content is required when using --block")
	}

	// Determine content to update
	var blocks map[string]string
	var err error

	if contentFile != "" {
		blocks, err = canvas.LoadBlocksFromFile(contentFile)
		if err != nil {
			return err
		}
	} else {
		blocks = map[string]string{
			componentName: componentContent,
		}
	}

	// Perform update
	result, err := canvas.Update(ctx, client, canvas.UpdateOptions{
		CanvasID: canvasID,
		Blocks:   blocks,
	})
	if err != nil {
		return err
	}

	// Print warnings about unknown blocks
	if len(result.Unknown) > 0 {
		fmt.Println("Warning: The following blocks were not found in the canvas:")
		for _, name := range result.Unknown {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Println()

		// Show available blocks
		availableBlocks, err := canvas.GetAvailableBlocks(ctx, client, canvasID)
		if err == nil {
			fmt.Println("Available blocks:")
			for _, name := range availableBlocks {
				fmt.Printf("  - %s\n", name)
			}
			fmt.Println()
		}
	}

	// Print errors
	for _, blockErr := range result.Errors {
		fmt.Printf("Error updating %q: %v\n", blockErr.BlockName, blockErr.Err)
	}

	// Print success/verbose info
	if verbose {
		fmt.Println("Update complete:")
	}

	fmt.Printf("Updated %d block(s)", result.SuccessCount)
	if result.ErrorCount > 0 {
		fmt.Printf(", %d error(s)", result.ErrorCount)
	}
	fmt.Println()

	return nil
}
