package render

import (
	"fmt"
	"strings"

	aha "github.com/grokify/aha-go"
)

// MermaidOptions configures Mermaid rendering.
type MermaidOptions struct {
	Theme     string // "default", "dark", "forest", "neutral"
	Direction string // "TB", "BT", "LR", "RL"
}

// DefaultMermaidOptions returns default Mermaid rendering options.
func DefaultMermaidOptions() *MermaidOptions {
	return &MermaidOptions{
		Theme:     "dark",
		Direction: "TB",
	}
}

// RenderCanvasMermaid renders a strategic model to Mermaid diagram language.
func RenderCanvasMermaid(sm *aha.StrategicModel, opts *MermaidOptions) string {
	if opts == nil {
		opts = DefaultMermaidOptions()
	}

	switch strings.ToLower(sm.Kind) {
	case "opportunity":
		return renderOpportunityMermaid(sm, opts)
	case "lean canvas":
		return renderLeanUXMermaid(sm, opts)
	case "business model":
		return renderBMCMermaid(sm, opts)
	default:
		return renderGenericMermaid(sm, opts)
	}
}

func renderOpportunityMermaid(sm *aha.StrategicModel, _ *MermaidOptions) string {
	var sb strings.Builder

	sb.WriteString("%%{init: {'theme': 'dark'}}%%\n")
	sb.WriteString("flowchart TB\n\n")
	_, _ = fmt.Fprintf(&sb, "    subgraph title[\"%s\"]\n", escapeMermaid(sm.Name))
	sb.WriteString("    direction TB\n\n")

	components := componentMap(sm)

	// Row 1
	sb.WriteString("    subgraph row1[\"Problem Space\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "users", "Users & Customers", components["Users & Customers"])
	writeMermaidNode(&sb, "problems", "Problems", components["Problems"])
	writeMermaidNode(&sb, "solution_ideas", "Solution Ideas", components["Solution Ideas"])
	sb.WriteString("    end\n\n")

	// Row 2
	sb.WriteString("    subgraph row2[\"Solution Space\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "solutions_today", "Solutions Today", components["Solutions Today"])
	writeMermaidNode(&sb, "user_value", "User Value", components["User Value"])
	writeMermaidNode(&sb, "adoption", "Adoption Strategy", components["Adoption Strategy"])
	sb.WriteString("    end\n\n")

	// Row 3
	sb.WriteString("    subgraph row3[\"Business\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "user_metrics", "User Metrics", components["User Metrics"])
	writeMermaidNode(&sb, "business_problem", "Business Problem", components["Business Problem"])
	writeMermaidNode(&sb, "business_metrics", "Business Metrics", components["Business Metrics"])
	sb.WriteString("    end\n\n")

	// Row 4
	sb.WriteString("    subgraph row4[\"Investment\"]\n")
	writeMermaidNode(&sb, "budget", "Budget", components["Budget"])
	sb.WriteString("    end\n\n")

	// Connections
	sb.WriteString("    row1 --> row2\n")
	sb.WriteString("    row2 --> row3\n")
	sb.WriteString("    row3 --> row4\n")
	sb.WriteString("    end\n")

	// Styling
	sb.WriteString("\n    classDef teal fill:#085041,stroke:#5dcaa5,color:#9fe1cb\n")
	sb.WriteString("    classDef blue fill:#1e3a5f,stroke:#60a5fa,color:#bfdbfe\n")
	sb.WriteString("    classDef amber fill:#633806,stroke:#ef9f27,color:#fac775\n")
	sb.WriteString("    classDef purple fill:#3c3489,stroke:#afa9ec,color:#cecbf6\n")

	sb.WriteString("\n    class users,problems,solution_ideas teal\n")
	sb.WriteString("    class solutions_today,user_value,adoption blue\n")
	sb.WriteString("    class user_metrics,business_problem,business_metrics amber\n")
	sb.WriteString("    class budget purple\n")

	return sb.String()
}

func renderLeanUXMermaid(sm *aha.StrategicModel, _ *MermaidOptions) string {
	var sb strings.Builder

	sb.WriteString("%%{init: {'theme': 'dark'}}%%\n")
	sb.WriteString("flowchart TB\n\n")
	_, _ = fmt.Fprintf(&sb, "    subgraph title[\"%s\"]\n", escapeMermaid(sm.Name))
	sb.WriteString("    direction TB\n\n")

	components := componentMap(sm)

	// Row 1
	sb.WriteString("    subgraph row1[\"Business\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "business_problem", "Business Problem", components["Business Problem"])
	writeMermaidNode(&sb, "business_outcomes", "Business Outcomes", components["Business Outcomes"])
	sb.WriteString("    end\n\n")

	// Row 2
	sb.WriteString("    subgraph row2[\"Users\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "users", "Users", components["Users"])
	writeMermaidNode(&sb, "benefits", "Benefits", components["Benefits"])
	writeMermaidNode(&sb, "solutions", "Solutions", components["Solutions"])
	sb.WriteString("    end\n\n")

	// Row 3
	sb.WriteString("    subgraph row3[\"Validation\"]\n")
	writeMermaidNode(&sb, "hypotheses", "Hypotheses", components["Hypotheses"])
	sb.WriteString("    end\n\n")

	// Row 4
	sb.WriteString("    subgraph row4[\"Experiment\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "riskiest", "Riskiest Assumption", components["Riskiest Assumption"])
	writeMermaidNode(&sb, "experiment", "Smallest Experiment", components["Smallest Experiment"])
	sb.WriteString("    end\n\n")

	// Connections
	sb.WriteString("    row1 --> row2\n")
	sb.WriteString("    row2 --> row3\n")
	sb.WriteString("    row3 --> row4\n")
	sb.WriteString("    end\n")

	// Styling
	sb.WriteString("\n    classDef purple fill:#3c3489,stroke:#afa9ec,color:#cecbf6\n")
	sb.WriteString("    classDef teal fill:#085041,stroke:#5dcaa5,color:#9fe1cb\n")
	sb.WriteString("    classDef amber fill:#633806,stroke:#ef9f27,color:#fac775\n")
	sb.WriteString("    classDef coral fill:#712b13,stroke:#f0997b,color:#f5c4b3\n")

	sb.WriteString("\n    class business_problem,business_outcomes purple\n")
	sb.WriteString("    class users,benefits,solutions teal\n")
	sb.WriteString("    class hypotheses amber\n")
	sb.WriteString("    class riskiest,experiment coral\n")

	return sb.String()
}

func renderBMCMermaid(sm *aha.StrategicModel, _ *MermaidOptions) string {
	var sb strings.Builder

	sb.WriteString("%%{init: {'theme': 'dark'}}%%\n")
	sb.WriteString("flowchart TB\n\n")
	_, _ = fmt.Fprintf(&sb, "    subgraph title[\"%s\"]\n", escapeMermaid(sm.Name))
	sb.WriteString("    direction TB\n\n")

	components := componentMap(sm)

	// Main content
	sb.WriteString("    subgraph main[\"Business Model\"]\n")
	sb.WriteString("    direction LR\n\n")

	// Left: Partners
	sb.WriteString("    subgraph partners[\"Partners\"]\n")
	writeMermaidNode(&sb, "key_partners", "Key Partners", components["Key Partners"])
	sb.WriteString("    end\n\n")

	// Activities and Resources
	sb.WriteString("    subgraph activities[\"Operations\"]\n")
	sb.WriteString("    direction TB\n")
	writeMermaidNode(&sb, "key_activities", "Key Activities", components["Key Activities"])
	writeMermaidNode(&sb, "key_resources", "Key Resources", components["Key Resources"])
	sb.WriteString("    end\n\n")

	// Value Propositions
	sb.WriteString("    subgraph value[\"Value\"]\n")
	writeMermaidNode(&sb, "value_props", "Value Propositions", components["Value Propositions"])
	sb.WriteString("    end\n\n")

	// Customer Relationships and Channels
	sb.WriteString("    subgraph delivery[\"Delivery\"]\n")
	sb.WriteString("    direction TB\n")
	writeMermaidNode(&sb, "cust_rel", "Customer Relationships", components["Customer Relationships"])
	writeMermaidNode(&sb, "channels", "Channels", components["Channels"])
	sb.WriteString("    end\n\n")

	// Customer Segments
	sb.WriteString("    subgraph customers[\"Customers\"]\n")
	writeMermaidNode(&sb, "cust_segments", "Customer Segments", components["Customer Segments"])
	sb.WriteString("    end\n\n")

	// Connections within main
	sb.WriteString("    partners --> activities\n")
	sb.WriteString("    activities --> value\n")
	sb.WriteString("    value --> delivery\n")
	sb.WriteString("    delivery --> customers\n")
	sb.WriteString("    end\n\n")

	// Bottom row
	sb.WriteString("    subgraph finances[\"Finances\"]\n")
	sb.WriteString("    direction LR\n")
	writeMermaidNode(&sb, "cost_structure", "Cost Structure", components["Cost Structure"])
	writeMermaidNode(&sb, "revenue_streams", "Revenue Streams", components["Revenue Streams"])
	sb.WriteString("    end\n\n")

	sb.WriteString("    main --> finances\n")
	sb.WriteString("    end\n")

	// Styling
	sb.WriteString("\n    classDef teal fill:#085041,stroke:#5dcaa5,color:#9fe1cb\n")
	sb.WriteString("    classDef blue fill:#1e3a5f,stroke:#60a5fa,color:#bfdbfe\n")
	sb.WriteString("    classDef purple fill:#3c3489,stroke:#afa9ec,color:#cecbf6\n")
	sb.WriteString("    classDef amber fill:#633806,stroke:#ef9f27,color:#fac775\n")
	sb.WriteString("    classDef coral fill:#712b13,stroke:#f0997b,color:#f5c4b3\n")

	sb.WriteString("\n    class key_partners,cost_structure,revenue_streams teal\n")
	sb.WriteString("    class key_activities,key_resources blue\n")
	sb.WriteString("    class value_props purple\n")
	sb.WriteString("    class cust_rel,channels amber\n")
	sb.WriteString("    class cust_segments coral\n")

	return sb.String()
}

func renderGenericMermaid(sm *aha.StrategicModel, _ *MermaidOptions) string {
	var sb strings.Builder

	sb.WriteString("%%{init: {'theme': 'dark'}}%%\n")
	sb.WriteString("flowchart TB\n\n")
	_, _ = fmt.Fprintf(&sb, "    subgraph title[\"%s\"]\n", escapeMermaid(sm.Name))

	for i, comp := range sm.Components {
		id := fmt.Sprintf("block_%d", i)
		writeMermaidNode(&sb, id, comp.Name, &comp)
	}

	// Connect blocks sequentially
	for i := 0; i < len(sm.Components)-1; i++ {
		_, _ = fmt.Fprintf(&sb, "    block_%d --> block_%d\n", i, i+1)
	}

	sb.WriteString("    end\n")

	return sb.String()
}

func writeMermaidNode(sb *strings.Builder, id, title string, comp *aha.StrategicModelComponent) {
	escapedTitle := escapeMermaid(title)
	if comp != nil && comp.Description != "" {
		content := stripHTML(comp.Description)
		if len(content) > 40 {
			content = content[:40] + "..."
		}
		_, _ = fmt.Fprintf(sb, "    %s[\"%s<br/><small>%s</small>\"]\n", id, escapedTitle, escapeMermaid(content))
	} else {
		_, _ = fmt.Fprintf(sb, "    %s[\"%s\"]\n", id, escapedTitle)
	}
}

func escapeMermaid(s string) string {
	// Escape special characters for Mermaid
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	// Allow <br/> and <small> tags
	s = strings.ReplaceAll(s, "&lt;br/&gt;", "<br/>")
	s = strings.ReplaceAll(s, "&lt;small&gt;", "<small>")
	s = strings.ReplaceAll(s, "&lt;/small&gt;", "</small>")
	return s
}
