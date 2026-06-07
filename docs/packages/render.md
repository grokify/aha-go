# render Package

The `render` package generates diagrams from Aha.io data in Mermaid, D2, and SVG formats.

```go
import "github.com/grokify/aha-go/render"
```

## Overview

Transform Aha.io entities into visual diagrams:

- **Mermaid** - Text-based diagrams for Markdown/docs
- **D2** - Declarative diagrams with high-quality output
- **SVG** - Vector graphics for web/presentations

## Mermaid Diagrams

### Roadmap Timeline

```go
releases, _ := client.ListReleases(ctx, "PRODUCT-KEY")

mermaid := render.ReleasesToMermaid(releases, render.MermaidOptions{
    Title:     "Product Roadmap",
    ShowDates: true,
})

fmt.Println(mermaid)
```

Output:

```
gantt
    title Product Roadmap
    dateFormat YYYY-MM-DD

    section Releases
    Q1 2024    :2024-01-01, 2024-03-31
    Q2 2024    :2024-04-01, 2024-06-30
```

### Flowchart

```go
features, _ := client.ListFeatures(ctx, "PRODUCT-KEY")

mermaid := render.FeaturesToFlowchart(features, render.MermaidOptions{
    Direction: "TB", // Top to Bottom
})
```

### Pie Chart

```go
mermaid := render.StatusDistributionToPie(features, render.PieOptions{
    Title: "Feature Status",
})
```

Output:

```
pie title Feature Status
    "Done" : 45
    "In Progress" : 30
    "Ready" : 15
    "Backlog" : 10
```

### MermaidOptions

```go
type MermaidOptions struct {
    Title        string // Diagram title
    Direction    string // TB, BT, LR, RL
    Theme        string // default, dark, forest, neutral
    ShowDates    bool   // Include dates
    ShowFeatures bool   // Include features in releases
    DateFormat   string // Date format (YYYY-MM-DD)
}
```

## D2 Diagrams

### Feature Dependencies

```go
features, _ := client.ListFeatures(ctx, "PRODUCT-KEY")

d2 := render.FeaturesToD2(features, render.D2Options{
    Title:     "Feature Dependencies",
    Direction: "right",
})
```

Output:

```
direction: right
title: Feature Dependencies

FEAT-1: User Auth {
    shape: rectangle
}
FEAT-2: Dashboard {
    shape: rectangle
}

FEAT-2 -> FEAT-1: depends on
```

### Initiative Hierarchy

```go
initiatives, _ := client.ListInitiatives(ctx, "PRODUCT-KEY")

d2 := render.InitiativesToD2(initiatives, render.D2Options{
    ShowGoals:    true,
    ShowFeatures: true,
})
```

### D2Options

```go
type D2Options struct {
    Title        string // Diagram title
    Direction    string // up, down, left, right
    Theme        string // Theme name
    Sketch       bool   // Hand-drawn style
    ShowGoals    bool   // Include goals
    ShowFeatures bool   // Include features
}
```

## SVG Output

### Render D2 to SVG

```go
d2Content := render.FeaturesToD2(features, render.D2Options{})

svg, err := render.D2ToSVG(d2Content)
if err != nil {
    log.Fatal(err)
}

os.WriteFile("diagram.svg", svg, 0644)
```

### Render Options

```go
svg, err := render.D2ToSVG(d2Content, render.SVGOptions{
    Width:  1200,
    Height: 800,
    Pad:    20,
})
```

### SVGOptions

```go
type SVGOptions struct {
    Width  int    // Output width
    Height int    // Output height
    Pad    int    // Padding
    Theme  string // Color theme
}
```

## Format Helpers

### Export Format

```go
type Format string

const (
    FormatMermaid Format = "mermaid"
    FormatD2      Format = "d2"
    FormatSVG     Format = "svg"
    FormatPNG     Format = "png"
)
```

### Auto-detect Format

```go
format := render.DetectFormat("output.svg") // Returns FormatSVG
format := render.DetectFormat("diagram.mmd") // Returns FormatMermaid
```

## Rendering Functions

### Releases

```go
// Mermaid Gantt chart
mermaid := render.ReleasesToMermaid(releases, opts)

// D2 timeline
d2 := render.ReleasesToD2(releases, opts)
```

### Features

```go
// Dependency graph
mermaid := render.FeatureDependenciesToMermaid(features, opts)
d2 := render.FeatureDependenciesToD2(features, opts)

// Status flowchart
mermaid := render.FeaturesToFlowchart(features, opts)
```

### Initiatives

```go
// Hierarchy diagram
mermaid := render.InitiativeHierarchyToMermaid(initiatives)
d2 := render.InitiativesToD2(initiatives, opts)
```

### Ideas

```go
// Vote distribution
mermaid := render.IdeasToVotePie(ideas, opts)

// Category breakdown
mermaid := render.IdeasByCategoryToBar(ideas, opts)
```

## Integration Examples

### Embed in Markdown

```go
func generateReadme(ctx context.Context, client *aha.Client) string {
    releases, _ := client.ListReleases(ctx, "PRODUCT")

    mermaid := render.ReleasesToMermaid(releases, render.MermaidOptions{})

    return fmt.Sprintf(`# Roadmap

%s%s%s
`, "```mermaid\n", mermaid, "\n```")
}
```

### HTML with SVG

```go
func generateHTML(features []aha.Feature) string {
    d2 := render.FeaturesToD2(features, render.D2Options{})
    svg, _ := render.D2ToSVG(d2)

    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body>%s</body>
</html>
`, string(svg))
}
```

### Write to File

```go
func saveDiagram(content string, filename string, format render.Format) error {
    switch format {
    case render.FormatSVG:
        svg, err := render.D2ToSVG(content)
        if err != nil {
            return err
        }
        return os.WriteFile(filename, svg, 0644)

    case render.FormatMermaid, render.FormatD2:
        return os.WriteFile(filename, []byte(content), 0644)

    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}
```

## API Reference

See [pkg.go.dev/github.com/grokify/aha-go/render](https://pkg.go.dev/github.com/grokify/aha-go/render)
