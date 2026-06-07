package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Feature represents an Aha feature.
type Feature struct {
	ID             string
	ReferenceNum   string
	Name           string
	Description    string
	ProductID      string
	URL            string
	Resource       string
	CommentsCount  int64
	ProgressSource string
	WorkUnits      int64
	StartDate      *time.Time
	DueDate        *time.Time
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	Tags           []string
	WorkflowStatus *WorkflowStatus
	Release        *Release
	AssignedTo     *User
	CustomFields   []CustomField
}

// WorkflowStatus represents a workflow status.
type WorkflowStatus struct {
	ID       string
	Name     string
	Position int64
	Complete bool
	Color    string
}

// CustomField represents a custom field value.
type CustomField struct {
	Key   string
	Name  string
	Value any
	Type  string
}

// FeatureList represents a paginated list of features.
type FeatureList struct {
	Features   []FeatureMeta
	Pagination Pagination
}

// FeatureMeta represents feature metadata in list responses.
type FeatureMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// Pagination represents pagination info.
type Pagination struct {
	TotalRecords int64
	TotalPages   int64
	CurrentPage  int64
}

// GetFeature retrieves a feature by ID or reference number.
func (c *Client) GetFeature(ctx context.Context, id string) (*Feature, error) {
	resp, err := c.apiClient.GetFeature(ctx, api.GetFeatureParams{
		FeatureID: id,
	})
	if err != nil {
		return nil, wrapError("GetFeature", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.FeatureResponse:
		if f, ok := r.Feature.Get(); ok {
			return featureFromAPI(f), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "feature not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "feature not found"}
	}
}

// ListFeaturesOptions configures ListFeatures.
type ListFeaturesOptions struct {
	Query          string
	AssignedToUser string
	Tag            string
	UpdatedSince   *time.Time
	Page           int
	PerPage        int
}

// ListFeaturesOption configures a ListFeatures call.
type ListFeaturesOption func(*ListFeaturesOptions)

// WithFeatureQuery filters features by search query.
func WithFeatureQuery(query string) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.Query = query }
}

// WithFeatureAssignee filters features by assigned user.
func WithFeatureAssignee(email string) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.AssignedToUser = email }
}

// WithFeatureTag filters features by tag.
func WithFeatureTag(tag string) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.Tag = tag }
}

// WithFeatureUpdatedSince filters features updated after the given time.
func WithFeatureUpdatedSince(t time.Time) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.UpdatedSince = &t }
}

// WithFeaturePage sets the page number for pagination.
func WithFeaturePage(page int) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.Page = page }
}

// WithFeaturePerPage sets the number of results per page.
func WithFeaturePerPage(perPage int) ListFeaturesOption {
	return func(o *ListFeaturesOptions) { o.PerPage = perPage }
}

// ListFeatures lists features with optional filtering.
func (c *Client) ListFeatures(ctx context.Context, opts ...ListFeaturesOption) (*FeatureList, error) {
	cfg := &ListFeaturesOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListFeaturesParams{}
	if cfg.Query != "" {
		params.Q = api.NewOptString(cfg.Query)
	}
	if cfg.AssignedToUser != "" {
		params.AssignedToUser = api.NewOptString(cfg.AssignedToUser)
	}
	if cfg.Tag != "" {
		params.Tag = api.NewOptString(cfg.Tag)
	}
	if cfg.UpdatedSince != nil {
		params.UpdatedSince = api.NewOptDateTime(*cfg.UpdatedSince)
	}
	if cfg.Page > 0 {
		params.Page = api.NewOptInt32(int32(cfg.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if cfg.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(cfg.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListFeatures(ctx, params)
	if err != nil {
		return nil, wrapError("ListFeatures", err)
	}

	return featureListFromAPI(resp), nil
}

// ListReleaseFeatures lists features in a release.
func (c *Client) ListReleaseFeatures(ctx context.Context, releaseID string, opts ...ListOption) (*FeatureList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListReleaseFeaturesParams{
		ReleaseID: releaseID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListReleaseFeatures(ctx, params)
	if err != nil {
		return nil, wrapError("ListReleaseFeatures", err)
	}

	return featureListFromAPI(resp), nil
}

// CreateFeatureOptions configures CreateFeature.
type CreateFeatureOptions struct {
	Name             string
	Description      string
	WorkflowStatus   string
	AssignedToUser   string
	Tags             string
	StartDate        *time.Time
	DueDate          *time.Time
	OriginalEstimate string
	Initiative       string
}

// CreateFeatureOption configures a CreateFeature call.
type CreateFeatureOption func(*CreateFeatureOptions)

// WithFeatureName sets the feature name.
func WithFeatureName(name string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.Name = name }
}

// WithFeatureDescription sets the feature description.
func WithFeatureDescription(desc string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.Description = desc }
}

// WithFeatureStatus sets the workflow status.
func WithFeatureStatus(status string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.WorkflowStatus = status }
}

// WithFeatureAssignedTo sets the assigned user.
func WithFeatureAssignedTo(email string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.AssignedToUser = email }
}

// WithFeatureTags sets the tags (comma-separated).
func WithFeatureTags(tags string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.Tags = tags }
}

// WithFeatureStartDate sets the start date.
func WithFeatureStartDate(t time.Time) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.StartDate = &t }
}

// WithFeatureDueDate sets the due date.
func WithFeatureDueDate(t time.Time) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.DueDate = &t }
}

// WithFeatureEstimate sets the original estimate (e.g., "2d", "4h").
func WithFeatureEstimate(estimate string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.OriginalEstimate = estimate }
}

// WithFeatureInitiative sets the initiative.
func WithFeatureInitiative(initiative string) CreateFeatureOption {
	return func(o *CreateFeatureOptions) { o.Initiative = initiative }
}

// CreateFeature creates a new feature in a release.
func (c *Client) CreateFeature(ctx context.Context, releaseID string, opts ...CreateFeatureOption) (*Feature, error) {
	cfg := &CreateFeatureOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "feature name is required"}
	}

	feature := api.FeatureCreate{
		Name: cfg.Name,
	}
	if cfg.Description != "" {
		feature.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		feature.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.AssignedToUser != "" {
		feature.AssignedToUser = api.NewOptString(cfg.AssignedToUser)
	}
	if cfg.Tags != "" {
		feature.Tags = api.NewOptString(cfg.Tags)
	}
	if cfg.StartDate != nil {
		feature.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.DueDate != nil {
		feature.DueDate = api.NewOptDate(*cfg.DueDate)
	}
	if cfg.OriginalEstimate != "" {
		feature.OriginalEstimateText = api.NewOptString(cfg.OriginalEstimate)
	}
	if cfg.Initiative != "" {
		feature.Initiative = api.NewOptString(cfg.Initiative)
	}

	req := &api.FeatureCreateRequest{
		Feature: feature,
	}

	resp, err := c.apiClient.CreateReleaseFeature(ctx, req, api.CreateReleaseFeatureParams{
		ReleaseID: releaseID,
	})
	if err != nil {
		return nil, wrapError("CreateFeature", err)
	}

	if f, ok := resp.Feature.Get(); ok {
		return featureFromAPI(f), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: feature not returned"}
}

// UpdateFeatureOptions configures UpdateFeature.
type UpdateFeatureOptions struct {
	Name              string
	Description       string
	WorkflowStatus    string
	AssignedToUser    string
	Tags              string
	StartDate         *time.Time
	DueDate           *time.Time
	Release           string
	OriginalEstimate  string
	RemainingEstimate string
	Initiative        string
	ReleasePhase      string
}

// UpdateFeatureOption configures an UpdateFeature call.
type UpdateFeatureOption func(*UpdateFeatureOptions)

// WithUpdateFeatureName sets the feature name.
func WithUpdateFeatureName(name string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.Name = name }
}

// WithUpdateFeatureDescription sets the feature description.
func WithUpdateFeatureDescription(desc string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.Description = desc }
}

// WithUpdateFeatureStatus sets the workflow status.
func WithUpdateFeatureStatus(status string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.WorkflowStatus = status }
}

// WithUpdateFeatureAssignedToUser sets the assigned user.
func WithUpdateFeatureAssignedToUser(email string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.AssignedToUser = email }
}

// WithUpdateFeatureRelease sets the release.
func WithUpdateFeatureRelease(release string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.Release = release }
}

// WithUpdateFeatureTags sets the tags.
func WithUpdateFeatureTags(tags string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.Tags = tags }
}

// WithUpdateFeatureStartDate sets the start date.
func WithUpdateFeatureStartDate(t time.Time) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.StartDate = &t }
}

// WithUpdateFeatureDueDate sets the due date.
func WithUpdateFeatureDueDate(t time.Time) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.DueDate = &t }
}

// WithUpdateFeatureInitiative sets the initiative.
func WithUpdateFeatureInitiative(initiative string) UpdateFeatureOption {
	return func(o *UpdateFeatureOptions) { o.Initiative = initiative }
}

// UpdateFeature updates an existing feature.
func (c *Client) UpdateFeature(ctx context.Context, id string, opts ...UpdateFeatureOption) (*Feature, error) {
	cfg := &UpdateFeatureOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	feature := api.FeatureUpdate{}
	if cfg.Name != "" {
		feature.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		feature.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		feature.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.AssignedToUser != "" {
		feature.AssignedToUser = api.NewOptString(cfg.AssignedToUser)
	}
	if cfg.Tags != "" {
		feature.Tags = api.NewOptString(cfg.Tags)
	}
	if cfg.StartDate != nil {
		feature.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.DueDate != nil {
		feature.DueDate = api.NewOptDate(*cfg.DueDate)
	}
	if cfg.Release != "" {
		feature.Release = api.NewOptString(cfg.Release)
	}
	if cfg.OriginalEstimate != "" {
		feature.OriginalEstimateText = api.NewOptString(cfg.OriginalEstimate)
	}
	if cfg.RemainingEstimate != "" {
		feature.RemainingEstimateText = api.NewOptString(cfg.RemainingEstimate)
	}
	if cfg.Initiative != "" {
		feature.Initiative = api.NewOptString(cfg.Initiative)
	}
	if cfg.ReleasePhase != "" {
		feature.ReleasePhase = api.NewOptString(cfg.ReleasePhase)
	}

	req := &api.FeatureUpdateRequest{
		Feature: feature,
	}

	resp, err := c.apiClient.UpdateFeature(ctx, req, api.UpdateFeatureParams{
		FeatureID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateFeature", err)
	}

	if f, ok := resp.Feature.Get(); ok {
		return featureFromAPI(f), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: feature not returned"}
}

// featureFromAPI converts an API feature to a domain feature.
func featureFromAPI(f api.Feature) *Feature {
	feature := &Feature{
		ID:           f.ID,
		ReferenceNum: f.ReferenceNum,
		Name:         f.Name,
		CreatedAt:    f.CreatedAt,
	}

	if v, ok := f.Description.Get(); ok {
		feature.Description = v
	}
	if v, ok := f.ProductID.Get(); ok {
		feature.ProductID = v
	}
	if v, ok := f.URL.Get(); ok {
		feature.URL = v
	}
	if v, ok := f.Resource.Get(); ok {
		feature.Resource = v
	}
	if v, ok := f.CommentsCount.Get(); ok {
		feature.CommentsCount = v
	}
	if v, ok := f.ProgressSource.Get(); ok {
		feature.ProgressSource = v
	}
	if v, ok := f.WorkUnits.Get(); ok {
		feature.WorkUnits = v
	}
	if v, ok := f.StartDate.Get(); ok {
		feature.StartDate = &v
	}
	if v, ok := f.DueDate.Get(); ok {
		feature.DueDate = &v
	}
	if v, ok := f.UpdatedAt.Get(); ok {
		feature.UpdatedAt = &v
	}
	feature.Tags = f.Tags
	if v, ok := f.WorkflowStatus.Get(); ok {
		feature.WorkflowStatus = workflowStatusFromAPI(v)
	}
	if v, ok := f.Release.Get(); ok {
		feature.Release = releaseFromAPI(v)
	}
	if v, ok := f.AssignedToUser.Get(); ok {
		feature.AssignedTo = userFromAPI(v)
	}
	feature.CustomFields = customFieldsFromAPI(f.CustomFields)

	return feature
}

// featureListFromAPI converts an API features response to a domain feature list.
func featureListFromAPI(resp *api.FeaturesResponse) *FeatureList {
	list := &FeatureList{}

	list.Features = make([]FeatureMeta, len(resp.Features))
	for i, f := range resp.Features {
		list.Features[i] = featureMetaFromAPI(f)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// featureMetaFromAPI converts an API feature meta to a domain feature meta.
func featureMetaFromAPI(f api.FeatureMeta) FeatureMeta {
	meta := FeatureMeta{}
	if v, ok := f.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := f.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := f.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := f.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := f.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := f.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}

// workflowStatusFromAPI converts an API workflow status.
func workflowStatusFromAPI(s api.WorkflowStatus) *WorkflowStatus {
	status := &WorkflowStatus{}
	if v, ok := s.ID.Get(); ok {
		status.ID = v
	}
	if v, ok := s.Name.Get(); ok {
		status.Name = v
	}
	if v, ok := s.Position.Get(); ok {
		status.Position = v
	}
	if v, ok := s.Complete.Get(); ok {
		status.Complete = v
	}
	if v, ok := s.Color.Get(); ok {
		status.Color = v
	}
	return status
}

// customFieldsFromAPI converts API custom fields.
func customFieldsFromAPI(fields []api.CustomField) []CustomField {
	result := make([]CustomField, len(fields))
	for i, f := range fields {
		result[i] = CustomField{}
		if v, ok := f.Key.Get(); ok {
			result[i].Key = v
		}
		if v, ok := f.Name.Get(); ok {
			result[i].Name = v
		}
		if v, ok := f.Type.Get(); ok {
			result[i].Type = v
		}
		result[i].Value = f.Value
	}
	return result
}

// paginationFromAPI converts API pagination.
func paginationFromAPI(p api.Pagination) Pagination {
	pg := Pagination{}
	if v, ok := p.TotalRecords.Get(); ok {
		pg.TotalRecords = v
	}
	if v, ok := p.TotalPages.Get(); ok {
		pg.TotalPages = v
	}
	if v, ok := p.CurrentPage.Get(); ok {
		pg.CurrentPage = v
	}
	return pg
}
