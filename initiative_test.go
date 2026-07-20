package aha

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
)

const getInitiativeFixture = `{
	"initiative": {
		"id": "7654719053155169681",
		"name": "Sample initiative",
		"reference_num": "SAVIN-S-76",
		"created_at": "2026-06-23T22:19:24.135Z",
		"custom_fields": [
			{
				"id": "7654721000724170303",
				"key": "aha_initiative_rank",
				"name": "Initiative Rank",
				"updatedAt": "2026-06-23T22:26:57Z",
				"type": "number",
				"value": "1.0"
			},
			{
				"id": "7662388581428353209",
				"key": "initiative_tags",
				"name": "Tags",
				"updatedAt": "2026-07-14T14:21:05Z",
				"type": "array",
				"value": ["Platform_Initiative"]
			}
		]
	}
}`

func TestGetInitiativeCustomFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getInitiativeFixture))
	})

	initiative, err := client.GetInitiative(t.Context(), "SAVIN-S-76")
	if err != nil {
		t.Fatalf("GetInitiative: %v", err)
	}

	if len(initiative.CustomFields) != 2 {
		t.Fatalf("len(CustomFields) = %d, want 2", len(initiative.CustomFields))
	}

	rank := initiative.CustomFields[0]
	if rank.Key != "aha_initiative_rank" {
		t.Errorf("CustomFields[0].Key = %q, want %q", rank.Key, "aha_initiative_rank")
	}
	if rank.Name != "Initiative Rank" {
		t.Errorf("CustomFields[0].Name = %q, want %q", rank.Name, "Initiative Rank")
	}
	if rank.Type != "number" {
		t.Errorf("CustomFields[0].Type = %q, want %q", rank.Type, "number")
	}
	// Value is preserved as raw JSON (jx.Raw), consistent with feature/idea
	// custom fields — callers unmarshal it into the type they expect.
	if got := rawValue(t, rank.Value); got != `"1.0"` {
		t.Errorf("CustomFields[0].Value = %s, want %s", got, `"1.0"`)
	}

	tags := initiative.CustomFields[1]
	if tags.Type != "array" {
		t.Errorf("CustomFields[1].Type = %q, want %q", tags.Type, "array")
	}
	if got := rawValue(t, tags.Value); got != `["Platform_Initiative"]` {
		t.Errorf("CustomFields[1].Value = %s, want %s", got, `["Platform_Initiative"]`)
	}
}

// rawValue asserts a custom field value is raw JSON and returns it as a string.
func rawValue(t *testing.T, v any) string {
	t.Helper()
	raw, ok := v.(jx.Raw)
	if !ok {
		t.Fatalf("custom field value type = %T, want jx.Raw", v)
	}
	return string(raw)
}
