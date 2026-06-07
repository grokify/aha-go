# Quick Start

This guide walks you through common operations with `aha-go`.

## Prerequisites

- [Installed](installation.md) the library or CLI
- [Configured](authentication.md) your API credentials

## Using the CLI

### List Products

```bash
aha product list
```

### List Features in a Product

```bash
aha feature list --product PRODUCT-KEY
```

### Get Feature Details

```bash
aha feature get FEAT-123
```

### Create a Feature

```bash
aha feature create --product PRODUCT-KEY \
  --name "New Feature" \
  --description "Feature description"
```

### Output Formats

```bash
# JSON output
aha feature list --product PRODUCT-KEY --output json

# Table output (default)
aha feature list --product PRODUCT-KEY --output table
```

## Using the Library

### Basic Client Setup

```go
package main

import (
    "context"
    "fmt"
    "log"

    aha "github.com/grokify/aha-go"
)

func main() {
    ctx := context.Background()

    // Create client from environment variables
    client, err := aha.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    // List products
    products, err := client.ListProducts(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, p := range products {
        fmt.Printf("%s: %s\n", p.ReferencePrefix, p.Name)
    }
}
```

### Working with Features

```go
// List features in a product
features, err := client.ListFeatures(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

// Get a specific feature
feature, err := client.GetFeature(ctx, "FEAT-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Feature: %s - %s\n", feature.ReferenceNum, feature.Name)

// Create a feature
newFeature, err := client.CreateFeature(ctx, "PRODUCT-KEY", aha.CreateFeatureRequest{
    Name:        "New Feature",
    Description: "Description here",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created: %s\n", newFeature.ReferenceNum)

// Update a feature
_, err = client.UpdateFeature(ctx, "FEAT-123", aha.UpdateFeatureRequest{
    Name: "Updated Feature Name",
})
```

### Working with Ideas

```go
// List ideas
ideas, err := client.ListIdeas(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

// Get top voted ideas
for _, idea := range ideas {
    if idea.Votes >= 10 {
        fmt.Printf("%s: %s (%d votes)\n", idea.ReferenceNum, idea.Name, idea.Votes)
    }
}
```

### Working with Releases

```go
// List releases
releases, err := client.ListReleases(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

// Find active releases
for _, r := range releases {
    if !r.Released {
        fmt.Printf("%s: %s (due: %s)\n", r.ReferenceNum, r.Name, r.ReleaseDate)
    }
}
```

### Adding Comments

```go
// Add a comment to a feature
comment, err := client.CreateComment(ctx, aha.CreateCommentRequest{
    CommentableType: "Feature",
    CommentableID:   "FEAT-123",
    Body:            "This is a comment from the API",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Comment added: %s\n", comment.ID)
```

## Error Handling

```go
features, err := client.ListFeatures(ctx, "PRODUCT-KEY")
if err != nil {
    // Check for specific error types
    if aha.IsNotFoundError(err) {
        fmt.Println("Product not found")
    } else if aha.IsUnauthorizedError(err) {
        fmt.Println("Invalid API key")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Pagination

For large result sets, use pagination options:

```go
// Get first page
features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithPage(1),
    aha.WithPerPage(50),
)

// Iterate through all pages
page := 1
for {
    features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
        aha.WithPage(page),
        aha.WithPerPage(100),
    )
    if err != nil {
        log.Fatal(err)
    }
    if len(features) == 0 {
        break
    }

    for _, f := range features {
        // Process feature
    }
    page++
}
```

## Next Steps

- [CLI Command Reference](../cli/commands.md)
- [Use Case Guides](../use-cases/index.md)
- [Package Documentation](../packages/index.md)
