package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Initiative represents an Aha initiative.
type Initiative struct {
	ID             string
	ReferenceNum   string
	Name           string
	Description    string
	Color          string
	Position       int64
	Value          float64
	Effort         float64
	Presented      bool
	StartDate      *time.Time
	EndDate        *time.Time
	Progress       float64
	ProgressSource string
	URL            string
	Resource       string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	WorkflowStatus *WorkflowStatus
	Epic           *EpicMeta
	Features       []FeatureMeta
}

// EpicMeta represents epic metadata.
type EpicMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// InitiativeList represents a paginated list of initiatives.
type InitiativeList struct {
	Initiatives []InitiativeMeta
	Pagination  Pagination
}

// InitiativeMeta represents initiative metadata in list responses.
type InitiativeMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// GetInitiative retrieves an initiative by ID or reference number.
func (c *Client) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	resp, err := c.apiClient.GetInitiative(ctx, api.GetInitiativeParams{
		InitiativeID: id,
	})
	if err != nil {
		return nil, wrapError("GetInitiative", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.InitiativeResponse:
		if i, ok := r.Initiative.Get(); ok {
			return initiativeFromAPI(i), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "initiative not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "initiative not found"}
	}
}

// ListInitiativesOptions configures ListInitiatives.
type ListInitiativesOptions struct {
	Query        string
	UpdatedSince *time.Time
	Page         int
	PerPage      int
}

// ListInitiativesOption configures a ListInitiatives call.
type ListInitiativesOption func(*ListInitiativesOptions)

// WithInitiativeQuery filters initiatives by search query.
func WithInitiativeQuery(query string) ListInitiativesOption {
	return func(o *ListInitiativesOptions) { o.Query = query }
}

// WithInitiativeUpdatedSince filters initiatives updated after the given time.
func WithInitiativeUpdatedSince(t time.Time) ListInitiativesOption {
	return func(o *ListInitiativesOptions) { o.UpdatedSince = &t }
}

// WithInitiativePage sets the page number for pagination.
func WithInitiativePage(page int) ListInitiativesOption {
	return func(o *ListInitiativesOptions) { o.Page = page }
}

// WithInitiativePerPage sets the number of results per page.
func WithInitiativePerPage(perPage int) ListInitiativesOption {
	return func(o *ListInitiativesOptions) { o.PerPage = perPage }
}

// ListInitiatives lists initiatives with optional filtering.
func (c *Client) ListInitiatives(ctx context.Context, opts ...ListInitiativesOption) (*InitiativeList, error) {
	cfg := &ListInitiativesOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListInitiativesParams{}
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

	resp, err := c.apiClient.ListInitiatives(ctx, params)
	if err != nil {
		return nil, wrapError("ListInitiatives", err)
	}

	return initiativeListFromAPI(resp), nil
}

// ListProductInitiatives lists initiatives for a product.
func (c *Client) ListProductInitiatives(ctx context.Context, productID string, opts ...ListOption) (*InitiativeList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductInitiativesParams{
		ProductID: productID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListProductInitiatives(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductInitiatives", err)
	}

	return initiativeListFromAPI(resp), nil
}

// CreateInitiativeOptions configures CreateInitiative.
type CreateInitiativeOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	EndDate        *time.Time
	Value          *float64
	Effort         *float64
	Color          string
	Presented      *bool
}

// CreateInitiativeOption configures a CreateInitiative call.
type CreateInitiativeOption func(*CreateInitiativeOptions)

// WithInitiativeName sets the initiative name.
func WithInitiativeName(name string) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.Name = name }
}

// WithInitiativeDescription sets the initiative description.
func WithInitiativeDescription(desc string) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.Description = desc }
}

// WithInitiativeStatus sets the workflow status.
func WithInitiativeStatus(status string) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.WorkflowStatus = status }
}

// WithInitiativeStartDate sets the start date.
func WithInitiativeStartDate(t time.Time) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.StartDate = &t }
}

// WithInitiativeEndDate sets the end date.
func WithInitiativeEndDate(t time.Time) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.EndDate = &t }
}

// WithInitiativeValue sets the value score.
func WithInitiativeValue(value float64) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.Value = &value }
}

// WithInitiativeEffort sets the effort score.
func WithInitiativeEffort(effort float64) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.Effort = &effort }
}

// WithInitiativeColor sets the color.
func WithInitiativeColor(color string) CreateInitiativeOption {
	return func(o *CreateInitiativeOptions) { o.Color = color }
}

// CreateInitiative creates a new initiative in a product.
func (c *Client) CreateInitiative(ctx context.Context, productID string, opts ...CreateInitiativeOption) (*Initiative, error) {
	cfg := &CreateInitiativeOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "initiative name is required"}
	}

	initiative := api.InitiativeCreate{
		Name: cfg.Name,
	}
	if cfg.Description != "" {
		initiative.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		initiative.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		initiative.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.EndDate != nil {
		initiative.EndDate = api.NewOptDate(*cfg.EndDate)
	}
	if cfg.Value != nil {
		initiative.Value = api.NewOptFloat64(*cfg.Value)
	}
	if cfg.Effort != nil {
		initiative.Effort = api.NewOptFloat64(*cfg.Effort)
	}
	if cfg.Color != "" {
		initiative.Color = api.NewOptString(cfg.Color)
	}
	if cfg.Presented != nil {
		initiative.Presented = api.NewOptBool(*cfg.Presented)
	}

	req := &api.InitiativeCreateRequest{
		Initiative: initiative,
	}

	resp, err := c.apiClient.CreateProductInitiative(ctx, req, api.CreateProductInitiativeParams{
		ProductID: productID,
	})
	if err != nil {
		return nil, wrapError("CreateInitiative", err)
	}

	if i, ok := resp.Initiative.Get(); ok {
		return initiativeFromAPI(i), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: initiative not returned"}
}

// UpdateInitiativeOptions configures UpdateInitiative.
type UpdateInitiativeOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	EndDate        *time.Time
	Value          *float64
	Effort         *float64
	Color          string
	Presented      *bool
}

// UpdateInitiativeOption configures an UpdateInitiative call.
type UpdateInitiativeOption func(*UpdateInitiativeOptions)

// UpdateInitiative updates an existing initiative.
func (c *Client) UpdateInitiative(ctx context.Context, id string, opts ...UpdateInitiativeOption) (*Initiative, error) {
	cfg := &UpdateInitiativeOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	initiative := api.InitiativeUpdate{}
	if cfg.Name != "" {
		initiative.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		initiative.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		initiative.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		initiative.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.EndDate != nil {
		initiative.EndDate = api.NewOptDate(*cfg.EndDate)
	}
	if cfg.Value != nil {
		initiative.Value = api.NewOptFloat64(*cfg.Value)
	}
	if cfg.Effort != nil {
		initiative.Effort = api.NewOptFloat64(*cfg.Effort)
	}
	if cfg.Color != "" {
		initiative.Color = api.NewOptString(cfg.Color)
	}
	if cfg.Presented != nil {
		initiative.Presented = api.NewOptBool(*cfg.Presented)
	}

	req := &api.InitiativeUpdateRequest{
		Initiative: initiative,
	}

	resp, err := c.apiClient.UpdateInitiative(ctx, req, api.UpdateInitiativeParams{
		InitiativeID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateInitiative", err)
	}

	if i, ok := resp.Initiative.Get(); ok {
		return initiativeFromAPI(i), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: initiative not returned"}
}

// initiativeFromAPI converts an API initiative to a domain initiative.
func initiativeFromAPI(i api.Initiative) *Initiative {
	initiative := &Initiative{
		ID:           i.ID,
		ReferenceNum: i.ReferenceNum,
		Name:         i.Name,
		CreatedAt:    i.CreatedAt,
	}

	if v, ok := i.Description.Get(); ok {
		initiative.Description = v
	}
	if v, ok := i.Color.Get(); ok {
		initiative.Color = v
	}
	if v, ok := i.Position.Get(); ok {
		initiative.Position = v
	}
	if v, ok := i.Value.Get(); ok {
		initiative.Value = v
	}
	if v, ok := i.Effort.Get(); ok {
		initiative.Effort = v
	}
	if v, ok := i.Presented.Get(); ok {
		initiative.Presented = v
	}
	if v, ok := i.StartDate.Get(); ok {
		initiative.StartDate = &v
	}
	if v, ok := i.EndDate.Get(); ok {
		initiative.EndDate = &v
	}
	if v, ok := i.Progress.Get(); ok {
		initiative.Progress = v
	}
	if v, ok := i.ProgressSource.Get(); ok {
		initiative.ProgressSource = v
	}
	if v, ok := i.URL.Get(); ok {
		initiative.URL = v
	}
	if v, ok := i.Resource.Get(); ok {
		initiative.Resource = v
	}
	if v, ok := i.UpdatedAt.Get(); ok {
		initiative.UpdatedAt = &v
	}
	if v, ok := i.WorkflowStatus.Get(); ok {
		initiative.WorkflowStatus = workflowStatusFromAPI(v)
	}
	if v, ok := i.Epic.Get(); ok {
		initiative.Epic = epicMetaFromAPI(v)
	}
	initiative.Features = make([]FeatureMeta, len(i.Features))
	for idx, f := range i.Features {
		initiative.Features[idx] = featureMetaFromAPI(f)
	}

	return initiative
}

// initiativeListFromAPI converts an API initiatives response to a domain initiative list.
func initiativeListFromAPI(resp *api.InitiativesResponse) *InitiativeList {
	list := &InitiativeList{}

	list.Initiatives = make([]InitiativeMeta, len(resp.Initiatives))
	for i, init := range resp.Initiatives {
		list.Initiatives[i] = initiativeMetaFromAPI(init)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// initiativeMetaFromAPI converts an API initiative meta to a domain initiative meta.
func initiativeMetaFromAPI(i api.InitiativeMeta) InitiativeMeta {
	meta := InitiativeMeta{}
	if v, ok := i.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := i.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := i.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := i.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := i.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := i.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}

// epicMetaFromAPI converts an API epic meta to a pointer.
func epicMetaFromAPI(e api.EpicMeta) *EpicMeta {
	epic := &EpicMeta{}
	if v, ok := e.ID.Get(); ok {
		epic.ID = v
	}
	if v, ok := e.ReferenceNum.Get(); ok {
		epic.ReferenceNum = v
	}
	if v, ok := e.Name.Get(); ok {
		epic.Name = v
	}
	if v, ok := e.URL.Get(); ok {
		epic.URL = v
	}
	if v, ok := e.Resource.Get(); ok {
		epic.Resource = v
	}
	if v, ok := e.CreatedAt.Get(); ok {
		epic.CreatedAt = v
	}
	return epic
}

// epicMetaValueFromAPI converts an API epic meta to a value.
func epicMetaValueFromAPI(e api.EpicMeta) EpicMeta {
	meta := EpicMeta{}
	if v, ok := e.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := e.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := e.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := e.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := e.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := e.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}
