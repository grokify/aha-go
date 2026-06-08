// Package example provides a handwritten GraphQL client for the Aha.io API.
// This package serves as example code for users learning to write GraphQL
// clients without code generation. For production use, see the generated
// client in github.com/grokify/aha-go/graphql/generated.
package example

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is a GraphQL client for Aha.io.
type Client struct {
	subdomain  string
	apiKey     string
	httpClient *http.Client
	endpoint   string // optional override for testing
}

// NewClient creates a new GraphQL client.
func NewClient(subdomain, apiKey string) *Client {
	return &Client{
		subdomain:  subdomain,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
}

// NewClientWithHTTP creates a new GraphQL client with a custom HTTP client.
func NewClientWithHTTP(subdomain, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		subdomain:  subdomain,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// Endpoint returns the GraphQL endpoint URL.
func (c *Client) Endpoint() string {
	if c.endpoint != "" {
		return c.endpoint
	}
	// Aha.io GraphQL API is at /api/v2/graphql (not /api/graphql)
	return fmt.Sprintf("https://%s.aha.io/api/v2/graphql", c.subdomain)
}

// SetEndpoint sets a custom endpoint URL (for testing).
func (c *Client) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

// Request represents a GraphQL request.
type Request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// Response represents a GraphQL response.
type Response struct {
	Data   json.RawMessage `json:"data"`
	Errors []Error         `json:"errors,omitempty"`
}

// Error represents a GraphQL error.
type Error struct {
	Message    string   `json:"message"`
	Path       []string `json:"path,omitempty"`
	Extensions any      `json:"extensions,omitempty"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message
}

// Do executes a GraphQL request and returns the raw response.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp Response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &resp, nil
}

// Query executes a GraphQL query and unmarshals the data into result.
func (c *Client) Query(ctx context.Context, query string, variables map[string]any, result any) error {
	req := &Request{
		Query:     query,
		Variables: variables,
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		return resp.Errors[0]
	}

	if result != nil && len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, result); err != nil {
			return fmt.Errorf("unmarshaling data: %w", err)
		}
	}

	return nil
}
