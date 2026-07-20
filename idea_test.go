package aha

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(
		WithSubdomain("test"),
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

const listIdeasFixture = `{
	"ideas": [
		{
			"id": "7142957303733180934",
			"name": "Increase number of customproperties available",
			"reference_num": "EIC-I-3558",
			"score": 0,
			"created_at": "2022-09-13T20:02:26.615Z",
			"updated_at": "2026-07-17T18:30:07.384Z",
			"votes": 205,
			"workflow_status": {
				"id": "7104817458763743132",
				"name": "Accepted",
				"position": 7,
				"complete": false,
				"color": "#64b80b"
			},
			"categories": [
				{
					"id": "7438762818542387377",
					"name": "Administration - Identity Repository",
					"parent_id": null,
					"project_id": "6956010895139277648",
					"created_at": "2024-11-18T23:18:55.029Z"
				}
			],
			"url": "https://test.aha.io/ideas/ideas/EIC-I-3558",
			"resource": "https://test.aha.io/api/v1/ideas/EIC-I-3558"
		}
	],
	"pagination": {
		"total_records": 1,
		"total_pages": 1,
		"current_page": 1
	}
}`

func TestListIdeasDefaultFields(t *testing.T) {
	var gotFields string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotFields = r.URL.Query().Get("fields")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listIdeasFixture))
	})

	list, err := client.ListIdeas(t.Context(), WithIdeaSort("popular"))
	if err != nil {
		t.Fatalf("ListIdeas: %v", err)
	}

	wantFields := "name,reference_num,created_at,updated_at,description,votes,categories,score,status_changed_at,workflow_status,feature,url,resource"
	if gotFields != wantFields {
		t.Errorf("fields = %q, want %q", gotFields, wantFields)
	}

	if len(list.Ideas) != 1 {
		t.Fatalf("len(Ideas) = %d, want 1", len(list.Ideas))
	}
	idea := list.Ideas[0]
	if idea.Votes != 205 {
		t.Errorf("Votes = %d, want 205", idea.Votes)
	}
	if idea.Score != 0 {
		t.Errorf("Score = %d, want 0", idea.Score)
	}
	if len(idea.Categories) != 1 || idea.Categories[0].Name != "Administration - Identity Repository" {
		t.Errorf("Categories = %+v, want one category named %q", idea.Categories, "Administration - Identity Repository")
	}
}

func TestListIdeasWithIdeaFields(t *testing.T) {
	var gotFields string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotFields = r.URL.Query().Get("fields")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listIdeasFixture))
	})

	if _, err := client.ListIdeas(t.Context(), WithIdeaFields("votes")); err != nil {
		t.Fatalf("ListIdeas: %v", err)
	}

	wantFields := "name,reference_num,created_at,updated_at,votes"
	if gotFields != wantFields {
		t.Errorf("fields = %q, want %q", gotFields, wantFields)
	}
}
