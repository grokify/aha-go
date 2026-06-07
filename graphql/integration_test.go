// +build integration

package graphql_test

import (
	"context"
	"os"
	"testing"

	"github.com/grokify/aha-go/graphql"
)

// TestSearchDocumentsIntegration tests the GraphQL search against a real Aha.io account.
// Run with: go test -tags=integration -v ./graphql/...
//
// Required environment variables:
//   - AHA_SUBDOMAIN: Your Aha.io subdomain
//   - AHA_API_KEY: Your Aha.io API key
func TestSearchDocumentsIntegration(t *testing.T) {
	subdomain := os.Getenv("AHA_SUBDOMAIN")
	if subdomain == "" {
		t.Skip("AHA_SUBDOMAIN not set, skipping integration test")
	}

	apiKey := os.Getenv("AHA_API_KEY")
	if apiKey == "" {
		t.Skip("AHA_API_KEY not set, skipping integration test")
	}

	client := graphql.NewClient(subdomain, apiKey)

	t.Logf("Testing GraphQL endpoint: %s", client.Endpoint())

	// Test search
	var result graphql.SearchDocumentsResponse
	err := client.Query(context.Background(), graphql.SearchDocumentsQuery, map[string]any{
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

// TestSearchDocumentsFeatures tests searching for features specifically.
func TestSearchDocumentsFeaturesIntegration(t *testing.T) {
	subdomain := os.Getenv("AHA_SUBDOMAIN")
	if subdomain == "" {
		t.Skip("AHA_SUBDOMAIN not set, skipping integration test")
	}

	apiKey := os.Getenv("AHA_API_KEY")
	if apiKey == "" {
		t.Skip("AHA_API_KEY not set, skipping integration test")
	}

	client := graphql.NewClient(subdomain, apiKey)

	var result graphql.SearchDocumentsResponse
	err := client.Query(context.Background(), graphql.SearchDocumentsQuery, map[string]any{
		"query":          "feature",
		"searchableType": []string{"Feature"},
	}, &result)

	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	t.Logf("Feature search results: %d total", result.SearchDocuments.TotalCount)
}
