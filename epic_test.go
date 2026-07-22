package aha

import (
	"net/http"
	"testing"
)

// getEpicFixture mirrors a real Aha API response where progress is null
// (no progress recorded yet) - previously broke decoding since the schema
// didn't mark it nullable.
const getEpicFixture = `{
	"epic": {
		"id": "6936046092427283054",
		"reference_num": "IN-E-9",
		"name": "Sample epic",
		"description": {
			"id": "6936046092427283055",
			"body": "",
			"notable_id": "6936046092427283054",
			"notable_type": "Epic",
			"editor_version": 1,
			"created_at": "2021-03-05T06:01:16.109Z",
			"updated_at": "2021-03-05T06:01:16.109Z",
			"attachments": []
		},
		"progress": null,
		"created_at": "2021-03-05T06:01:16.109Z",
		"url": "https://test.aha.io/epics/IN-E-9",
		"resource": "https://test.aha.io/api/v1/epics/IN-E-9"
	}
}`

func TestGetEpicNullProgress(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getEpicFixture))
	})

	epic, err := client.GetEpic(t.Context(), "IN-E-9")
	if err != nil {
		t.Fatalf("GetEpic: %v", err)
	}
	if epic.Progress != 0 {
		t.Errorf("Progress = %v, want 0", epic.Progress)
	}
}
