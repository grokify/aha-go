# Use Cases

This section provides narrative guides for common workflows and integration patterns with `aha-go`.

## Overview

| Use Case | Description |
|----------|-------------|
| [Querying Features](features.md) | Search, filter, and analyze features |
| [Managing Ideas](ideas.md) | Track customer feedback and votes |
| [Strategic Canvases](canvases.md) | Create and manage strategic canvases |
| [Generating Diagrams](diagrams.md) | Visualize roadmaps and relationships |

## Integration Patterns

### Automation Scripts

Use `aha-go` to automate repetitive tasks:

- Generate weekly status reports
- Sync features with external systems
- Archive completed items
- Bulk update feature statuses

### CI/CD Integration

Integrate with your development workflow:

- Update feature status when PRs merge
- Link commits to requirements
- Generate release notes from completed features

### Data Export

Export data for analysis or backup:

- Export features to CSV/JSON
- Archive ideas with vote counts
- Snapshot release state

### Slack/Teams Bots

Build notification bots:

- Alert on high-voted ideas
- Notify when releases are published
- Daily feature progress summaries

## Example: Weekly Status Report

```go
package main

import (
    "context"
    "fmt"
    "time"

    aha "github.com/grokify/aha-go"
)

func main() {
    ctx := context.Background()
    client, _ := aha.NewClient()

    // Get features updated this week
    weekAgo := time.Now().AddDate(0, 0, -7)

    features, _ := client.ListFeatures(ctx, "PLATFORM",
        aha.WithUpdatedSince(weekAgo),
    )

    // Group by status
    byStatus := make(map[string][]aha.Feature)
    for _, f := range features {
        byStatus[f.Status] = append(byStatus[f.Status], f)
    }

    // Generate report
    fmt.Println("# Weekly Status Report")
    fmt.Printf("Period: %s to %s\n\n",
        weekAgo.Format("Jan 2"),
        time.Now().Format("Jan 2"))

    for status, feats := range byStatus {
        fmt.Printf("## %s (%d)\n", status, len(feats))
        for _, f := range feats {
            fmt.Printf("- %s: %s\n", f.ReferenceNum, f.Name)
        }
        fmt.Println()
    }
}
```

## Example: Idea Voting Alerts

```go
// Monitor ideas and alert when they reach vote thresholds
func monitorIdeas(ctx context.Context, client *aha.Client, threshold int) {
    ideas, _ := client.ListIdeas(ctx, "PLATFORM")

    for _, idea := range ideas {
        if idea.Votes >= threshold {
            // Send alert (Slack, email, etc.)
            fmt.Printf("High-voted idea: %s (%d votes)\n",
                idea.Name, idea.Votes)
        }
    }
}
```

## Next Steps

Explore detailed guides for specific use cases:

- [Querying Features](features.md) - Advanced filtering and analysis
- [Managing Ideas](ideas.md) - Customer feedback workflows
- [Strategic Canvases](canvases.md) - Business model tools
- [Generating Diagrams](diagrams.md) - Visual roadmaps
