package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/aha-go/canvas"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	exportOutput string
	exportWidth  int
	exportHeight int
)

// canvasExportCmd exports a strategic canvas.
var canvasExportCmd = &cobra.Command{
	Use:   "export <canvas-id>",
	Short: "Export a strategic canvas",
	Long: `Export a strategic canvas (strategic model) from Aha.io.

Supported formats:
  - svg      SVG image with grid layout (default)
  - json     Raw JSON data
  - d2       D2 diagram language (https://d2lang.com)
  - mermaid  Mermaid diagram syntax (https://mermaid.js.org)

The SVG output uses a dark theme with color-coded blocks matching
the canvas type's standard layout.

Examples:
  # Export to SVG (default)
  aha canvas export PROD-SM-1 -o canvas.svg

  # Export to JSON
  aha canvas export PROD-SM-1 --format json -o canvas.json

  # Export to D2 diagram format
  aha canvas export PROD-SM-1 --format d2 -o canvas.d2

  # Export to Mermaid diagram format
  aha canvas export PROD-SM-1 --format mermaid -o canvas.mmd

  # Export with custom dimensions (SVG only)
  aha canvas export PROD-SM-1 --width 1200 --height 800 -o canvas.svg

  # Output to stdout
  aha canvas export PROD-SM-1`,
	Args: cobra.ExactArgs(1),
	RunE: runExportCanvas,
}

func init() {
	canvasCmd.AddCommand(canvasExportCmd)

	canvasExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "svg", "Output format (svg, json, d2, mermaid)")
	canvasExportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
	canvasExportCmd.Flags().IntVar(&exportWidth, "width", 800, "SVG width")
	canvasExportCmd.Flags().IntVar(&exportHeight, "height", 600, "SVG height")
}

func runExportCanvas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	canvasID := args[0]

	// Map format string to enum
	var format canvas.ExportFormat
	switch exportFormat {
	case "svg":
		format = canvas.FormatSVG
	case "json":
		format = canvas.FormatJSON
	case "d2":
		format = canvas.FormatD2
	case "mermaid", "mmd":
		format = canvas.FormatMermaid
	default:
		return fmt.Errorf("unsupported format: %s (use svg, json, d2, or mermaid)", exportFormat)
	}

	// Export
	output, err := canvas.Export(ctx, client, canvas.ExportOptions{
		CanvasID: canvasID,
		Format:   format,
		Width:    exportWidth,
		Height:   exportHeight,
	})
	if err != nil {
		return err
	}

	// Write output
	if exportOutput == "" {
		// Write to stdout
		fmt.Print(string(output))
	} else {
		// Write to file
		if err := os.WriteFile(exportOutput, output, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Exported to %s (%d bytes)\n", exportOutput, len(output))
	}

	return nil
}
