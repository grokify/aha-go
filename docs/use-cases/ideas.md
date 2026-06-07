# Managing Ideas

Ideas in Aha.io capture customer feedback and feature requests. This guide covers patterns for managing ideas, tracking votes, and prioritizing based on customer input.

## Understanding Ideas

Ideas represent customer requests that can be:

- Voted on by customers/stakeholders
- Categorized for triage
- Promoted to features
- Linked to existing features

## Listing Ideas

### Basic List

```go
ideas, err := client.ListIdeas(ctx, "PRODUCT-KEY")
if err != nil {
    log.Fatal(err)
}

for _, idea := range ideas {
    fmt.Printf("%s: %s (%d votes)\n", idea.ReferenceNum, idea.Name, idea.Votes)
}
```

### With CLI

```bash
aha idea list --product PRODUCT-KEY
```

## Sorting by Votes

### Top Voted Ideas

```go
ideas, err := client.ListIdeas(ctx, "PRODUCT-KEY",
    aha.WithSort("votes"),
    aha.WithSortDirection("desc"),
    aha.WithLimit(10),
)

fmt.Println("Top 10 Ideas by Votes:")
for i, idea := range ideas {
    fmt.Printf("%d. %s - %d votes\n", i+1, idea.Name, idea.Votes)
}
```

### CLI

```bash
aha idea list --product PRODUCT-KEY --sort votes --limit 10
```

## Analyzing Ideas

### Vote Distribution

```go
func analyzeVotes(ideas []aha.Idea) {
    var total, max int
    buckets := map[string]int{
        "0":     0,
        "1-5":   0,
        "6-10":  0,
        "11-50": 0,
        "50+":   0,
    }

    for _, idea := range ideas {
        total += idea.Votes
        if idea.Votes > max {
            max = idea.Votes
        }

        switch {
        case idea.Votes == 0:
            buckets["0"]++
        case idea.Votes <= 5:
            buckets["1-5"]++
        case idea.Votes <= 10:
            buckets["6-10"]++
        case idea.Votes <= 50:
            buckets["11-50"]++
        default:
            buckets["50+"]++
        }
    }

    fmt.Printf("Total ideas: %d\n", len(ideas))
    fmt.Printf("Total votes: %d\n", total)
    fmt.Printf("Max votes: %d\n", max)
    fmt.Printf("Average: %.1f\n", float64(total)/float64(len(ideas)))
    fmt.Println("\nDistribution:")
    for bucket, count := range buckets {
        fmt.Printf("  %s votes: %d ideas\n", bucket, count)
    }
}
```

### By Category

```go
func groupByCategory(ideas []aha.Idea) map[string][]aha.Idea {
    result := make(map[string][]aha.Idea)
    for _, idea := range ideas {
        for _, cat := range idea.Categories {
            result[cat.Name] = append(result[cat.Name], idea)
        }
    }
    return result
}

ideas, _ := client.ListIdeas(ctx, "PRODUCT-KEY")
byCategory := groupByCategory(ideas)

for category, categoryIdeas := range byCategory {
    totalVotes := 0
    for _, idea := range categoryIdeas {
        totalVotes += idea.Votes
    }
    fmt.Printf("%s: %d ideas, %d total votes\n",
        category, len(categoryIdeas), totalVotes)
}
```

## Idea Triage Workflow

### Find Untriaged Ideas

```go
// Ideas without a workflow status indicating triage
func findUntriaged(ideas []aha.Idea) []aha.Idea {
    var untriaged []aha.Idea
    for _, idea := range ideas {
        if idea.WorkflowStatus == "" || idea.WorkflowStatus == "New" {
            untriaged = append(untriaged, idea)
        }
    }
    return untriaged
}
```

### High-Priority Ideas

```go
// Ideas with high votes that need attention
func findHighPriority(ideas []aha.Idea, voteThreshold int) []aha.Idea {
    var highPriority []aha.Idea
    for _, idea := range ideas {
        if idea.Votes >= voteThreshold && idea.FeatureID == "" {
            highPriority = append(highPriority, idea)
        }
    }
    return highPriority
}

// Find ideas with 10+ votes not linked to features
highPriority := findHighPriority(ideas, 10)
fmt.Printf("Found %d high-priority ideas needing features\n", len(highPriority))
```

## Monitoring and Alerts

### Vote Threshold Alerts

```go
func checkVoteThresholds(ctx context.Context, client *aha.Client, product string) {
    ideas, _ := client.ListIdeas(ctx, product)

    thresholds := []int{10, 25, 50, 100}

    for _, threshold := range thresholds {
        var crossing []aha.Idea
        for _, idea := range ideas {
            // Check if votes just crossed threshold (within last day)
            if idea.Votes >= threshold && idea.Votes < threshold+5 {
                crossing = append(crossing, idea)
            }
        }

        if len(crossing) > 0 {
            fmt.Printf("\nIdeas crossing %d votes:\n", threshold)
            for _, idea := range crossing {
                fmt.Printf("  - %s: %s (%d votes)\n",
                    idea.ReferenceNum, idea.Name, idea.Votes)
            }
        }
    }
}
```

### Weekly Idea Summary

```go
func weeklyIdeaSummary(ctx context.Context, client *aha.Client, product string) {
    weekAgo := time.Now().AddDate(0, 0, -7)

    ideas, _ := client.ListIdeas(ctx, product,
        aha.WithCreatedSince(weekAgo),
    )

    fmt.Printf("New ideas this week: %d\n", len(ideas))

    // Total new votes
    totalVotes := 0
    for _, idea := range ideas {
        totalVotes += idea.Votes
    }
    fmt.Printf("Total votes on new ideas: %d\n", totalVotes)

    // Top new idea
    if len(ideas) > 0 {
        top := ideas[0]
        for _, idea := range ideas[1:] {
            if idea.Votes > top.Votes {
                top = idea
            }
        }
        fmt.Printf("Top new idea: %s (%d votes)\n", top.Name, top.Votes)
    }
}
```

## Exporting Ideas

### To CSV

```go
func exportIdeasCSV(ideas []aha.Idea, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Header
    writer.Write([]string{"Reference", "Name", "Votes", "Status", "Created"})

    // Data
    for _, idea := range ideas {
        writer.Write([]string{
            idea.ReferenceNum,
            idea.Name,
            fmt.Sprintf("%d", idea.Votes),
            idea.WorkflowStatus,
            idea.CreatedAt,
        })
    }

    return nil
}
```

### To JSON with Analysis

```go
type IdeaReport struct {
    GeneratedAt string         `json:"generated_at"`
    Product     string         `json:"product"`
    TotalIdeas  int            `json:"total_ideas"`
    TotalVotes  int            `json:"total_votes"`
    TopIdeas    []aha.Idea     `json:"top_ideas"`
    ByCategory  map[string]int `json:"by_category"`
}

func generateIdeaReport(ctx context.Context, client *aha.Client, product string) IdeaReport {
    ideas, _ := client.ListIdeas(ctx, product)

    report := IdeaReport{
        GeneratedAt: time.Now().Format(time.RFC3339),
        Product:     product,
        TotalIdeas:  len(ideas),
        ByCategory:  make(map[string]int),
    }

    // Calculate totals
    for _, idea := range ideas {
        report.TotalVotes += idea.Votes
        for _, cat := range idea.Categories {
            report.ByCategory[cat.Name]++
        }
    }

    // Top 10 by votes
    sort.Slice(ideas, func(i, j int) bool {
        return ideas[i].Votes > ideas[j].Votes
    })
    if len(ideas) > 10 {
        report.TopIdeas = ideas[:10]
    } else {
        report.TopIdeas = ideas
    }

    return report
}
```

## Comments on Ideas

### List Comments

```go
comments, err := client.ListComments(ctx,
    aha.WithCommentableType("Idea"),
    aha.WithCommentableID("IDEA-123"),
)

for _, c := range comments {
    fmt.Printf("%s: %s\n", c.User.Name, c.Body)
}
```

### Add Response Comment

```go
_, err := client.CreateComment(ctx, aha.CreateCommentRequest{
    CommentableType: "Idea",
    CommentableID:   "IDEA-123",
    Body:            "Thanks for this suggestion! We're reviewing it for our Q2 roadmap.",
})
```
