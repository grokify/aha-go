# API Reference

Complete API documentation is available on pkg.go.dev.

## Package Documentation

| Package | pkg.go.dev |
|---------|------------|
| `aha-go` (root) | [pkg.go.dev/github.com/grokify/aha-go](https://pkg.go.dev/github.com/grokify/aha-go) |
| `canvas` | [pkg.go.dev/github.com/grokify/aha-go/canvas](https://pkg.go.dev/github.com/grokify/aha-go/canvas) |
| `browser` | [pkg.go.dev/github.com/grokify/aha-go/browser](https://pkg.go.dev/github.com/grokify/aha-go/browser) |
| `render` | [pkg.go.dev/github.com/grokify/aha-go/render](https://pkg.go.dev/github.com/grokify/aha-go/render) |

## Core Types

### Client

```go
type Client struct {
    // unexported fields
}

func NewClient() (*Client, error)
func NewClientWithConfig(cfg Config) (*Client, error)
```

### Config

```go
type Config struct {
    Subdomain string
    APIKey    string
    BaseURL   string // Optional, defaults to https://{subdomain}.aha.io
}

func ConfigFromEnv() Config
```

## Entity Types

### Feature

```go
type Feature struct {
    ID              string    `json:"id"`
    ReferenceNum    string    `json:"reference_num"`
    Name            string    `json:"name"`
    Description     string    `json:"description"`
    Status          string    `json:"workflow_status"`
    AssignedTo      *User     `json:"assigned_to_user"`
    Release         *Release  `json:"release"`
    StartDate       string    `json:"start_date"`
    DueDate         string    `json:"due_date"`
    Tags            []Tag     `json:"tags"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Idea

```go
type Idea struct {
    ID             string     `json:"id"`
    ReferenceNum   string     `json:"reference_num"`
    Name           string     `json:"name"`
    Description    string     `json:"description"`
    Votes          int        `json:"votes"`
    WorkflowStatus string     `json:"workflow_status"`
    Categories     []Category `json:"idea_categories"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}
```

### Release

```go
type Release struct {
    ID           string    `json:"id"`
    ReferenceNum string    `json:"reference_num"`
    Name         string    `json:"name"`
    StartDate    string    `json:"start_date"`
    ReleaseDate  string    `json:"release_date"`
    Released     bool      `json:"released"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Initiative

```go
type Initiative struct {
    ID           string    `json:"id"`
    ReferenceNum string    `json:"reference_num"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Status       string    `json:"workflow_status"`
    Value        float64   `json:"value"`
    Effort       float64   `json:"effort"`
    Progress     float64   `json:"progress"`
    StartDate    string    `json:"start_date"`
    EndDate      string    `json:"end_date"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Goal

```go
type Goal struct {
    ID           string    `json:"id"`
    ReferenceNum string    `json:"reference_num"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Status       string    `json:"workflow_status"`
    Progress     float64   `json:"progress"`
    StartDate    string    `json:"start_date"`
    EndDate      string    `json:"end_date"`
    CreatedAt    time.Time `json:"created_at"`
}
```

## Client Methods

### Products

```go
func (c *Client) ListProducts(ctx context.Context) ([]Product, error)
func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error)
```

### Features

```go
func (c *Client) ListFeatures(ctx context.Context, productID string, opts ...Option) ([]Feature, error)
func (c *Client) GetFeature(ctx context.Context, ref string) (*Feature, error)
func (c *Client) CreateFeature(ctx context.Context, productID string, req CreateFeatureRequest) (*Feature, error)
func (c *Client) UpdateFeature(ctx context.Context, ref string, req UpdateFeatureRequest) (*Feature, error)
```

### Ideas

```go
func (c *Client) ListIdeas(ctx context.Context, productID string, opts ...Option) ([]Idea, error)
func (c *Client) GetIdea(ctx context.Context, ref string) (*Idea, error)
```

### Releases

```go
func (c *Client) ListReleases(ctx context.Context, productID string, opts ...Option) ([]Release, error)
func (c *Client) GetRelease(ctx context.Context, ref string) (*Release, error)
```

### Comments

```go
func (c *Client) ListComments(ctx context.Context, opts ...Option) ([]Comment, error)
func (c *Client) GetComment(ctx context.Context, id string) (*Comment, error)
func (c *Client) CreateComment(ctx context.Context, req CreateCommentRequest) (*Comment, error)
func (c *Client) DeleteComment(ctx context.Context, id string) error
```

## Options

```go
// Pagination
func WithPage(page int) Option
func WithPerPage(perPage int) Option

// Filtering
func WithStatus(status string) Option
func WithRelease(release string) Option
func WithAssignedTo(email string) Option

// Sorting
func WithSort(field string) Option
func WithSortDirection(dir string) Option

// Date filtering
func WithUpdatedSince(t time.Time) Option
func WithCreatedSince(t time.Time) Option
```

## Error Handling

```go
func IsNotFoundError(err error) bool
func IsUnauthorizedError(err error) bool
func IsRateLimitError(err error) bool
```

## Aha.io API Documentation

For details on the underlying Aha.io API:

- [Aha.io API Documentation](https://www.aha.io/api)
- [Aha.io API Reference](https://www.aha.io/api/reference)
