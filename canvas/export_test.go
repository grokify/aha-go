package canvas

import (
	"strings"
	"testing"

	aha "github.com/grokify/aha-go"
)

func TestExportModelSVG(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:           "123",
		ReferenceNum: "PROD-SM-1",
		Name:         "Test Canvas",
		Kind:         "Opportunity",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Users & Customers", Description: "<p>Test users</p>"},
			{ID: "2", Name: "Problems", Description: "<p>Test problems</p>"},
		},
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatSVG,
		Width:  800,
		Height: 600,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	svg := string(output)

	// Check SVG structure
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG should start with <svg tag")
	}
	if !strings.HasSuffix(svg, "</svg>") {
		t.Error("SVG should end with </svg> tag")
	}
	if !strings.Contains(svg, "Test Canvas") {
		t.Error("SVG should contain canvas name")
	}
	if !strings.Contains(svg, "viewBox") {
		t.Error("SVG should have viewBox attribute")
	}
}

func TestExportModelJSON(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:           "123",
		ReferenceNum: "PROD-SM-1",
		Name:         "Test Canvas",
		Kind:         "Opportunity",
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatJSON,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	json := string(output)

	// Check JSON structure
	if !strings.Contains(json, `"ID": "123"`) {
		t.Error("JSON should contain ID field")
	}
	if !strings.Contains(json, `"Name": "Test Canvas"`) {
		t.Error("JSON should contain Name field")
	}
	if !strings.Contains(json, `"Kind": "Opportunity"`) {
		t.Error("JSON should contain Kind field")
	}
}

func TestExportModelUnsupportedFormat(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "123",
		Name: "Test Canvas",
	}

	_, err := ExportModel(sm, ExportOptions{
		Format: "xml",
	})
	if err == nil {
		t.Error("ExportModel() should return error for unsupported format")
	}
}

func TestExportOptionsDefaults(t *testing.T) {
	opts := DefaultExportOptions()

	if opts.Format != FormatSVG {
		t.Errorf("Default format = %s, want %s", opts.Format, FormatSVG)
	}
	if opts.Width != 800 {
		t.Errorf("Default width = %d, want 800", opts.Width)
	}
	if opts.Height != 600 {
		t.Errorf("Default height = %d, want 600", opts.Height)
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := SupportedFormats()

	if len(formats) != 4 {
		t.Errorf("SupportedFormats() returned %d formats, want 4", len(formats))
	}

	hasFormat := func(f ExportFormat) bool {
		for _, format := range formats {
			if format == f {
				return true
			}
		}
		return false
	}

	if !hasFormat(FormatSVG) {
		t.Error("SupportedFormats() should include SVG")
	}
	if !hasFormat(FormatJSON) {
		t.Error("SupportedFormats() should include JSON")
	}
	if !hasFormat(FormatD2) {
		t.Error("SupportedFormats() should include D2")
	}
	if !hasFormat(FormatMermaid) {
		t.Error("SupportedFormats() should include Mermaid")
	}
}

func TestExportModelD2(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "123",
		Name: "Test Canvas",
		Kind: "Opportunity",
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatD2,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	d2 := string(output)

	if !strings.Contains(d2, "Opportunity Canvas") {
		t.Error("D2 should contain canvas type")
	}
	if !strings.Contains(d2, "Test Canvas") {
		t.Error("D2 should contain canvas name")
	}
	if !strings.Contains(d2, "grid-") {
		t.Error("D2 should contain grid layout")
	}
}

func TestExportModelMermaid(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "123",
		Name: "Test Canvas",
		Kind: "Opportunity",
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatMermaid,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	mmd := string(output)

	if !strings.Contains(mmd, "flowchart TB") {
		t.Error("Mermaid should contain flowchart directive")
	}
	if !strings.Contains(mmd, "Test Canvas") {
		t.Error("Mermaid should contain canvas name")
	}
}

func TestExportModelLeanUXSVG(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "123",
		Name: "Lean UX Test",
		Kind: "Lean Canvas",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Business Problem"},
			{ID: "2", Name: "Business Outcomes"},
		},
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatSVG,
		Width:  800,
		Height: 600,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	svg := string(output)
	if !strings.Contains(svg, "Lean UX Canvas") {
		t.Error("SVG should contain Lean UX Canvas descriptor")
	}
}

func TestExportModelBMCSVG(t *testing.T) {
	sm := &aha.StrategicModel{
		ID:   "123",
		Name: "BMC Test",
		Kind: "Business Model",
		Components: []aha.StrategicModelComponent{
			{ID: "1", Name: "Key Partners"},
			{ID: "2", Name: "Value Propositions"},
		},
	}

	output, err := ExportModel(sm, ExportOptions{
		Format: FormatSVG,
		Width:  800,
		Height: 600,
	})
	if err != nil {
		t.Fatalf("ExportModel() error = %v", err)
	}

	svg := string(output)
	if !strings.Contains(svg, "Business Model Canvas") {
		t.Error("SVG should contain Business Model Canvas descriptor")
	}
}
