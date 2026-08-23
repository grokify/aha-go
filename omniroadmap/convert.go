package omniroadmap

import (
	"fmt"
	"strings"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/omniroadmap-core/provider"
)

func canonicalID(sourceID string) string {
	return fmt.Sprintf("%s:%s", providerName, sourceID)
}

// itemFromFeatureMeta converts an aha.FeatureMeta (as returned by list
// endpoints) to a canonical Item. List responses only carry meta-level
// fields — Description/Status/Progress/etc. are unavailable here; call
// GetItem for full fidelity on a specific record.
func itemFromFeatureMeta(m aha.FeatureMeta) provider.Item {
	item := provider.Item{
		ID:        canonicalID(m.ID),
		Provider:  providerName,
		SourceID:  m.ID,
		SourceRef: m.ReferenceNum,
		SourceURL: m.URL,
		Kind:      provider.ItemKindFeature,
		Name:      m.Name,
	}
	if !m.CreatedAt.IsZero() {
		item.CreatedAt = &m.CreatedAt
	}
	return item
}

// itemFromFeature converts a full aha.Feature (as returned by GetFeature) to
// a canonical Item.
func itemFromFeature(f *aha.Feature) *provider.Item {
	item := &provider.Item{
		ID:          canonicalID(f.ID),
		Provider:    providerName,
		SourceID:    f.ID,
		SourceRef:   f.ReferenceNum,
		SourceURL:   f.URL,
		Kind:        provider.ItemKindFeature,
		Name:        f.Name,
		Description: f.Description,
		Status:      statusFromWorkflowStatus(f.WorkflowStatus),
		Progress:    f.Progress,
		StartDate:   f.StartDate,
		DueDate:     f.DueDate,
		UpdatedAt:   f.UpdatedAt,
		Tags:        f.Tags,
		Metadata: map[string]any{
			"aha.comments_count":  f.CommentsCount,
			"aha.progress_source": f.ProgressSource,
			"aha.work_units":      f.WorkUnits,
		},
	}
	if !f.CreatedAt.IsZero() {
		item.CreatedAt = &f.CreatedAt
	}
	if f.Release != nil {
		item.ReleaseID = canonicalReleaseID(f.Release.ID)
	}
	if f.AssignedTo != nil {
		item.Owner = &provider.Person{
			ID:    f.AssignedTo.ID,
			Name:  f.AssignedTo.Name(),
			Email: f.AssignedTo.Email,
		}
	}
	item.CustomFields = customFieldsFromAha(f.CustomFields)
	return item
}

// itemFromInitiativeMeta converts an aha.InitiativeMeta to a canonical Item.
func itemFromInitiativeMeta(m aha.InitiativeMeta) provider.Item {
	item := provider.Item{
		ID:        canonicalID(m.ID),
		Provider:  providerName,
		SourceID:  m.ID,
		SourceRef: m.ReferenceNum,
		SourceURL: m.URL,
		Kind:      provider.ItemKindInitiative,
		Name:      m.Name,
	}
	if !m.CreatedAt.IsZero() {
		item.CreatedAt = &m.CreatedAt
	}
	return item
}

// itemFromInitiative converts a full aha.Initiative to a canonical Item.
func itemFromInitiative(i *aha.Initiative) *provider.Item {
	item := &provider.Item{
		ID:          canonicalID(i.ID),
		Provider:    providerName,
		SourceID:    i.ID,
		SourceRef:   i.ReferenceNum,
		SourceURL:   i.URL,
		Kind:        provider.ItemKindInitiative,
		Name:        i.Name,
		Description: i.Description,
		Status:      statusFromWorkflowStatus(i.WorkflowStatus),
		StartDate:   i.StartDate,
		DueDate:     i.EndDate,
		UpdatedAt:   i.UpdatedAt,
		Metadata: map[string]any{
			"aha.value":     i.Value,
			"aha.effort":    i.Effort,
			"aha.presented": i.Presented,
			"aha.position":  i.Position,
			"aha.color":     i.Color,
		},
	}
	if i.Progress != 0 {
		progress := i.Progress
		item.Progress = &progress
	}
	if !i.CreatedAt.IsZero() {
		item.CreatedAt = &i.CreatedAt
	}
	if i.Epic != nil {
		item.ParentID = canonicalID(i.Epic.ID)
	}
	item.CustomFields = customFieldsFromAha(i.CustomFields)
	return item
}

// itemFromEpicMeta converts an aha.EpicMeta to a canonical Item.
func itemFromEpicMeta(m aha.EpicMeta) provider.Item {
	item := provider.Item{
		ID:        canonicalID(m.ID),
		Provider:  providerName,
		SourceID:  m.ID,
		SourceRef: m.ReferenceNum,
		SourceURL: m.URL,
		Kind:      provider.ItemKindEpic,
		Name:      m.Name,
	}
	if !m.CreatedAt.IsZero() {
		item.CreatedAt = &m.CreatedAt
	}
	return item
}

// itemFromEpic converts a full aha.Epic to a canonical Item.
func itemFromEpic(e *aha.Epic) *provider.Item {
	item := &provider.Item{
		ID:          canonicalID(e.ID),
		Provider:    providerName,
		SourceID:    e.ID,
		SourceRef:   e.ReferenceNum,
		SourceURL:   e.URL,
		Kind:        provider.ItemKindEpic,
		Name:        e.Name,
		Description: e.Description,
		Status:      statusFromWorkflowStatus(e.WorkflowStatus),
		StartDate:   e.StartDate,
		DueDate:     e.DueDate,
		UpdatedAt:   e.UpdatedAt,
		Tags:        e.Tags,
		Metadata: map[string]any{
			"aha.comments_count":  e.CommentsCount,
			"aha.progress_source": e.ProgressSource,
			"aha.position":        e.Position,
			"aha.color":           e.Color,
		},
	}
	if e.Progress != 0 {
		progress := e.Progress
		item.Progress = &progress
	}
	if !e.CreatedAt.IsZero() {
		item.CreatedAt = &e.CreatedAt
	}
	if e.Release != nil {
		item.ReleaseID = canonicalReleaseID(e.Release.ID)
	}
	if e.Initiative != nil {
		item.ParentID = canonicalID(e.Initiative.ID)
	}
	return item
}

func canonicalReleaseID(sourceID string) string {
	return fmt.Sprintf("%s:%s", providerName, sourceID)
}

// releaseFromAha converts a full aha.Release to a canonical Release.
func releaseFromAha(r *aha.Release) provider.Release {
	rel := provider.Release{
		ID:          canonicalReleaseID(r.ID),
		Provider:    providerName,
		SourceID:    r.ID,
		SourceRef:   r.ReferenceNum,
		SourceURL:   r.URL,
		Name:        r.Name,
		StartDate:   r.StartDate,
		ReleaseDate: r.ReleaseDate,
		Released:    r.Released,
		Metadata: map[string]any{
			"aha.parking_lot":           r.ParkingLot,
			"aha.theme":                 r.Theme,
			"aha.external_release_date": r.ExternalReleaseDate,
			"aha.progress_source":       r.ProgressSource,
		},
	}
	if r.Progress != nil {
		progress := *r.Progress
		rel.Progress = &progress
	}
	rel.Status = statusFromWorkflowStatus(r.WorkflowStatus)
	return rel
}

// statusFromWorkflowStatus converts an aha.WorkflowStatus to a canonical
// Status, normalizing into a StatusCategory. Aha's own workflow statuses
// don't carry a done/in-progress/todo category, only a Complete flag and a
// Position — so the heuristic below is a best-effort normalization:
// Complete is authoritative for "done"; otherwise, position 0 is treated as
// "todo" (Aha's first workflow status is conventionally "the backlog");
// anything else in between is "in_progress". Cancellation isn't
// distinguishable from "done" this way — Aha has no separate cancelled
// flag — so canceled statuses will be categorized as done.
func statusFromWorkflowStatus(ws *aha.WorkflowStatus) *provider.Status {
	if ws == nil {
		return nil
	}
	category := provider.StatusCategoryInProgress
	switch {
	case ws.Complete:
		category = provider.StatusCategoryDone
	case ws.Position == 0:
		category = provider.StatusCategoryTodo
	case strings.Contains(strings.ToLower(ws.Name), "cancel"):
		category = provider.StatusCategoryCanceled
	}
	return &provider.Status{
		ID:       ws.ID,
		Name:     ws.Name,
		Category: category,
		Complete: ws.Complete,
		Color:    ws.Color,
		Position: ws.Position,
	}
}

// customFieldsFromAha converts aha-go's per-record custom field values to
// canonical CustomFields (already a near-identical shape).
func customFieldsFromAha(fields []aha.CustomField) []provider.CustomField {
	if len(fields) == 0 {
		return nil
	}
	out := make([]provider.CustomField, len(fields))
	for i, f := range fields {
		out[i] = provider.CustomField{
			Key:   f.Key,
			Name:  f.Name,
			Value: f.Value,
			Type:  f.Type,
		}
	}
	return out
}

// customFieldDefinitionsFromAha converts custom field *definitions*
// (schema/metadata, from ListCustomFieldDefinitions — distinct from
// per-record CustomField values) to canonical CustomFieldDefinitions.
func customFieldDefinitionsFromAha(defs []aha.CustomFieldDefinition) []provider.CustomFieldDefinition {
	out := make([]provider.CustomFieldDefinition, len(defs))
	for i, d := range defs {
		out[i] = provider.CustomFieldDefinition{
			ID:   d.ID,
			Key:  d.Key,
			Name: d.Name,
			Type: d.Type,
			Metadata: map[string]any{
				"aha.customfieldable_type": d.CustomFieldableType,
				"aha.internal_name":        d.InternalName,
				"aha.position":             d.Position,
				"aha.api_type":             d.APIType,
				"aha.allows_other_option":  d.AllowsOtherOption,
			},
		}
	}
	return out
}
