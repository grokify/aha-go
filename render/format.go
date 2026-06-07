// Package render provides rendering utilities for Aha canvases.
package render

// Format identifies the output format for canvas rendering.
type Format string

const (
	FormatSVG     Format = "svg"
	FormatD2      Format = "d2"
	FormatMermaid Format = "mermaid"
)

// SupportedFormats returns all supported render formats.
func SupportedFormats() []Format {
	return []Format{FormatSVG, FormatD2, FormatMermaid}
}

// FileExtension returns the file extension for a format.
func (f Format) FileExtension() string {
	switch f {
	case FormatSVG:
		return ".svg"
	case FormatD2:
		return ".d2"
	case FormatMermaid:
		return ".mmd"
	default:
		return ".txt"
	}
}

// MimeType returns the MIME type for a format.
func (f Format) MimeType() string {
	switch f {
	case FormatSVG:
		return "image/svg+xml"
	case FormatD2:
		return "text/plain"
	case FormatMermaid:
		return "text/plain"
	default:
		return "text/plain"
	}
}
