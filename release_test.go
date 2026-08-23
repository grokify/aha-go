package aha

import (
	"net/http"
	"testing"
	"time"
)

const getReleaseFixture = `{
	"release": {
		"id": "6863840104778168921",
		"reference_num": "SE-R-2",
		"name": "MVP3",
		"release_date": "2021-07-30",
		"parking_lot": false,
		"theme": {
			"id": "6935912103773849090",
			"body": "<p>Ship the MVP3 milestone.</p>",
			"notable_id": "6863840104778168921",
			"notable_type": "Release",
			"editor_version": 1,
			"created_at": "2021-03-04T21:21:19.444Z",
			"updated_at": "2021-03-04T21:21:19.444Z",
			"attachments": []
		},
		"url": "https://test.aha.io/releases/SE-R-2",
		"resource": "https://test.aha.io/api/v1/releases/SE-R-2"
	}
}`

func TestGetReleaseTheme(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getReleaseFixture))
	})

	release, err := client.GetRelease(t.Context(), "SE-R-2")
	if err != nil {
		t.Fatalf("GetRelease: %v", err)
	}

	want := "<p>Ship the MVP3 milestone.</p>"
	if release.Theme != want {
		t.Errorf("Theme = %q, want %q", release.Theme, want)
	}
}

func TestUpdateReleaseTheme(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getReleaseFixture, &gotBody)

	if _, err := client.UpdateRelease(t.Context(), "SE-R-2", WithReleaseTheme("<p>Updated theme.</p>")); err != nil {
		t.Fatalf("UpdateRelease: %v", err)
	}

	release, ok := gotBody["release"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing release object: %v", gotBody)
	}
	if theme, _ := release["theme"].(string); theme != "<p>Updated theme.</p>" {
		t.Errorf("request theme = %q, want %q", theme, "<p>Updated theme.</p>")
	}
}

const getReleaseProgressFixture = `{
	"release": {
		"id": "6863840104778168921",
		"reference_num": "SE-R-2",
		"name": "MVP3",
		"progress_source": "progress_manual",
		"progress": 60
	}
}`

func TestGetReleaseProgress(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getReleaseProgressFixture))
	})

	release, err := client.GetRelease(t.Context(), "SE-R-2")
	if err != nil {
		t.Fatalf("GetRelease: %v", err)
	}

	if release.ProgressSource != "progress_manual" {
		t.Errorf("ProgressSource = %q, want %q", release.ProgressSource, "progress_manual")
	}
	if release.Progress == nil || *release.Progress != 60 {
		t.Errorf("Progress = %v, want 60", release.Progress)
	}
}

func TestUpdateReleaseProgress(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getReleaseFixture, &gotBody)

	if _, err := client.UpdateRelease(t.Context(), "SE-R-2",
		WithReleaseProgressSource("progress_manual"),
		WithReleaseProgress(90),
	); err != nil {
		t.Fatalf("UpdateRelease: %v", err)
	}

	release, ok := gotBody["release"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing release object: %v", gotBody)
	}
	if v, _ := release["progress_source"].(string); v != "progress_manual" {
		t.Errorf("request progress_source = %q, want %q", v, "progress_manual")
	}
	if v, _ := release["progress"].(float64); v != 90 {
		t.Errorf("request progress = %v, want 90", v)
	}
}

const getReleaseWorkflowStatusFixture = `{
	"release": {
		"id": "6863840104778168921",
		"reference_num": "SE-R-2",
		"name": "MVP3",
		"released": true,
		"workflow_status": {
			"id": "6935912103773849099",
			"name": "Shipped",
			"position": 5,
			"complete": true,
			"color": "#00FF00"
		}
	}
}`

func TestGetReleaseWorkflowStatus(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getReleaseWorkflowStatusFixture))
	})

	release, err := client.GetRelease(t.Context(), "SE-R-2")
	if err != nil {
		t.Fatalf("GetRelease: %v", err)
	}

	if !release.Released {
		t.Error("Released = false, want true")
	}
	if release.WorkflowStatus == nil {
		t.Fatal("WorkflowStatus = nil, want non-nil")
	}
	if release.WorkflowStatus.Name != "Shipped" {
		t.Errorf("WorkflowStatus.Name = %q, want %q", release.WorkflowStatus.Name, "Shipped")
	}
	if !release.WorkflowStatus.Complete {
		t.Error("WorkflowStatus.Complete = false, want true")
	}
}

func TestUpdateReleaseWorkflowStatus(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getReleaseWorkflowStatusFixture, &gotBody)

	if _, err := client.UpdateRelease(t.Context(), "SE-R-2", WithReleaseWorkflowStatus("Shipped")); err != nil {
		t.Fatalf("UpdateRelease: %v", err)
	}

	release, ok := gotBody["release"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing release object: %v", gotBody)
	}
	if v, _ := release["workflow_status"].(string); v != "Shipped" {
		t.Errorf("request workflow_status = %q, want %q", v, "Shipped")
	}
}

func TestUpdateReleaseExternalReleaseDate(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getReleaseFixture, &gotBody)

	date := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	if _, err := client.UpdateRelease(t.Context(), "SE-R-2", WithReleaseExternalReleaseDate(date)); err != nil {
		t.Fatalf("UpdateRelease: %v", err)
	}

	release, ok := gotBody["release"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing release object: %v", gotBody)
	}
	if v, _ := release["external_release_date"].(string); v != "2026-09-01" {
		t.Errorf("request external_release_date = %q, want %q", v, "2026-09-01")
	}
}

func TestUpdateReleaseDevelopmentStartedOn(t *testing.T) {
	var gotBody map[string]any
	client := newUpdateCaptureClient(t, getReleaseFixture, &gotBody)

	date := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	if _, err := client.UpdateRelease(t.Context(), "SE-R-2", WithReleaseDevelopmentStartedOn(date)); err != nil {
		t.Fatalf("UpdateRelease: %v", err)
	}

	release, ok := gotBody["release"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing release object: %v", gotBody)
	}
	if v, _ := release["development_started_on"].(string); v != "2026-07-15" {
		t.Errorf("request development_started_on = %q, want %q", v, "2026-07-15")
	}
}
