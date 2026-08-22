package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// IdeaEndorsement represents a vote ("endorsement") on an Aha idea, with the
// identity of the voter. Aha's UI calls these "votes"; the API calls them
// "endorsements".
type IdeaEndorsement struct {
	ID                   string
	IdeaID               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Value                string
	Link                 string
	Weight               int
	EndorsedByPortalUser *EndorsementPortalUser
	EndorsedByIdeaUser   *EndorsementIdeaUser
}

// EndorsementPortalUser is the portal-side identity of an idea's voter.
type EndorsementPortalUser struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

// EndorsementIdeaUser is Aha's internal contact record for an idea's voter.
// Distinct from EndorsementPortalUser: same identity, different Aha record.
type EndorsementIdeaUser struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	Title     string
}

// IdeaEndorsementList represents a paginated list of idea endorsements.
type IdeaEndorsementList struct {
	IdeaEndorsements []IdeaEndorsement
	Pagination       Pagination
}

// ListIdeaEndorsements lists the endorsements (votes) on an idea, including
// voter identity (name, email).
func (c *Client) ListIdeaEndorsements(ctx context.Context, ideaID string, opts ...ListOption) (*IdeaEndorsementList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListIdeaEndorsementsParams{
		IdeaID: ideaID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListIdeaEndorsements(ctx, params)
	if err != nil {
		return nil, wrapError("ListIdeaEndorsements", err)
	}

	return ideaEndorsementListFromAPI(resp), nil
}

// ideaEndorsementListFromAPI converts an API idea endorsements response to a
// domain idea endorsement list.
func ideaEndorsementListFromAPI(resp *api.IdeaEndorsementsResponse) *IdeaEndorsementList {
	list := &IdeaEndorsementList{}

	list.IdeaEndorsements = make([]IdeaEndorsement, len(resp.IdeaEndorsements))
	for i, e := range resp.IdeaEndorsements {
		list.IdeaEndorsements[i] = ideaEndorsementFromAPI(e)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// ideaEndorsementFromAPI converts an API idea endorsement to a domain idea endorsement.
func ideaEndorsementFromAPI(e api.IdeaEndorsement) IdeaEndorsement {
	endorsement := IdeaEndorsement{}
	if v, ok := e.ID.Get(); ok {
		endorsement.ID = v
	}
	if v, ok := e.IdeaID.Get(); ok {
		endorsement.IdeaID = v
	}
	if v, ok := e.CreatedAt.Get(); ok {
		endorsement.CreatedAt = v
	}
	if v, ok := e.UpdatedAt.Get(); ok {
		endorsement.UpdatedAt = v
	}
	if e.Value.Set && !e.Value.Null {
		endorsement.Value = e.Value.Value
	}
	if e.Link.Set && !e.Link.Null {
		endorsement.Link = e.Link.Value
	}
	if v, ok := e.Weight.Get(); ok {
		endorsement.Weight = v
	}
	if v, ok := e.EndorsedByPortalUser.Get(); ok {
		endorsement.EndorsedByPortalUser = endorsementPortalUserFromAPI(v)
	}
	if v, ok := e.EndorsedByIdeaUser.Get(); ok {
		endorsement.EndorsedByIdeaUser = endorsementIdeaUserFromAPI(v)
	}
	return endorsement
}

func endorsementPortalUserFromAPI(u api.EndorsementPortalUser) *EndorsementPortalUser {
	user := &EndorsementPortalUser{}
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
	return user
}

func endorsementIdeaUserFromAPI(u api.EndorsementIdeaUser) *EndorsementIdeaUser {
	user := &EndorsementIdeaUser{}
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
	if u.Title.Set && !u.Title.Null {
		user.Title = u.Title.Value
	}
	return user
}
