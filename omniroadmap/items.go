package omniroadmap

import (
	"context"
	"slices"

	aha "github.com/grokify/aha-go"
	omniroadmap "github.com/grokify/omniroadmap-core"
	"github.com/grokify/omniroadmap-core/provider"
)

// ListItems lists items across the kinds requested (or all supported kinds
// if req.Kinds is empty). Aha has no single cross-kind list endpoint, so
// this issues one call per requested kind and concatenates the results.
// Pagination (req.Page/PerPage) is applied identically to every underlying
// call — TotalCount in the response is a sum across kinds, not a single
// authoritative total from one endpoint.
//
// Results only carry meta-level fields (ID, name, reference, URL,
// created-at) since Aha's list endpoints return lightweight records —
// Description/Status/Progress/etc. require GetItem.
func (p *Provider) ListItems(ctx context.Context, req *provider.ListItemsRequest) (*provider.ListItemsResponse, error) {
	kinds := req.Kinds
	if len(kinds) == 0 {
		kinds = p.Capabilities().Kinds
	}

	resp := &provider.ListItemsResponse{Page: req.Page, PerPage: req.PerPage}
	for _, kind := range kinds {
		items, total, err := p.listItemsByKind(ctx, kind, req)
		if err != nil {
			return nil, err
		}
		resp.Items = append(resp.Items, items...)
		resp.TotalCount += total
	}
	return resp, nil
}

func (p *Provider) listItemsByKind(ctx context.Context, kind provider.ItemKind, req *provider.ListItemsRequest) ([]provider.Item, int, error) {
	switch kind {
	case provider.ItemKindFeature:
		opts := []aha.ListFeaturesOption{}
		if req.Page > 0 {
			opts = append(opts, aha.WithFeaturePage(req.Page))
		}
		if req.PerPage > 0 {
			opts = append(opts, aha.WithFeaturePerPage(req.PerPage))
		}
		list, err := p.client.ListFeatures(ctx, opts...)
		if err != nil {
			return nil, 0, wrapErr("ListFeatures", err)
		}
		items := make([]provider.Item, len(list.Features))
		for i, m := range list.Features {
			items[i] = itemFromFeatureMeta(m)
		}
		return items, int(list.Pagination.TotalRecords), nil

	case provider.ItemKindInitiative:
		opts := []aha.ListInitiativesOption{}
		if req.Page > 0 {
			opts = append(opts, aha.WithInitiativePage(req.Page))
		}
		if req.PerPage > 0 {
			opts = append(opts, aha.WithInitiativePerPage(req.PerPage))
		}
		list, err := p.client.ListInitiatives(ctx, opts...)
		if err != nil {
			return nil, 0, wrapErr("ListInitiatives", err)
		}
		items := make([]provider.Item, len(list.Initiatives))
		for i, m := range list.Initiatives {
			items[i] = itemFromInitiativeMeta(m)
		}
		return items, int(list.Pagination.TotalRecords), nil

	case provider.ItemKindEpic:
		opts := []aha.ListEpicsOption{}
		if req.Page > 0 {
			opts = append(opts, aha.WithEpicPage(req.Page))
		}
		if req.PerPage > 0 {
			opts = append(opts, aha.WithEpicPerPage(req.PerPage))
		}
		list, err := p.client.ListEpics(ctx, opts...)
		if err != nil {
			return nil, 0, wrapErr("ListEpics", err)
		}
		items := make([]provider.Item, len(list.Epics))
		for i, m := range list.Epics {
			items[i] = itemFromEpicMeta(m)
		}
		return items, int(list.Pagination.TotalRecords), nil

	default:
		if !slices.Contains(p.Capabilities().Kinds, kind) {
			return nil, 0, omniroadmap.ErrUnsupportedOperation
		}
		return nil, 0, nil
	}
}

// GetItem fetches a single, fully-populated item by kind and source ID.
func (p *Provider) GetItem(ctx context.Context, req *provider.GetItemRequest) (*provider.Item, error) {
	switch req.Kind {
	case provider.ItemKindFeature:
		f, err := p.client.GetFeature(ctx, req.ID)
		if err != nil {
			return nil, wrapErr("GetFeature", err)
		}
		return itemFromFeature(f), nil

	case provider.ItemKindInitiative:
		i, err := p.client.GetInitiative(ctx, req.ID)
		if err != nil {
			return nil, wrapErr("GetInitiative", err)
		}
		return itemFromInitiative(i), nil

	case provider.ItemKindEpic:
		e, err := p.client.GetEpic(ctx, req.ID)
		if err != nil {
			return nil, wrapErr("GetEpic", err)
		}
		return itemFromEpic(e), nil

	default:
		return nil, omniroadmap.ErrUnsupportedOperation
	}
}

// wrapErr wraps an aha-go client error, preserving the underlying error for
// errors.Is/As while identifying which call failed.
func wrapErr(op string, err error) error {
	return &opError{op: op, err: err}
}

type opError struct {
	op  string
	err error
}

func (e *opError) Error() string { return "omniroadmap/aha: " + e.op + ": " + e.err.Error() }
func (e *opError) Unwrap() error { return e.err }
