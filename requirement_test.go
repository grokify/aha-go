package aha

import (
	"net/http"
	"testing"
)

// getRequirementFixture mirrors a real Aha API response where
// original_estimate, remaining_estimate, work_done, and assigned_to_user are
// all null - previously broke decoding since the schema didn't mark them
// nullable.
const getRequirementFixture = `{
	"requirement": {
		"id": "6944094681093364153",
		"reference_num": "SE-3-1",
		"name": "Sample requirement",
		"description": {
			"id": "6944094681093364154",
			"body": "<p>Add the required fields.</p>",
			"notable_id": "6944094681093364153",
			"notable_type": "Requirement",
			"editor_version": 1,
			"created_at": "2021-03-05T06:01:16.109Z",
			"updated_at": "2021-03-05T06:01:16.109Z",
			"attachments": []
		},
		"original_estimate": null,
		"remaining_estimate": null,
		"work_done": null,
		"assigned_to_user": null,
		"created_at": "2021-03-05T06:01:16.109Z",
		"url": "https://test.aha.io/requirements/SE-3-1",
		"resource": "https://test.aha.io/api/v1/requirements/SE-3-1"
	}
}`

func TestGetRequirementNullEstimatesAndAssignee(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getRequirementFixture))
	})

	req, err := client.GetRequirement(t.Context(), "SE-3-1")
	if err != nil {
		t.Fatalf("GetRequirement: %v", err)
	}

	want := "<p>Add the required fields.</p>"
	if req.Description != want {
		t.Errorf("Description = %q, want %q", req.Description, want)
	}
	if req.OriginalEstimate != 0 || req.RemainingEstimate != 0 || req.WorkDone != 0 {
		t.Errorf("estimates = (%v, %v, %v), want all 0", req.OriginalEstimate, req.RemainingEstimate, req.WorkDone)
	}
}
