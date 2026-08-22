package aha

import (
	"net/http"
	"testing"
)

const listIdeaEndorsementsFixture = `{
	"idea_endorsements": [
		{
			"id": "7000000000000000001",
			"idea_id": "7000000000000000000",
			"created_at": "2026-08-19T14:50:08.917Z",
			"updated_at": "2026-08-19T14:50:08.917Z",
			"value": null,
			"link": null,
			"weight": 1,
			"endorsed_by_portal_user": {
				"id": "7000000000000000010",
				"name": "Ada Lovelace",
				"email": "ada@example.com",
				"created_at": "2025-07-25T16:17:52.240Z"
			},
			"endorsed_by_idea_user": {
				"id": "7000000000000000011",
				"name": "Ada Lovelace",
				"email": "ada@example.com",
				"created_at": "2025-07-25T16:17:52.225Z",
				"title": null
			}
		}
	],
	"pagination": {
		"total_records": 1,
		"total_pages": 1,
		"current_page": 1
	}
}`

func TestListIdeaEndorsements(t *testing.T) {
	var gotPath, gotPage, gotPerPage string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotPage = r.URL.Query().Get("page")
		gotPerPage = r.URL.Query().Get("per_page")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listIdeaEndorsementsFixture))
	})

	list, err := client.ListIdeaEndorsements(t.Context(), "7000000000000000000", WithPage(2), WithPerPage(50))
	if err != nil {
		t.Fatalf("ListIdeaEndorsements: %v", err)
	}

	wantPath := "/ideas/7000000000000000000/endorsements"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if gotPage != "2" {
		t.Errorf("page query param = %q, want %q", gotPage, "2")
	}
	if gotPerPage != "50" {
		t.Errorf("per_page query param = %q, want %q", gotPerPage, "50")
	}

	if len(list.IdeaEndorsements) != 1 {
		t.Fatalf("len(IdeaEndorsements) = %d, want 1", len(list.IdeaEndorsements))
	}
	e := list.IdeaEndorsements[0]

	if e.ID != "7000000000000000001" {
		t.Errorf("ID = %q, want %q", e.ID, "7000000000000000001")
	}
	if e.IdeaID != "7000000000000000000" {
		t.Errorf("IdeaID = %q, want %q", e.IdeaID, "7000000000000000000")
	}
	if e.Weight != 1 {
		t.Errorf("Weight = %d, want 1", e.Weight)
	}
	if e.Value != "" {
		t.Errorf("Value = %q, want empty (null in fixture)", e.Value)
	}

	if e.EndorsedByPortalUser == nil {
		t.Fatal("EndorsedByPortalUser is nil")
	}
	if e.EndorsedByPortalUser.Email != "ada@example.com" {
		t.Errorf("EndorsedByPortalUser.Email = %q, want %q", e.EndorsedByPortalUser.Email, "ada@example.com")
	}

	if e.EndorsedByIdeaUser == nil {
		t.Fatal("EndorsedByIdeaUser is nil")
	}
	if e.EndorsedByIdeaUser.Email != "ada@example.com" {
		t.Errorf("EndorsedByIdeaUser.Email = %q, want %q", e.EndorsedByIdeaUser.Email, "ada@example.com")
	}
	if e.EndorsedByIdeaUser.Title != "" {
		t.Errorf("EndorsedByIdeaUser.Title = %q, want empty (null in fixture)", e.EndorsedByIdeaUser.Title)
	}

	if list.Pagination.TotalRecords != 1 {
		t.Errorf("Pagination.TotalRecords = %d, want 1", list.Pagination.TotalRecords)
	}
}
