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
	ProductLine       bool
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
	CreatedAt       time.Time
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

// ListProducts lists all products (workspaces).
func (c *Client) ListProducts(ctx context.Context, opts ...ListOption) (*ProductList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductsParams{}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListProducts(ctx, params)
	if err != nil {
		return nil, wrapError("ListProducts", err)
	}

	return productListFromAPI(resp), nil
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
	if v, ok := p.ProductLine.Get(); ok {
		product.ProductLine = v
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
	if v, ok := p.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	return meta
}
