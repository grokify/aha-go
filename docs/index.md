# aha-go

Go SDK and CLI for the [Aha.io](https://www.aha.io/) API.

## Overview

`aha-go` provides a complete Go client library and command-line interface for interacting with Aha.io's product management platform. Whether you're automating workflows, building integrations, or querying product data, this library gives you programmatic access to all major Aha.io entities.

## Features

- **Full API Coverage** - Access features, ideas, epics, releases, initiatives, goals, requirements, comments, users, and products
- **CLI Tool** - Command-line interface for quick interactions and scripting
- **Strategic Canvases** - Create and export Business Model Canvas, Lean Canvas, and custom templates
- **Diagram Generation** - Render roadmaps and relationships as Mermaid, D2, or SVG diagrams
- **Type Safety** - Fully typed Go structs for all API responses

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    aha "github.com/grokify/aha-go"
)

func main() {
    // Create client (uses AHA_SUBDOMAIN and AHA_API_KEY env vars)
    client, err := aha.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    // List features in a product
    features, err := client.ListFeatures(context.Background(), "PRODUCT-KEY")
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range features {
        fmt.Printf("%s: %s (%s)\n", f.ReferenceNum, f.Name, f.Status)
    }
}
```

## Documentation

| Section | Description |
|---------|-------------|
| [Getting Started](getting-started/installation.md) | Installation, authentication, and first steps |
| [CLI Reference](cli/index.md) | Command-line tool usage and examples |
| [Use Cases](use-cases/index.md) | Common workflows and integration patterns |
| [Packages](packages/index.md) | Package documentation for canvas, browser, render |
| [API Reference](api-reference.md) | Links to pkg.go.dev documentation |

## Installation

```bash
# As a library
go get github.com/grokify/aha-go

# CLI tool
go install github.com/grokify/aha-go/cmd/aha@latest
```

## Links

- [GitHub Repository](https://github.com/grokify/aha-go)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/grokify/aha-go)
- [Aha.io API Documentation](https://www.aha.io/api)
