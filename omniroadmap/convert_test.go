package omniroadmap

import (
	"testing"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/omniroadmap-core/provider"
)

func TestStatusFromWorkflowStatus(t *testing.T) {
	tests := []struct {
		name string
		ws   *aha.WorkflowStatus
		want provider.StatusCategory
	}{
		{"nil", nil, ""},
		{"complete", &aha.WorkflowStatus{Name: "Shipped", Complete: true, Position: 5}, provider.StatusCategoryDone},
		{"position zero", &aha.WorkflowStatus{Name: "Backlog", Position: 0}, provider.StatusCategoryTodo},
		{"cancelled name", &aha.WorkflowStatus{Name: "Cancelled", Position: 3}, provider.StatusCategoryCanceled},
		{"in progress", &aha.WorkflowStatus{Name: "In Development", Position: 2}, provider.StatusCategoryInProgress},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusFromWorkflowStatus(tt.ws)
			if tt.ws == nil {
				if got != nil {
					t.Fatalf("statusFromWorkflowStatus(nil) = %+v, want nil", got)
				}
				return
			}
			if got.Category != tt.want {
				t.Errorf("Category = %q, want %q", got.Category, tt.want)
			}
			if got.Name != tt.ws.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.ws.Name)
			}
		})
	}
}

func TestItemFromFeature(t *testing.T) {
	f := &aha.Feature{
		ID:           "123",
		ReferenceNum: "SE-1",
		Name:         "Test feature",
		Description:  "A description",
		URL:          "https://test.aha.io/features/SE-1",
		WorkflowStatus: &aha.WorkflowStatus{
			ID: "ws-1", Name: "Shipped", Complete: true,
		},
		Release: &aha.Release{ID: "rel-1"},
		Tags:    []string{"backend"},
	}

	item := itemFromFeature(f)

	if item.ID != "aha:123" {
		t.Errorf("ID = %q, want %q", item.ID, "aha:123")
	}
	if item.Provider != "aha" {
		t.Errorf("Provider = %q, want %q", item.Provider, "aha")
	}
	if item.Kind != provider.ItemKindFeature {
		t.Errorf("Kind = %q, want %q", item.Kind, provider.ItemKindFeature)
	}
	if item.SourceRef != "SE-1" {
		t.Errorf("SourceRef = %q, want %q", item.SourceRef, "SE-1")
	}
	if item.Status == nil || item.Status.Category != provider.StatusCategoryDone {
		t.Errorf("Status = %+v, want Category=done", item.Status)
	}
	if item.ReleaseID != "aha:rel-1" {
		t.Errorf("ReleaseID = %q, want %q", item.ReleaseID, "aha:rel-1")
	}
	if len(item.Tags) != 1 || item.Tags[0] != "backend" {
		t.Errorf("Tags = %v, want [backend]", item.Tags)
	}
}

func TestItemFromFeatureMeta(t *testing.T) {
	m := aha.FeatureMeta{ID: "123", ReferenceNum: "SE-1", Name: "Test feature"}
	item := itemFromFeatureMeta(m)

	if item.ID != "aha:123" {
		t.Errorf("ID = %q, want %q", item.ID, "aha:123")
	}
	if item.Kind != provider.ItemKindFeature {
		t.Errorf("Kind = %q, want %q", item.Kind, provider.ItemKindFeature)
	}
	if item.Description != "" {
		t.Errorf("Description = %q, want empty (list responses are meta-only)", item.Description)
	}
}

func TestCustomFieldDefinitionsFromAha(t *testing.T) {
	defs := []aha.CustomFieldDefinition{
		{ID: "1", Key: "priority", Name: "Priority", Type: "string", CustomFieldableType: "Feature"},
	}
	got := customFieldDefinitionsFromAha(defs)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Key != "priority" {
		t.Errorf("Key = %q, want %q", got[0].Key, "priority")
	}
	if got[0].Metadata["aha.customfieldable_type"] != "Feature" {
		t.Errorf("Metadata[aha.customfieldable_type] = %v, want %q", got[0].Metadata["aha.customfieldable_type"], "Feature")
	}
}
