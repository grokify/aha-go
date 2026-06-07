package aha

import (
	"net/http"
	"time"
)

// Option configures the Aha client.
type Option func(*Config)

// WithSubdomain sets the Aha account subdomain.
func WithSubdomain(subdomain string) Option {
	return func(c *Config) {
		c.Subdomain = subdomain
	}
}

// WithAPIKey sets the Aha API key.
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithBaseURL overrides the default API URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// ListOptions configures list operations.
type ListOptions struct {
	Page    int
	PerPage int
}

// ListOption configures a list operation.
type ListOption func(*ListOptions)

// WithPage sets the page number for pagination.
func WithPage(page int) ListOption {
	return func(o *ListOptions) {
		o.Page = page
	}
}

// WithPerPage sets the number of results per page.
func WithPerPage(perPage int) ListOption {
	return func(o *ListOptions) {
		o.PerPage = perPage
	}
}

// applyListOptions applies list options and returns the configured options.
func applyListOptions(opts ...ListOption) *ListOptions {
	o := &ListOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
