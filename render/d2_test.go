package render

import (
	"strings"
	"testing"

	aha "github.com/grokify/aha-go"
)

func TestRenderCanvasD2(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d2 := RenderCanvasD2(tt.model, nil)

			if !strings.Contains(d2, tt.wantDesc) {
				t.Errorf("D2 output should contain %q", tt.wantDesc)
			}
			if !strings.Contains(d2, tt.model.Name) {
				t.Errorf("D2 output should contain canvas name %q", tt.model.Name)
			}
			if !strings.Contains(d2, "grid-") {
				t.Error("D2 output should contain grid layout")
			}
		})
	}
}

func TestRenderOpportunityD2Blocks(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Opportunity",
	}

	d2 := RenderCanvasD2(sm, nil)

	expectedBlocks := []string{
		"Users & Customers",
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
		if !strings.Contains(d2, block) {
			t.Errorf("D2 output should contain block %q", block)
		}
	}
}

func TestRenderBMCD2Blocks(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Business Model",
	}

	d2 := RenderCanvasD2(sm, nil)

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
		if !strings.Contains(d2, block) {
			t.Errorf("D2 output should contain block %q", block)
		}
	}
}

func TestRenderD2WithContent(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test Canvas",
		Kind: "Opportunity",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Users & Customers", Description: "<p>Enterprise customers</p>"},
		},
	}

	d2 := RenderCanvasD2(sm, nil)

	if !strings.Contains(d2, "Enterprise customers") {
		t.Error("D2 output should contain component content")
	}
	if !strings.Contains(d2, "tooltip:") {
		t.Error("D2 output should contain tooltip for content")
	}
}

func TestDefaultD2Options(t *testing.T) {
	opts := DefaultD2Options()

	if opts.Theme != 0 {
		t.Errorf("Default theme = %d, want 0", opts.Theme)
	}
	if opts.Layout != "dagre" {
		t.Errorf("Default layout = %s, want dagre", opts.Layout)
	}
	if opts.Direction != "right" {
		t.Errorf("Default direction = %s, want right", opts.Direction)
	}
}
