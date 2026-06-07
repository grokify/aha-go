# Strategic Canvases

Aha.io supports strategic canvases for business modeling. The `canvas` package provides tools to create, manage, and export canvases programmatically.

## Canvas Types

| Canvas | Description |
|--------|-------------|
| Business Model Canvas | 9-block business model framework |
| Lean Canvas | Startup-focused adaptation of BMC |
| Lean UX Canvas | UX hypothesis and experiment planning |
| Opportunity Canvas | Opportunity assessment framework |

## Listing Canvases

### With Code

```go
import "github.com/grokify/aha-go/canvas"

canvases, err := client.ListCanvases(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

for _, c := range canvases {
    fmt.Printf("%s: %s (%s)\n", c.ID, c.Name, c.Type)
}
```

### With CLI

```bash
aha canvas list --product PRODUCT-KEY
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
        "Strategic roadmapping",
    },
    Channels: []string{
        "Direct sales",
        "Partner network",
    },
    CustomerRelationships: []string{
        "Dedicated support",
        "Self-service portal",
    },
    RevenueStreams: []string{
        "SaaS subscriptions",
        "Professional services",
    },
    KeyResources: []string{
        "Engineering team",
        "Customer success team",
    },
    KeyActivities: []string{
        "Product development",
        "Customer onboarding",
    },
    KeyPartners: []string{
        "Integration partners",
        "Resellers",
    },
    CostStructure: []string{
        "Engineering salaries",
        "Cloud infrastructure",
    },
}

created, err := client.CreateBusinessModelCanvas(ctx, "PRODUCT-KEY", bmc)
```

### From JSON File

```bash
# Create from JSON template
aha canvas create bmc --product PRODUCT-KEY --name "My BMC" --file bmc.json
```

Example `bmc.json`:

```json
{
  "customer_segments": [
    "Enterprise product teams",
    "Startup founders"
  ],
  "value_propositions": [
    "Unified product management"
  ],
  "channels": ["Direct sales"],
  "customer_relationships": ["Dedicated support"],
  "revenue_streams": ["SaaS subscriptions"],
  "key_resources": ["Engineering team"],
  "key_activities": ["Product development"],
  "key_partners": ["Integration partners"],
  "cost_structure": ["Engineering salaries"]
}
```

### Lean UX Canvas

```go
leanux := canvas.LeanUXCanvas{
    Name: "Feature Hypothesis",
    BusinessProblem: "Users struggle to find relevant features",
    BusinessOutcomes: []string{
        "Increase feature discovery by 25%",
        "Reduce support tickets about features",
    },
    Users: []string{
        "New users (first 30 days)",
        "Product managers",
    },
    UserOutcomes: []string{
        "Quickly find features they need",
        "Understand feature capabilities",
    },
    Solutions: []string{
        "Contextual feature suggestions",
        "Interactive feature tour",
    },
    Hypotheses: []string{
        "If we add contextual suggestions, users will discover 25% more features",
    },
    Experiments: []string{
        "A/B test suggestions in dashboard",
        "User interviews on feature discovery",
    },
}

created, err := client.CreateLeanUXCanvas(ctx, "PRODUCT-KEY", leanux)
```

### Opportunity Canvas

```go
opp := canvas.OpportunityCanvas{
    Name: "Mobile App Opportunity",
    Problem: "Users can't access product data on mobile",
    Solution: "Native mobile app for iOS and Android",
    KeyMetrics: []string{
        "Mobile DAU",
        "Mobile session duration",
    },
    UniqueValueProposition: "Product management on the go",
    UnfairAdvantage: "Deep integration with existing platform",
    Channels: []string{
        "App stores",
        "In-app promotion",
    },
    CustomerSegments: []string{
        "Mobile-first product managers",
        "Executives needing quick updates",
    },
    CostStructure: []string{
        "Mobile development team",
        "App store fees",
    },
    RevenueStreams: []string{
        "Premium mobile features",
        "Increased retention",
    },
}

created, err := client.CreateOpportunityCanvas(ctx, "PRODUCT-KEY", opp)
```

## Exporting Canvases

### To PNG

```go
import "github.com/grokify/aha-go/canvas"

err := canvas.Export(ctx, client, canvasID, canvas.ExportOptions{
    Format: canvas.FormatPNG,
    Output: "canvas.png",
})
```

### CLI Export

```bash
# Export to PNG
aha canvas export CANVAS-ID --format png --output canvas.png

# Export to PDF
aha canvas export CANVAS-ID --format pdf --output canvas.pdf

# Export to SVG
aha canvas export CANVAS-ID --format svg --output canvas.svg
```

## Updating Canvases

```go
// Update specific fields
err := client.UpdateCanvas(ctx, canvasID, canvas.UpdateRequest{
    Name: "Updated Canvas Name",
    Data: map[string]any{
        "customer_segments": []string{
            "New segment 1",
            "New segment 2",
        },
    },
})
```

## Canvas Templates

### Loading from JSON

```go
// Load canvas data from JSON file
data, err := os.ReadFile("templates/bmc-template.json")
if err != nil {
    log.Fatal(err)
}

var bmc canvas.BusinessModelCanvas
if err := json.Unmarshal(data, &bmc); err != nil {
    log.Fatal(err)
}

bmc.Name = "New Product BMC"
created, err := client.CreateBusinessModelCanvas(ctx, "PRODUCT-KEY", bmc)
```

### Saving as Template

```go
// Fetch existing canvas
existing, err := client.GetCanvas(ctx, canvasID)
if err != nil {
    log.Fatal(err)
}

// Save as template
template := existing.Data
template["name"] = "" // Clear name for reuse

data, _ := json.MarshalIndent(template, "", "  ")
os.WriteFile("templates/my-template.json", data, 0644)
```

## Best Practices

### Version Canvases

Track canvas changes over time:

```go
func snapshotCanvas(ctx context.Context, client *aha.Client, canvasID string) {
    canvas, _ := client.GetCanvas(ctx, canvasID)

    timestamp := time.Now().Format("2006-01-02-150405")
    filename := fmt.Sprintf("canvas-%s-%s.json", canvasID, timestamp)

    data, _ := json.MarshalIndent(canvas, "", "  ")
    os.WriteFile(filename, data, 0644)
}
```

### Batch Canvas Creation

```go
// Create multiple canvases from templates
templates := []string{"bmc.json", "leanux.json", "opportunity.json"}

for _, tmpl := range templates {
    data, _ := os.ReadFile(tmpl)
    var canvas map[string]any
    json.Unmarshal(data, &canvas)

    _, err := client.CreateCanvas(ctx, "PRODUCT-KEY", canvas)
    if err != nil {
        log.Printf("Failed to create %s: %v", tmpl, err)
    }
}
```
