package aha

import (
	"net/http"
	"testing"
)

const listIdeaUsersFixture = `{
	"idea_users": [
		{
			"id": "7000000000000000020",
			"name": "Grace Hopper",
			"email": "grace@example.com",
			"created_at": "2025-07-25T16:17:52.225Z",
			"idea_organizations": [
				{
					"id": "7000000000000000030",
					"name": "Example Corp",
					"created_at": "2026-02-19T14:02:26.435Z",
					"url": "https://example.aha.io/ideas/idea_organizations/ACCOUNT-O-1",
					"resource": "https://example.aha.io/api/v1/idea_organizations/ACCOUNT-O-1"
				}
			]
		}
	],
	"pagination": {
		"total_records": 1,
		"total_pages": 1,
		"current_page": 1
	}
}`

const getIdeaUserFixture = `{
	"idea_user": {
		"id": "7000000000000000020",
		"name": "Grace Hopper",
		"email": "grace@example.com",
		"created_at": "2025-07-25T16:17:52.225Z",
		"idea_organizations": []
	}
}`

func TestListIdeaUsers(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listIdeaUsersFixture))
	})

	list, err := client.ListIdeaUsers(t.Context())
	if err != nil {
		t.Fatalf("ListIdeaUsers: %v", err)
	}

	if len(list.IdeaUsers) != 1 {
		t.Fatalf("len(IdeaUsers) = %d, want 1", len(list.IdeaUsers))
	}
	u := list.IdeaUsers[0]
	if u.Email != "grace@example.com" {
		t.Errorf("Email = %q, want %q", u.Email, "grace@example.com")
	}
	if len(u.IdeaOrganizations) != 1 || u.IdeaOrganizations[0].Name != "Example Corp" {
		t.Errorf("IdeaOrganizations = %+v, want one named Example Corp", u.IdeaOrganizations)
	}
	if list.Pagination.TotalRecords != 1 {
		t.Errorf("Pagination.TotalRecords = %d, want 1", list.Pagination.TotalRecords)
	}
}

func TestGetIdeaUser(t *testing.T) {
	var gotPath string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getIdeaUserFixture))
	})

	u, err := client.GetIdeaUser(t.Context(), "7000000000000000020")
	if err != nil {
		t.Fatalf("GetIdeaUser: %v", err)
	}

	wantPath := "/idea_users/7000000000000000020"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if u.Name != "Grace Hopper" {
		t.Errorf("Name = %q, want %q", u.Name, "Grace Hopper")
	}
	if len(u.IdeaOrganizations) != 0 {
		t.Errorf("IdeaOrganizations = %+v, want empty", u.IdeaOrganizations)
	}
}
