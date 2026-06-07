package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// StrategicModel represents an Aha strategic model (canvas).
type StrategicModel struct {
	ID           string
	ReferenceNum string
	Name         string
	Kind         string // e.g., "Opportunity", "Lean Canvas", "Business Model"
	Description  string
	URL          string
	Resource     string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Project      *ProductMeta
	Components   []StrategicModelComponent
}

// StrategicModelComponent represents a block within a strategic model.
type StrategicModelComponent struct {
	ID          string
	Name        string // e.g., "Users & Customers", "Problems", "Solutions"
	Description string // Content of the block (HTML)
	Position    int64
}

// StrategicModelList represents a paginated list of strategic models.
type StrategicModelList struct {
	StrategicModels []StrategicModelMeta
	Pagination      Pagination
}

// StrategicModelMeta represents strategic model metadata in list responses.
type StrategicModelMeta struct {
	ID           string
	ReferenceNum string
	Name         string
	Kind         string
	URL          string
	Resource     string
	CreatedAt    time.Time
}

// GetStrategicModel retrieves a strategic model by ID or reference number.
func (c *Client) GetStrategicModel(ctx context.Context, id string) (*StrategicModel, error) {
	resp, err := c.apiClient.GetStrategicModel(ctx, api.GetStrategicModelParams{
		StrategyModelID: id,
	})
	if err != nil {
		return nil, wrapError("GetStrategicModel", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.StrategicModelResponse:
		if sm, ok := r.StrategyModel.Get(); ok {
			return strategicModelFromAPI(sm), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "strategic model not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "strategic model not found"}
	}
}

// ListStrategicModelsOptions configures ListStrategicModels.
type ListStrategicModelsOptions struct {
	Kind    string // Filter by model kind
	Page    int
	PerPage int
}

// ListStrategicModelsOption configures a ListStrategicModels call.
type ListStrategicModelsOption func(*ListStrategicModelsOptions)

// WithStrategicModelKind filters strategic models by kind.
func WithStrategicModelKind(kind string) ListStrategicModelsOption {
	return func(o *ListStrategicModelsOptions) { o.Kind = kind }
}

// ListStrategicModels lists strategic models with optional filtering.
func (c *Client) ListStrategicModels(ctx context.Context, opts ...ListStrategicModelsOption) (*StrategicModelList, error) {
	cfg := &ListStrategicModelsOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListStrategicModelsParams{}
	if cfg.Kind != "" {
		params.Kind = api.NewOptString(cfg.Kind)
	}
	if cfg.Page > 0 {
		params.Page = api.NewOptInt32(int32(cfg.Page))
	}
	if cfg.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(cfg.PerPage))
	}

	resp, err := c.apiClient.ListStrategicModels(ctx, params)
	if err != nil {
		return nil, wrapError("ListStrategicModels", err)
	}

	return strategicModelListFromAPI(resp), nil
}

// ListProductStrategicModels lists strategic models for a product.
func (c *Client) ListProductStrategicModels(ctx context.Context, productID string, opts ...ListStrategicModelsOption) (*StrategicModelList, error) {
	cfg := &ListStrategicModelsOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListProductStrategicModelsParams{
		ProductID: productID,
	}
	if cfg.Kind != "" {
		params.Kind = api.NewOptString(cfg.Kind)
	}
	if cfg.Page > 0 {
		params.Page = api.NewOptInt32(int32(cfg.Page))
	}
	if cfg.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(cfg.PerPage))
	}

	resp, err := c.apiClient.ListProductStrategicModels(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductStrategicModels", err)
	}

	return strategicModelListFromAPI(resp), nil
}

// CreateStrategicModelOptions configures CreateStrategicModel.
type CreateStrategicModelOptions struct {
	Name        string
	Kind        string // Required: e.g., "Opportunity", "Lean Canvas"
	Description string
}

// CreateStrategicModelOption configures a CreateStrategicModel call.
type CreateStrategicModelOption func(*CreateStrategicModelOptions)

// WithStrategicModelName sets the strategic model name.
func WithStrategicModelName(name string) CreateStrategicModelOption {
	return func(o *CreateStrategicModelOptions) { o.Name = name }
}

// WithStrategicModelDescription sets the strategic model description.
func WithStrategicModelDescription(desc string) CreateStrategicModelOption {
	return func(o *CreateStrategicModelOptions) { o.Description = desc }
}

// CreateStrategicModel creates a new strategic model in a product.
func (c *Client) CreateStrategicModel(ctx context.Context, productID string, kind string, opts ...CreateStrategicModelOption) (*StrategicModel, error) {
	cfg := &CreateStrategicModelOptions{Kind: kind}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "strategic model name is required"}
	}
	if cfg.Kind == "" {
		return nil, &APIError{StatusCode: 400, Message: "strategic model kind is required"}
	}

	model := api.StrategicModelCreate{
		Name: cfg.Name,
		Kind: cfg.Kind,
	}
	if cfg.Description != "" {
		model.Description = api.NewOptString(cfg.Description)
	}

	req := &api.StrategicModelCreateRequest{
		StrategyModel: model,
	}

	resp, err := c.apiClient.CreateProductStrategicModel(ctx, req, api.CreateProductStrategicModelParams{
		ProductID: productID,
	})
	if err != nil {
		return nil, wrapError("CreateStrategicModel", err)
	}

	if sm, ok := resp.StrategyModel.Get(); ok {
		return strategicModelFromAPI(sm), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: strategic model not returned"}
}

// UpdateStrategicModelOptions configures UpdateStrategicModel.
type UpdateStrategicModelOptions struct {
	Name        string
	Description string
}

// UpdateStrategicModelOption configures an UpdateStrategicModel call.
type UpdateStrategicModelOption func(*UpdateStrategicModelOptions)

// UpdateStrategicModel updates an existing strategic model.
func (c *Client) UpdateStrategicModel(ctx context.Context, id string, opts ...UpdateStrategicModelOption) (*StrategicModel, error) {
	cfg := &UpdateStrategicModelOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	model := api.StrategicModelUpdate{}
	if cfg.Name != "" {
		model.Name = api.NewOptString(cfg.Name)
	}
	if cfg.Description != "" {
		model.Description = api.NewOptString(cfg.Description)
	}

	req := &api.StrategicModelUpdateRequest{
		StrategyModel: model,
	}

	resp, err := c.apiClient.UpdateStrategicModel(ctx, req, api.UpdateStrategicModelParams{
		StrategyModelID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateStrategicModel", err)
	}

	if sm, ok := resp.StrategyModel.Get(); ok {
		return strategicModelFromAPI(sm), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: strategic model not returned"}
}

// UpdateStrategicModelComponent updates a component within a strategic model.
func (c *Client) UpdateStrategicModelComponent(ctx context.Context, modelID, componentID, description string) (*StrategicModelComponent, error) {
	component := api.StrategicModelComponentUpdate{}
	if description != "" {
		component.Description = api.NewOptString(description)
	}

	req := &api.StrategicModelComponentUpdateRequest{
		Component: component,
	}

	resp, err := c.apiClient.UpdateStrategicModelComponent(ctx, req, api.UpdateStrategicModelComponentParams{
		StrategyModelID: modelID,
		ComponentID:     componentID,
	})
	if err != nil {
		return nil, wrapError("UpdateStrategicModelComponent", err)
	}

	if comp, ok := resp.Component.Get(); ok {
		return strategicModelComponentFromAPI(comp), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: component not returned"}
}

// strategicModelFromAPI converts an API strategic model to a domain strategic model.
func strategicModelFromAPI(sm api.StrategicModel) *StrategicModel {
	model := &StrategicModel{
		ID:           sm.ID,
		ReferenceNum: sm.ReferenceNum,
		Name:         sm.Name,
		Kind:         sm.Kind,
		CreatedAt:    sm.CreatedAt,
	}

	if v, ok := sm.Description.Get(); ok {
		model.Description = v
	}
	if v, ok := sm.URL.Get(); ok {
		model.URL = v
	}
	if v, ok := sm.Resource.Get(); ok {
		model.Resource = v
	}
	if v, ok := sm.UpdatedAt.Get(); ok {
		model.UpdatedAt = &v
	}
	if v, ok := sm.Project.Get(); ok {
		pm := productMetaFromAPI(v)
		model.Project = &pm
	}

	model.Components = make([]StrategicModelComponent, len(sm.Components))
	for i, comp := range sm.Components {
		if c := strategicModelComponentFromAPI(comp); c != nil {
			model.Components[i] = *c
		}
	}

	return model
}

// strategicModelListFromAPI converts an API strategic models response to a domain list.
func strategicModelListFromAPI(resp *api.StrategicModelsResponse) *StrategicModelList {
	list := &StrategicModelList{}

	list.StrategicModels = make([]StrategicModelMeta, len(resp.StrategyModels))
	for i, sm := range resp.StrategyModels {
		list.StrategicModels[i] = strategicModelMetaFromAPI(sm)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// strategicModelMetaFromAPI converts an API strategic model meta to a domain meta.
func strategicModelMetaFromAPI(sm api.StrategicModelMeta) StrategicModelMeta {
	meta := StrategicModelMeta{}
	if v, ok := sm.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := sm.ReferenceNum.Get(); ok {
		meta.ReferenceNum = v
	}
	if v, ok := sm.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := sm.Kind.Get(); ok {
		meta.Kind = v
	}
	if v, ok := sm.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := sm.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := sm.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}

// strategicModelComponentFromAPI converts an API component to a domain component.
func strategicModelComponentFromAPI(comp api.StrategicModelComponent) *StrategicModelComponent {
	c := &StrategicModelComponent{}
	if v, ok := comp.ID.Get(); ok {
		c.ID = v
	}
	if v, ok := comp.Name.Get(); ok {
		c.Name = v
	}
	if v, ok := comp.Description.Get(); ok {
		c.Description = v
	}
	if v, ok := comp.Position.Get(); ok {
		c.Position = v
	}
	return c
}
