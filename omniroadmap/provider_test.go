package omniroadmap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/omniroadmap-core/provider"
	"github.com/grokify/omniroadmap-core/provider/providertest"
)

func newTestProvider(t *testing.T, handler http.HandlerFunc) *Provider {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := aha.NewClient(
		aha.WithSubdomain("test"),
		aha.WithAPIKey("test-key"),
		aha.WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("aha.NewClient: %v", err)
	}
	return NewProvider(client)
}

func TestConformance(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"features":[],"pagination":{"total_records":0}}`))
	})

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: true,
	})
}

func TestListItems_UnsupportedKind(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected for an unsupported kind")
	})

	_, err := p.ListItems(t.Context(), &provider.ListItemsRequest{
		Kinds: []provider.ItemKind{provider.ItemKindObjective},
	})
	if err == nil {
		t.Fatal("expected an error for an unsupported kind, got nil")
	}
}

func TestListReleases_RequiresProductID(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected without a configured productID")
	})

	_, err := p.ListReleases(t.Context(), &provider.ListReleasesRequest{})
	if err == nil {
		t.Fatal("expected ErrUnsupportedOperation without a productID, got nil")
	}
}

func TestListStatuses_RequiresProductID(t *testing.T) {
	p := newTestProvider(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no HTTP call expected without a configured productID")
	})

	_, err := p.ListStatuses(t.Context(), &provider.ListStatusesRequest{})
	if err == nil {
		t.Fatal("expected ErrUnsupportedOperation without a productID, got nil")
	}
}
