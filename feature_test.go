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

func TestUpdateFeatureOriginalEstimate(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getFeatureFixture, &gotBody)

	if _, err := client.UpdateFeature(t.Context(), "SE-2", WithUpdateFeatureOriginalEstimate("2d")); err != nil {
		t.Fatalf("UpdateFeature: %v", err)
	}

	feature, ok := gotBody["feature"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing feature object: %v", gotBody)
	}
	if v, _ := feature["original_estimate_text"].(string); v != "2d" {
		t.Errorf("request original_estimate_text = %q, want %q", v, "2d")
	}
}

func TestUpdateFeatureRemainingEstimate(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getFeatureFixture, &gotBody)

	if _, err := client.UpdateFeature(t.Context(), "SE-2", WithUpdateFeatureRemainingEstimate("4h")); err != nil {
		t.Fatalf("UpdateFeature: %v", err)
	}

	feature, ok := gotBody["feature"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing feature object: %v", gotBody)
	}
	if v, _ := feature["remaining_estimate_text"].(string); v != "4h" {
		t.Errorf("request remaining_estimate_text = %q, want %q", v, "4h")
	}
}

const getFeatureProgressFixture = `{
	"feature": {
		"id": "6863830667240702272",
		"reference_num": "SE-2",
		"name": "Authenticate via External emails",
		"created_at": "2021-03-04T20:25:26.282Z",
		"progress_source": "progress_manual",
		"progress": 42.5
	}
}`

func TestGetFeatureProgress(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getFeatureProgressFixture))
	})

	feature, err := client.GetFeature(t.Context(), "SE-2")
	if err != nil {
		t.Fatalf("GetFeature: %v", err)
	}

	if feature.ProgressSource != "progress_manual" {
		t.Errorf("ProgressSource = %q, want %q", feature.ProgressSource, "progress_manual")
	}
	if feature.Progress == nil || *feature.Progress != 42.5 {
		t.Errorf("Progress = %v, want 42.5", feature.Progress)
	}
}

func TestUpdateFeatureProgress(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getFeatureFixture, &gotBody)

	if _, err := client.UpdateFeature(t.Context(), "SE-2",
		WithUpdateFeatureProgressSource("progress_manual"),
		WithUpdateFeatureProgress(75),
	); err != nil {
		t.Fatalf("UpdateFeature: %v", err)
	}

	feature, ok := gotBody["feature"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing feature object: %v", gotBody)
	}
	if v, _ := feature["progress_source"].(string); v != "progress_manual" {
		t.Errorf("request progress_source = %q, want %q", v, "progress_manual")
	}
	if v, _ := feature["progress"].(float64); v != 75 {
		t.Errorf("request progress = %v, want 75", v)
	}
}

func TestUpdateFeatureEpic(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getFeatureFixture, &gotBody)

	if _, err := client.UpdateFeature(t.Context(), "SE-2", WithUpdateFeatureEpic("SE-E-1")); err != nil {
		t.Fatalf("UpdateFeature: %v", err)
	}

	feature, ok := gotBody["feature"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing feature object: %v", gotBody)
	}
	if v, _ := feature["epic"].(string); v != "SE-E-1" {
		t.Errorf("request epic = %q, want %q", v, "SE-E-1")
	}
}

func TestUpdateFeatureReleasePhase(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getFeatureFixture, &gotBody)

	if _, err := client.UpdateFeature(t.Context(), "SE-2", WithUpdateFeatureReleasePhase("Development")); err != nil {
		t.Fatalf("UpdateFeature: %v", err)
	}

	feature, ok := gotBody["feature"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing feature object: %v", gotBody)
	}
	if v, _ := feature["release_phase"].(string); v != "Development" {
		t.Errorf("request release_phase = %q, want %q", v, "Development")
	}
}
