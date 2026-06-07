# graphql

The `graphql` package provides a typed client for Aha.io's GraphQL API.

## Installation

```go
import "github.com/grokify/aha-go/graphql"
```

## Overview

The GraphQL API provides access to functionality not available in the REST API, including full-text search across documents.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/grokify/aha-go/graphql"
)

func main() {
    // Create client
    client := graphql.NewClient("mycompany", "your-api-key")

    // Search for documents
    var result graphql.SearchDocumentsResponse
    err := client.Query(context.Background(), graphql.SearchDocumentsQuery, map[string]any{
        "query":          "onboarding",
        "searchableType": []string{"Page"},
    }, &result)
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    for _, node := range result.SearchDocuments.Nodes {
        fmt.Printf("%s: %s\n", node.Name, node.URL)
    }
}
```

## Client

### Creating a Client

```go
// Basic client
client := graphql.NewClient("subdomain", "api-key")

// With custom HTTP client
httpClient := &http.Client{Timeout: 30 * time.Second}
client := graphql.NewClientWithHTTP("subdomain", "api-key", httpClient)
```

### Executing Queries

```go
// Using Query method (recommended)
var result SomeResponseType
err := client.Query(ctx, queryString, variables, &result)

// Using Do method (for raw access)
req := &graphql.Request{
    Query:     queryString,
    Variables: variables,
}
resp, err := client.Do(ctx, req)
```

## Predefined Queries

### SearchDocumentsQuery

Search across pages, features, and other document types.

```go
var result graphql.SearchDocumentsResponse
err := client.Query(ctx, graphql.SearchDocumentsQuery, map[string]any{
    "query":          "search term",
    "searchableType": []string{"Page", "Feature"},
}, &result)

// Access results
for _, node := range result.SearchDocuments.Nodes {
    fmt.Printf("Name: %s\n", node.Name)
    fmt.Printf("URL: %s\n", node.URL)
    fmt.Printf("Type: %s\n", node.SearchableType)
    fmt.Printf("ID: %s\n", node.SearchableID)
}

// Pagination info
fmt.Printf("Page %d of %d\n", result.SearchDocuments.CurrentPage, result.SearchDocuments.TotalPages)
fmt.Printf("Total: %d\n", result.SearchDocuments.TotalCount)
```

### GetPageQuery

Retrieve a page by reference number.

```go
var result graphql.PageResponse
err := client.Query(ctx, graphql.GetPageQuery, map[string]any{
    "id":            "PAGE-123",
    "includeParent": true,
}, &result)

if result.Page != nil {
    fmt.Printf("Name: %s\n", result.Page.Name)
    fmt.Printf("Content: %s\n", result.Page.Description.MarkdownBody)
}
```

### GetFeatureQuery

Retrieve a feature by reference number.

```go
var result graphql.FeatureResponse
err := client.Query(ctx, graphql.GetFeatureQuery, map[string]any{
    "id": "FEAT-123",
}, &result)

if result.Feature != nil {
    fmt.Printf("Name: %s\n", result.Feature.Name)
}
```

### GetRequirementQuery

Retrieve a requirement by reference number.

```go
var result graphql.RequirementResponse
err := client.Query(ctx, graphql.GetRequirementQuery, map[string]any{
    "id": "REQ-123",
}, &result)

if result.Requirement != nil {
    fmt.Printf("Name: %s\n", result.Requirement.Name)
}
```

## Response Types

### SearchDocumentsResponse

```go
type SearchDocumentsResponse struct {
    SearchDocuments SearchResults `json:"searchDocuments"`
}

type SearchResults struct {
    Nodes       []DocumentNode `json:"nodes"`
    CurrentPage int            `json:"currentPage"`
    TotalCount  int            `json:"totalCount"`
    TotalPages  int            `json:"totalPages"`
    IsLastPage  bool           `json:"isLastPage"`
}

type DocumentNode struct {
    Name           string `json:"name"`
    URL            string `json:"url"`
    SearchableID   string `json:"searchableId"`
    SearchableType string `json:"searchableType"`
}
```

### PageResponse

```go
type PageResponse struct {
    Page *Page `json:"page"`
}

type Page struct {
    Name        string       `json:"name"`
    Description *Description `json:"description"`
    Children    []PageRef    `json:"children"`
    Parent      *PageRef     `json:"parent"`
}
```

## Error Handling

```go
var result graphql.SearchDocumentsResponse
err := client.Query(ctx, graphql.SearchDocumentsQuery, variables, &result)
if err != nil {
    // Check if it's a GraphQL error
    if gqlErr, ok := err.(graphql.Error); ok {
        fmt.Printf("GraphQL error: %s\n", gqlErr.Message)
        fmt.Printf("Path: %v\n", gqlErr.Path)
    } else {
        // Network or other error
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Testing

The package provides a `SetEndpoint` method for testing with mock servers:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"data":{"searchDocuments":{"nodes":[]}}}`))
}))
defer server.Close()

client := graphql.NewClient("test", "key")
client.SetEndpoint(server.URL)

// Now queries go to mock server
```

## API Reference

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/grokify/aha-go/graphql).
