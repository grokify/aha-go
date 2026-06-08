// Package graphql provides GraphQL clients for the Aha.io GraphQL API.
//
// This package provides two ways to interact with the Aha.io GraphQL API:
//
// 1. Manual client (Client) - A simple client for executing raw GraphQL queries
// 2. Generated client (via genqlient) - Type-safe generated functions
//
// For the generated client, use NewGenqlientClient and the functions in the
// generated subpackage:
//
//	import (
//		"github.com/grokify/aha-go/graphql"
//		"github.com/grokify/aha-go/graphql/generated"
//	)
//
//	client := graphql.NewGenqlientClient("mycompany", "api-key")
//	resp, err := generated.GetFeature(ctx, client, "FEAT-123")
//
// For the manual client, use NewClient and the Query method:
//
//	client := graphql.NewClient("mycompany", "api-key")
//	var result MyResponse
//	err := client.Query(ctx, myQuery, variables, &result)
package graphql

import (
	"net/http"

	genql "github.com/Khan/genqlient/graphql"
)

// NewGenqlientClient creates a genqlient-compatible client for use with
// the generated query functions in the generated subpackage.
//
// Example:
//
//	client := graphql.NewGenqlientClient("mycompany", "my-api-key")
//	resp, err := generated.GetFeature(ctx, client, "FEAT-123")
func NewGenqlientClient(subdomain, apiKey string) genql.Client {
	endpoint := "https://" + subdomain + ".aha.io/api/v2/graphql"
	return genql.NewClient(endpoint, &authenticatedHTTPClient{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	})
}

// NewGenqlientClientWithHTTP creates a genqlient-compatible client with a
// custom HTTP client for use with the generated query functions.
func NewGenqlientClientWithHTTP(subdomain, apiKey string, httpClient *http.Client) genql.Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	endpoint := "https://" + subdomain + ".aha.io/api/v2/graphql"
	return genql.NewClient(endpoint, &authenticatedHTTPClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	})
}

// authenticatedHTTPClient wraps an http.Client to add Aha.io authentication.
type authenticatedHTTPClient struct {
	apiKey     string
	httpClient *http.Client
}

// Do implements the genql.Doer interface, adding authentication headers.
func (c *authenticatedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	return c.httpClient.Do(req) //nolint:gosec // G704: URL is constructed from user-provided subdomain, expected for SDK
}
