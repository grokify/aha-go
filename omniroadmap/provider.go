// Package omniroadmap implements omniroadmap-core's provider.Provider for
// Aha!, wrapping aha-go's own *aha.Client. It follows the pattern used by
// elevenlabs-go/omnivoice and opik-go/integrations/omnillm: since aha-go is
// itself a from-scratch SDK (Aha publishes no official Go SDK), the adapter
// lives inside this repo rather than a separate omni-aha repo.
package omniroadmap

import (
	aha "github.com/grokify/aha-go"
	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// providerName is the registry key and Provider.Name() value for this adapter.
const providerName = "aha"

// Provider implements provider.Provider using aha-go's *aha.Client.
type Provider struct {
	client    *aha.Client
	productID string
}

var _ provider.Provider = (*Provider)(nil)

// Option configures a Provider.
type Option func(*Provider)

// WithProductID scopes product-level-only operations (ListStatuses,
// product-scoped custom field definitions) to a single Aha product
// (workspace). Aha's workflow/status endpoints have no unscoped "all
// products" variant, so ListStatuses returns
// omniroadmap.ErrUnsupportedOperation when this isn't set.
func WithProductID(productID string) Option {
	return func(p *Provider) { p.productID = productID }
}

// NewProvider wraps client as a provider.Provider.
func NewProvider(client *aha.Client, opts ...Option) *Provider {
	p := &Provider{client: client}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Provider) Name() string { return providerName }

// Close is a no-op: aha.Client owns no closable resources beyond its HTTP
// client, which the standard library manages for us.
func (p *Provider) Close() error { return nil }

func (p *Provider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Kinds: []provider.ItemKind{
			provider.ItemKindFeature,
			provider.ItemKindEpic,
			provider.ItemKindInitiative,
		},
		SupportsReleases:     true,
		SupportsObjectives:   false,
		SupportsCustomFields: true,
		SupportsWrite:        false,
	}
}

func init() {
	_ = omniroadmap.RegisterProvider(providerName, func(config any) (provider.Provider, error) {
		client, ok := config.(*aha.Client)
		if !ok {
			return nil, omniroadmap.NewAPIError(providerName, 0, "invalid_config",
				"omniroadmap/aha: expected *aha.Client config")
		}
		return NewProvider(client), nil
	})
}
