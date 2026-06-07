# browser Package

The `browser` package provides an HTTP server for browsing Aha.io data in a web interface.

```go
import "github.com/grokify/aha-go/browser"
```

## Overview

The browser package creates a local web server that renders Aha.io entities as HTML pages. Useful for:

- Quick data exploration
- Stakeholder demos
- Offline browsing (with cached data)
- Custom dashboards

## Quick Start

```go
package main

import (
    "log"

    aha "github.com/grokify/aha-go"
    "github.com/grokify/aha-go/browser"
)

func main() {
    client, err := aha.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    server := browser.New(client, browser.Options{
        Port: 8080,
    })

    log.Println("Starting browser at http://localhost:8080")
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Options

```go
type Options struct {
    // Port to listen on (default: 8080)
    Port int

    // Host to bind to (default: "localhost")
    Host string

    // Product to browse (optional, shows selector if empty)
    Product string

    // Title for the browser (default: "Aha Browser")
    Title string

    // Custom CSS file path
    CustomCSS string

    // Enable caching
    EnableCache bool

    // Cache TTL in seconds (default: 300)
    CacheTTL int
}
```

## Starting the Server

### Basic

```go
server := browser.New(client, browser.Options{})
server.Start() // Blocks
```

### Background

```go
server := browser.New(client, browser.Options{
    Port: 8080,
})

go server.Start()

// Do other work...

server.Stop()
```

### Custom Host

```go
server := browser.New(client, browser.Options{
    Host: "0.0.0.0",  // Listen on all interfaces
    Port: 3000,
})
```

## Routes

The browser serves these routes:

| Route | Description |
|-------|-------------|
| `/` | Product selector or product overview |
| `/products` | List all products |
| `/products/{id}` | Product detail |
| `/features` | Feature list |
| `/features/{ref}` | Feature detail |
| `/ideas` | Idea list |
| `/ideas/{ref}` | Idea detail |
| `/releases` | Release list |
| `/releases/{ref}` | Release detail |
| `/initiatives` | Initiative list |
| `/initiatives/{ref}` | Initiative detail |

## Customization

### Custom Templates

```go
server := browser.New(client, browser.Options{})

// Override a template
server.SetTemplate("feature", `
<!DOCTYPE html>
<html>
<head><title>{{.Name}}</title></head>
<body>
    <h1>{{.ReferenceNum}}: {{.Name}}</h1>
    <p>Status: {{.Status}}</p>
    <div>{{.Description}}</div>
</body>
</html>
`)
```

### Custom CSS

```go
server := browser.New(client, browser.Options{
    CustomCSS: "/path/to/custom.css",
})
```

Or inline:

```go
server.SetCSS(`
    body { font-family: system-ui; }
    .feature { border: 1px solid #ccc; padding: 1rem; }
    .status-done { color: green; }
    .status-progress { color: blue; }
`)
```

### Adding Routes

```go
server := browser.New(client, browser.Options{})

// Add custom handler
server.HandleFunc("/custom", func(w http.ResponseWriter, r *http.Request) {
    features, _ := client.ListFeatures(r.Context(), "PRODUCT")
    // Custom rendering
})
```

## Caching

Enable caching for better performance:

```go
server := browser.New(client, browser.Options{
    EnableCache: true,
    CacheTTL:    300, // 5 minutes
})
```

### Manual Cache Control

```go
// Clear all cache
server.ClearCache()

// Clear specific entity
server.ClearCacheFor("features")
```

## Templates

### Default Template Variables

All templates receive these variables:

| Variable | Description |
|----------|-------------|
| `.Title` | Page title |
| `.Product` | Current product |
| `.Data` | Entity data |
| `.Error` | Error message if any |

### Feature Template

```html
{{define "feature"}}
<div class="feature">
    <h1>{{.Data.ReferenceNum}}: {{.Data.Name}}</h1>
    <span class="status status-{{.Data.Status | lower}}">
        {{.Data.Status}}
    </span>
    <div class="description">
        {{.Data.Description | safeHTML}}
    </div>
    {{if .Data.Requirements}}
    <h2>Requirements</h2>
    <ul>
        {{range .Data.Requirements}}
        <li>{{.Name}}</li>
        {{end}}
    </ul>
    {{end}}
</div>
{{end}}
```

## Middleware

```go
server := browser.New(client, browser.Options{})

// Add authentication middleware
server.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check auth
        if !isAuthorized(r) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
})
```

## API Reference

See [pkg.go.dev/github.com/grokify/aha-go/browser](https://pkg.go.dev/github.com/grokify/aha-go/browser)
