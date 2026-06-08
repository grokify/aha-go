//go:build integration

package example_test

import (
	"context"
	"testing"

	"github.com/grokify/aha-go/graphql"
	"github.com/grokify/aha-go/graphql/example"
)

// TestSearchDocumentsIntegration tests the handwritten GraphQL client.
// Run with: go test -tags=integration -v ./graphql/example/...
//
// Credentials can be provided via:
//   - AHA_SUBDOMAIN + AHA_API_KEY (direct)
//   - GOAUTH_CREDENTIALS_FILE + GOAUTH_ACCOUNT (goauth file)
func TestSearchDocumentsIntegration(t *testing.T) {
	creds, err := graphql.LoadTestCredentials()
	if err != nil {
		t.Skip(graphql.SkipReason())
	}

	client := example.NewClient(creds.Subdomain, creds.APIKey)

	t.Logf("Testing GraphQL endpoint: %s", client.Endpoint())

	// Test search
	var result example.SearchDocumentsResponse
	err = client.Query(context.Background(), example.SearchDocumentsQuery, map[string]any{
		"query":          "test",
		"searchableType": []string{"Page"},
	}, &result)

	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	t.Logf("Search results: %d total, page %d of %d",
		result.SearchDocuments.TotalCount,
		result.SearchDocuments.CurrentPage,
		result.SearchDocuments.TotalPages)

	for i, node := range result.SearchDocuments.Nodes {
		t.Logf("  [%d] %s (%s): %s", i+1, node.Name, node.SearchableType, node.URL)
	}
}

// TestSearchDocumentsFeaturesIntegration tests searching for features specifically.
func TestSearchDocumentsFeaturesIntegration(t *testing.T) {
	creds, err := graphql.LoadTestCredentials()
	if err != nil {
		t.Skip(graphql.SkipReason())
	}

	client := example.NewClient(creds.Subdomain, creds.APIKey)

	var result example.SearchDocumentsResponse
	err = client.Query(context.Background(), example.SearchDocumentsQuery, map[string]any{
		"query":          "feature",
		"searchableType": []string{"Feature"},
	}, &result)

	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	t.Logf("Feature search results: %d total", result.SearchDocuments.TotalCount)
}
