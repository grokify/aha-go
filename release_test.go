package aha

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
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
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body: %v", err)
		}
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("unmarshaling request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getReleaseFixture))
	})

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
