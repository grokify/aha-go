# Querying Features

Features are the core work items in Aha.io. This guide covers common patterns for querying, filtering, and analyzing features.

## Listing Features

### Basic List

```go
features, err := client.ListFeatures(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

for _, f := range features {
    fmt.Printf("%s: %s\n", f.ReferenceNum, f.Name)
}
```

### With CLI

```bash
aha feature list --product PRODUCT-KEY
```

## Filtering Features

### By Status

```go
// Get in-progress features
features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithStatus("In Progress"),
)
```

```bash
aha feature list --product PRODUCT-KEY --status "In Progress"
```

### By Release

```go
// Features in a specific release
features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithRelease("REL-2024Q1"),
)
```

### By Assignee

```go
// Features assigned to a user
features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithAssignedTo("user@company.com"),
)
```

### By Date Range

```go
// Features updated in the last 7 days
weekAgo := time.Now().AddDate(0, 0, -7)
features, err := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithUpdatedSince(weekAgo),
)
```

## Analyzing Features

### Group by Status

```go
func groupByStatus(features []aha.Feature) map[string][]aha.Feature {
    result := make(map[string][]aha.Feature)
    for _, f := range features {
        result[f.Status] = append(result[f.Status], f)
    }
    return result
}

features, _ := client.ListFeatures(ctx, "PRODUCT-KEY")
byStatus := groupByStatus(features)

for status, feats := range byStatus {
    fmt.Printf("%s: %d features\n", status, len(feats))
}
```

### Calculate Completion Rate

```go
func completionRate(features []aha.Feature) float64 {
    if len(features) == 0 {
        return 0
    }

    done := 0
    for _, f := range features {
        if f.Status == "Done" || f.Status == "Shipped" {
            done++
        }
    }

    return float64(done) / float64(len(features)) * 100
}

features, _ := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithRelease("REL-2024Q1"),
)
fmt.Printf("Release completion: %.1f%%\n", completionRate(features))
```

### Find Overdue Features

```go
func findOverdue(features []aha.Feature) []aha.Feature {
    var overdue []aha.Feature
    now := time.Now()

    for _, f := range features {
        if f.DueDate != "" && f.Status != "Done" {
            due, _ := time.Parse("2006-01-02", f.DueDate)
            if due.Before(now) {
                overdue = append(overdue, f)
            }
        }
    }
    return overdue
}
```

## Creating Features

### REST API

```go
feature, err := client.CreateFeature(ctx, "PRODUCT-KEY", aha.CreateFeatureRequest{
    Name:        "User Authentication",
    Description: "Implement OAuth2 login flow",
})
```

### With All Options (REST)

```go
feature, err := client.CreateFeature(ctx, "PRODUCT-KEY", aha.CreateFeatureRequest{
    Name:         "User Authentication",
    Description:  "Implement OAuth2 login flow with SSO support",
    ReleaseID:    "release-id",
    AssignedToID: "user-id",
    Status:       "Ready for Development",
    Tags:         []string{"security", "authentication"},
})
```

### GraphQL API

The GraphQL API provides type-safe mutations for creating features:

```go
import (
    "github.com/grokify/aha-go/graphql"
    "github.com/grokify/aha-go/graphql/generated"
)

client := graphql.NewGenqlientClient("mycompany", "your-api-key")

// Create a feature in a release
resp, err := generated.CreateFeatureWithRelease(ctx, client,
    "User Authentication",
    "RELEASE-ID",
    "<p>Implement OAuth2 login flow</p>",
    "security, authentication",
    nil,
)
fmt.Printf("Created: %s\n", resp.CreateFeature.Feature.ReferenceNum)
```

### With Custom Fields (GraphQL)

```go
// First create the feature
resp, err := generated.CreateFeatureWithRelease(ctx, client,
    "Feature Name", "RELEASE-ID", "<p>Description</p>", "", nil,
)

// Then set custom fields
_, err = generated.SetCustomFieldValues(ctx, client,
    resp.CreateFeature.Feature.Id,
    generated.CustomFieldableTypeEnumFeature,
    []generated.CustomFieldValueInput{
        {Key: "priority_score", Value: 85},
        {Key: "customer_segment", Value: "Enterprise"},
    },
)
```

## Updating Features

### Update Status

```go
_, err := client.UpdateFeature(ctx, "FEAT-123", aha.UpdateFeatureRequest{
    Status: "In Progress",
})
```

### Bulk Update

```go
// Move all "Ready" features to "In Progress"
features, _ := client.ListFeatures(ctx, "PRODUCT-KEY",
    aha.WithStatus("Ready"),
)

for _, f := range features {
    _, err := client.UpdateFeature(ctx, f.ReferenceNum, aha.UpdateFeatureRequest{
        Status: "In Progress",
    })
    if err != nil {
        log.Printf("Failed to update %s: %v", f.ReferenceNum, err)
    }
}
```

## Working with Requirements

Features can have child requirements (sub-tasks).

### List Requirements

```go
requirements, err := client.ListRequirements(ctx, "FEAT-123")
for _, r := range requirements {
    fmt.Printf("  - %s: %s\n", r.ReferenceNum, r.Name)
}
```

### Add Requirement

```go
req, err := client.CreateRequirement(ctx, aha.CreateRequirementRequest{
    FeatureID:   "feature-id",
    Name:        "API endpoint",
    Description: "Create REST API endpoint",
})
```

## Export to JSON

```go
features, _ := client.ListFeatures(ctx, "PRODUCT-KEY")

data, _ := json.MarshalIndent(features, "", "  ")
os.WriteFile("features.json", data, 0644)
```

## Common Patterns

### Progress Dashboard Data

```go
type ProgressData struct {
    Total      int
    Done       int
    InProgress int
    Blocked    int
}

func getProgress(ctx context.Context, client *aha.Client, product string) ProgressData {
    features, _ := client.ListFeatures(ctx, product)

    data := ProgressData{Total: len(features)}
    for _, f := range features {
        switch f.Status {
        case "Done", "Shipped":
            data.Done++
        case "In Progress":
            data.InProgress++
        case "Blocked":
            data.Blocked++
        }
    }
    return data
}
```

### Release Notes Generator

```go
func generateReleaseNotes(ctx context.Context, client *aha.Client, releaseRef string) string {
    features, _ := client.ListFeatures(ctx, "PRODUCT-KEY",
        aha.WithRelease(releaseRef),
        aha.WithStatus("Done"),
    )

    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("# Release %s\n\n", releaseRef))
    sb.WriteString("## Features\n\n")

    for _, f := range features {
        sb.WriteString(fmt.Sprintf("- **%s**: %s\n", f.ReferenceNum, f.Name))
    }

    return sb.String()
}
```
