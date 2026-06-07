package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Requirement represents an Aha requirement.
type Requirement struct {
	ID                string
	ReferenceNum      string
	Name              string
	Description       string
	Position          int64
	OriginalEstimate  float64
	RemainingEstimate float64
	WorkDone          float64
	URL               string
	Resource          string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	WorkflowStatus    *WorkflowStatus
	AssignedToUser    *User
	Feature           *FeatureMeta
}

// RequirementList represents a paginated list of requirements.
type RequirementList struct {
	Requirements []RequirementMeta
	Pagination   Pagination
}

// RequirementMeta represents requirement metadata in list responses.
type RequirementMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// GetRequirement retrieves a requirement by ID or reference number.
func (c *Client) GetRequirement(ctx context.Context, id string) (*Requirement, error) {
	resp, err := c.apiClient.GetRequirement(ctx, api.GetRequirementParams{
		RequirementID: id,
	})
	if err != nil {
		return nil, wrapError("GetRequirement", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.RequirementResponse:
		if req, ok := r.Requirement.Get(); ok {
			return requirementFromAPI(req), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "requirement not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "requirement not found"}
	}
}

// ListFeatureRequirements lists requirements for a feature.
func (c *Client) ListFeatureRequirements(ctx context.Context, featureID string, opts ...ListOption) (*RequirementList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListFeatureRequirementsParams{
		FeatureID: featureID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListFeatureRequirements(ctx, params)
	if err != nil {
		return nil, wrapError("ListFeatureRequirements", err)
	}

	return requirementListFromAPI(resp), nil
}

// CreateRequirementOptions configures CreateRequirement.
type CreateRequirementOptions struct {
	Name             string
	Description      string
	WorkflowStatus   string
	AssignedToUser   string
	OriginalEstimate *float64
}

// CreateRequirementOption configures a CreateRequirement call.
type CreateRequirementOption func(*CreateRequirementOptions)

// WithRequirementName sets the requirement name.
func WithRequirementName(name string) CreateRequirementOption {
	return func(o *CreateRequirementOptions) { o.Name = name }
}

// WithRequirementDescription sets the requirement description.
func WithRequirementDescription(desc string) CreateRequirementOption {
	return func(o *CreateRequirementOptions) { o.Description = desc }
}

// WithRequirementStatus sets the workflow status.
func WithRequirementStatus(status string) CreateRequirementOption {
	return func(o *CreateRequirementOptions) { o.WorkflowStatus = status }
}

// WithRequirementAssignedTo sets the assigned user.
func WithRequirementAssignedTo(user string) CreateRequirementOption {
	return func(o *CreateRequirementOptions) { o.AssignedToUser = user }
}

// WithRequirementEstimate sets the original estimate.
func WithRequirementEstimate(estimate float64) CreateRequirementOption {
	return func(o *CreateRequirementOptions) { o.OriginalEstimate = &estimate }
}

// CreateRequirement creates a new requirement for a feature.
func (c *Client) CreateRequirement(ctx context.Context, featureID string, opts ...CreateRequirementOption) (*Requirement, error) {
	cfg := &CreateRequirementOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "requirement name is required"}
	}

	req := api.RequirementCreate{
		Name: cfg.Name,
	}
	if cfg.Description != "" {
		req.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		req.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.AssignedToUser != "" {
		req.AssignedToUser = api.NewOptString(cfg.AssignedToUser)
	}
	if cfg.OriginalEstimate != nil {
		req.OriginalEstimate = api.NewOptFloat32(float32(*cfg.OriginalEstimate))
	}

	reqBody := &api.RequirementCreateRequest{
		Requirement: req,
	}

	resp, err := c.apiClient.CreateFeatureRequirement(ctx, reqBody, api.CreateFeatureRequirementParams{
		FeatureID: featureID,
	})
	if err != nil {
		return nil, wrapError("CreateRequirement", err)
	}

	if r, ok := resp.Requirement.Get(); ok {
		return requirementFromAPI(r), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: requirement not returned"}
}

// UpdateRequirementOptions configures UpdateRequirement.
type UpdateRequirementOptions struct {
	Name              string
	Description       string
	WorkflowStatus    string
	AssignedToUser    string
	OriginalEstimate  *float64
	RemainingEstimate *float64
	WorkDone          *float64
}

// UpdateRequirementOption configures an UpdateRequirement call.
type UpdateRequirementOption func(*UpdateRequirementOptions)

// WithUpdateRequirementName sets the requirement name.
func WithUpdateRequirementName(name string) UpdateRequirementOption {
	return func(o *UpdateRequirementOptions) { o.Name = name }
}

// WithUpdateRequirementDescription sets the requirement description.
func WithUpdateRequirementDescription(desc string) UpdateRequirementOption {
	return func(o *UpdateRequirementOptions) { o.Description = desc }
}

// WithUpdateRequirementStatus sets the workflow status.
func WithUpdateRequirementStatus(status string) UpdateRequirementOption {
	return func(o *UpdateRequirementOptions) { o.WorkflowStatus = status }
}

// WithUpdateRequirementWorkDone sets the work done.
func WithUpdateRequirementWorkDone(done float64) UpdateRequirementOption {
	return func(o *UpdateRequirementOptions) { o.WorkDone = &done }
}

// UpdateRequirement updates an existing requirement.
func (c *Client) UpdateRequirement(ctx context.Context, id string, opts ...UpdateRequirementOption) (*Requirement, error) {
	cfg := &UpdateRequirementOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	req := api.RequirementUpdate{}
	if cfg.Name != "" {
		req.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		req.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		req.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.AssignedToUser != "" {
		req.AssignedToUser = api.NewOptString(cfg.AssignedToUser)
	}
	if cfg.OriginalEstimate != nil {
		req.OriginalEstimate = api.NewOptFloat32(float32(*cfg.OriginalEstimate))
	}
	if cfg.RemainingEstimate != nil {
		req.RemainingEstimate = api.NewOptFloat32(float32(*cfg.RemainingEstimate))
	}
	if cfg.WorkDone != nil {
		req.WorkDone = api.NewOptFloat32(float32(*cfg.WorkDone))
	}

	reqBody := &api.RequirementUpdateRequest{
		Requirement: req,
	}

	resp, err := c.apiClient.UpdateRequirement(ctx, reqBody, api.UpdateRequirementParams{
		RequirementID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateRequirement", err)
	}

	if r, ok := resp.Requirement.Get(); ok {
		return requirementFromAPI(r), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: requirement not returned"}
}

// DeleteRequirement deletes a requirement.
func (c *Client) DeleteRequirement(ctx context.Context, id string) error {
	err := c.apiClient.DeleteRequirement(ctx, api.DeleteRequirementParams{
		RequirementID: id,
	})
	if err != nil {
		return wrapError("DeleteRequirement", err)
	}
	return nil
}

// requirementFromAPI converts an API requirement to a domain requirement.
func requirementFromAPI(r api.Requirement) *Requirement {
	req := &Requirement{
		ID:           r.ID,
		ReferenceNum: r.ReferenceNum,
		Name:         r.Name,
		CreatedAt:    r.CreatedAt,
	}

	if v, ok := r.Description.Get(); ok {
		req.Description = v
	}
	if v, ok := r.Position.Get(); ok {
		req.Position = v
	}
	if v, ok := r.OriginalEstimate.Get(); ok {
		req.OriginalEstimate = float64(v)
	}
	if v, ok := r.RemainingEstimate.Get(); ok {
		req.RemainingEstimate = float64(v)
	}
	if v, ok := r.WorkDone.Get(); ok {
		req.WorkDone = float64(v)
	}
	if v, ok := r.URL.Get(); ok {
		req.URL = v
	}
	if v, ok := r.Resource.Get(); ok {
		req.Resource = v
	}
	if v, ok := r.UpdatedAt.Get(); ok {
		req.UpdatedAt = &v
	}
	if v, ok := r.WorkflowStatus.Get(); ok {
		req.WorkflowStatus = workflowStatusFromAPI(v)
	}
	if v, ok := r.AssignedToUser.Get(); ok {
		req.AssignedToUser = userFromAPI(v)
	}
	if v, ok := r.Feature.Get(); ok {
		req.Feature = &FeatureMeta{
			ID:           v.ID.Or(""),
			ReferenceNum: v.ReferenceNum.Or(""),
			Name:         v.Name.Or(""),
			URL:          v.URL.Or(""),
			Resource:     v.Resource.Or(""),
		}
		if ca, ok := v.CreatedAt.Get(); ok {
			req.Feature.CreatedAt = ca
		}
	}

	return req
}

// requirementListFromAPI converts an API requirements response to a domain requirement list.
func requirementListFromAPI(resp *api.RequirementsResponse) *RequirementList {
	list := &RequirementList{}

	list.Requirements = make([]RequirementMeta, len(resp.Requirements))
	for i, r := range resp.Requirements {
		list.Requirements[i] = requirementMetaFromAPI(r)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// requirementMetaFromAPI converts an API requirement meta to a domain requirement meta.
func requirementMetaFromAPI(r api.RequirementMeta) RequirementMeta {
	meta := RequirementMeta{}
	if v, ok := r.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := r.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := r.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := r.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := r.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := r.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}
