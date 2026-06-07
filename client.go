package aha

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/grokify/aha-go/internal/api"
)

const (
	// SDKVersion is the version of this SDK.
	SDKVersion = "0.1.0"

	// SDKName is the name of this SDK.
	SDKName = "aha-go"
)

// Client provides access to the Aha.io API.
type Client struct {
	config    *Config
	apiClient *api.Client
}

// NewClient creates a new Aha client with the given options.
//
// Configuration is loaded in the following order (later values override earlier):
//  1. Default values
//  2. Environment variables (AHA_SUBDOMAIN, AHA_API_KEY)
//  3. Options passed to NewClient
func NewClient(opts ...Option) (*Client, error) {
	cfg := &Config{}
	cfg.loadDefaults()
	cfg.loadEnv()

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// Create security source for authentication
	secSource := &securitySource{apiKey: cfg.APIKey}

	// Create custom HTTP client with SDK headers
	httpClient := &sdkHTTPClient{
		client: cfg.HTTPClient,
	}

	// Build base URL
	baseURL := cfg.buildBaseURL()

	// Create ogen client
	apiClient, err := api.NewClient(baseURL, secSource, api.WithClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("creating API client: %w", err)
	}

	return &Client{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

// Subdomain returns the configured Aha subdomain.
func (c *Client) Subdomain() string {
	return c.config.Subdomain
}

// BaseURL returns the API base URL.
func (c *Client) BaseURL() string {
	return c.config.buildBaseURL()
}

// API returns the low-level ogen-generated client for advanced use.
// This allows access to API operations not covered by the high-level wrapper.
func (c *Client) API() *api.Client {
	return c.apiClient
}

// HTTPClient returns the underlying HTTP client for raw requests.
func (c *Client) HTTPClient() *http.Client {
	return c.config.HTTPClient
}

// APIKey returns the configured API key.
func (c *Client) APIKey() string {
	return c.config.APIKey
}

// DoRaw performs a raw HTTP request to the Aha API.
// The path should be a relative path like "/api/v1/features/123".
// This is useful for API operations not yet covered by typed methods.
func (c *Client) DoRaw(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	url := c.BaseURL() + path

	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add authentication
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// Add SDK headers
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", SDKName, SDKVersion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.config.HTTPClient.Do(req)
}

// securitySource implements api.SecuritySource for bearer token auth.
type securitySource struct {
	apiKey string
}

// BearerAuth returns the bearer token for authentication.
func (s *securitySource) BearerAuth(ctx context.Context, operationName api.OperationName) (api.BearerAuth, error) {
	return api.BearerAuth{
		Token: s.apiKey,
	}, nil
}

// sdkHTTPClient wraps an HTTP client to add SDK headers.
type sdkHTTPClient struct {
	client *http.Client
}

// Do executes an HTTP request with SDK headers.
func (c *sdkHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add SDK identification headers
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", SDKName, SDKVersion))
	req.Header.Set("X-Aha-SDK-Version", SDKVersion)
	req.Header.Set("X-Aha-SDK-Lang", "go")

	// Add content type if not set
	if req.Header.Get("Content-Type") == "" && req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add accept header if not set
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	return c.client.Do(req)
}
