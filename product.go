package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Product represents an Aha product (workspace).
type Product struct {
	ID                string
	ReferencePrefix   string
	Name              string
	Description       string
	ProductLine       bool
	ParentID          string
	WorkspaceType     string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	URL               string
	Resource          string
	HasIdeas          bool
	HasMasterFeatures bool
}

// ProductList represents a paginated list of products.
type ProductList struct {
	Products   []ProductMeta
	Pagination Pagination
}

// ProductMeta represents product metadata in list responses.
type ProductMeta struct {
	ID              string
	ReferencePrefix string
	Name            string
	ProductLine     bool
	WorkspaceType   string
	CreatedAt       time.Time
}

// ListProductsOptions configures ListProducts.
type ListProductsOptions struct {
	UpdatedSince    *time.Time
	WithIdeaPortals bool
	Page            int
	PerPage         int
}

// ListProductsOption configures a ListProducts call.
type ListProductsOption func(*ListProductsOptions)

// WithUpdatedSince filters to products updated after the given time.
func WithUpdatedSince(t time.Time) ListProductsOption {
	return func(o *ListProductsOptions) { o.UpdatedSince = &t }
}

// WithIdeaPortals filters to only products with idea portals.
func WithIdeaPortals() ListProductsOption {
	return func(o *ListProductsOptions) { o.WithIdeaPortals = true }
}

// WithProductsPage sets the page number for listing.
func WithProductsPage(page int) ListProductsOption {
	return func(o *ListProductsOptions) { o.Page = page }
}

// WithProductsPerPage sets the results per page.
func WithProductsPerPage(perPage int) ListProductsOption {
	return func(o *ListProductsOptions) { o.PerPage = perPage }
}

// GetProduct retrieves a product by ID or reference prefix.
func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	resp, err := c.apiClient.GetProduct(ctx, api.GetProductParams{
		ProductID: id,
	})
	if err != nil {
		return nil, wrapError("GetProduct", err)
	}

	if p, ok := resp.Product.Get(); ok {
		return productFromAPI(p), nil
	}
	return nil, &APIError{StatusCode: 404, Message: "product not found"}
}

// ListProducts lists all products (workspaces) including Aha! Develop teams.
func (c *Client) ListProducts(ctx context.Context, opts ...ListProductsOption) (*ProductList, error) {
	options := &ListProductsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := api.ListProductsParams{}
	if options.UpdatedSince != nil {
		params.UpdatedSince = api.NewOptDateTime(*options.UpdatedSince)
	}
	if options.WithIdeaPortals {
		params.WithIdeaPortals = api.NewOptBool(true)
	}
	if options.Page > 0 {
		params.Page = api.NewOptInt32(int32(options.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if options.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(options.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListProducts(ctx, params)
	if err != nil {
		return nil, wrapError("ListProducts", err)
	}

	return productListFromAPI(resp), nil
}

// CreateProductOptions configures CreateProduct.
type CreateProductOptions struct {
	Name            string
	ReferencePrefix string
	Description     string
	ParentID        string
	WorkspaceType   string
	ProductLine     bool
	ProductLineType string
}

// CreateProductOption configures a CreateProduct call.
type CreateProductOption func(*CreateProductOptions)

// WithProductName sets the product name.
func WithProductName(name string) CreateProductOption {
	return func(o *CreateProductOptions) { o.Name = name }
}

// WithProductReferencePrefix sets the reference prefix.
func WithProductReferencePrefix(prefix string) CreateProductOption {
	return func(o *CreateProductOptions) { o.ReferencePrefix = prefix }
}

// WithProductDescription sets the product description.
func WithProductDescription(desc string) CreateProductOption {
	return func(o *CreateProductOptions) { o.Description = desc }
}

// WithProductParentID sets the parent product line ID.
func WithProductParentID(parentID string) CreateProductOption {
	return func(o *CreateProductOptions) { o.ParentID = parentID }
}

// WithProductWorkspaceType sets the workspace type.
func WithProductWorkspaceType(workspaceType string) CreateProductOption {
	return func(o *CreateProductOptions) { o.WorkspaceType = workspaceType }
}

// WithProductLine marks this as a product line.
func WithProductLine(productLineType string) CreateProductOption {
	return func(o *CreateProductOptions) {
		o.ProductLine = true
		o.ProductLineType = productLineType
	}
}

// CreateProduct creates a new product.
func (c *Client) CreateProduct(ctx context.Context, name, referencePrefix string, opts ...CreateProductOption) (*Product, error) {
	options := &CreateProductOptions{
		Name:            name,
		ReferencePrefix: referencePrefix,
	}
	for _, opt := range opts {
		opt(options)
	}

	productCreate := api.ProductCreate{
		Name:            options.Name,
		ReferencePrefix: options.ReferencePrefix,
	}
	if options.Description != "" {
		productCreate.Description = api.NewOptString(options.Description)
	}
	if options.ParentID != "" {
		productCreate.ParentID = api.NewOptString(options.ParentID)
	}
	if options.WorkspaceType != "" {
		productCreate.WorkspaceType = api.NewOptString(options.WorkspaceType)
	}
	if options.ProductLine {
		productCreate.ProductLine = api.NewOptBool(true)
		if options.ProductLineType != "" {
			productCreate.ProductLineType = api.NewOptString(options.ProductLineType)
		}
	}

	req := &api.ProductCreateRequest{
		Product: productCreate,
	}

	resp, err := c.apiClient.CreateProduct(ctx, req)
	if err != nil {
		return nil, wrapError("CreateProduct", err)
	}

	if p, ok := resp.Product.Get(); ok {
		return productFromAPI(p), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: product not returned"}
}

// CreateProductLine creates a new product line.
func (c *Client) CreateProductLine(ctx context.Context, name, referencePrefix, productLineType string, opts ...CreateProductOption) (*Product, error) {
	allOpts := append([]CreateProductOption{WithProductLine(productLineType)}, opts...)
	return c.CreateProduct(ctx, name, referencePrefix, allOpts...)
}

// UpdateProductOptions configures UpdateProduct.
type UpdateProductOptions struct {
	Name            string
	ReferencePrefix string
	Description     string
	ParentID        string
	WorkspaceType   string
}

// UpdateProductOption configures an UpdateProduct call.
type UpdateProductOption func(*UpdateProductOptions)

// WithUpdateProductName sets the product name.
func WithUpdateProductName(name string) UpdateProductOption {
	return func(o *UpdateProductOptions) { o.Name = name }
}

// WithUpdateProductReferencePrefix sets the reference prefix.
func WithUpdateProductReferencePrefix(prefix string) UpdateProductOption {
	return func(o *UpdateProductOptions) { o.ReferencePrefix = prefix }
}

// WithUpdateProductDescription sets the product description.
func WithUpdateProductDescription(desc string) UpdateProductOption {
	return func(o *UpdateProductOptions) { o.Description = desc }
}

// WithUpdateProductParentID sets the parent product line ID.
func WithUpdateProductParentID(parentID string) UpdateProductOption {
	return func(o *UpdateProductOptions) { o.ParentID = parentID }
}

// WithUpdateProductWorkspaceType sets the workspace type.
func WithUpdateProductWorkspaceType(workspaceType string) UpdateProductOption {
	return func(o *UpdateProductOptions) { o.WorkspaceType = workspaceType }
}

// UpdateProduct updates an existing product.
func (c *Client) UpdateProduct(ctx context.Context, productID string, opts ...UpdateProductOption) (*Product, error) {
	options := &UpdateProductOptions{}
	for _, opt := range opts {
		opt(options)
	}

	productUpdate := api.ProductUpdate{}
	if options.Name != "" {
		productUpdate.Name = api.NewOptString(options.Name)
	}
	if options.ReferencePrefix != "" {
		productUpdate.ReferencePrefix = api.NewOptString(options.ReferencePrefix)
	}
	if options.Description != "" {
		productUpdate.Description = api.NewOptString(options.Description)
	}
	if options.ParentID != "" {
		productUpdate.ParentID = api.NewOptString(options.ParentID)
	}
	if options.WorkspaceType != "" {
		productUpdate.WorkspaceType = api.NewOptString(options.WorkspaceType)
	}

	req := &api.ProductUpdateRequest{
		Product: productUpdate,
	}

	resp, err := c.apiClient.UpdateProduct(ctx, req, api.UpdateProductParams{
		ProductID: productID,
	})
	if err != nil {
		return nil, wrapError("UpdateProduct", err)
	}

	if p, ok := resp.Product.Get(); ok {
		return productFromAPI(p), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: product not returned"}
}

// productFromAPI converts an API product to a domain product.
func productFromAPI(p api.Product) *Product {
	product := &Product{}
	if v, ok := p.ID.Get(); ok {
		product.ID = v
	}
	if v, ok := p.ReferencePrefix.Get(); ok {
		product.ReferencePrefix = v
	}
	if v, ok := p.Name.Get(); ok {
		product.Name = v
	}
	if v, ok := p.Description.Get(); ok {
		product.Description = v
	}
	if v, ok := p.ProductLine.Get(); ok {
		product.ProductLine = v
	}
	if v, ok := p.ParentID.Get(); ok {
		product.ParentID = v
	}
	if v, ok := p.WorkspaceType.Get(); ok {
		product.WorkspaceType = v
	}
	if v, ok := p.CreatedAt.Get(); ok {
		product.CreatedAt = v
	}
	if v, ok := p.UpdatedAt.Get(); ok {
		product.UpdatedAt = &v
	}
	if v, ok := p.URL.Get(); ok {
		product.URL = v
	}
	if v, ok := p.Resource.Get(); ok {
		product.Resource = v
	}
	if v, ok := p.HasIdeas.Get(); ok {
		product.HasIdeas = v
	}
	if v, ok := p.HasMasterFeatures.Get(); ok {
		product.HasMasterFeatures = v
	}
	return product
}

// productListFromAPI converts an API products response to a domain product list.
func productListFromAPI(resp *api.ProductsResponse) *ProductList {
	list := &ProductList{}

	list.Products = make([]ProductMeta, len(resp.Products))
	for i, p := range resp.Products {
		list.Products[i] = productMetaFromAPI(p)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// productMetaFromAPI converts an API product meta to a domain product meta.
func productMetaFromAPI(p api.ProductMeta) ProductMeta {
	meta := ProductMeta{}
	if v, ok := p.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := p.ReferencePrefix.Get(); ok {
		meta.ReferencePrefix = v
	}
	if v, ok := p.Name.Get(); ok {
		meta.Name = v
	}
	if v, ok := p.ProductLine.Get(); ok {
		meta.ProductLine = v
	}
	if v, ok := p.WorkspaceType.Get(); ok {
		meta.WorkspaceType = v
	}
	if v, ok := p.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}
