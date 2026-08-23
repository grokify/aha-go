package graphql_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	genql "github.com/Khan/genqlient/graphql"

	"github.com/grokify/aha-go/graphql"
	"github.com/grokify/aha-go/graphql/generated"
)

func testClient(t *testing.T, handler http.HandlerFunc) genql.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return genql.NewClient(srv.URL, http.DefaultClient)
}

func TestSetCustomFieldValues_Success(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"setCustomFieldValues": {
					"customFieldValues": [
						{"id": "1", "key": "priority", "value": "High", "humanValue": "High"}
					],
					"errors": {"attributes": []}
				}
			}
		}`))
	})

	got, err := graphql.SetCustomFieldValues(context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		map[string]any{"priority": "High"})
	if err != nil {
		t.Fatalf("SetCustomFieldValues() error = %v", err)
	}
	if len(got) != 1 || got[0].Key != "priority" || got[0].Value != "High" || got[0].HumanValue != "High" {
		t.Errorf("SetCustomFieldValues() = %+v", got)
	}
}

func TestSetCustomFieldValues_PayloadErrorsSurfaced(t *testing.T) {
	// Transport-level success (HTTP 200, valid GraphQL response), but the
	// mutation payload itself reports a field-level validation error - this
	// must NOT be treated as success.
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"setCustomFieldValues": {
					"customFieldValues": [],
					"errors": {
						"attributes": [
							{"name": "priority", "fullMessages": ["Priority is not a valid custom field key"]}
						]
					}
				}
			}
		}`))
	})

	_, err := graphql.SetCustomFieldValues(context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		map[string]any{"priority": "High"})
	if err == nil {
		t.Fatal("SetCustomFieldValues() error = nil, want non-nil for payload-level validation error")
	}
}

func TestSetCustomFieldValues_TransportError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := graphql.SetCustomFieldValues(context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		map[string]any{"priority": "High"})
	if err == nil {
		t.Fatal("SetCustomFieldValues() error = nil, want non-nil for HTTP 500")
	}
}

func TestSetCustomFieldValues_MixedScalarTypes(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"setCustomFieldValues": {
					"customFieldValues": [
						{"id": "1", "key": "priority", "value": "High", "humanValue": "High"},
						{"id": "2", "key": "story_points", "value": 8, "humanValue": "8"},
						{"id": "3", "key": "is_blocked", "value": false, "humanValue": "No"}
					],
					"errors": {"attributes": []}
				}
			}
		}`))
	})

	got, err := graphql.SetCustomFieldValues(context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		map[string]any{"priority": "High", "story_points": 8, "is_blocked": false})
	if err != nil {
		t.Fatalf("SetCustomFieldValues() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 values, got %d: %+v", len(got), got)
	}
}
