// Package graphql provides a type-safe GraphQL client for the Aha.io API.
//
// Use NewGenqlientClient and the functions in the generated subpackage:
//
//	import (
//		"github.com/grokify/aha-go/graphql"
//		"github.com/grokify/aha-go/graphql/generated"
//	)
//
//	client := graphql.NewGenqlientClient("mycompany", "api-key")
//	resp, err := generated.GetFeature(ctx, client, "FEAT-123")
//
// For example code showing how to write a GraphQL client without code
// generation, see the github.com/grokify/aha-go/graphql/example package.
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
