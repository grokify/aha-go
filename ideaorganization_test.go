package aha

import (
	"net/http"
	"testing"
)

const listIdeaOrganizationsFixture = `{
	"idea_organizations": [
		{
			"id": "7000000000000000030",
			"name": "Example Corp",
			"created_at": "2026-02-19T14:02:26.435Z",
			"url": "https://example.aha.io/ideas/idea_organizations/ACCOUNT-O-1",
			"resource": "https://example.aha.io/api/v1/idea_organizations/ACCOUNT-O-1"
		}
	],
	"pagination": {
		"total_records": 1,
		"total_pages": 1,
		"current_page": 1
	}
}`

const getIdeaOrganizationFixture = `{
	"idea_organization": {
		"id": "7000000000000000030",
		"name": "Example Corp",
		"reference_num": "ACCOUNT-O-1",
		"url": "https://example.aha.io/ideas/idea_organizations/ACCOUNT-O-1",
		"created_at": "2026-02-19T14:02:26.435Z",
		"updated_at": "2026-02-19T14:03:44.444Z",
		"endorsements_count": 4,
		"email_domains": "example.com",
		"revenue": 250000.50
	}
}`

const getIdeaOrganizationNoRevenueFixture = `{
	"idea_organization": {
		"id": "7000000000000000031",
		"name": "No Revenue Corp",
		"reference_num": "ACCOUNT-O-2",
		"endorsements_count": 0,
		"email_domains": null,
		"revenue": null
	}
}`

func TestListIdeaOrganizations(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listIdeaOrganizationsFixture))
	})

	list, err := client.ListIdeaOrganizations(t.Context())
	if err != nil {
		t.Fatalf("ListIdeaOrganizations: %v", err)
	}

	if len(list.IdeaOrganizations) != 1 || list.IdeaOrganizations[0].Name != "Example Corp" {
		t.Errorf("IdeaOrganizations = %+v, want one named Example Corp", list.IdeaOrganizations)
	}
}

func TestGetIdeaOrganization(t *testing.T) {
	var gotPath string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getIdeaOrganizationFixture))
	})

	org, err := client.GetIdeaOrganization(t.Context(), "7000000000000000030")
	if err != nil {
		t.Fatalf("GetIdeaOrganization: %v", err)
	}

	wantPath := "/idea_organizations/7000000000000000030"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if org.EmailDomains != "example.com" {
		t.Errorf("EmailDomains = %q, want %q", org.EmailDomains, "example.com")
	}
	if org.Revenue == nil || *org.Revenue != 250000.50 {
		t.Errorf("Revenue = %v, want 250000.50", org.Revenue)
	}
	if org.EndorsementsCount != 4 {
		t.Errorf("EndorsementsCount = %d, want 4", org.EndorsementsCount)
	}
	if org.ReferenceNum != "ACCOUNT-O-1" {
		t.Errorf("ReferenceNum = %q, want %q", org.ReferenceNum, "ACCOUNT-O-1")
	}
}

func TestGetIdeaOrganizationNullFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(getIdeaOrganizationNoRevenueFixture))
	})

	org, err := client.GetIdeaOrganization(t.Context(), "7000000000000000031")
	if err != nil {
		t.Fatalf("GetIdeaOrganization: %v", err)
	}

	if org.Revenue != nil {
		t.Errorf("Revenue = %v, want nil", org.Revenue)
	}
	if org.EmailDomains != "" {
		t.Errorf("EmailDomains = %q, want empty (null in fixture)", org.EmailDomains)
	}
}
