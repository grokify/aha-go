package omniroadmap

import (
	"context"

	aha "github.com/grokify/aha-go"
	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListReleases lists releases across all products the client can access.
// Unlike ListItems, aha-go's ListProductReleases returns fully-populated
// Release objects (no meta-vs-full distinction for releases), but it's
// product-scoped — without a configured productID (see WithProductID) this
// returns omniroadmap.ErrUnsupportedOperation.
func (p *Provider) ListReleases(ctx context.Context, req *provider.ListReleasesRequest) (*provider.ListReleasesResponse, error) {
	if p.productID == "" {
		return nil, omniroadmap.ErrUnsupportedOperation
	}

	opts := []aha.ListOption{}
	if req.Page > 0 {
		opts = append(opts, aha.WithPage(req.Page))
	}
	if req.PerPage > 0 {
		opts = append(opts, aha.WithPerPage(req.PerPage))
	}

	list, err := p.client.ListProductReleases(ctx, p.productID, opts...)
	if err != nil {
		return nil, wrapErr("ListProductReleases", err)
	}

	releases := make([]provider.Release, len(list.Releases))
	for i := range list.Releases {
		releases[i] = releaseFromAha(&list.Releases[i])
	}
	return &provider.ListReleasesResponse{
		Releases:   releases,
		Page:       req.Page,
		PerPage:    req.PerPage,
		TotalCount: int(list.Pagination.TotalRecords),
	}, nil
}
