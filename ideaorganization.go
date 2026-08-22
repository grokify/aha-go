package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// IdeaOrganization is the full-detail shape of an Aha idea organization
// (customer/account record), available only via GetIdeaOrganization.
// EmailDomains is Aha's own authoritative domain-to-organization mapping,
// sourced from Aha's CRM data.
type IdeaOrganization struct {
	ID                string
	Name              string
	ReferenceNum      string
	URL               string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EndorsementsCount int
	EmailDomains      string
	Revenue           *float64
}

// IdeaOrganizationList represents a paginated list of idea organizations.
// The list response is the light IdeaOrganizationRef shape (no
// email_domains/revenue/endorsements_count) -- use GetIdeaOrganization for
// full detail.
type IdeaOrganizationList struct {
	IdeaOrganizations []IdeaOrganizationRef
	Pagination        Pagination
}

// ListIdeaOrganizations lists idea organizations (customer/account
// records). Account-wide, not scoped to a single idea.
func (c *Client) ListIdeaOrganizations(ctx context.Context, opts ...ListOption) (*IdeaOrganizationList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListIdeaOrganizationsParams{}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListIdeaOrganizations(ctx, params)
	if err != nil {
		return nil, wrapError("ListIdeaOrganizations", err)
	}

	return ideaOrganizationListFromAPI(resp), nil
}

// GetIdeaOrganization retrieves an idea organization by ID, with full
// detail (email_domains, revenue, endorsements_count) not present in the
// list response.
func (c *Client) GetIdeaOrganization(ctx context.Context, id string) (*IdeaOrganization, error) {
	resp, err := c.apiClient.GetIdeaOrganization(ctx, api.GetIdeaOrganizationParams{IdeaOrganizationID: id})
	if err != nil {
		return nil, wrapError("GetIdeaOrganization", err)
	}

	o, ok := resp.IdeaOrganization.Get()
	if !ok {
		return &IdeaOrganization{}, nil
	}
	return ideaOrganizationFromAPI(o), nil
}

func ideaOrganizationListFromAPI(resp *api.IdeaOrganizationsResponse) *IdeaOrganizationList {
	list := &IdeaOrganizationList{}

	list.IdeaOrganizations = make([]IdeaOrganizationRef, len(resp.IdeaOrganizations))
	for i, o := range resp.IdeaOrganizations {
		list.IdeaOrganizations[i] = ideaOrganizationRefFromAPI(o)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

func ideaOrganizationFromAPI(o api.IdeaOrganization) *IdeaOrganization {
	org := &IdeaOrganization{}
	if v, ok := o.ID.Get(); ok {
		org.ID = v
	}
	if v, ok := o.Name.Get(); ok {
		org.Name = v
	}
	if v, ok := o.ReferenceNum.Get(); ok {
		org.ReferenceNum = v
	}
	if v, ok := o.URL.Get(); ok {
		org.URL = v
	}
	if v, ok := o.CreatedAt.Get(); ok {
		org.CreatedAt = v
	}
	if v, ok := o.UpdatedAt.Get(); ok {
		org.UpdatedAt = v
	}
	if v, ok := o.EndorsementsCount.Get(); ok {
		org.EndorsementsCount = v
	}
	if o.EmailDomains.Set && !o.EmailDomains.Null {
		org.EmailDomains = o.EmailDomains.Value
	}
	if o.Revenue.Set && !o.Revenue.Null {
		v := o.Revenue.Value
		org.Revenue = &v
	}
	return org
}
