package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// TemplateBlock defines a block in a strategic model template.
type TemplateBlock struct {
	Name        string // Block name/label
	Description string // Block description/placeholder text
	Row         int    // Row position (0-indexed)
	Column      int    // Column position (0-indexed)
	RowSpan     int    // Number of rows to span (default 1)
	ColSpan     int    // Number of columns to span (default 1)
}

// TemplateConfig defines a strategic model template configuration.
type TemplateConfig struct {
	Name        string          // Template name
	Description string          // Template description
	Rows        int             // Number of rows in grid
	Columns     int             // Number of columns in grid
	Blocks      []TemplateBlock // Block definitions
}

// PredefinedTemplates contains common canvas template configurations.
var PredefinedTemplates = map[string]TemplateConfig{
	"capability-stack": {
		Name:        "Capability Stack",
		Description: "Layered capability model for product/platform planning",
		Rows:        4,
		Columns:     1,
		Blocks: []TemplateBlock{
			{Name: "Experience Layer", Description: "User-facing capabilities and interfaces", Row: 0, Column: 0},
			{Name: "Application Layer", Description: "Business logic and application services", Row: 1, Column: 0},
			{Name: "Platform Layer", Description: "Shared platform capabilities and APIs", Row: 2, Column: 0},
			{Name: "Infrastructure Layer", Description: "Infrastructure and foundational services", Row: 3, Column: 0},
		},
	},
	"maturity-model": {
		Name:        "Maturity Model",
		Description: "Capability maturity assessment grid",
		Rows:        5,
		Columns:     5,
		Blocks: []TemplateBlock{
			// Headers (could be implemented as special blocks or labels)
			{Name: "Domain", Description: "Capability domain", Row: 0, Column: 0},
			{Name: "M1 - Initial", Description: "Ad-hoc, reactive", Row: 0, Column: 1},
			{Name: "M2 - Managed", Description: "Documented, repeatable", Row: 0, Column: 2},
			{Name: "M3 - Defined", Description: "Standardized, proactive", Row: 0, Column: 3},
			{Name: "M4 - Measured", Description: "Quantified, controlled", Row: 0, Column: 4},
			// Row blocks for domains would be added dynamically
		},
	},
	"opportunity-patton": {
		Name:        "Opportunity Canvas (Patton)",
		Description: "Jeff Patton's 10-block opportunity assessment",
		Rows:        3,
		Columns:     3,
		Blocks: []TemplateBlock{
			// Row 1: Problem Space
			{Name: "Problems", Description: "What problems are we solving?", Row: 0, Column: 0},
			{Name: "Users", Description: "Who has these problems?", Row: 0, Column: 1},
			{Name: "Current Solutions", Description: "How do they solve it today?", Row: 0, Column: 2},
			// Row 2: Solution Space
			{Name: "Value Proposition", Description: "What's our unique solution?", Row: 1, Column: 0},
			{Name: "User Value", Description: "What value do users get?", Row: 1, Column: 1},
			{Name: "Business Value", Description: "What value does the business get?", Row: 1, Column: 2},
			// Row 3: Validation
			{Name: "Assumptions", Description: "What must be true?", Row: 2, Column: 0},
			{Name: "Risks", Description: "What could go wrong?", Row: 2, Column: 1},
			{Name: "Budget/Timeline", Description: "What resources are needed?", Row: 2, Column: 2},
		},
	},
	"feature-canvas": {
		Name:        "Feature Canvas (Efimov)",
		Description: "Nikita Efimov's feature planning canvas",
		Rows:        4,
		Columns:     2,
		Blocks: []TemplateBlock{
			{Name: "Idea Statement", Description: "One-sentence feature description", Row: 0, Column: 0, ColSpan: 2},
			// Problem Area (left)
			{Name: "Situations", Description: "When/where does the problem occur?", Row: 1, Column: 0},
			{Name: "Problems", Description: "What specific problems exist?", Row: 2, Column: 0},
			{Name: "Value", Description: "What value does solving this provide?", Row: 3, Column: 0},
			// Solution Area (right)
			{Name: "Capabilities", Description: "What must the solution do?", Row: 1, Column: 1},
			{Name: "Restrictions", Description: "What constraints exist?", Row: 2, Column: 1},
			{Name: "Limitations", Description: "What's out of scope?", Row: 3, Column: 1},
		},
	},
}

// CreateTemplate creates a new strategic model template in Aha via browser automation.
// This requires the user to be logged in first.
func (c *Client) CreateTemplate(ctx context.Context, config TemplateConfig) error {
	if c.browser == nil {
		return fmt.Errorf("browser not connected, call Connect first")
	}

	page := c.browser.MustPage()
	defer func() { _ = page.Close() }()

	// Navigate to strategic model templates settings
	templatesURL := fmt.Sprintf("%s/settings/account/strategic_model_templates", c.BaseURL())
	if err := page.Timeout(c.timeout).Navigate(templatesURL); err != nil {
		return fmt.Errorf("failed to navigate to templates page: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page load: %w", err)
	}

	// Wait for page to stabilize
	time.Sleep(2 * time.Second)

	// Click "Add strategic model template" or similar button
	// Note: Exact selectors will need adjustment based on Aha's actual UI
	addBtn, err := page.Element("a[href*='new'], button:has-text('Add'), .add-template-btn")
	if err != nil {
		return fmt.Errorf("failed to find add template button: %w", err)
	}
	if err := addBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click add button: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for form load: %w", err)
	}

	// Fill in template name
	nameInput, err := page.Element("input[name*='name'], input#template_name")
	if err != nil {
		return fmt.Errorf("failed to find name input: %w", err)
	}
	if err := nameInput.Input(config.Name); err != nil {
		return fmt.Errorf("failed to input template name: %w", err)
	}

	// Fill in description if field exists
	descInput, err := page.Element("textarea[name*='description'], input[name*='description']")
	if err == nil && config.Description != "" {
		if err := descInput.Input(config.Description); err != nil {
			return fmt.Errorf("failed to input description: %w", err)
		}
	}

	// Configure grid dimensions
	// Note: This depends on Aha's UI for template configuration
	// May need to set rows/columns via dropdown or number inputs

	// Add blocks
	for _, block := range config.Blocks {
		if err := c.addTemplateBlock(page, block); err != nil {
			return fmt.Errorf("failed to add block %q: %w", block.Name, err)
		}
	}

	// Save template
	saveBtn, err := page.Element("button[type='submit'], input[type='submit'], .save-btn")
	if err != nil {
		return fmt.Errorf("failed to find save button: %w", err)
	}
	if err := saveBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click save: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for save completion: %w", err)
	}

	return nil
}

// addTemplateBlock adds a single block to the template being edited.
func (c *Client) addTemplateBlock(page *rod.Page, block TemplateBlock) error {
	// Click add block button
	addBlockBtn, err := page.Element(".add-block-btn, button:has-text('Add block')")
	if err != nil {
		return fmt.Errorf("failed to find add block button: %w", err)
	}
	if err := addBlockBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click add block: %w", err)
	}

	// Wait for block form/modal
	time.Sleep(500 * time.Millisecond)

	// Fill in block name
	blockNameInput, err := page.Element(".block-name-input, input[name*='block'][name*='name']")
	if err != nil {
		return fmt.Errorf("failed to find block name input: %w", err)
	}
	if err := blockNameInput.Input(block.Name); err != nil {
		return fmt.Errorf("failed to input block name: %w", err)
	}

	// Fill in block description/placeholder
	if block.Description != "" {
		blockDescInput, err := page.Element(".block-description-input, textarea[name*='block']")
		if err == nil {
			if err := blockDescInput.Input(block.Description); err != nil {
				return fmt.Errorf("failed to input block description: %w", err)
			}
		}
	}

	// Set position if UI supports it
	// This would need to be implemented based on Aha's actual block positioning UI

	// Confirm/save block
	confirmBtn, err := page.Element(".confirm-block-btn, button:has-text('Add'), button:has-text('Save')")
	if err == nil {
		if err := confirmBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
			return fmt.Errorf("failed to confirm block: %w", err)
		}
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

// CreatePredefinedTemplate creates one of the predefined templates.
func (c *Client) CreatePredefinedTemplate(ctx context.Context, templateName string) error {
	config, ok := PredefinedTemplates[templateName]
	if !ok {
		return fmt.Errorf("unknown predefined template: %s", templateName)
	}
	return c.CreateTemplate(ctx, config)
}

// ListPredefinedTemplates returns the names of all predefined templates.
func ListPredefinedTemplates() []string {
	names := make([]string, 0, len(PredefinedTemplates))
	for name := range PredefinedTemplates {
		names = append(names, name)
	}
	return names
}
