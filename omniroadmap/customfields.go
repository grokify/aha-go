package omniroadmap

import (
	"context"

	"github.com/grokify/omniroadmap-core/provider"
)

// ListCustomFieldDefinitions lists custom field definitions. Unlike
// ListStatuses, aha-go's ListCustomFieldDefinitions has an unscoped
// variant covering all products the client can access, so no productID
// configuration is required. req.Kind is currently ignored — Aha's API
// doesn't let us filter definitions by CustomFieldableType server-side in
// the unscoped call, so all definitions across all kinds are returned.
func (p *Provider) ListCustomFieldDefinitions(ctx context.Context, req *provider.ListCustomFieldDefinitionsRequest) (*provider.ListCustomFieldDefinitionsResponse, error) {
	defs, err := p.client.ListCustomFieldDefinitions(ctx)
	if err != nil {
		return nil, wrapErr("ListCustomFieldDefinitions", err)
	}
	return &provider.ListCustomFieldDefinitionsResponse{
		Definitions: customFieldDefinitionsFromAha(defs),
	}, nil
}
