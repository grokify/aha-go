# canvas Package

The `canvas` package provides tools for creating, managing, and exporting strategic canvases in Aha.io.

```go
import "github.com/grokify/aha-go/canvas"
```

## Overview

Strategic canvases are visual frameworks for business planning. This package supports:

- Business Model Canvas
- Lean Canvas
- Lean UX Canvas
- Opportunity Canvas
- Custom canvas templates

## Canvas Types

### BusinessModelCanvas

The classic 9-block business model framework.

```go
type BusinessModelCanvas struct {
    Name                  string   `json:"name"`
    CustomerSegments      []string `json:"customer_segments"`
    ValuePropositions     []string `json:"value_propositions"`
    Channels              []string `json:"channels"`
    CustomerRelationships []string `json:"customer_relationships"`
    RevenueStreams        []string `json:"revenue_streams"`
    KeyResources          []string `json:"key_resources"`
    KeyActivities         []string `json:"key_activities"`
    KeyPartners           []string `json:"key_partners"`
    CostStructure         []string `json:"cost_structure"`
}
```

### LeanCanvas

Startup-focused adaptation emphasizing problem-solution fit.

```go
type LeanCanvas struct {
    Name               string   `json:"name"`
    Problem            []string `json:"problem"`
    Solution           []string `json:"solution"`
    UniqueValue        string   `json:"unique_value_proposition"`
    UnfairAdvantage    string   `json:"unfair_advantage"`
    CustomerSegments   []string `json:"customer_segments"`
    KeyMetrics         []string `json:"key_metrics"`
    Channels           []string `json:"channels"`
    CostStructure      []string `json:"cost_structure"`
    RevenueStreams     []string `json:"revenue_streams"`
}
```

### LeanUXCanvas

UX hypothesis and experiment planning.

```go
type LeanUXCanvas struct {
    Name             string   `json:"name"`
    BusinessProblem  string   `json:"business_problem"`
    BusinessOutcomes []string `json:"business_outcomes"`
    Users            []string `json:"users"`
    UserOutcomes     []string `json:"user_outcomes"`
    Solutions        []string `json:"solutions"`
    Hypotheses       []string `json:"hypotheses"`
    Experiments      []string `json:"experiments"`
}
```

## Creating Canvases

### Business Model Canvas

```go
bmc := canvas.BusinessModelCanvas{
    Name: "Platform Business Model",
    CustomerSegments: []string{
        "Enterprise product teams",
        "Startup founders",
    },
    ValuePropositions: []string{
        "Unified product management",
    },
    // ... other fields
}

created, err := canvas.CreateBusinessModelCanvas(ctx, client, "PRODUCT-KEY", bmc)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created canvas: %s\n", created.ID)
```

### From JSON

```go
data, _ := os.ReadFile("templates/bmc.json")

var bmc canvas.BusinessModelCanvas
json.Unmarshal(data, &bmc)

created, err := canvas.CreateBusinessModelCanvas(ctx, client, "PRODUCT-KEY", bmc)
```

## Listing Canvases

```go
canvases, err := canvas.List(ctx, client, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

for _, c := range canvases {
    fmt.Printf("%s: %s (%s)\n", c.ID, c.Name, c.Type)
}
```

## Getting a Canvas

```go
c, err := canvas.Get(ctx, client, canvasID)
if err != nil {
    log.Fatal(err)
}

// Access canvas data
fmt.Printf("Name: %s\n", c.Name)
fmt.Printf("Type: %s\n", c.Type)
```

## Exporting Canvases

### Export Options

```go
type ExportOptions struct {
    Format ExportFormat // PNG, PDF, SVG
    Output string       // Output file path
    Width  int          // Optional width
    Height int          // Optional height
}
```

### Export to PNG

```go
err := canvas.Export(ctx, client, canvasID, canvas.ExportOptions{
    Format: canvas.FormatPNG,
    Output: "canvas.png",
    Width:  1920,
})
```

### Export to PDF

```go
err := canvas.Export(ctx, client, canvasID, canvas.ExportOptions{
    Format: canvas.FormatPDF,
    Output: "canvas.pdf",
})
```

### Export to Bytes

```go
data, err := canvas.ExportBytes(ctx, client, canvasID, canvas.FormatPNG)
if err != nil {
    log.Fatal(err)
}
// Use data directly
```

## Updating Canvases

```go
err := canvas.Update(ctx, client, canvasID, canvas.UpdateRequest{
    Name: "Updated Canvas Name",
    Data: map[string]any{
        "customer_segments": []string{"New segment"},
    },
})
```

## Templates

### Predefined Templates

```go
// Get empty Business Model Canvas template
template := canvas.BusinessModelCanvasTemplate()

// Get Lean Canvas template
template := canvas.LeanCanvasTemplate()
```

### Custom Templates

```go
// Save canvas as template
func saveAsTemplate(c *canvas.Canvas, filename string) error {
    data, err := json.MarshalIndent(c.Data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

// Load template
func loadTemplate(filename string) (map[string]any, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var template map[string]any
    err = json.Unmarshal(data, &template)
    return template, err
}
```

## API Reference

See [pkg.go.dev/github.com/grokify/aha-go/canvas](https://pkg.go.dev/github.com/grokify/aha-go/canvas)
