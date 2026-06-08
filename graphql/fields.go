package graphql

import (
	"context"
	"fmt"

	genql "github.com/Khan/genqlient/graphql"
	"github.com/grokify/aha-go/graphql/generated"
)

// FieldRequirement represents a field's requirement status.
type FieldRequirement struct {
	FieldID         string
	Required        bool
	ReadOnly        bool
	IsCustomField   bool
	CustomFieldKey  string
	CustomFieldName string
	CustomFieldType string
}

// ScreenRequirements holds all field requirements for a screen/record type.
type ScreenRequirements struct {
	RecordType string
	Fields     []FieldRequirement
}

// RequiredFieldIDs returns the IDs of all required fields.
func (sr *ScreenRequirements) RequiredFieldIDs() []string {
	var ids []string
	for _, f := range sr.Fields {
		if f.Required {
			ids = append(ids, f.FieldID)
		}
	}
	return ids
}

// RequiredCustomFieldKeys returns the keys of all required custom fields.
func (sr *ScreenRequirements) RequiredCustomFieldKeys() []string {
	var keys []string
	for _, f := range sr.Fields {
		if f.Required && f.IsCustomField {
			keys = append(keys, f.CustomFieldKey)
		}
	}
	return keys
}

// IsRequired checks if a field is required by ID or custom field key.
func (sr *ScreenRequirements) IsRequired(fieldIDOrKey string) bool {
	for _, f := range sr.Fields {
		if f.FieldID == fieldIDOrKey || f.CustomFieldKey == fieldIDOrKey {
			return f.Required
		}
	}
	return false
}

// GetFeatureRequirements fetches the field requirements for features in a project.
// It requires an existing feature ID from the project to get the screen definition.
func GetFeatureRequirements(ctx context.Context, client genql.Client, featureID string) (*ScreenRequirements, error) {
	resp, err := generated.GetFeatureScreenDefinition(ctx, client, featureID)
	if err != nil {
		return nil, fmt.Errorf("fetching feature screen definition: %w", err)
	}

	if resp.Feature.ScreenDefinition == nil {
		return nil, fmt.Errorf("no screen definition found for feature %s", featureID)
	}

	sr := &ScreenRequirements{
		RecordType: string(resp.Feature.ScreenDefinition.ScreenableType),
	}

	for _, control := range resp.Feature.ScreenDefinition.ScreenDefinitionControls {
		fieldID := ""
		if control.FieldId != nil {
			fieldID = *control.FieldId
		}
		fr := FieldRequirement{
			FieldID:  fieldID,
			Required: control.Required,
			ReadOnly: control.ReadOnly,
		}

		if control.CustomFieldDefinition != nil {
			fr.IsCustomField = true
			fr.CustomFieldKey = control.CustomFieldDefinition.Key
			fr.CustomFieldName = control.CustomFieldDefinition.Name
			fr.CustomFieldType = string(control.CustomFieldDefinition.Type)
		}

		sr.Fields = append(sr.Fields, fr)
	}

	return sr, nil
}

// CustomFieldDefinition represents a custom field definition.
type CustomFieldDefinition struct {
	ID      string
	Key     string
	Name    string
	Type    string
	Options []CustomFieldOption
}

// CustomFieldOption represents an option for choice-based custom fields.
type CustomFieldOption struct {
	ID       string
	Name     string
	Color    int
	Position int
}

// GetProjectCustomFieldDefinitions fetches all custom field definitions for a project.
func GetProjectCustomFieldDefinitions(ctx context.Context, client genql.Client, projectID string) ([]CustomFieldDefinition, error) {
	resp, err := generated.GetProjectCustomFields(ctx, client, projectID)
	if err != nil {
		return nil, fmt.Errorf("fetching project custom fields: %w", err)
	}

	var defs []CustomFieldDefinition
	for _, cf := range resp.Project.CustomFieldsRelatedToTeam {
		def := CustomFieldDefinition{
			ID:   cf.Id,
			Key:  cf.Key,
			Name: cf.Name,
			Type: string(cf.Type),
		}

		for _, opt := range cf.CustomFieldOptions {
			name := ""
			if opt.Name != nil {
				name = *opt.Name
			}
			color := 0
			if opt.Color != nil {
				color = *opt.Color
			}
			position := 0
			if opt.Position != nil {
				position = *opt.Position
			}
			def.Options = append(def.Options, CustomFieldOption{
				ID:       opt.Id,
				Name:     name,
				Color:    color,
				Position: position,
			})
		}

		defs = append(defs, def)
	}

	return defs, nil
}

// ValidateRequiredFields checks if all required fields are provided.
// Returns a list of missing required field IDs/keys.
func ValidateRequiredFields(requirements *ScreenRequirements, providedFields map[string]any) []string {
	var missing []string

	for _, f := range requirements.Fields {
		if !f.Required {
			continue
		}

		// Check by field ID
		if _, ok := providedFields[f.FieldID]; ok {
			continue
		}

		// Check by custom field key
		if f.IsCustomField {
			if _, ok := providedFields[f.CustomFieldKey]; ok {
				continue
			}
		}

		// Field is missing
		if f.IsCustomField {
			missing = append(missing, f.CustomFieldKey)
		} else {
			missing = append(missing, f.FieldID)
		}
	}

	return missing
}
