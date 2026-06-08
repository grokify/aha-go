//go:build integration

package generated_test

import (
	"context"
	"testing"

	"github.com/grokify/aha-go/graphql"
	"github.com/grokify/aha-go/graphql/generated"
)

// TestGenqlientQueries tests the genqlient-generated GraphQL client.
// Run with: go test -tags=integration -v ./graphql/generated/...
//
// Credentials can be provided via:
//   - AHA_SUBDOMAIN + AHA_API_KEY (direct)
//   - GOAUTH_CREDENTIALS_FILE + GOAUTH_ACCOUNT (goauth file)
func TestGenqlientQueries(t *testing.T) {
	creds, err := graphql.LoadTestCredentials()
	if err != nil {
		t.Skip(graphql.SkipReason())
	}

	client := graphql.NewGenqlientClient(creds.Subdomain, creds.APIKey)
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
