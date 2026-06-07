//nolint:dupl // Render tests follow consistent pattern across formats
package render

import (
	"strings"
	"testing"

	aha "github.com/grokify/aha-go"
)

func TestRenderCanvasMermaid(t *testing.T) {
	tests := []struct {
		name     string
		model    *aha.StrategicModel
		wantType string
	}{
		{
			name: "Opportunity Canvas",
			model: &aha.StrategicModel{
				ID:   "1",
				Name: "Test Opportunity",
				Kind: "Opportunity",
			},
			wantType: "flowchart TB",
		},
		{
			name: "Lean UX Canvas",
			model: &aha.StrategicModel{
				ID:   "2",
				Name: "Test Lean UX",
				Kind: "Lean Canvas",
			},
			wantType: "flowchart TB",
		},
		{
			name: "Business Model Canvas",
			model: &aha.StrategicModel{
				ID:   "3",
				Name: "Test BMC",
				Kind: "Business Model",
			},
			wantType: "flowchart TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmd := RenderCanvasMermaid(tt.model, nil)

			if !strings.Contains(mmd, tt.wantType) {
				t.Errorf("Mermaid output should contain %q", tt.wantType)
			}
			if !strings.Contains(mmd, tt.model.Name) {
				t.Errorf("Mermaid output should contain canvas name %q", tt.model.Name)
			}
			if !strings.Contains(mmd, "%%{init:") {
				t.Error("Mermaid output should contain init directive")
			}
		})
	}
}

func TestRenderOpportunityMermaidBlocks(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Opportunity",
	}

	mmd := RenderCanvasMermaid(sm, nil)

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
		if !strings.Contains(mmd, block) {
			t.Errorf("Mermaid output should contain block %q", block)
		}
	}
}

func TestRenderLeanUXMermaidBlocks(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Lean Canvas",
	}

	mmd := RenderCanvasMermaid(sm, nil)

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
		if !strings.Contains(mmd, block) {
			t.Errorf("Mermaid output should contain block %q", block)
		}
	}
}

func TestRenderBMCMermaidBlocks(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Business Model",
	}

	mmd := RenderCanvasMermaid(sm, nil)

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
		if !strings.Contains(mmd, block) {
			t.Errorf("Mermaid output should contain block %q", block)
		}
	}
}

func TestRenderMermaidWithContent(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test Canvas",
		Kind: "Opportunity",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Users & Customers", Description: "<p>Enterprise customers</p>"},
		},
	}

	mmd := RenderCanvasMermaid(sm, nil)

	if !strings.Contains(mmd, "Enterprise customers") {
		t.Error("Mermaid output should contain component content")
	}
}

func TestMermaidStyling(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "1",
		Name: "Test",
		Kind: "Opportunity",
	}

	mmd := RenderCanvasMermaid(sm, nil)

	// Check for class definitions
	if !strings.Contains(mmd, "classDef teal") {
		t.Error("Mermaid output should contain teal class definition")
	}
	if !strings.Contains(mmd, "classDef blue") {
		t.Error("Mermaid output should contain blue class definition")
	}
	if !strings.Contains(mmd, "classDef amber") {
		t.Error("Mermaid output should contain amber class definition")
	}
}

func TestDefaultMermaidOptions(t *testing.T) {
	opts := DefaultMermaidOptions()

	if opts.Theme != "dark" {
		t.Errorf("Default theme = %s, want dark", opts.Theme)
	}
	if opts.Direction != "TB" {
		t.Errorf("Default direction = %s, want TB", opts.Direction)
	}
}
