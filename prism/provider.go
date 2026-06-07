// Package prism provides a CanvasProvider implementation for Aha.io Strategic Models.
// This allows exporting prism-roadmap canvases to Aha.io.
package prism

import (
	"context"
	"fmt"
	"strings"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/export"
)

// Provider implements export.CanvasProvider for Aha.io Strategic Models.
type Provider struct {
	client    *aha.Client
	productID string // Default product ID to use
}

// NewProvider creates a new Aha.io canvas provider.
func NewProvider(client *aha.Client, productID string) *Provider {
	return &Provider{
		client:    client,
		productID: productID,
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return "aha"
}

// SupportedTypes returns the canvas types supported by Aha Strategic Models.
func (p *Provider) SupportedTypes() []canvas.CanvasType {
	return []canvas.CanvasType{
		canvas.CanvasTypeOpportunity,
		canvas.CanvasTypeLeanUX,
		canvas.CanvasTypeBMC,
	}
}

// CreateCanvas creates a new canvas in Aha.io as a Strategic Model.
func (p *Provider) CreateCanvas(ctx context.Context, c *canvas.Canvas) (string, error) {
	if c == nil {
		return "", export.WrapError("aha", "CreateCanvas", fmt.Errorf("canvas is nil"))
	}

	kind, err := canvasTypeToAhaKind(c.Type)
	if err != nil {
		return "", export.WrapError("aha", "CreateCanvas", err)
	}

	meta := c.GetMetadata()
	name := "Untitled Canvas"
	if meta != nil && meta.Title != "" {
		name = meta.Title
	}

	// Create the strategic model
	sm, err := p.client.CreateStrategicModel(ctx, p.productID, kind,
		aha.WithStrategicModelName(name),
	)
	if err != nil {
		return "", export.WrapError("aha", "CreateCanvas", err)
	}

	// Update components with canvas content
	if err := p.updateComponents(ctx, sm.ID, c); err != nil {
		return sm.ID, export.WrapError("aha", "CreateCanvas", err)
	}

	return sm.ID, nil
}

// UpdateCanvas updates an existing canvas in Aha.io.
func (p *Provider) UpdateCanvas(ctx context.Context, externalID string, c *canvas.Canvas) error {
	if c == nil {
		return export.WrapError("aha", "UpdateCanvas", fmt.Errorf("canvas is nil"))
	}

	// Get the existing model to get component IDs
	sm, err := p.client.GetStrategicModel(ctx, externalID)
	if err != nil {
		return export.WrapError("aha", "UpdateCanvas", err)
	}

	if sm == nil {
		return export.WrapError("aha", "UpdateCanvas", export.ErrNotFound)
	}

	// Update components with new content
	return p.updateComponents(ctx, externalID, c)
}

// GetCanvas retrieves a canvas from Aha.io.
func (p *Provider) GetCanvas(ctx context.Context, externalID string) (*canvas.Canvas, error) {
	sm, err := p.client.GetStrategicModel(ctx, externalID)
	if err != nil {
		return nil, export.WrapError("aha", "GetCanvas", err)
	}
	if sm == nil {
		return nil, export.WrapError("aha", "GetCanvas", export.ErrNotFound)
	}

	return strategicModelToCanvas(sm)
}

// DeleteCanvas removes a canvas from Aha.io.
// Note: Aha API may not support deletion; this returns an error if unsupported.
func (p *Provider) DeleteCanvas(ctx context.Context, externalID string) error {
	// Aha Strategic Models API does not support deletion via API
	return export.WrapError("aha", "DeleteCanvas", fmt.Errorf("deletion not supported by Aha Strategic Models API"))
}

// canvasTypeToAhaKind maps canvas types to Aha Strategic Model kinds.
func canvasTypeToAhaKind(ct canvas.CanvasType) (string, error) {
	switch ct {
	case canvas.CanvasTypeOpportunity:
		return "Opportunity", nil
	case canvas.CanvasTypeLeanUX:
		return "Lean Canvas", nil
	case canvas.CanvasTypeBMC:
		return "Business Model", nil
	default:
		return "", export.ErrUnsupportedType
	}
}

// ahaKindToCanvasType maps Aha Strategic Model kinds to canvas types.
func ahaKindToCanvasType(kind string) (canvas.CanvasType, error) {
	switch strings.ToLower(kind) {
	case "opportunity":
		return canvas.CanvasTypeOpportunity, nil
	case "lean canvas":
		return canvas.CanvasTypeLeanUX, nil
	case "business model":
		return canvas.CanvasTypeBMC, nil
	default:
		return "", export.ErrUnsupportedType
	}
}

// updateComponents updates Aha Strategic Model components with canvas content.
func (p *Provider) updateComponents(ctx context.Context, modelID string, c *canvas.Canvas) error {
	// Get existing model to find component IDs
	sm, err := p.client.GetStrategicModel(ctx, modelID)
	if err != nil {
		return err
	}

	// Build a map of component name to ID
	componentIDs := make(map[string]string)
	for _, comp := range sm.Components {
		componentIDs[comp.Name] = comp.ID
	}

	// Generate content for each block based on canvas type
	blocks := canvasToBlocks(c)

	// Update each block
	for name, content := range blocks {
		compID, ok := componentIDs[name]
		if !ok {
			// Component doesn't exist in this model type, skip
			continue
		}

		_, err := p.client.UpdateStrategicModelComponent(ctx, modelID, compID, content)
		if err != nil {
			return fmt.Errorf("failed to update component %q: %w", name, err)
		}
	}

	return nil
}

// canvasToBlocks extracts block content from a canvas.
func canvasToBlocks(c *canvas.Canvas) map[string]string {
	blocks := make(map[string]string)

	switch c.Type {
	case canvas.CanvasTypeOpportunity:
		if c.Opportunity != nil {
			blocks = opportunityToBlocks(c.Opportunity)
		}
	case canvas.CanvasTypeLeanUX:
		if c.LeanUX != nil {
			blocks = leanUXToBlocks(c.LeanUX)
		}
	case canvas.CanvasTypeBMC:
		if c.BMC != nil {
			blocks = bmcToBlocks(c.BMC)
		}
	}

	return blocks
}

// opportunityToBlocks converts OpportunityCanvas to block content.
func opportunityToBlocks(oc *canvas.OpportunityCanvas) map[string]string {
	blocks := make(map[string]string)

	// Users & Customers
	var users []string
	for _, u := range oc.Users {
		users = append(users, fmt.Sprintf("<p>%s: %s</p>", u.Name, u.Description))
	}
	blocks["Users & Customers"] = strings.Join(users, "\n")

	// Problems
	var problems []string
	for _, prob := range oc.Problems {
		problems = append(problems, fmt.Sprintf("<p>%s</p>", prob.Description))
	}
	blocks["Problems"] = strings.Join(problems, "\n")

	// Solution Ideas
	blocks["Solution Ideas"] = formatList(oc.SolutionIdeas)

	// Solutions Today
	var solutions []string
	for _, s := range oc.CurrentSolutions {
		solutions = append(solutions, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", s.Name, s.Description))
	}
	blocks["Solutions Today"] = strings.Join(solutions, "\n")

	// User Value
	blocks["User Value"] = formatList(oc.UserValue)

	// Adoption Strategy
	blocks["Adoption Strategy"] = formatList(oc.AdoptionStrategy)

	// User Metrics
	blocks["User Metrics"] = formatList(oc.UserMetrics)

	// Business Problem
	blocks["Business Problem"] = fmt.Sprintf("<p>%s</p>", oc.BusinessProblem)

	// Business Metrics
	blocks["Business Metrics"] = formatList(oc.BusinessMetrics)

	// Budget
	if oc.Budget != nil {
		var budgetParts []string
		if oc.Budget.TimeEstimate != "" {
			budgetParts = append(budgetParts, fmt.Sprintf("Time: %s", oc.Budget.TimeEstimate))
		}
		if oc.Budget.TeamSize != "" {
			budgetParts = append(budgetParts, fmt.Sprintf("Team: %s", oc.Budget.TeamSize))
		}
		if oc.Budget.CostEstimate != "" {
			budgetParts = append(budgetParts, fmt.Sprintf("Cost: %s", oc.Budget.CostEstimate))
		}
		blocks["Budget"] = formatList(budgetParts)
	}

	return blocks
}

// leanUXToBlocks converts LeanUXCanvas to block content.
func leanUXToBlocks(lc *canvas.LeanUXCanvas) map[string]string {
	blocks := make(map[string]string)

	// Business Problem
	blocks["Business Problem"] = fmt.Sprintf("<p>%s</p>", lc.BusinessProblem)

	// Business Outcomes
	var outcomes []string
	for _, o := range lc.BusinessOutcomes {
		outcomes = append(outcomes, fmt.Sprintf("<p>%s</p>", o.Description))
	}
	blocks["Business Outcomes"] = strings.Join(outcomes, "\n")

	// Users
	var users []string
	for _, u := range lc.Users {
		users = append(users, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", u.Name, u.Description))
	}
	blocks["Users"] = strings.Join(users, "\n")

	// Benefits (User Outcomes)
	var benefits []string
	for _, o := range lc.UserOutcomes {
		benefits = append(benefits, fmt.Sprintf("<p>%s</p>", o.Description))
	}
	blocks["Benefits"] = strings.Join(benefits, "\n")

	// Solutions
	var solutions []string
	for _, s := range lc.Solutions {
		solutions = append(solutions, fmt.Sprintf("<p>%s</p>", s.Description))
	}
	blocks["Solutions"] = strings.Join(solutions, "\n")

	// Hypotheses
	var hypotheses []string
	for _, h := range lc.Hypotheses {
		hypotheses = append(hypotheses, fmt.Sprintf("<p><strong>We believe:</strong> %s</p><p><strong>Will result in:</strong> %s</p>", h.WeBelieve, h.WillResultIn))
	}
	blocks["Hypotheses"] = strings.Join(hypotheses, "\n")

	// Riskiest Assumption
	blocks["Riskiest Assumption"] = fmt.Sprintf("<p>%s</p>", lc.RiskiestAssumption)

	// Smallest Experiment
	if lc.Experiment != nil {
		blocks["Smallest Experiment"] = fmt.Sprintf("<p><strong>%s</strong></p><p>Method: %s</p><p>Success Criteria: %s</p>",
			lc.Experiment.Description,
			lc.Experiment.Method,
			lc.Experiment.SuccessCriteria)
	}

	return blocks
}

// bmcToBlocks converts BusinessModelCanvas to block content.
func bmcToBlocks(bmc *canvas.BusinessModelCanvas) map[string]string {
	blocks := make(map[string]string)

	// Key Partners
	var partners []string
	for _, p := range bmc.KeyPartnerships {
		partners = append(partners, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", p.Partner, p.Description))
	}
	blocks["Key Partners"] = strings.Join(partners, "\n")

	// Key Activities
	var activities []string
	for _, a := range bmc.KeyActivities {
		activities = append(activities, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", a.Name, a.Description))
	}
	blocks["Key Activities"] = strings.Join(activities, "\n")

	// Key Resources
	var resources []string
	for _, r := range bmc.KeyResources {
		resources = append(resources, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", r.Name, r.Description))
	}
	blocks["Key Resources"] = strings.Join(resources, "\n")

	// Value Propositions
	var valueProps []string
	for _, vp := range bmc.ValuePropositions {
		valueProps = append(valueProps, fmt.Sprintf("<p>%s</p>", vp.Description))
	}
	blocks["Value Propositions"] = strings.Join(valueProps, "\n")

	// Customer Relationships
	var relationships []string
	for _, r := range bmc.CustomerRelationships {
		relationships = append(relationships, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", r.Type, r.Description))
	}
	blocks["Customer Relationships"] = strings.Join(relationships, "\n")

	// Channels
	var channels []string
	for _, ch := range bmc.Channels {
		channels = append(channels, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", ch.Name, ch.Description))
	}
	blocks["Channels"] = strings.Join(channels, "\n")

	// Customer Segments
	var segments []string
	for _, seg := range bmc.CustomerSegments {
		segments = append(segments, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", seg.Name, seg.Description))
	}
	blocks["Customer Segments"] = strings.Join(segments, "\n")

	// Cost Structure
	var costs []string
	for _, c := range bmc.CostStructure {
		label := c.Category
		if label == "" {
			label = c.Type
		}
		costs = append(costs, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", label, c.Description))
	}
	blocks["Cost Structure"] = strings.Join(costs, "\n")

	// Revenue Streams
	var revenues []string
	for _, r := range bmc.RevenueStreams {
		revenues = append(revenues, fmt.Sprintf("<p><strong>%s</strong>: %s</p>", r.Type, r.Description))
	}
	blocks["Revenue Streams"] = strings.Join(revenues, "\n")

	return blocks
}

// strategicModelToCanvas converts an Aha Strategic Model to a Canvas.
func strategicModelToCanvas(sm *aha.StrategicModel) (*canvas.Canvas, error) {
	canvasType, err := ahaKindToCanvasType(sm.Kind)
	if err != nil {
		return nil, err
	}

	// Build a map of component name to content
	componentContent := make(map[string]string)
	for _, comp := range sm.Components {
		componentContent[comp.Name] = comp.Description
	}

	meta := canvas.Metadata{
		ID:    sm.ID,
		Title: sm.Name,
	}

	switch canvasType {
	case canvas.CanvasTypeOpportunity:
		oc := &canvas.OpportunityCanvas{
			Metadata:        meta,
			BusinessProblem: stripHTML(componentContent["Business Problem"]),
		}
		return canvas.NewOpportunity(oc), nil

	case canvas.CanvasTypeLeanUX:
		lc := &canvas.LeanUXCanvas{
			Metadata:           meta,
			BusinessProblem:    stripHTML(componentContent["Business Problem"]),
			RiskiestAssumption: stripHTML(componentContent["Riskiest Assumption"]),
		}
		return canvas.NewLeanUX(lc), nil

	case canvas.CanvasTypeBMC:
		bmc := &canvas.BusinessModelCanvas{
			Metadata: meta,
		}
		return canvas.NewBMC(bmc), nil
	}

	return nil, export.ErrUnsupportedType
}

// formatList formats a slice of strings as an HTML list.
func formatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("<ul>")
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li>%s</li>", item))
	}
	sb.WriteString("</ul>")
	return sb.String()
}

// stripHTML removes HTML tags from a string (basic implementation).
func stripHTML(s string) string {
	// Basic implementation - remove tags
	result := s
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return strings.TrimSpace(result)
}
