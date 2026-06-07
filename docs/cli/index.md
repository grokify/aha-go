# CLI Overview

The `aha` command-line tool provides quick access to Aha.io data from your terminal.

## Installation

```bash
go install github.com/grokify/aha-go/cmd/aha@latest
```

## Configuration

Set environment variables:

```bash
export AHA_SUBDOMAIN="yourcompany"
export AHA_API_KEY="your-api-key"
```

## Basic Usage

```bash
# Get help
aha --help

# List products
aha product list

# List features in a product
aha feature list --product PRODUCT-KEY

# Get feature details
aha feature get FEAT-123
```

## Output Formats

All commands support multiple output formats:

```bash
# Table format (default)
aha feature list --product PLATFORM

# JSON format
aha feature list --product PLATFORM --output json

# JSON output is useful for scripting
aha feature get FEAT-123 --output json | jq '.name'
```

## Command Categories

| Category | Description |
|----------|-------------|
| `product` | List and view products |
| `feature` | Manage features |
| `idea` | View ideas and votes |
| `epic` | View epics |
| `release` | View releases |
| `initiative` | View strategic initiatives |
| `goal` | View goals |
| `requirement` | Manage feature requirements |
| `comment` | Manage comments |
| `user` | View users |
| `canvas` | Create and export strategic canvases |
| `completion` | Generate shell completions |

## Shell Completion

Enable tab completion for your shell:

=== "Bash"

    ```bash
    aha completion bash > /etc/bash_completion.d/aha
    ```

=== "Zsh"

    ```bash
    aha completion zsh > "${fpath[1]}/_aha"
    ```

=== "Fish"

    ```bash
    aha completion fish > ~/.config/fish/completions/aha.fish
    ```

## Common Workflows

### Daily Standup Check

```bash
# Features in progress
aha feature list --product PLATFORM --status "In Progress"

# Recent comments
aha comment list --product PLATFORM --since "24h"
```

### Release Planning

```bash
# List releases
aha release list --product PLATFORM

# Features in a release
aha feature list --product PLATFORM --release REL-2024Q1
```

### Idea Triage

```bash
# Top voted ideas
aha idea list --product PLATFORM --sort votes --limit 10

# Ideas without features
aha idea list --product PLATFORM --no-feature
```

## Next Steps

- [Complete Command Reference](commands.md)
- [Use Case Examples](../use-cases/index.md)
