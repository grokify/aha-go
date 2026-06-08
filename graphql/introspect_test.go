//go:build integration

package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/grokify/aha-go/graphql"
)

// TestIntrospection tests if the Aha.io GraphQL API supports introspection.
// Run with: go test -tags=integration -v -run TestIntrospection ./graphql/...
//
// Credentials can be provided via:
//   - AHA_SUBDOMAIN + AHA_API_KEY (direct)
//   - GOAUTH_CREDENTIALS_FILE + GOAUTH_ACCOUNT (goauth file)
func TestIntrospection(t *testing.T) {
	creds, err := graphql.LoadTestCredentials()
	if err != nil {
		t.Skip(graphql.SkipReason())
	}

	// Introspection query to get schema
	introspectionQuery := `{
		__schema {
			types {
				name
				kind
				description
			}
			queryType {
				name
				fields {
					name
					description
				}
			}
		}
	}`

	endpoint := fmt.Sprintf("https://%s.aha.io/api/v2/graphql", creds.Subdomain)
	t.Logf("Testing introspection at: %s", endpoint)

	reqBody := map[string]any{
		"query": introspectionQuery,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+creds.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	t.Logf("Status: %d", resp.StatusCode)
	t.Logf("Response (first 2000 chars): %s", string(respBody)[:min(2000, len(respBody))])

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Introspection failed with status %d", resp.StatusCode)
	}

	// Check if introspection is disabled
	if bytes.Contains(respBody, []byte("introspection")) && bytes.Contains(respBody, []byte("disabled")) {
		t.Log("Introspection appears to be disabled")
	}

	// Parse response to see what we got
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if errors, ok := result["errors"]; ok {
		t.Logf("GraphQL errors: %v", errors)
	}

	if data, ok := result["data"]; ok {
		t.Logf("Got data: introspection is supported!")
		if schema, ok := data.(map[string]any)["__schema"]; ok {
			if types, ok := schema.(map[string]any)["types"].([]any); ok {
				t.Logf("Found %d types in schema", len(types))
			}
		}
	}
}
