package canvas

import (
	"context"
	"encoding/json"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/render"
)

// ExportFormat represents an export format.
type ExportFormat string

// Export formats.
const (
	FormatSVG     ExportFormat = "svg"
	FormatJSON    ExportFormat = "json"
	FormatD2      ExportFormat = "d2"
	FormatMermaid ExportFormat = "mermaid"
)

// ExportOptions configures canvas export.
type ExportOptions struct {
	CanvasID string
	Format   ExportFormat
	Width    int
	Height   int
}

// DefaultExportOptions returns default export options.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format: FormatSVG,
		Width:  800,
		Height: 600,
	}
}

// Export exports a canvas to the specified format.
func Export(ctx context.Context, client *aha.Client, opts ExportOptions) ([]byte, error) {
	if opts.CanvasID == "" {
		return nil, fmt.Errorf("canvas ID is required")
	}

	// Get the canvas
	sm, err := client.GetStrategicModel(ctx, opts.CanvasID)
	if err != nil {
		return nil, fmt.Errorf("failed to get canvas: %w", err)
	}

	return ExportModel(sm, opts)
}

// ExportModel exports a strategic model to the specified format.
// This function is useful for testing without an API call.
func ExportModel(sm *aha.StrategicModel, opts ExportOptions) ([]byte, error) {
	switch opts.Format {
	case FormatSVG:
		svgOpts := &render.SVGOptions{
			Width:      opts.Width,
			Height:     opts.Height,
			FontFamily: `-apple-system, "system-ui", "Segoe UI", sans-serif`,
			Theme:      "dark",
		}
		if svgOpts.Width == 0 {
			svgOpts.Width = 800
		}
		if svgOpts.Height == 0 {
			svgOpts.Height = 600
		}
		return []byte(render.RenderCanvasSVG(sm, svgOpts)), nil

	case FormatJSON:
		return json.MarshalIndent(sm, "", "  ")

	case FormatD2:
		return []byte(render.RenderCanvasD2(sm, nil)), nil

	case FormatMermaid:
		return []byte(render.RenderCanvasMermaid(sm, nil)), nil

	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

// SupportedFormats returns the list of supported export formats.
func SupportedFormats() []ExportFormat {
	return []ExportFormat{FormatSVG, FormatJSON, FormatD2, FormatMermaid}
}
