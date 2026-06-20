package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Release represents an Aha release.
type Release struct {
	ID                  string
	ReferenceNum        string
	Name                string
	StartDate           *time.Time
	ReleaseDate         *time.Time
	ExternalReleaseDate *time.Time
	Released            bool
	ParkingLot          bool
	URL                 string
	Resource            string
}

// ReleaseList represents a paginated list of releases.
type ReleaseList struct {
	Releases   []Release
	Pagination Pagination
}

// GetRelease retrieves a release by ID or reference number.
func (c *Client) GetRelease(ctx context.Context, id string) (*Release, error) {
	resp, err := c.apiClient.GetRelease(ctx, api.GetReleaseParams{
		ReleaseID: id,
	})
	if err != nil {
		return nil, wrapError("GetRelease", err)
	}

	if r, ok := resp.Release.Get(); ok {
		return releaseFromAPI(r), nil
	}
	return nil, &APIError{StatusCode: 404, Message: "release not found"}
}

// ListProductReleases lists releases for a product.
func (c *Client) ListProductReleases(ctx context.Context, productID string, opts ...ListOption) (*ReleaseList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductReleasesParams{
		ProductID: productID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListProductReleases(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductReleases", err)
	}

	return releaseListFromAPI(resp), nil
}

// UpdateReleaseOptions configures UpdateRelease.
type UpdateReleaseOptions struct {
	Name                 string
	StartDate            *time.Time
	ReleaseDate          *time.Time
	ExternalReleaseDate  *time.Time
	DevelopmentStartedOn *time.Time
	ParkingLot           *bool
}

// UpdateReleaseOption configures an UpdateRelease call.
type UpdateReleaseOption func(*UpdateReleaseOptions)

// WithReleaseName sets the release name.
func WithReleaseName(name string) UpdateReleaseOption {
	return func(o *UpdateReleaseOptions) { o.Name = name }
}

// WithReleaseStartDate sets the start date.
func WithReleaseStartDate(t time.Time) UpdateReleaseOption {
	return func(o *UpdateReleaseOptions) { o.StartDate = &t }
}

// WithReleaseDate sets the release date.
func WithReleaseDate(t time.Time) UpdateReleaseOption {
	return func(o *UpdateReleaseOptions) { o.ReleaseDate = &t }
}

// WithReleaseParkingLot sets whether this is a parking lot release.
func WithReleaseParkingLot(parkingLot bool) UpdateReleaseOption {
	return func(o *UpdateReleaseOptions) { o.ParkingLot = &parkingLot }
}

// UpdateRelease updates an existing release.
func (c *Client) UpdateRelease(ctx context.Context, id string, opts ...UpdateReleaseOption) (*Release, error) {
	cfg := &UpdateReleaseOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	release := api.ReleaseUpdate{}
	if cfg.Name != "" {
		release.Name = api.NewOptString(cfg.Name)
	}
	if cfg.StartDate != nil {
		release.StartDate = api.NewOptNilDate(*cfg.StartDate)
	}
	if cfg.ReleaseDate != nil {
		release.ReleaseDate = api.NewOptNilDate(*cfg.ReleaseDate)
	}
	if cfg.ExternalReleaseDate != nil {
		release.ExternalReleaseDate = api.NewOptNilDate(*cfg.ExternalReleaseDate)
	}
	if cfg.DevelopmentStartedOn != nil {
		release.DevelopmentStartedOn = api.NewOptDate(*cfg.DevelopmentStartedOn)
	}
	if cfg.ParkingLot != nil {
		release.ParkingLot = api.NewOptBool(*cfg.ParkingLot)
	}

	req := &api.ReleaseUpdateRequest{
		Release: release,
	}

	resp, err := c.apiClient.UpdateRelease(ctx, req, api.UpdateReleaseParams{
		ReleaseID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateRelease", err)
	}

	if r, ok := resp.Release.Get(); ok {
		return releaseFromAPI(r), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: release not returned"}
}

// releaseFromAPI converts an API release to a domain release.
func releaseFromAPI(r api.Release) *Release {
	release := &Release{}
	if v, ok := r.ID.Get(); ok {
		release.ID = v
	}
	if v, ok := r.ReferenceNum.Get(); ok {
		release.ReferenceNum = v
	}
	if v, ok := r.Name.Get(); ok {
		release.Name = v
	}
	if v, ok := r.StartDate.Get(); ok {
		release.StartDate = &v
	}
	if v, ok := r.ReleaseDate.Get(); ok {
		release.ReleaseDate = &v
	}
	if v, ok := r.ExternalReleaseDate.Get(); ok {
		release.ExternalReleaseDate = &v
	}
	if v, ok := r.Released.Get(); ok {
		release.Released = v
	}
	if v, ok := r.ParkingLot.Get(); ok {
		release.ParkingLot = v
	}
	if v, ok := r.URL.Get(); ok {
		release.URL = v
	}
	if v, ok := r.Resource.Get(); ok {
		release.Resource = v
	}
	return release
}

// releaseListFromAPI converts an API releases response to a domain release list.
func releaseListFromAPI(resp *api.ReleasesResponse) *ReleaseList {
	list := &ReleaseList{}

	list.Releases = make([]Release, len(resp.Releases))
	for i, r := range resp.Releases {
		list.Releases[i] = *releaseFromAPI(r)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}
