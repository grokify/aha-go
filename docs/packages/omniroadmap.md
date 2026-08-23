# omniroadmap Package

The `omniroadmap` package adapts `*aha.Client` to
[`omniroadmap-core`](https://github.com/grokify/omniroadmap-core)'s
`provider.Provider` interface, so tools that aggregate roadmap data across
multiple product-management systems (e.g. `plexusone/dashforge`) can talk
to Aha through the same contract they use for every other provider.

```go
import "github.com/grokify/aha-go/omniroadmap"
```

Since Aha publishes no official Go SDK, aha-go is itself a from-scratch
client — so the adapter lives inside this repo rather than a separate
`omni-aha` repo, following the pattern used by `elevenlabs-go/omnivoice`
and `opik-go/integrations/omnillm`.

## Creating a Provider

```go
import (
    aha "github.com/grokify/aha-go"
    "github.com/grokify/aha-go/omniroadmap"
)

client, err := aha.NewClient()
if err != nil {
    log.Fatal(err)
}

// WithProductID scopes product-level-only operations (ListStatuses,
// product-scoped releases) to a single Aha product. Required for
// ListStatuses/ListReleases; without it, those return
// omniroadmap.ErrUnsupportedOperation.
p := omniroadmap.NewProvider(client, omniroadmap.WithProductID("PROD"))
```

The provider self-registers under the name `"aha"` via `init()`, so
`omniroadmap.RegisterProvider`-based factories can construct it from an
`*aha.Client` config without an explicit import-time call.

## Capabilities

```go
p.Capabilities()
// provider.Capabilities{
//     Kinds:                []provider.ItemKind{Feature, Epic, Initiative},
//     SupportsReleases:     true,
//     SupportsObjectives:   false,
//     SupportsCustomFields: true,
//     SupportsWrite:        false,
// }
```

## Listing and Fetching Items

```go
import "github.com/grokify/omniroadmap-core/provider"

// List across all supported kinds (or a subset via req.Kinds)
resp, err := p.ListItems(ctx, &provider.ListItemsRequest{
    Page:    1,
    PerPage: 50,
})

// List results only carry meta-level fields (ID, name, reference, URL,
// created-at) - fetch a single item for full details:
item, err := p.GetItem(ctx, &provider.GetItemRequest{
    Kind: provider.ItemKindFeature,
    ID:   "PROD-123",
})
```

Aha has no single cross-kind list endpoint, so `ListItems` issues one call
per requested kind and concatenates the results; `TotalCount` in the
response is a sum across kinds, not a single authoritative total.

## Releases, Statuses, and Custom Fields

```go
// Releases and workflow statuses require WithProductID (Aha's
// release/workflow endpoints are product-scoped, with no unscoped
// "all products" variant)
releases, err := p.ListReleases(ctx, &provider.ListReleasesRequest{})
statuses, err := p.ListStatuses(ctx, &provider.ListStatusesRequest{})

// Custom field definitions have no product-scoping requirement
defs, err := p.ListCustomFieldDefinitions(ctx, &provider.ListCustomFieldDefinitionsRequest{})
```

## Error Handling

Every underlying `aha.Client` error is wrapped to identify which call
failed, while remaining unwrappable via `errors.Is`/`errors.As`:

```go
_, err := p.GetItem(ctx, &provider.GetItemRequest{Kind: provider.ItemKindFeature, ID: "BAD-ID"})
if aha.IsNotFound(err) {
    // still detectable through the wrapper
}
```

## API Reference

See [pkg.go.dev/github.com/grokify/aha-go/omniroadmap](https://pkg.go.dev/github.com/grokify/aha-go/omniroadmap)
