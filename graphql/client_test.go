package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("mycompany", "secret123")

	if client.subdomain != "mycompany" {
		t.Errorf("subdomain = %q, want %q", client.subdomain, "mycompany")
	}
	if client.apiKey != "secret123" {
		t.Errorf("apiKey = %q, want %q", client.apiKey, "secret123")
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestNewClientWithHTTP(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
		wantNil    bool
	}{
		{
			name:       "nil http client uses default",
			httpClient: nil,
			wantNil:    false,
		},
		{
			name:       "custom http client preserved",
			httpClient: &http.Client{},
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClientWithHTTP("sub", "key", tt.httpClient)
			if client.httpClient == nil {
				t.Error("httpClient should not be nil")
			}
		})
	}
}

func TestClientEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		subdomain string
		want      string
	}{
		{
			name:      "standard subdomain",
			subdomain: "mycompany",
			want:      "https://mycompany.aha.io/api/v2/graphql",
		},
		{
			name:      "empty subdomain",
			subdomain: "",
			want:      "https://.aha.io/api/v2/graphql",
		},
		{
			name:      "hyphenated subdomain",
			subdomain: "my-company",
			want:      "https://my-company.aha.io/api/v2/graphql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.subdomain, "key")
			got := client.Endpoint()
			if got != tt.want {
				t.Errorf("Endpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClientDo(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		wantErr        bool
		wantData       string
	}{
		{
			name:           "successful response",
			serverResponse: `{"data":{"test":"value"},"errors":null}`,
			serverStatus:   http.StatusOK,
			wantErr:        false,
			wantData:       `{"test":"value"}`,
		},
		{
			name:           "server error status",
			serverResponse: `{"error":"internal error"}`,
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "response with graphql errors",
			serverResponse: `{"data":null,"errors":[{"message":"query failed"}]}`,
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("Authorization") != "Bearer testkey" {
					t.Error("missing or incorrect Authorization header")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing or incorrect Content-Type header")
				}
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}

				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create client that points to test server
			client := &Client{
				subdomain:  "test",
				apiKey:     "testkey",
				httpClient: server.Client(),
			}
			client.SetEndpoint(server.URL)

			req := &Request{
				Query: "query { test }",
			}

			resp, err := client.Do(context.Background(), req)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.wantData != "" {
				if string(resp.Data) != tt.wantData {
					t.Errorf("Data = %s, want %s", string(resp.Data), tt.wantData)
				}
			}
		})
	}
}

func TestClientQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request to verify variables
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Return mock search results
		response := `{
			"data": {
				"searchDocuments": {
					"nodes": [{"name": "Test Doc", "url": "https://example.com", "searchableId": "123", "searchableType": "Feature"}],
					"currentPage": 1,
					"totalCount": 1,
					"totalPages": 1,
					"isLastPage": true
				}
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client := &Client{
		subdomain:  "test",
		apiKey:     "testkey",
		httpClient: server.Client(),
	}
	client.SetEndpoint(server.URL)

	var result SearchDocumentsResponse
	err := client.Query(context.Background(), SearchDocumentsQuery, map[string]any{
		"query":          "test",
		"searchableType": []string{"Feature"},
	}, &result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(result.SearchDocuments.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(result.SearchDocuments.Nodes))
	}

	if result.SearchDocuments.Nodes[0].Name != "Test Doc" {
		t.Errorf("node name = %q, want %q", result.SearchDocuments.Nodes[0].Name, "Test Doc")
	}
}

func TestClientQueryWithErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"data": null,
			"errors": [{"message": "Record not found", "path": ["page"]}]
		}`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client := &Client{
		subdomain:  "test",
		apiKey:     "testkey",
		httpClient: server.Client(),
	}
	client.SetEndpoint(server.URL)

	var result PageResponse
	err := client.Query(context.Background(), GetPageQuery, map[string]any{
		"id":            "INVALID",
		"includeParent": true,
	}, &result)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	graphqlErr, ok := err.(Error)
	if !ok {
		t.Errorf("expected Error type, got %T", err)
		return
	}

	if graphqlErr.Message != "Record not found" {
		t.Errorf("error message = %q, want %q", graphqlErr.Message, "Record not found")
	}
}

func TestError(t *testing.T) {
	err := Error{
		Message: "test error message",
		Path:    []string{"query", "field"},
	}

	if err.Error() != "test error message" {
		t.Errorf("Error() = %q, want %q", err.Error(), "test error message")
	}
}
