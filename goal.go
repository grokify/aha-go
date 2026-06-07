package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Goal represents an Aha goal.
type Goal struct {
	ID             string
	ReferenceNum   string
	Name           string
	Description    string
	Progress       float64
	ProgressSource string
	Status         string
	StartDate      *time.Time
	EndDate        *time.Time
	URL            string
	Resource       string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	TimeFrame      *TimeFrame
	WorkflowStatus *WorkflowStatus
	CustomFields   []CustomField
}

// TimeFrame represents a goal time frame.
type TimeFrame struct {
	ID   string
	Name string
}

// GoalList represents a paginated list of goals.
type GoalList struct {
	Goals      []GoalMeta
	Pagination Pagination
}

// GoalMeta represents goal metadata in list responses.
type GoalMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// GetGoal retrieves a goal by ID or reference number.
func (c *Client) GetGoal(ctx context.Context, id string) (*Goal, error) {
	resp, err := c.apiClient.GetGoal(ctx, api.GetGoalParams{
		GoalID: id,
	})
	if err != nil {
		return nil, wrapError("GetGoal", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.GoalResponse:
		if g, ok := r.Goal.Get(); ok {
			return goalFromAPI(g), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "goal not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "goal not found"}
	}
}

// ListGoalsOptions configures ListGoals.
type ListGoalsOptions struct {
	Query        string
	UpdatedSince *time.Time
	Page         int
	PerPage      int
}

// ListGoalsOption configures a ListGoals call.
type ListGoalsOption func(*ListGoalsOptions)

// WithGoalQuery filters goals by search query.
func WithGoalQuery(query string) ListGoalsOption {
	return func(o *ListGoalsOptions) { o.Query = query }
}

// WithGoalUpdatedSince filters goals updated after the given time.
func WithGoalUpdatedSince(t time.Time) ListGoalsOption {
	return func(o *ListGoalsOptions) { o.UpdatedSince = &t }
}

// WithGoalPage sets the page number for pagination.
func WithGoalPage(page int) ListGoalsOption {
	return func(o *ListGoalsOptions) { o.Page = page }
}

// WithGoalPerPage sets the number of results per page.
func WithGoalPerPage(perPage int) ListGoalsOption {
	return func(o *ListGoalsOptions) { o.PerPage = perPage }
}

// ListGoals lists goals with optional filtering.
func (c *Client) ListGoals(ctx context.Context, opts ...ListGoalsOption) (*GoalList, error) {
	cfg := &ListGoalsOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListGoalsParams{}
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

	resp, err := c.apiClient.ListGoals(ctx, params)
	if err != nil {
		return nil, wrapError("ListGoals", err)
	}

	return goalListFromAPI(resp), nil
}

// ListProductGoals lists goals for a product.
func (c *Client) ListProductGoals(ctx context.Context, productID string, opts ...ListOption) (*GoalList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductGoalsParams{
		ProductID: productID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListProductGoals(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductGoals", err)
	}

	return goalListFromAPI(resp), nil
}

// CreateGoalOptions configures CreateGoal.
type CreateGoalOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	EndDate        *time.Time
}

// CreateGoalOption configures a CreateGoal call.
type CreateGoalOption func(*CreateGoalOptions)

// WithGoalName sets the goal name.
func WithGoalName(name string) CreateGoalOption {
	return func(o *CreateGoalOptions) { o.Name = name }
}

// WithGoalDescription sets the goal description.
func WithGoalDescription(desc string) CreateGoalOption {
	return func(o *CreateGoalOptions) { o.Description = desc }
}

// WithGoalStatus sets the workflow status.
func WithGoalStatus(status string) CreateGoalOption {
	return func(o *CreateGoalOptions) { o.WorkflowStatus = status }
}

// WithGoalStartDate sets the start date.
func WithGoalStartDate(t time.Time) CreateGoalOption {
	return func(o *CreateGoalOptions) { o.StartDate = &t }
}

// WithGoalEndDate sets the end date.
func WithGoalEndDate(t time.Time) CreateGoalOption {
	return func(o *CreateGoalOptions) { o.EndDate = &t }
}

// CreateGoal creates a new goal in a product.
func (c *Client) CreateGoal(ctx context.Context, productID string, opts ...CreateGoalOption) (*Goal, error) {
	cfg := &CreateGoalOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "goal name is required"}
	}

	goal := api.GoalCreate{
		Name: cfg.Name,
	}
	if cfg.Description != "" {
		goal.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		goal.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		goal.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.EndDate != nil {
		goal.EndDate = api.NewOptDate(*cfg.EndDate)
	}

	req := &api.GoalCreateRequest{
		Goal: goal,
	}

	resp, err := c.apiClient.CreateProductGoal(ctx, req, api.CreateProductGoalParams{
		ProductID: productID,
	})
	if err != nil {
		return nil, wrapError("CreateGoal", err)
	}

	if g, ok := resp.Goal.Get(); ok {
		return goalFromAPI(g), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: goal not returned"}
}

// UpdateGoalOptions configures UpdateGoal.
type UpdateGoalOptions struct {
	Name           string
	Description    string
	WorkflowStatus string
	StartDate      *time.Time
	EndDate        *time.Time
	Progress       *float64
}

// UpdateGoalOption configures an UpdateGoal call.
type UpdateGoalOption func(*UpdateGoalOptions)

// WithUpdateGoalName sets the goal name.
func WithUpdateGoalName(name string) UpdateGoalOption {
	return func(o *UpdateGoalOptions) { o.Name = name }
}

// WithUpdateGoalDescription sets the goal description.
func WithUpdateGoalDescription(desc string) UpdateGoalOption {
	return func(o *UpdateGoalOptions) { o.Description = desc }
}

// WithUpdateGoalStatus sets the workflow status.
func WithUpdateGoalStatus(status string) UpdateGoalOption {
	return func(o *UpdateGoalOptions) { o.WorkflowStatus = status }
}

// WithUpdateGoalProgress sets the progress (when progress_source is manual).
func WithUpdateGoalProgress(progress float64) UpdateGoalOption {
	return func(o *UpdateGoalOptions) { o.Progress = &progress }
}

// UpdateGoal updates an existing goal.
func (c *Client) UpdateGoal(ctx context.Context, id string, opts ...UpdateGoalOption) (*Goal, error) {
	cfg := &UpdateGoalOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	goal := api.GoalUpdate{}
	if cfg.Name != "" {
		goal.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		goal.Description = api.NewOptString(cfg.Description)
	}
	if cfg.WorkflowStatus != "" {
		goal.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.StartDate != nil {
		goal.StartDate = api.NewOptDate(*cfg.StartDate)
	}
	if cfg.EndDate != nil {
		goal.EndDate = api.NewOptDate(*cfg.EndDate)
	}
	if cfg.Progress != nil {
		goal.Progress = api.NewOptFloat32(float32(*cfg.Progress))
	}

	req := &api.GoalUpdateRequest{
		Goal: goal,
	}

	resp, err := c.apiClient.UpdateGoal(ctx, req, api.UpdateGoalParams{
		GoalID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateGoal", err)
	}

	if g, ok := resp.Goal.Get(); ok {
		return goalFromAPI(g), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: goal not returned"}
}

// goalFromAPI converts an API goal to a domain goal.
func goalFromAPI(g api.Goal) *Goal {
	goal := &Goal{
		ID:           g.ID,
		ReferenceNum: g.ReferenceNum,
		Name:         g.Name,
		CreatedAt:    g.CreatedAt,
	}

	if v, ok := g.Description.Get(); ok {
		goal.Description = v
	}
	if v, ok := g.Progress.Get(); ok {
		goal.Progress = float64(v)
	}
	if v, ok := g.ProgressSource.Get(); ok {
		goal.ProgressSource = v
	}
	if v, ok := g.Status.Get(); ok {
		goal.Status = v
	}
	if v, ok := g.StartDate.Get(); ok {
		goal.StartDate = &v
	}
	if v, ok := g.EndDate.Get(); ok {
		goal.EndDate = &v
	}
	if v, ok := g.URL.Get(); ok {
		goal.URL = v
	}
	if v, ok := g.Resource.Get(); ok {
		goal.Resource = v
	}
	if v, ok := g.UpdatedAt.Get(); ok {
		goal.UpdatedAt = &v
	}
	if v, ok := g.TimeFrame.Get(); ok {
		goal.TimeFrame = timeFrameFromAPI(v)
	}
	if v, ok := g.WorkflowStatus.Get(); ok {
		goal.WorkflowStatus = workflowStatusFromAPI(v)
	}

	// Convert custom fields
	goal.CustomFields = customFieldsFromAPI(g.CustomFields)

	return goal
}

// goalListFromAPI converts an API goals response to a domain goal list.
func goalListFromAPI(resp *api.GoalsResponse) *GoalList {
	list := &GoalList{}

	list.Goals = make([]GoalMeta, len(resp.Goals))
	for i, g := range resp.Goals {
		list.Goals[i] = goalMetaFromAPI(g)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// goalMetaFromAPI converts an API goal meta to a domain goal meta.
func goalMetaFromAPI(g api.GoalMeta) GoalMeta {
	meta := GoalMeta{}
	if v, ok := g.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := g.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := g.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := g.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := g.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := g.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}

// timeFrameFromAPI converts an API time frame.
func timeFrameFromAPI(tf api.TimeFrame) *TimeFrame {
	timeFrame := &TimeFrame{}
	if v, ok := tf.ID.Get(); ok {
		timeFrame.ID = v
	}
	if v, ok := tf.Name.Get(); ok {
		timeFrame.Name = v
	}
	return timeFrame
}
