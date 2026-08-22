package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// IdeaUser represents a voter identity, account-wide (not scoped to a
// single idea). Richer than the endorsement-embedded EndorsementIdeaUser
// stub: carries the voter's known organization affiliations.
type IdeaUser struct {
	ID                string
	Name              string
	Email             string
	CreatedAt         time.Time
	IdeaOrganizations []IdeaOrganizationRef
}

// IdeaOrganizationRef is the light shape of an idea organization, as
// embedded on IdeaUser and returned by ListIdeaOrganizations. For full
// detail (email_domains, revenue, endorsements_count), use GetIdeaOrganization.
type IdeaOrganizationRef struct {
	ID        string
	Name      string
	CreatedAt time.Time
	URL       string
	Resource  string
}

// IdeaUserList represents a paginated list of idea users.
type IdeaUserList struct {
	IdeaUsers  []IdeaUser
	Pagination Pagination
}

// ListIdeaUsers lists idea users (voter identities). Account-wide, not
// scoped to a single idea.
func (c *Client) ListIdeaUsers(ctx context.Context, opts ...ListOption) (*IdeaUserList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListIdeaUsersParams{}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListIdeaUsers(ctx, params)
	if err != nil {
		return nil, wrapError("ListIdeaUsers", err)
	}

	return ideaUserListFromAPI(resp), nil
}

// GetIdeaUser retrieves an idea user by ID.
func (c *Client) GetIdeaUser(ctx context.Context, id string) (*IdeaUser, error) {
	resp, err := c.apiClient.GetIdeaUser(ctx, api.GetIdeaUserParams{IdeaUserID: id})
	if err != nil {
		return nil, wrapError("GetIdeaUser", err)
	}

	u, ok := resp.IdeaUser.Get()
	if !ok {
		return &IdeaUser{}, nil
	}
	return ideaUserFromAPI(u), nil
}

func ideaUserListFromAPI(resp *api.IdeaUsersResponse) *IdeaUserList {
	list := &IdeaUserList{}

	list.IdeaUsers = make([]IdeaUser, len(resp.IdeaUsers))
	for i, u := range resp.IdeaUsers {
		list.IdeaUsers[i] = *ideaUserFromAPI(u)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

func ideaUserFromAPI(u api.IdeaUser) *IdeaUser {
	user := &IdeaUser{}
	if v, ok := u.ID.Get(); ok {
		user.ID = v
	}
	if v, ok := u.Name.Get(); ok {
		user.Name = v
	}
	if v, ok := u.Email.Get(); ok {
		user.Email = v
	}
	if v, ok := u.CreatedAt.Get(); ok {
		user.CreatedAt = v
	}
	user.IdeaOrganizations = make([]IdeaOrganizationRef, len(u.IdeaOrganizations))
	for i, o := range u.IdeaOrganizations {
		user.IdeaOrganizations[i] = ideaOrganizationRefFromAPI(o)
	}
	return user
}

func ideaOrganizationRefFromAPI(o api.IdeaOrganizationRef) IdeaOrganizationRef {
	ref := IdeaOrganizationRef{}
	if v, ok := o.ID.Get(); ok {
		ref.ID = v
	}
	if v, ok := o.Name.Get(); ok {
		ref.Name = v
	}
	if v, ok := o.CreatedAt.Get(); ok {
		ref.CreatedAt = v
	}
	if v, ok := o.URL.Get(); ok {
		ref.URL = v
	}
	if v, ok := o.Resource.Get(); ok {
		ref.Resource = v
	}
	return ref
}
