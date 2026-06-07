package render

import (
	"strings"
	"testing"

	aha "github.com/grokify/aha-go"
)

func TestRenderCanvasSVG(t *testing.T) {
	tests := []struct {
		name     string
		model    *aha.StrategicModel
		wantDesc string
	}{
		{
			name: "Opportunity Canvas",
			model: &aha.StrategicModel{
				ID:   "1",
				Name: "Test Opportunity",
				Kind: "Opportunity",
			},
			wantDesc: "Opportunity Canvas",
		},
		{
			name: "Lean UX Canvas",
			model: &aha.StrategicModel{
				ID:   "2",
				Name: "Test Lean UX",
				Kind: "Lean Canvas",
			},
			wantDesc: "Lean UX Canvas",
		},
		{
			name: "Business Model Canvas",
			model: &aha.StrategicModel{
				ID:   "3",
				Name: "Test BMC",
				Kind: "Business Model",
			},
			wantDesc: "Business Model Canvas",
		},
		{
			name: "Unknown Canvas",
			model: &aha.StrategicModel{
				ID:   "4",
				Name: "Test Unknown",
				Kind: "Custom Type",
			},
			wantDesc: "Custom Type Canvas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svg := RenderCanvasSVG(tt.model, nil)

			if !strings.HasPrefix(svg, "<svg") {
				t.Error("SVG should start with <svg tag")
			}
			if !strings.HasSuffix(svg, "</svg>") {
				t.Error("SVG should end with </svg> tag")
			}
			if !strings.Contains(svg, tt.model.Name) {
				t.Errorf("SVG should contain canvas name %q", tt.model.Name)
			}
			if !strings.Contains(svg, tt.wantDesc) {
				t.Errorf("SVG should contain desc %q", tt.wantDesc)
			}
		})
	}
}

func TestRenderCanvasSVGWithOptions(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Opportunity",
	}

	opts := &SVGOptions{
		Width:      1200,
		Height:     800,
		FontFamily: "Arial",
		Theme:      "dark",
	}

	svg := RenderCanvasSVG(sm, opts)

	if !strings.Contains(svg, "1200") {
		t.Error("SVG should contain width 1200")
	}
	if !strings.Contains(svg, "800") {
		t.Error("SVG should contain height 800")
	}
	if !strings.Contains(svg, "Arial") {
		t.Error("SVG should contain font family Arial")
	}
}

func TestRenderCanvasSVGWithComponents(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test Canvas",
		Kind: "Opportunity",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Users & Customers", Description: "<p>Enterprise customers</p>"},
			{ID: "2", Name: "Problems", Description: "<p>Slow deployment times</p>"},
		},
	}

	svg := RenderCanvasSVG(sm, nil)

	// Check blocks are rendered
	if !strings.Contains(svg, "Users &amp; Customers") {
		t.Error("SVG should contain 'Users & Customers' block")
	}
	if !strings.Contains(svg, "Problems") {
		t.Error("SVG should contain 'Problems' block")
	}
	// Content should be included (stripped of HTML)
	if !strings.Contains(svg, "Enterprise customers") {
		t.Error("SVG should contain component content")
	}
}

func TestDefaultSVGOptions(t *testing.T) {
	opts := DefaultSVGOptions()

	if opts.Width != 800 {
		t.Errorf("Default width = %d, want 800", opts.Width)
	}
	if opts.Height != 600 {
		t.Errorf("Default height = %d, want 600", opts.Height)
	}
	if opts.Theme != "dark" {
		t.Errorf("Default theme = %s, want dark", opts.Theme)
	}
	if opts.FontFamily == "" {
		t.Error("Default font family should not be empty")
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<p>Hello</p>", "Hello"},
		{"<ul><li>Item 1</li><li>Item 2</li></ul>", "Item 1  Item 2"},
		{"No tags here", "No tags here"},
		{"<strong>Bold</strong> and <em>italic</em>", "Bold  and  italic"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripHTML(tt.input)
			// Normalize whitespace for comparison
			got = strings.Join(strings.Fields(got), " ")
			want := strings.Join(strings.Fields(tt.want), " ")
			if got != want {
				t.Errorf("stripHTML(%q) = %q, want %q", tt.input, got, want)
			}
		})
	}
}

func TestRenderOpportunitySVGBlockCount(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Opportunity",
	}

	svg := RenderCanvasSVG(sm, nil)

	// Opportunity canvas has 10 blocks
	expectedBlocks := []string{
		"Users &amp; Customers",
		"Problems",
		"Solution Ideas",
		"Solutions Today",
		"User Value",
		"Adoption Strategy",
		"User Metrics",
		"Business Problem",
		"Business Metrics",
		"Budget",
	}

	for _, block := range expectedBlocks {
		if !strings.Contains(svg, block) {
			t.Errorf("SVG should contain block %q", block)
		}
	}
}

func TestRenderLeanUXSVGBlockCount(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Lean Canvas",
	}

	svg := RenderCanvasSVG(sm, nil)

	// Lean UX canvas has 8 blocks
	expectedBlocks := []string{
		"Business Problem",
		"Business Outcomes",
		"Users",
		"Benefits",
		"Solutions",
		"Hypotheses",
		"Riskiest Assumption",
		"Smallest Experiment",
	}

	for _, block := range expectedBlocks {
		if !strings.Contains(svg, block) {
			t.Errorf("SVG should contain block %q", block)
		}
	}
}

func TestRenderBMCSVGBlockCount(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Business Model",
	}

	svg := RenderCanvasSVG(sm, nil)

	// BMC canvas has 9 blocks
	expectedBlocks := []string{
		"Key Partners",
		"Key Activities",
		"Key Resources",
		"Value Propositions",
		"Customer Relationships",
		"Channels",
		"Customer Segments",
		"Cost Structure",
		"Revenue Streams",
	}

	for _, block := range expectedBlocks {
		if !strings.Contains(svg, block) {
			t.Errorf("SVG should contain block %q", block)
		}
	}
}
