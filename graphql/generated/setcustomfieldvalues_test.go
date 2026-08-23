package generated_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	genql "github.com/Khan/genqlient/graphql"

	"github.com/grokify/aha-go/graphql/generated"
)

// TestSetCustomFieldValues_ScalarResponseValue is the regression test for
// the genqlient.yaml JSON->map[string]any global scalar binding, which was
// wrong for CustomFieldValueInput.value/CustomFieldValue.value (custom
// field values are usually plain scalars, not JSON objects). Before the
// @genqlient(for: "...", bind: "any") directives were added to
// mutations.graphql, a response containing a bare string for "value"
// (the common case) failed to unmarshal, making a successful mutation
// look like a hard failure. This proves that no longer happens.
func TestSetCustomFieldValues_ScalarResponseValue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer srv.Close()

	client := genql.NewClient(srv.URL, http.DefaultClient)
	resp, err := generated.SetCustomFieldValues(
		context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		[]generated.CustomFieldValueInput{
			{Key: "priority", Value: "High"},
			{Key: "story_points", Value: 8},
			{Key: "is_blocked", Value: false},
		},
	)
	if err != nil {
		t.Fatalf("SetCustomFieldValues failed to unmarshal a scalar response value: %v", err)
	}

	values := resp.SetCustomFieldValues.CustomFieldValues
	if len(values) != 3 {
		t.Fatalf("expected 3 custom field values, got %d", len(values))
	}
	if got, want := values[0].Value, any("High"); got != want {
		t.Errorf("values[0].Value = %#v, want %#v", got, want)
	}
	// JSON numbers decode as float64 through encoding/json into `any`.
	if got, want := values[1].Value, any(float64(8)); got != want {
		t.Errorf("values[1].Value = %#v, want %#v", got, want)
	}
	if got, want := values[2].Value, any(false); got != want {
		t.Errorf("values[2].Value = %#v, want %#v", got, want)
	}
}

// TestSetCustomFieldValues_RequestBody confirms the *request* side also
// accepts bare scalars (not just the response side) - marshals the
// variables and checks the JSON shape sent over the wire.
func TestSetCustomFieldValues_RequestBody(t *testing.T) {
	var captured map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data": {"setCustomFieldValues": {"customFieldValues": [], "errors": {"attributes": []}}}}`))
	}))
	defer srv.Close()

	client := genql.NewClient(srv.URL, http.DefaultClient)
	_, err := generated.SetCustomFieldValues(
		context.Background(), client,
		"FEAT-123", generated.CustomFieldableTypeEnumFeature,
		[]generated.CustomFieldValueInput{{Key: "priority", Value: "High"}},
	)
	if err != nil {
		t.Fatalf("SetCustomFieldValues: %v", err)
	}

	vars, ok := captured["variables"].(map[string]any)
	if !ok {
		t.Fatalf("request body missing variables: %#v", captured)
	}
	cfv, ok := vars["customFieldValues"].([]any)
	if !ok || len(cfv) != 1 {
		t.Fatalf("customFieldValues variable malformed: %#v", vars["customFieldValues"])
	}
	entry, ok := cfv[0].(map[string]any)
	if !ok || entry["value"] != "High" {
		t.Fatalf("expected bare scalar value \"High\" on the wire, got %#v", cfv[0])
	}
}
