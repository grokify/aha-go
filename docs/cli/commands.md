# CLI Commands

Complete reference for all `aha` CLI commands.

## Global Flags

These flags work with all commands:

| Flag | Description |
|------|-------------|
| `--output`, `-o` | Output format: `table`, `json` |
| `--help`, `-h` | Show help |

## product

Manage products (workspaces).

### product list

List all products in your account.

```bash
aha product list
aha product list --with-idea-portals
aha product list --updated-since 2024-01-01
aha product list --output json
```

| Flag | Description |
|------|-------------|
| `--with-idea-portals` | Only list products with idea portals |
| `--updated-since` | Filter to products updated after date (ISO8601 or YYYY-MM-DD) |
| `--tree` | Show as tree structure |

### product get

Get details for a specific product.

```bash
aha product get PRODUCT-KEY
aha product get PRODUCT-KEY --output json
```

### product create

Create a new product or product line.

```bash
# Create a product
aha product create --name "My Product" --prefix PROD

# Create with description
aha product create --name "My Product" --prefix PROD --description "Product description"

# Create under a product line
aha product create --name "My Product" --prefix PROD --parent PORT

# Create a product line
aha product create --name "My Portfolio" --prefix PORT --product-line --product-line-type portfolio
```

| Flag | Description |
|------|-------------|
| `--name`, `-n` | Product name (required) |
| `--prefix`, `-p` | Reference prefix (required) |
| `--description`, `-d` | Product description (HTML allowed) |
| `--parent` | Parent product line ID or prefix |
| `--workspace-type` | Workspace type (product_workspace, it_workspace, etc.) |
| `--product-line` | Create as a product line |
| `--product-line-type` | Product line type (required with --product-line) |

### product update

Update an existing product.

```bash
aha product update PROD --name "New Name"
aha product update PROD --description "Updated description"
aha product update PROD --parent PORT
```

| Flag | Description |
|------|-------------|
| `--name`, `-n` | New product name |
| `--prefix`, `-p` | New reference prefix |
| `--description`, `-d` | New description (HTML allowed) |
| `--parent` | New parent product line ID or prefix |
| `--workspace-type` | New workspace type |

## feature

Manage features.

### feature list

List features in a product.

```bash
aha feature list --product PRODUCT-KEY [flags]
```

| Flag | Description |
|------|-------------|
| `--product`, `-p` | Product key (required) |
| `--status` | Filter by status |
| `--release` | Filter by release |
| `--assigned-to` | Filter by assignee email |
| `--limit` | Maximum results |

### feature get

Get a feature by reference number.

```bash
aha feature get FEAT-123
```

### feature create

Create a new feature.

```bash
aha feature create --product PRODUCT-KEY --name "Feature Name" [flags]
```

| Flag | Description |
|------|-------------|
| `--product`, `-p` | Product key (required) |
| `--name` | Feature name (required) |
| `--description` | Feature description |
| `--release` | Release to assign |
| `--status` | Initial status |

### feature update

Update an existing feature.

```bash
aha feature update FEAT-123 [flags]
```

| Flag | Description |
|------|-------------|
| `--name` | New name |
| `--description` | New description |
| `--status` | New status |
| `--release` | New release |

## idea

View customer ideas.

### idea list

List ideas in a product.

```bash
aha idea list --product PRODUCT-KEY [flags]
```

| Flag | Description |
|------|-------------|
| `--product`, `-p` | Product key (required) |
| `--sort` | Sort by: `votes`, `created_at`, `updated_at` |
| `--limit` | Maximum results |

### idea get

Get an idea by reference number.

```bash
aha idea get IDEA-456
```

## epic

View epics.

### epic list

List epics in a product.

```bash
aha epic list --product PRODUCT-KEY
```

### epic get

Get an epic by reference number.

```bash
aha epic get EPIC-789
```

## release

View releases.

### release list

List releases in a product.

```bash
aha release list --product PRODUCT-KEY [flags]
```

| Flag | Description |
|------|-------------|
| `--product`, `-p` | Product key (required) |
| `--active` | Only show unreleased |

### release get

Get a release by reference number.

```bash
aha release get REL-2024Q1
```

## initiative

View strategic initiatives.

### initiative list

List initiatives.

```bash
aha initiative list --product PRODUCT-KEY
```

### initiative get

Get an initiative by reference number.

```bash
aha initiative get INIT-123
```

## goal

View goals.

### goal list

List goals.

```bash
aha goal list --product PRODUCT-KEY
```

### goal get

Get a goal by reference number.

```bash
aha goal get GOAL-456
```

## requirement

Manage feature requirements.

### requirement list

List requirements for a feature.

```bash
aha requirement list --feature FEAT-123
```

### requirement get

Get a requirement by ID.

```bash
aha requirement get REQ-789
```

### requirement create

Create a requirement on a feature.

```bash
aha requirement create --feature FEAT-123 --name "Requirement" [flags]
```

| Flag | Description |
|------|-------------|
| `--feature` | Feature reference (required) |
| `--name` | Requirement name (required) |
| `--description` | Description |

### requirement update

Update a requirement.

```bash
aha requirement update REQ-789 --name "Updated Name"
```

### requirement delete

Delete a requirement.

```bash
aha requirement delete REQ-789
```

## comment

Manage comments on entities.

### comment list

List comments.

```bash
aha comment list --feature FEAT-123
aha comment list --idea IDEA-456
aha comment list --product PRODUCT-KEY
```

### comment get

Get a comment by ID.

```bash
aha comment get COMMENT-ID
```

### comment create

Add a comment.

```bash
aha comment create --feature FEAT-123 --body "Comment text"
aha comment create --idea IDEA-456 --body "Comment text"
```

### comment delete

Delete a comment.

```bash
aha comment delete COMMENT-ID
```

## user

View users.

### user list

List users in the account.

```bash
aha user list
```

### user get

Get a user by ID or email.

```bash
aha user get USER-ID
aha user get user@company.com
```

### user me

Get the current authenticated user.

```bash
aha user me
```

## canvas

Create and export strategic canvases.

### canvas list

List canvases in a product.

```bash
aha canvas list --product PRODUCT-KEY
```

### canvas get

Get a canvas by ID.

```bash
aha canvas get CANVAS-ID
```

### canvas create bmc

Create a Business Model Canvas.

```bash
aha canvas create bmc --product PRODUCT-KEY --name "My BMC" --file bmc.json
```

### canvas create leanux

Create a Lean UX Canvas.

```bash
aha canvas create leanux --product PRODUCT-KEY --name "My Lean UX" --file leanux.json
```

### canvas create opportunity

Create an Opportunity Canvas.

```bash
aha canvas create opportunity --product PRODUCT-KEY --name "Opportunity" --file opp.json
```

### canvas export

Export a canvas to various formats.

```bash
aha canvas export CANVAS-ID --format png --output canvas.png
aha canvas export CANVAS-ID --format pdf --output canvas.pdf
```

## completion

Generate shell completion scripts.

```bash
# Bash
aha completion bash

# Zsh
aha completion zsh

# Fish
aha completion fish

# PowerShell
aha completion powershell
```
