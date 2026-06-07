package render

import (
	"fmt"
	"strings"

	aha "github.com/grokify/aha-go"
)

// D2Options configures D2 rendering.
type D2Options struct {
	Theme     int    // D2 theme number (0-8)
	Layout    string // "dagre", "elk", "tala"
	Direction string // "right", "down", "left", "up"
}

// DefaultD2Options returns default D2 rendering options.
func DefaultD2Options() *D2Options {
	return &D2Options{
		Theme:     0,
		Layout:    "dagre",
		Direction: "right",
	}
}

// RenderCanvasD2 renders a strategic model to D2 diagram language.
func RenderCanvasD2(sm *aha.StrategicModel, opts *D2Options) string {
	if opts == nil {
		opts = DefaultD2Options()
	}

	switch strings.ToLower(sm.Kind) {
	case "opportunity":
		return renderOpportunityD2(sm, opts)
	case "lean canvas":
		return renderLeanUXD2(sm, opts)
	case "business model":
		return renderBMCD2(sm, opts)
	default:
		return renderGenericD2(sm, opts)
	}
}

func renderOpportunityD2(sm *aha.StrategicModel, _ *D2Options) string {
	var sb strings.Builder

	sb.WriteString("# Opportunity Canvas\n")
	sb.WriteString("# Jeff Patton's 10-block structure\n\n")
	_, _ = fmt.Fprintf(&sb, "title: %s {shape: text; near: top-center}\n\n", escapeD2(sm.Name))

	// Use grid layout
	sb.WriteString("grid-rows: 4\n")
	sb.WriteString("grid-columns: 3\n")
	sb.WriteString("grid-gap: 10\n\n")

	components := componentMap(sm)

	// Row 1
	writeD2Block(&sb, "users", "Users & Customers", components["Users & Customers"], "teal")
	writeD2Block(&sb, "problems", "Problems", components["Problems"], "teal")
	writeD2Block(&sb, "solution_ideas", "Solution Ideas", components["Solution Ideas"], "teal")

	// Row 2
	writeD2Block(&sb, "solutions_today", "Solutions Today", components["Solutions Today"], "blue")
	writeD2Block(&sb, "user_value", "User Value", components["User Value"], "blue")
	writeD2Block(&sb, "adoption", "Adoption Strategy", components["Adoption Strategy"], "blue")

	// Row 3
	writeD2Block(&sb, "user_metrics", "User Metrics", components["User Metrics"], "amber")
	writeD2Block(&sb, "business_problem", "Business Problem", components["Business Problem"], "amber")
	writeD2Block(&sb, "business_metrics", "Business Metrics", components["Business Metrics"], "amber")

	// Row 4 - Budget spans all columns
	sb.WriteString("budget: Budget {\n")
	sb.WriteString("  grid-column: 1 / 4\n")
	if comp := components["Budget"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#3c3489\"\n")
	sb.WriteString("}\n")

	return sb.String()
}

func renderLeanUXD2(sm *aha.StrategicModel, _ *D2Options) string {
	var sb strings.Builder

	sb.WriteString("# Lean UX Canvas\n")
	sb.WriteString("# Jeff Gothelf's 8-block structure\n\n")
	_, _ = fmt.Fprintf(&sb, "title: %s {shape: text; near: top-center}\n\n", escapeD2(sm.Name))

	sb.WriteString("grid-rows: 4\n")
	sb.WriteString("grid-columns: 3\n")
	sb.WriteString("grid-gap: 10\n\n")

	components := componentMap(sm)

	// Row 1 - Business Problem and Outcomes span columns
	sb.WriteString("business_problem: Business Problem {\n")
	sb.WriteString("  grid-column: 1 / 2\n")
	if comp := components["Business Problem"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#3c3489\"\n")
	sb.WriteString("}\n\n")

	sb.WriteString("business_outcomes: Business Outcomes {\n")
	sb.WriteString("  grid-column: 2 / 4\n")
	if comp := components["Business Outcomes"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#3c3489\"\n")
	sb.WriteString("}\n\n")

	// Row 2
	writeD2Block(&sb, "users", "Users", components["Users"], "teal")
	writeD2Block(&sb, "benefits", "Benefits", components["Benefits"], "teal")
	writeD2Block(&sb, "solutions", "Solutions", components["Solutions"], "teal")

	// Row 3 - Hypotheses spans all columns
	sb.WriteString("hypotheses: Hypotheses {\n")
	sb.WriteString("  grid-column: 1 / 4\n")
	if comp := components["Hypotheses"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#633806\"\n")
	sb.WriteString("}\n\n")

	// Row 4
	sb.WriteString("riskiest: Riskiest Assumption {\n")
	sb.WriteString("  grid-column: 1 / 2\n")
	if comp := components["Riskiest Assumption"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#712b13\"\n")
	sb.WriteString("}\n\n")

	sb.WriteString("experiment: Smallest Experiment {\n")
	sb.WriteString("  grid-column: 2 / 4\n")
	if comp := components["Smallest Experiment"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#712b13\"\n")
	sb.WriteString("}\n")

	return sb.String()
}

func renderBMCD2(sm *aha.StrategicModel, _ *D2Options) string {
	var sb strings.Builder

	sb.WriteString("# Business Model Canvas\n")
	sb.WriteString("# Osterwalder's 9-block structure\n\n")
	_, _ = fmt.Fprintf(&sb, "title: %s {shape: text; near: top-center}\n\n", escapeD2(sm.Name))

	sb.WriteString("grid-rows: 3\n")
	sb.WriteString("grid-columns: 5\n")
	sb.WriteString("grid-gap: 10\n\n")

	components := componentMap(sm)

	// Row 1-2: Top section with varying heights
	// Col 1: Key Partners (spans 2 rows)
	sb.WriteString("key_partners: Key Partners {\n")
	sb.WriteString("  grid-row: 1 / 3\n")
	if comp := components["Key Partners"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#085041\"\n")
	sb.WriteString("}\n\n")

	// Col 2: Key Activities, Key Resources (stacked)
	writeD2Block(&sb, "key_activities", "Key Activities", components["Key Activities"], "blue")
	writeD2Block(&sb, "key_resources", "Key Resources", components["Key Resources"], "blue")

	// Col 3: Value Propositions (spans 2 rows)
	sb.WriteString("value_props: Value Propositions {\n")
	sb.WriteString("  grid-row: 1 / 3\n")
	if comp := components["Value Propositions"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#3c3489\"\n")
	sb.WriteString("}\n\n")

	// Col 4: Customer Relationships, Channels (stacked)
	writeD2Block(&sb, "cust_rel", "Customer Relationships", components["Customer Relationships"], "amber")
	writeD2Block(&sb, "channels", "Channels", components["Channels"], "amber")

	// Col 5: Customer Segments (spans 2 rows)
	sb.WriteString("cust_segments: Customer Segments {\n")
	sb.WriteString("  grid-row: 1 / 3\n")
	if comp := components["Customer Segments"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#712b13\"\n")
	sb.WriteString("}\n\n")

	// Row 3: Bottom section
	sb.WriteString("cost_structure: Cost Structure {\n")
	sb.WriteString("  grid-column: 1 / 3\n")
	if comp := components["Cost Structure"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#085041\"\n")
	sb.WriteString("}\n\n")

	sb.WriteString("revenue_streams: Revenue Streams {\n")
	sb.WriteString("  grid-column: 3 / 6\n")
	if comp := components["Revenue Streams"]; comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(&sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	sb.WriteString("  style.fill: \"#085041\"\n")
	sb.WriteString("}\n")

	return sb.String()
}

func renderGenericD2(sm *aha.StrategicModel, _ *D2Options) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "# %s Canvas\n\n", sm.Kind)
	_, _ = fmt.Fprintf(&sb, "title: %s {shape: text; near: top-center}\n\n", escapeD2(sm.Name))

	sb.WriteString("grid-rows: 3\n")
	sb.WriteString("grid-columns: 3\n")
	sb.WriteString("grid-gap: 10\n\n")

	colors := []string{"teal", "blue", "amber", "purple", "coral"}

	for i, comp := range sm.Components {
		id := fmt.Sprintf("block_%d", i)
		color := colors[i%len(colors)]
		writeD2Block(&sb, id, comp.Name, &comp, color)
	}

	return sb.String()
}

func writeD2Block(sb *strings.Builder, id, title string, comp *aha.StrategicModelComponent, color string) {
	colorMap := map[string]string{
		"teal":   "#085041",
		"blue":   "#1e3a5f",
		"amber":  "#633806",
		"purple": "#3c3489",
		"coral":  "#712b13",
	}

	fillColor := colorMap[color]
	if fillColor == "" {
		fillColor = "#333333"
	}

	_, _ = fmt.Fprintf(sb, "%s: %s {\n", id, escapeD2(title))
	if comp != nil && comp.Description != "" {
		_, _ = fmt.Fprintf(sb, "  tooltip: %s\n", escapeD2(stripHTML(comp.Description)))
	}
	_, _ = fmt.Fprintf(sb, "  style.fill: \"%s\"\n", fillColor)
	sb.WriteString("}\n\n")
}

func componentMap(sm *aha.StrategicModel) map[string]*aha.StrategicModelComponent {
	components := make(map[string]*aha.StrategicModelComponent)
	for i := range sm.Components {
		components[sm.Components[i].Name] = &sm.Components[i]
	}
	return components
}

func escapeD2(s string) string {
	// Escape special characters for D2
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
