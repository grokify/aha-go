package aha

import (
	"context"

	"github.com/grokify/aha-go/internal/api"
)

// CustomFieldDefinition represents a custom field definition.
type CustomFieldDefinition struct {
	ID                  string
	Name                string
	Key                 string
	Type                string
	CustomFieldableType string // Feature, Initiative, Epic, etc.
	InternalName        string
	Position            int64
	APIType             string
	AllowsOtherOption   bool
}

// CustomFieldOption represents an option for a select custom field.
type CustomFieldOption struct {
	ID       string
	Value    string
	Position int64
	Color    string
}

// ListCustomFieldDefinitions returns all custom field definitions.
func (c *Client) ListCustomFieldDefinitions(ctx context.Context) ([]CustomFieldDefinition, error) {
	resp, err := c.apiClient.ListCustomFieldDefinitions(ctx)
	if err != nil {
		return nil, wrapError("ListCustomFieldDefinitions", err)
	}

	var defs []CustomFieldDefinition
	for _, d := range resp.CustomFieldDefinitions {
		defs = append(defs, customFieldDefinitionFromAPI(d))
	}
	return defs, nil
}

// ListProductCustomFieldDefinitions returns custom field definitions for a product.
func (c *Client) ListProductCustomFieldDefinitions(ctx context.Context, productID string) ([]CustomFieldDefinition, error) {
	resp, err := c.apiClient.ListProductCustomFieldDefinitions(ctx, api.ListProductCustomFieldDefinitionsParams{
		ProductID: productID,
	})
	if err != nil {
		return nil, wrapError("ListProductCustomFieldDefinitions", err)
	}

	var defs []CustomFieldDefinition
	for _, d := range resp.CustomFieldDefinitions {
		defs = append(defs, customFieldDefinitionFromAPI(d))
	}
	return defs, nil
}

// ListCustomFieldOptions returns options for a select custom field.
func (c *Client) ListCustomFieldOptions(ctx context.Context, fieldID string) ([]CustomFieldOption, error) {
	resp, err := c.apiClient.ListCustomFieldOptions(ctx, api.ListCustomFieldOptionsParams{
		ID: fieldID,
	})
	if err != nil {
		return nil, wrapError("ListCustomFieldOptions", err)
	}

	var opts []CustomFieldOption
	for _, o := range resp.CustomFieldOptions {
		opts = append(opts, customFieldOptionFromAPI(o))
	}
	return opts, nil
}

func customFieldDefinitionFromAPI(d api.CustomFieldDefinition) CustomFieldDefinition {
	def := CustomFieldDefinition{
		ID:                  d.ID,
		Name:                d.Name,
		Key:                 d.Key,
		Type:                d.Type,
		CustomFieldableType: d.CustomFieldableType,
	}

	if v, ok := d.InternalName.Get(); ok {
		def.InternalName = v
	}
	if v, ok := d.Position.Get(); ok {
		def.Position = int64(v)
	}
	if v, ok := d.APIType.Get(); ok {
		def.APIType = v
	}
	if v, ok := d.AllowsOtherOption.Get(); ok {
		def.AllowsOtherOption = v
	}

	return def
}

func customFieldOptionFromAPI(o api.CustomFieldOption) CustomFieldOption {
	opt := CustomFieldOption{}

	if v, ok := o.ID.Get(); ok {
		opt.ID = v
	}
	if v, ok := o.Value.Get(); ok {
		opt.Value = v
	}
	if v, ok := o.Position.Get(); ok {
		opt.Position = int64(v)
	}
	if v, ok := o.Color.Get(); ok {
		opt.Color = v
	}

	return opt
}
