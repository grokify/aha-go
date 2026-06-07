// Package aha provides a Go client for the Aha.io product management API.
//
// This package provides an ergonomic wrapper around the Aha.io REST API,
// enabling Go applications to interact with Aha workspaces, features, ideas,
// releases, and other product management resources.
//
// # Quick Start
//
// Create a client and retrieve a feature:
//
//	client, err := aha.NewClient(
//	    aha.WithSubdomain("mycompany"),
//	    aha.WithAPIKey("your-api-key"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	feature, err := client.GetFeature(ctx, "PROD-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(feature.Name)
//
// # Configuration
//
// The client can be configured via options or environment variables:
//
//   - AHA_SUBDOMAIN: Your Aha account subdomain
//   - AHA_API_KEY: Your Aha API key
//
// Options take precedence over environment variables.
//
// # Error Handling
//
// API errors are returned as *APIError, which includes the HTTP status code
// and error message. Helper functions are provided for common error checks:
//
//	feature, err := client.GetFeature(ctx, "INVALID")
//	if aha.IsNotFound(err) {
//	    // Handle not found
//	}
//
// # Low-Level API Access
//
// For advanced use cases, the underlying ogen-generated client can be accessed:
//
//	apiClient := client.API()
//	// Use apiClient for operations not covered by the high-level API
package aha
