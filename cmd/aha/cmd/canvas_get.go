package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// canvasGetCmd retrieves a specific strategic canvas.
var canvasGetCmd = &cobra.Command{
	Use:   "get <canvas-id>",
	Short: "Get a strategic canvas by ID",
	Long: `Get details of a specific strategic canvas (strategic model) in Aha.io.

The canvas ID can be the numeric ID or the reference number.

Examples:
  # Get canvas by ID
  aha canvas get 12345678

  # Get canvas by reference number
  aha canvas get PROD-SM-1`,
	Args: cobra.ExactArgs(1),
	RunE: runGetCanvas,
}

func init() {
	canvasCmd.AddCommand(canvasGetCmd)
}

func runGetCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	canvasID := args[0]

	sm, err := client.GetStrategicModel(ctx, canvasID)
	if err != nil {
		return fmt.Errorf("failed to get canvas: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "ID:\t%s\n", sm.ID)
	_, _ = fmt.Fprintf(w, "Reference:\t%s\n", sm.ReferenceNum)
	_, _ = fmt.Fprintf(w, "Name:\t%s\n", sm.Name)
	_, _ = fmt.Fprintf(w, "Kind:\t%s\n", sm.Kind)
	if sm.Description != "" {
		_, _ = fmt.Fprintf(w, "Description:\t%s\n", sm.Description)
	}
	_, _ = fmt.Fprintf(w, "URL:\t%s\n", sm.URL)
	_, _ = fmt.Fprintf(w, "Created:\t%s\n", sm.CreatedAt.Format("2006-01-02 15:04:05"))
	if sm.UpdatedAt != nil {
		_, _ = fmt.Fprintf(w, "Updated:\t%s\n", sm.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	_ = w.Flush()

	if len(sm.Components) > 0 {
		fmt.Println()
		fmt.Println("Components (blocks):")
		fmt.Println()
		for _, comp := range sm.Components {
			fmt.Printf("  [%s] %s\n", comp.ID, comp.Name)
			if comp.Description != "" && verbose {
				// Show content preview (first 100 chars)
				content := comp.Description
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				fmt.Printf("      %s\n", content)
			}
		}
	}

	return nil
}
