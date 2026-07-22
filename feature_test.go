package aha

import (
	"net/http"
	"testing"
)

// getFeatureFixture mirrors a real Aha API response: description is a
// DescriptionObject (not a plain string), and assigned_to_user is null when
// the feature has no assignee. Both previously broke decoding.
const getFeatureFixture = `{
	"feature": {
		"id": "6863830667240702272",
		"reference_num": "SE-2",
		"name": "Authenticate via External emails",
		"description": {
			"id": "7329927923491557026",
			"body": "<p>Support login via external identity providers.</p>",
			"notable_id": "6863830667240702272",
			"notable_type": "Feature",
			"editor_version": 1,
			"created_at": "2024-01-30T16:23:56.359Z",
			"updated_at": "2024-01-30T16:23:56.359Z",
			"attachments": []
		},
		"created_at": "2021-03-04T20:25:26.282Z",
		"assigned_to_user": null,
		"url": "https://test.aha.io/features/SE-2",
		"resource": "https://test.aha.io/api/v1/features/SE-2"
	}
}`

func TestGetFeatureDescriptionAndNullAssignee(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getFeatureFixture))
	})

	feature, err := client.GetFeature(t.Context(), "SE-2")
	if err != nil {
		t.Fatalf("GetFeature: %v", err)
	}

	want := "<p>Support login via external identity providers.</p>"
	if feature.Description != want {
		t.Errorf("Description = %q, want %q", feature.Description, want)
	}
	if feature.AssignedTo != nil {
		t.Errorf("AssignedTo = %+v, want nil", feature.AssignedTo)
	}
}
