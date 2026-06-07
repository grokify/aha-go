// Package canvas provides high-level operations for Aha.io strategic canvases.
package canvas

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	aha "github.com/grokify/aha-go"
)

// Kind represents a canvas type.
type Kind string

// Canvas kinds.
const (
	KindOpportunity Kind = "Opportunity"
	KindLeanUX      Kind = "Lean Canvas"
	KindBMC         Kind = "Business Model"
)

// OpportunityBlocks defines the standard block names for Opportunity Canvas.
var OpportunityBlocks = []string{
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

// LeanUXBlocks defines the standard block names for Lean UX Canvas.
var LeanUXBlocks = []string{
	"Business Problem",
	"Business Outcomes",
	"Users",
	"Benefits",
	"Solutions",
	"Hypotheses",
	"Riskiest Assumption",
	"Smallest Experiment",
}

// BMCBlocks defines the standard block names for Business Model Canvas.
var BMCBlocks = []string{
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

// CreateOptions configures canvas creation.
type CreateOptions struct {
	ProductID   string
	Name        string
	Description string
	Kind        Kind
}

// Create creates a new strategic canvas in Aha.io.
func Create(ctx context.Context, client *aha.Client, opts CreateOptions) (*aha.StrategicModel, error) {
	if opts.ProductID == "" {
		return nil, fmt.Errorf("product ID is required")
	}
	if opts.Name == "" {
		return nil, fmt.Errorf("canvas name is required")
	}
	if opts.Kind == "" {
		return nil, fmt.Errorf("canvas kind is required")
	}

	var createOpts []aha.CreateStrategicModelOption
	createOpts = append(createOpts, aha.WithStrategicModelName(opts.Name))
	if opts.Description != "" {
		createOpts = append(createOpts, aha.WithStrategicModelDescription(opts.Description))
	}

	return client.CreateStrategicModel(ctx, opts.ProductID, string(opts.Kind), createOpts...)
}

// UpdateOptions configures canvas update.
type UpdateOptions struct {
	CanvasID string
	Blocks   map[string]string // Block name -> HTML content
}

// UpdateResult contains the result of an update operation.
type UpdateResult struct {
	SuccessCount int
	ErrorCount   int
	Errors       []BlockError
	Unknown      []string // Block names not found in canvas
}

// BlockError represents an error updating a specific block.
type BlockError struct {
	BlockName string
	Err       error
}

// Update updates canvas blocks with new content.
func Update(ctx context.Context, client *aha.Client, opts UpdateOptions) (*UpdateResult, error) {
	if opts.CanvasID == "" {
		return nil, fmt.Errorf("canvas ID is required")
	}
	if len(opts.Blocks) == 0 {
		return nil, fmt.Errorf("no blocks to update")
	}

	// Get the canvas to find component IDs
	sm, err := client.GetStrategicModel(ctx, opts.CanvasID)
	if err != nil {
		return nil, fmt.Errorf("failed to get canvas: %w", err)
	}

	// Build component name -> ID map
	componentIDs := make(map[string]string)
	for _, comp := range sm.Components {
		componentIDs[comp.Name] = comp.ID
	}

	result := &UpdateResult{}

	// Find unknown blocks
	for name := range opts.Blocks {
		if _, ok := componentIDs[name]; !ok {
			result.Unknown = append(result.Unknown, name)
		}
	}

	// Update each block
	for name, content := range opts.Blocks {
		compID, ok := componentIDs[name]
		if !ok {
			continue // Skip unknown blocks
		}

		_, err := client.UpdateStrategicModelComponent(ctx, opts.CanvasID, compID, content)
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, BlockError{BlockName: name, Err: err})
			continue
		}
		result.SuccessCount++
	}

	return result, nil
}

// LoadBlocksFromFile loads block content from a JSON file.
func LoadBlocksFromFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var blocks map[string]string
	if err := json.Unmarshal(data, &blocks); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return blocks, nil
}

// GetAvailableBlocks returns the available block names for a canvas.
func GetAvailableBlocks(ctx context.Context, client *aha.Client, canvasID string) ([]string, error) {
	sm, err := client.GetStrategicModel(ctx, canvasID)
	if err != nil {
		return nil, err
	}

	blocks := make([]string, len(sm.Components))
	for i, comp := range sm.Components {
		blocks[i] = comp.Name
	}
	return blocks, nil
}

// BlocksForKind returns the standard block names for a canvas kind.
func BlocksForKind(kind Kind) []string {
	switch kind {
	case KindOpportunity:
		return OpportunityBlocks
	case KindLeanUX:
		return LeanUXBlocks
	case KindBMC:
		return BMCBlocks
	default:
		return nil
	}
}
