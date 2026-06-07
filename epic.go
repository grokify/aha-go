package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Epic represents an Aha epic.
type Epic struct {
	ID             string
	ReferenceNum   string
	Name           string
	Description    string
	Progress       float64
	ProgressSource string
	Position       int64
	Color          string
	StartDate      *time.Time
	DueDate        *time.Time
	URL            string
	Resource       string
	CommentsCount  int64
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	Tags           []string
	WorkflowStatus *WorkflowStatus
	Release        *Release
	Initiative     *InitiativeMeta
}

// EpicList represents a paginated list of epics.
type EpicList struct {
	Epics      []EpicMeta
	Pagination Pagination
}

// EpicMeta is defined in initiative.go

// GetEpic retrieves an epic by ID or reference number.
func (c *Client) GetEpic(ctx context.Context, id string) (*Epic, error) {
	resp, err := c.apiClient.GetEpic(ctx, api.GetEpicParams{
		EpicID: id,
	})
	if err != nil {
		return nil, wrapError("GetEpic", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.EpicResponse:
		if e, ok := r.Epic.Get(); ok {
			return epicFromAPI(e), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "epic not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "epic not found"}
	}
}

// ListEpicsOptions configures ListEpics.
type ListEpicsOptions struct {
	Query        string
	UpdatedSince *time.Time
	Page         int
	PerPage      int
}

// ListEpicsOption configures a ListEpics call.
type ListEpicsOption func(*ListEpicsOptions)

// WithEpicQuery filters epics by search query.
func WithEpicQuery(query string) ListEpicsOption {
	return func(o *ListEpicsOptions) { o.Query = query }
}

// WithEpicUpdatedSince filters epics updated after the given time.
func WithEpicUpdatedSince(t time.Time) ListEpicsOption {
	return func(o *ListEpicsOptions) { o.UpdatedSince = &t }
}

// WithEpicPage sets the page number for pagination.
func WithEpicPage(page int) ListEpicsOption {
	return func(o *ListEpicsOptions) { o.Page = page }
}

// WithEpicPerPage sets the number of results per page.
func WithEpicPerPage(perPage int) ListEpicsOption {
	return func(o *ListEpicsOptions) { o.PerPage = perPage }
}

// ListEpics lists epics with optional filtering.
func (c *Client) ListEpics(ctx context.Context, opts ...ListEpicsOption) (*EpicList, error) {
	cfg := &ListEpicsOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListEpicsParams{}
	if cfg.Query != "" {
		params.Q = api.NewOptString(cfg.Query)
	}
	if cfg.UpdatedSince != nil {
		params.UpdatedSince = api.NewOptDateTime(*cfg.UpdatedSince)
	}
	if cfg.Page > 0 {
		params.Page = api.NewOptInt32(int32(cfg.Page))
	}
	if cfg.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(cfg.PerPage))
	}

	resp, err := c.apiClient.ListEpics(ctx, params)
	if err != nil {
		return nil, wrapError("ListEpics", err)
	}

	return epicListFromAPI(resp), nil
}

// ListProductEpics lists epics for a product.
func (c *Client) ListProductEpics(ctx context.Context, productID string, opts ...ListOption) (*EpicList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductEpicsParams{
		ProductID: productID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListProductEpics(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductEpics", err)
	}

	return epicListFromAPI(resp), nil
}

// CreateEpicOptions configures CreateEpic.
type CreateEpicOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	DueDate        *time.Time
	Color          string
	Initiative     string
}

// CreateEpicOption configures a CreateEpic call.
type CreateEpicOption func(*CreateEpicOptions)

// WithEpicName sets the epic name.
func WithEpicName(name string) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.Name = name }
}

// WithEpicDescription sets the epic description.
func WithEpicDescription(desc string) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.Description = desc }
}

// WithEpicStatus sets the workflow status.
func WithEpicStatus(status string) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.WorkflowStatus = status }
}

// WithEpicStartDate sets the start date.
func WithEpicStartDate(t time.Time) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.StartDate = &t }
}

// WithEpicDueDate sets the due date.
func WithEpicDueDate(t time.Time) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.DueDate = &t }
}

// WithEpicColor sets the color.
func WithEpicColor(color string) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.Color = color }
}

// WithEpicInitiative sets the initiative.
func WithEpicInitiative(initiative string) CreateEpicOption {
	return func(o *CreateEpicOptions) { o.Initiative = initiative }
}

// CreateEpic creates a new epic in a release.
func (c *Client) CreateEpic(ctx context.Context, releaseID string, opts ...CreateEpicOption) (*Epic, error) {
	cfg := &CreateEpicOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "epic name is required"}
	}

	epic := api.EpicCreate{
		Name: cfg.Name,
	}
	if cfg.Description != "" {
		epic.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		epic.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		epic.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.DueDate != nil {
		epic.DueDate = api.NewOptDate(*cfg.DueDate)
	}
	if cfg.Color != "" {
		epic.Color = api.NewOptString(cfg.Color)
	}
	if cfg.Initiative != "" {
		epic.Initiative = api.NewOptString(cfg.Initiative)
	}

	req := &api.EpicCreateRequest{
		Epic: epic,
	}

	resp, err := c.apiClient.CreateReleaseEpic(ctx, req, api.CreateReleaseEpicParams{
		ReleaseID: releaseID,
	})
	if err != nil {
		return nil, wrapError("CreateEpic", err)
	}

	if e, ok := resp.Epic.Get(); ok {
		return epicFromAPI(e), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: epic not returned"}
}

// UpdateEpicOptions configures UpdateEpic.
type UpdateEpicOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	DueDate        *time.Time
	Progress       *float64
	Color          string
	Initiative     string
}

// UpdateEpicOption configures an UpdateEpic call.
type UpdateEpicOption func(*UpdateEpicOptions)

// WithUpdateEpicName sets the epic name.
func WithUpdateEpicName(name string) UpdateEpicOption {
	return func(o *UpdateEpicOptions) { o.Name = name }
}

// WithUpdateEpicDescription sets the epic description.
func WithUpdateEpicDescription(desc string) UpdateEpicOption {
	return func(o *UpdateEpicOptions) { o.Description = desc }
}

// WithUpdateEpicStatus sets the workflow status.
func WithUpdateEpicStatus(status string) UpdateEpicOption {
	return func(o *UpdateEpicOptions) { o.WorkflowStatus = status }
}

// WithUpdateEpicProgress sets the progress (when progress_source is manual).
func WithUpdateEpicProgress(progress float64) UpdateEpicOption {
	return func(o *UpdateEpicOptions) { o.Progress = &progress }
}

// UpdateEpic updates an existing epic.
func (c *Client) UpdateEpic(ctx context.Context, id string, opts ...UpdateEpicOption) (*Epic, error) {
	cfg := &UpdateEpicOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	epic := api.EpicUpdate{}
	if cfg.Name != "" {
		epic.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		epic.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		epic.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		epic.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.DueDate != nil {
		epic.DueDate = api.NewOptDate(*cfg.DueDate)
	}
	if cfg.Progress != nil {
		epic.Progress = api.NewOptFloat32(float32(*cfg.Progress))
	}
	if cfg.Color != "" {
		epic.Color = api.NewOptString(cfg.Color)
	}
	if cfg.Initiative != "" {
		epic.Initiative = api.NewOptString(cfg.Initiative)
	}

	req := &api.EpicUpdateRequest{
		Epic: epic,
	}

	resp, err := c.apiClient.UpdateEpic(ctx, req, api.UpdateEpicParams{
		EpicID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateEpic", err)
	}

	if e, ok := resp.Epic.Get(); ok {
		return epicFromAPI(e), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: epic not returned"}
}

// epicFromAPI converts an API epic to a domain epic.
func epicFromAPI(e api.Epic) *Epic {
	epic := &Epic{
		ID:           e.ID,
		ReferenceNum: e.ReferenceNum,
		Name:         e.Name,
		CreatedAt:    e.CreatedAt,
	}

	if v, ok := e.Description.Get(); ok {
		epic.Description = v
	}
	if v, ok := e.Progress.Get(); ok {
		epic.Progress = float64(v)
	}
	if v, ok := e.ProgressSource.Get(); ok {
		epic.ProgressSource = v
	}
	if v, ok := e.Position.Get(); ok {
		epic.Position = v
	}
	if v, ok := e.Color.Get(); ok {
		epic.Color = v
	}
	if v, ok := e.StartDate.Get(); ok {
		epic.StartDate = &v
	}
	if v, ok := e.DueDate.Get(); ok {
		epic.DueDate = &v
	}
	if v, ok := e.URL.Get(); ok {
		epic.URL = v
	}
	if v, ok := e.Resource.Get(); ok {
		epic.Resource = v
	}
	if v, ok := e.CommentsCount.Get(); ok {
		epic.CommentsCount = v
	}
	if v, ok := e.UpdatedAt.Get(); ok {
		epic.UpdatedAt = &v
	}
	epic.Tags = e.Tags
	if v, ok := e.WorkflowStatus.Get(); ok {
		epic.WorkflowStatus = workflowStatusFromAPI(v)
	}
	if v, ok := e.Release.Get(); ok {
		epic.Release = releaseFromAPI(v)
	}
	if v, ok := e.Initiative.Get(); ok {
		epic.Initiative = &InitiativeMeta{
			ID:           v.ID.Or(""),
			ReferenceNum: v.ReferenceNum.Or(""),
			Name:         v.Name.Or(""),
			URL:          v.URL.Or(""),
			Resource:     v.Resource.Or(""),
		}
		if ca, ok := v.CreatedAt.Get(); ok {
			epic.Initiative.CreatedAt = ca
		}
	}

	return epic
}

// epicListFromAPI converts an API epics response to a domain epic list.
func epicListFromAPI(resp *api.EpicsResponse) *EpicList {
	list := &EpicList{}

	list.Epics = make([]EpicMeta, len(resp.Epics))
	for i, e := range resp.Epics {
		list.Epics[i] = epicMetaValueFromAPI(e)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}
