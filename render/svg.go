// Package render provides rendering utilities for Aha canvases.
package render

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	aha "github.com/grokify/aha-go"
)

// SVGOptions configures SVG rendering.
type SVGOptions struct {
	Width      int
	Height     int
	FontFamily string
	Theme      string // "dark" or "light"
}

// DefaultSVGOptions returns default SVG options.
func DefaultSVGOptions() *SVGOptions {
	return &SVGOptions{
		Width:      800,
		Height:     600,
		FontFamily: `-apple-system, "system-ui", "Segoe UI", sans-serif`,
		Theme:      "dark",
	}
}

// colorScheme defines colors for a block.
type colorScheme struct {
	fill         string
	stroke       string
	titleFill    string
	subtitleFill string
}

var (
	// Dark theme color schemes
	tealScheme = colorScheme{
		fill:         "rgb(8, 80, 65)",
		stroke:       "rgb(93, 202, 165)",
		titleFill:    "rgb(159, 225, 203)",
		subtitleFill: "rgb(93, 202, 165)",
	}
	blueScheme = colorScheme{
		fill:         "rgb(30, 58, 95)",
		stroke:       "rgb(96, 165, 250)",
		titleFill:    "rgb(191, 219, 254)",
		subtitleFill: "rgb(96, 165, 250)",
	}
	amberScheme = colorScheme{
		fill:         "rgb(99, 56, 6)",
		stroke:       "rgb(239, 159, 39)",
		titleFill:    "rgb(250, 199, 117)",
		subtitleFill: "rgb(239, 159, 39)",
	}
	purpleScheme = colorScheme{
		fill:         "rgb(60, 52, 137)",
		stroke:       "rgb(175, 169, 236)",
		titleFill:    "rgb(206, 203, 246)",
		subtitleFill: "rgb(175, 169, 236)",
	}
	coralScheme = colorScheme{
		fill:         "rgb(113, 43, 19)",
		stroke:       "rgb(240, 153, 123)",
		titleFill:    "rgb(245, 196, 179)",
		subtitleFill: "rgb(240, 153, 123)",
	}
)

// RenderCanvasSVG renders a strategic model to SVG.
func RenderCanvasSVG(sm *aha.StrategicModel, opts *SVGOptions) string {
	if opts == nil {
		opts = DefaultSVGOptions()
	}

	switch strings.ToLower(sm.Kind) {
	case "opportunity":
		return renderOpportunitySVG(sm, opts)
	case "lean canvas":
		return renderLeanUXSVG(sm, opts)
	case "business model":
		return renderBMCSVG(sm, opts)
	default:
		return renderGenericSVG(sm, opts)
	}
}

// renderOpportunitySVG renders an Opportunity Canvas to SVG.
func renderOpportunitySVG(sm *aha.StrategicModel, opts *SVGOptions) string {
	var sb strings.Builder

	// SVG header
	_, _ = fmt.Fprintf(&sb, `<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, opts.Width, opts.Height)
	_, _ = fmt.Fprintf(&sb, `<title>%s</title>`, html.EscapeString(sm.Name))
	sb.WriteString(`<desc>Opportunity Canvas - Jeff Patton's 10-block structure</desc>`)

	// Build component map
	components := make(map[string]*aha.StrategicModelComponent)
	for i := range sm.Components {
		components[sm.Components[i].Name] = &sm.Components[i]
	}

	// Grid layout: 3 columns, 4 rows (last row spans all)
	const (
		startX  = 40
		startY  = 30
		gap     = 10
		cornerR = 8
	)

	totalWidth := opts.Width - 2*startX
	totalHeight := opts.Height - 2*startY
	colWidth := (totalWidth - 2*gap) / 3
	rowHeight := (totalHeight - 3*gap) / 4

	// Row 1: Users & Customers, Problems, Solution Ideas
	row1Blocks := []string{"Users & Customers", "Problems", "Solution Ideas"}
	for i, name := range row1Blocks {
		x := startX + i*(colWidth+gap)
		renderBlock(&sb, x, startY, colWidth, rowHeight, name, components[name], tealScheme, opts)
	}

	// Row 2: Solutions Today, User Value, Adoption Strategy
	row2Blocks := []string{"Solutions Today", "User Value", "Adoption Strategy"}
	y2 := startY + rowHeight + gap
	for i, name := range row2Blocks {
		x := startX + i*(colWidth+gap)
		renderBlock(&sb, x, y2, colWidth, rowHeight, name, components[name], blueScheme, opts)
	}

	// Row 3: User Metrics, Business Problem, Business Metrics
	row3Blocks := []string{"User Metrics", "Business Problem", "Business Metrics"}
	y3 := startY + 2*(rowHeight+gap)
	for i, name := range row3Blocks {
		x := startX + i*(colWidth+gap)
		renderBlock(&sb, x, y3, colWidth, rowHeight, name, components[name], amberScheme, opts)
	}

	// Row 4: Budget (full width)
	y4 := startY + 3*(rowHeight+gap)
	renderBlock(&sb, startX, y4, totalWidth, rowHeight, "Budget", components["Budget"], purpleScheme, opts)

	sb.WriteString(`</svg>`)
	return sb.String()
}

// renderLeanUXSVG renders a Lean UX Canvas to SVG.
func renderLeanUXSVG(sm *aha.StrategicModel, opts *SVGOptions) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, `<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, opts.Width, opts.Height)
	_, _ = fmt.Fprintf(&sb, `<title>%s</title>`, html.EscapeString(sm.Name))
	sb.WriteString(`<desc>Lean UX Canvas - Jeff Gothelf's 8-block structure</desc>`)

	components := make(map[string]*aha.StrategicModelComponent)
	for i := range sm.Components {
		components[sm.Components[i].Name] = &sm.Components[i]
	}

	const (
		startX = 40
		startY = 30
		gap    = 10
	)

	totalWidth := opts.Width - 2*startX
	totalHeight := opts.Height - 2*startY
	halfWidth := (totalWidth - gap) / 2
	thirdWidth := (totalWidth - 2*gap) / 3
	rowHeight := (totalHeight - 3*gap) / 4

	// Row 1: Business Problem, Business Outcomes (2 columns)
	renderBlock(&sb, startX, startY, halfWidth, rowHeight, "Business Problem", components["Business Problem"], purpleScheme, opts)
	renderBlock(&sb, startX+halfWidth+gap, startY, halfWidth, rowHeight, "Business Outcomes", components["Business Outcomes"], purpleScheme, opts)

	// Row 2: Users, Benefits, Solutions (3 columns)
	y2 := startY + rowHeight + gap
	renderBlock(&sb, startX, y2, thirdWidth, rowHeight, "Users", components["Users"], tealScheme, opts)
	renderBlock(&sb, startX+thirdWidth+gap, y2, thirdWidth, rowHeight, "Benefits", components["Benefits"], tealScheme, opts)
	renderBlock(&sb, startX+2*(thirdWidth+gap), y2, thirdWidth, rowHeight, "Solutions", components["Solutions"], tealScheme, opts)

	// Row 3: Hypotheses (full width)
	y3 := startY + 2*(rowHeight+gap)
	renderBlock(&sb, startX, y3, totalWidth, rowHeight, "Hypotheses", components["Hypotheses"], amberScheme, opts)

	// Row 4: Riskiest Assumption, Smallest Experiment (2 columns)
	y4 := startY + 3*(rowHeight+gap)
	renderBlock(&sb, startX, y4, halfWidth, rowHeight, "Riskiest Assumption", components["Riskiest Assumption"], coralScheme, opts)
	renderBlock(&sb, startX+halfWidth+gap, y4, halfWidth, rowHeight, "Smallest Experiment", components["Smallest Experiment"], coralScheme, opts)

	sb.WriteString(`</svg>`)
	return sb.String()
}

// renderBMCSVG renders a Business Model Canvas to SVG.
func renderBMCSVG(sm *aha.StrategicModel, opts *SVGOptions) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, `<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, opts.Width, opts.Height)
	_, _ = fmt.Fprintf(&sb, `<title>%s</title>`, html.EscapeString(sm.Name))
	sb.WriteString(`<desc>Business Model Canvas - Osterwalder's 9-block structure</desc>`)

	components := make(map[string]*aha.StrategicModelComponent)
	for i := range sm.Components {
		components[sm.Components[i].Name] = &sm.Components[i]
	}

	const (
		startX = 40
		startY = 30
		gap    = 10
	)

	totalWidth := opts.Width - 2*startX
	totalHeight := opts.Height - 2*startY
	colWidth := (totalWidth - 4*gap) / 5
	halfHeight := (totalHeight - 2*gap) / 3
	topHeight := halfHeight * 2

	// Top section (2/3 height)
	// Col 1: Key Partners (full height)
	renderBlock(&sb, startX, startY, colWidth, topHeight, "Key Partners", components["Key Partners"], tealScheme, opts)

	// Col 2: Key Activities (top), Key Resources (bottom)
	x2 := startX + colWidth + gap
	renderBlock(&sb, x2, startY, colWidth, halfHeight-gap/2, "Key Activities", components["Key Activities"], blueScheme, opts)
	renderBlock(&sb, x2, startY+halfHeight+gap/2, colWidth, halfHeight-gap/2, "Key Resources", components["Key Resources"], blueScheme, opts)

	// Col 3: Value Propositions (full height)
	x3 := startX + 2*(colWidth+gap)
	renderBlock(&sb, x3, startY, colWidth, topHeight, "Value Propositions", components["Value Propositions"], purpleScheme, opts)

	// Col 4: Customer Relationships (top), Channels (bottom)
	x4 := startX + 3*(colWidth+gap)
	renderBlock(&sb, x4, startY, colWidth, halfHeight-gap/2, "Customer Relationships", components["Customer Relationships"], amberScheme, opts)
	renderBlock(&sb, x4, startY+halfHeight+gap/2, colWidth, halfHeight-gap/2, "Channels", components["Channels"], amberScheme, opts)

	// Col 5: Customer Segments (full height)
	x5 := startX + 4*(colWidth+gap)
	renderBlock(&sb, x5, startY, colWidth, topHeight, "Customer Segments", components["Customer Segments"], coralScheme, opts)

	// Bottom section (1/3 height)
	y2 := startY + topHeight + gap
	halfWidth := (totalWidth - gap) / 2
	renderBlock(&sb, startX, y2, halfWidth, halfHeight, "Cost Structure", components["Cost Structure"], tealScheme, opts)
	renderBlock(&sb, startX+halfWidth+gap, y2, halfWidth, halfHeight, "Revenue Streams", components["Revenue Streams"], tealScheme, opts)

	sb.WriteString(`</svg>`)
	return sb.String()
}

// renderGenericSVG renders an unknown canvas type as a simple grid.
func renderGenericSVG(sm *aha.StrategicModel, opts *SVGOptions) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, `<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, opts.Width, opts.Height)
	_, _ = fmt.Fprintf(&sb, `<title>%s</title>`, html.EscapeString(sm.Name))
	_, _ = fmt.Fprintf(&sb, `<desc>%s Canvas</desc>`, html.EscapeString(sm.Kind))

	const (
		startX = 40
		startY = 30
		gap    = 10
		cols   = 3
	)

	totalWidth := opts.Width - 2*startX
	colWidth := (totalWidth - (cols-1)*gap) / cols

	rows := (len(sm.Components) + cols - 1) / cols
	if rows == 0 {
		rows = 1
	}
	totalHeight := opts.Height - 2*startY
	rowHeight := (totalHeight - (rows-1)*gap) / rows

	schemes := []colorScheme{tealScheme, blueScheme, amberScheme, purpleScheme, coralScheme}

	for i, comp := range sm.Components {
		row := i / cols
		col := i % cols
		x := startX + col*(colWidth+gap)
		y := startY + row*(rowHeight+gap)
		scheme := schemes[row%len(schemes)]
		renderBlock(&sb, x, y, colWidth, rowHeight, comp.Name, &comp, scheme, opts)
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

// renderBlock renders a single block.
func renderBlock(sb *strings.Builder, x, y, width, height int, title string, comp *aha.StrategicModelComponent, scheme colorScheme, opts *SVGOptions) {
	sb.WriteString(`<g>`)

	// Rectangle
	_, _ = fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" rx="8" style="fill:%s;stroke:%s;stroke-width:0.5"/>`,
		x, y, width, height, scheme.fill, scheme.stroke)

	// Title
	centerX := x + width/2
	_, _ = fmt.Fprintf(sb, `<text x="%d" y="%d" text-anchor="middle" style="fill:%s;font-size:12px;font-weight:600;font-family:%s">%s</text>`,
		centerX, y+18, scheme.titleFill, opts.FontFamily, html.EscapeString(title))

	// Content preview (if available)
	if comp != nil && comp.Description != "" {
		content := stripHTML(comp.Description)
		// Truncate and show first line
		if len(content) > 50 {
			content = content[:50] + "..."
		}
		_, _ = fmt.Fprintf(sb, `<text x="%d" y="%d" text-anchor="start" style="fill:%s;font-size:10px;font-family:%s">%s</text>`,
			x+12, y+38, scheme.subtitleFill, opts.FontFamily, html.EscapeString(content))
	}

	sb.WriteString(`</g>`)
}

// stripHTML removes HTML tags from a string.
func stripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}
