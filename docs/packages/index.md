# Packages

`aha-go` is organized into focused packages for specific functionality.

## Package Overview

| Package | Import | Description |
|---------|--------|-------------|
| `aha` | `github.com/grokify/aha-go` | Core client and entity types |
| `canvas` | `github.com/grokify/aha-go/canvas` | Strategic canvas operations |
| `browser` | `github.com/grokify/aha-go/browser` | HTML entity browser |
| `render` | `github.com/grokify/aha-go/render` | Diagram generation |

## Core Package

The root package provides the main client and all entity types.

```go
import aha "github.com/grokify/aha-go"

client, err := aha.NewClient()

// Entity operations
features, _ := client.ListFeatures(ctx, "PRODUCT")
ideas, _ := client.ListIdeas(ctx, "PRODUCT")
releases, _ := client.ListReleases(ctx, "PRODUCT")
```

## Sub-packages

### canvas

Create and export strategic canvases.

```go
import "github.com/grokify/aha-go/canvas"

// Create Business Model Canvas
bmc := canvas.BusinessModelCanvas{
    Name: "Product BMC",
    // ...
}
```

[Full canvas documentation](canvas.md)

### browser

Serve Aha data as a browsable HTML interface.

```go
import "github.com/grokify/aha-go/browser"

server := browser.New(client, browser.Options{
    Port: 8080,
})
server.Start()
```

[Full browser documentation](browser.md)

### render

Generate diagrams from Aha data.

```go
import "github.com/grokify/aha-go/render"

mermaid := render.ReleasesToMermaid(releases, render.MermaidOptions{})
```

[Full render documentation](render.md)

## API Reference

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/grokify/aha-go).
