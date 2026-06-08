//go:build integration

package generated_test

import (
	"context"
	"os"
	"testing"

	"github.com/grokify/aha-go/graphql"
	"github.com/grokify/aha-go/graphql/generated"
)

// TestGenqlientQueries tests the genqlient-generated GraphQL client.
// Run with: go test -tags=integration -v ./graphql/generated/...
//
// Required environment variables (per official Aha API docs):
//   - AHA_SUBDOMAIN: Your Aha! subdomain (e.g., "yourcompany" for yourcompany.aha.io)
//   - AHA_API_KEY: Your Aha! API key
func TestGenqlientQueries(t *testing.T) {
	subdomain := os.Getenv("AHA_SUBDOMAIN")
	if subdomain == "" {
		t.Skip("AHA_SUBDOMAIN not set")
	}
	apiKey := os.Getenv("AHA_API_KEY")
	if apiKey == "" {
		t.Skip("AHA_API_KEY not set")
	}

	client := graphql.NewGenqlientClient(subdomain, apiKey)
	ctx := context.Background()

	// Test GetAccount (no ID = current account)
	t.Run("GetAccount", func(t *testing.T) {
		resp, err := generated.GetAccount(ctx, client, nil)
		if err != nil {
			t.Fatalf("GetAccount failed: %v", err)
		}
		t.Logf("Account: %s (domain: %s)", resp.Account.Name, resp.Account.Domain)
		if resp.Account.Name == "" {
			t.Error("Expected account name to be non-empty")
		}
	})

	// Test SearchDocuments
	t.Run("SearchDocuments", func(t *testing.T) {
		resp, err := generated.SearchDocuments(ctx, client, "product", []string{"Page", "Feature"})
		if err != nil {
			t.Fatalf("SearchDocuments failed: %v", err)
		}
		t.Logf("Search results: %d total", resp.SearchDocuments.TotalCount)
		for i, node := range resp.SearchDocuments.Nodes {
			if i >= 3 {
				break
			}
			name := ""
			if node.Name != nil {
				name = *node.Name
			}
			t.Logf("  - %s (%s)", name, node.SearchableType)
		}
	})
}
