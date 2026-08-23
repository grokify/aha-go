package omniroadmap

import (
	"context"

	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListStatuses lists workflow statuses for the configured product. Aha's
// workflow endpoint (ListProductWorkflows) has no unscoped "all products"
// variant, so this requires a productID to have been set via WithProductID
// — without one, it returns omniroadmap.ErrUnsupportedOperation.
func (p *Provider) ListStatuses(ctx context.Context, req *provider.ListStatusesRequest) (*provider.ListStatusesResponse, error) {
	if p.productID == "" {
		return nil, omniroadmap.ErrUnsupportedOperation
	}

	workflows, err := p.client.ListProductWorkflows(ctx, p.productID)
	if err != nil {
		return nil, wrapErr("ListProductWorkflows", err)
	}

	seen := map[string]bool{}
	var statuses []provider.Status
	for _, wf := range workflows.Workflows {
		for _, ws := range wf.Statuses {
			if seen[ws.ID] {
				continue
			}
			seen[ws.ID] = true
			s := statusFromWorkflowStatus(&ws)
			statuses = append(statuses, *s)
		}
	}
	return &provider.ListStatusesResponse{Statuses: statuses}, nil
}
