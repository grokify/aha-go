# Canvas Content Examples

This directory contains example JSON files with content for strategic canvases.

## Files

| File | Description | Blocks |
|------|-------------|--------|
| `opportunity.json` | Opportunity Canvas (Jeff Patton) | 10 blocks |
| `lean-ux.json` | Lean UX Canvas (Jeff Gothelf) | 8 blocks |
| `bmc.json` | Business Model Canvas (Osterwalder) | 9 blocks |

## Usage

Use these files with the `aha canvas update` command:

```bash
# Update an existing Opportunity canvas
aha canvas update PROD-SM-1 --file examples/canvas/opportunity.json

# Update a Lean UX canvas
aha canvas update PROD-SM-2 --file examples/canvas/lean-ux.json

# Update a Business Model Canvas
aha canvas update PROD-SM-3 --file examples/canvas/bmc.json
```

## JSON Format

The JSON format maps block names to HTML content:

```json
{
  "Block Name": "<p>Content with <strong>formatting</strong></p>",
  "Another Block": "<ul><li>List item 1</li><li>List item 2</li></ul>"
}
```

Supported HTML tags:

- `<p>` - Paragraphs
- `<ul>`, `<ol>`, `<li>` - Lists
- `<strong>`, `<b>` - Bold
- `<em>`, `<i>` - Italic
- `<br>` - Line breaks

## Block Names by Canvas Type

### Opportunity Canvas

1. Users & Customers
2. Problems
3. Solution Ideas
4. Solutions Today
5. User Value
6. Adoption Strategy
7. User Metrics
8. Business Problem
9. Business Metrics
10. Budget

### Lean UX Canvas

1. Business Problem
2. Business Outcomes
3. Users
4. Benefits
5. Solutions
6. Hypotheses
7. Riskiest Assumption
8. Smallest Experiment

### Business Model Canvas

1. Key Partners
2. Key Activities
3. Key Resources
4. Value Propositions
5. Customer Relationships
6. Channels
7. Customer Segments
8. Cost Structure
9. Revenue Streams
