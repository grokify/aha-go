package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Idea represents an Aha idea.
type Idea struct {
	ID              string
	ReferenceNum    string
	Name            string
	Description     string
	Votes           int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StatusChangedAt *time.Time
	WorkflowStatus  *WorkflowStatus
	Categories      []Category
	Feature         *IdeaFeature
}

// Category represents an idea category.
type Category struct {
	ID        string
	Name      string
	ParentID  string
	ProjectID string
	CreatedAt time.Time
}

// IdeaFeature represents a feature linked to an idea.
type IdeaFeature struct {
	ID           string
	ReferenceNum string
	Name         string
	URL          string
	Resource     string
	ProductID    string
	CreatedAt    time.Time
}

// IdeaList represents a paginated list of ideas.
type IdeaList struct {
	Ideas      []Idea
	Pagination Pagination
}

// GetIdea retrieves an idea by ID or reference number.
func (c *Client) GetIdea(ctx context.Context, id string) (*Idea, error) {
	resp, err := c.apiClient.GetIdea(ctx, api.GetIdeaParams{
		IdeaID: id,
	})
	if err != nil {
		return nil, wrapError("GetIdea", err)
	}

	if i, ok := resp.Idea.Get(); ok {
		return ideaFromAPI(i), nil
	}
	return nil, &APIError{StatusCode: 404, Message: "idea not found"}
}

// ListIdeasOptions configures ListIdeas.
type ListIdeasOptions struct {
	Query          string
	WorkflowStatus string
	Sort           string
	Spam           *bool
	Tag            string
	UserID         string
	IdeaUserID     string
	CreatedBefore  *time.Time
	CreatedSince   *time.Time
	UpdatedSince   *time.Time
	Page           int
	PerPage        int
}

// ListIdeasOption configures a ListIdeas call.
type ListIdeasOption func(*ListIdeasOptions)

// WithIdeaQuery filters ideas by search query.
func WithIdeaQuery(query string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.Query = query }
}

// WithIdeaStatus filters ideas by workflow status.
func WithIdeaStatus(status string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.WorkflowStatus = status }
}

// WithIdeaSort sets the sort order (recent, trending, popular).
func WithIdeaSort(sort string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.Sort = sort }
}

// WithIdeaTag filters ideas by tag.
func WithIdeaTag(tag string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.Tag = tag }
}

// WithIdeaCreatedSince filters ideas created after the given time.
func WithIdeaCreatedSince(t time.Time) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.CreatedSince = &t }
}

// WithIdeaUpdatedSince filters ideas updated after the given time.
func WithIdeaUpdatedSince(t time.Time) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.UpdatedSince = &t }
}

// WithIdeaPage sets the page number for pagination.
func WithIdeaPage(page int) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.Page = page }
}

// WithIdeaPerPage sets the number of results per page.
func WithIdeaPerPage(perPage int) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.PerPage = perPage }
}

// WithIdeaSpam filters ideas by spam status.
func WithIdeaSpam(spam bool) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.Spam = &spam }
}

// WithIdeaUserID filters ideas by creator user ID.
func WithIdeaUserID(userID string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.UserID = userID }
}

// WithIdeaIdeaUserID filters ideas by idea user ID.
func WithIdeaIdeaUserID(ideaUserID string) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.IdeaUserID = ideaUserID }
}

// WithIdeaCreatedBefore filters ideas created before the given time.
func WithIdeaCreatedBefore(t time.Time) ListIdeasOption {
	return func(o *ListIdeasOptions) { o.CreatedBefore = &t }
}

// ListIdeas lists ideas with optional filtering.
func (c *Client) ListIdeas(ctx context.Context, opts ...ListIdeasOption) (*IdeaList, error) {
	cfg := &ListIdeasOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	params := api.ListIdeasParams{}
	if cfg.Query != "" {
		params.Q = api.NewOptString(cfg.Query)
	}
	if cfg.WorkflowStatus != "" {
		params.WorkflowStatus = api.NewOptString(cfg.WorkflowStatus)
	}
	if cfg.Sort != "" {
		params.Sort = api.NewOptListIdeasSort(api.ListIdeasSort(cfg.Sort))
	}
	if cfg.Spam != nil {
		params.Spam = api.NewOptBool(*cfg.Spam)
	}
	if cfg.Tag != "" {
		params.Tag = api.NewOptString(cfg.Tag)
	}
	if cfg.UserID != "" {
		params.UserID = api.NewOptString(cfg.UserID)
	}
	if cfg.IdeaUserID != "" {
		params.IdeaUserID = api.NewOptString(cfg.IdeaUserID)
	}
	if cfg.CreatedBefore != nil {
		params.CreatedBefore = api.NewOptDateTime(*cfg.CreatedBefore)
	}
	if cfg.CreatedSince != nil {
		params.CreatedSince = api.NewOptDateTime(*cfg.CreatedSince)
	}
	if cfg.UpdatedSince != nil {
		params.UpdatedSince = api.NewOptDateTime(*cfg.UpdatedSince)
	}
	if cfg.Page > 0 {
		params.Page = api.NewOptInt32(int32(cfg.Page))
	}
	if cfg.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(cfg.PerPage))
	}

	resp, err := c.apiClient.ListIdeas(ctx, params)
	if err != nil {
		return nil, wrapError("ListIdeas", err)
	}

	return ideaListFromAPI(resp), nil
}

// ideaFromAPI converts an API idea to a domain idea.
func ideaFromAPI(i api.Idea) *Idea {
	idea := &Idea{
		ID:           i.ID,
		ReferenceNum: i.ReferenceNum,
		Name:         i.Name,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}

	if v, ok := i.Votes.Get(); ok {
		idea.Votes = v
	}
	if v, ok := i.Description.Get(); ok {
		idea.Description = v
	}
	if v, ok := i.StatusChangedAt.Get(); ok {
		idea.StatusChangedAt = &v
	}
	if v, ok := i.WorkflowStatus.Get(); ok {
		idea.WorkflowStatus = workflowStatusFromAPI(v)
	}
	idea.Categories = categoriesFromAPI(i.Categories)
	if v, ok := i.Feature.Get(); ok {
		idea.Feature = ideaFeatureFromAPI(v)
	}

	return idea
}

// ideaListFromAPI converts an API ideas response to a domain idea list.
func ideaListFromAPI(resp *api.IdeasResponse) *IdeaList {
	list := &IdeaList{}

	list.Ideas = make([]Idea, len(resp.Ideas))
	for i, idea := range resp.Ideas {
		list.Ideas[i] = *ideaFromAPI(idea)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// categoriesFromAPI converts API categories.
func categoriesFromAPI(cats []api.Category) []Category {
	result := make([]Category, len(cats))
	for i, c := range cats {
		result[i] = Category{}
		if v, ok := c.ID.Get(); ok {
			result[i].ID = v
		}
		if v, ok := c.Name.Get(); ok {
			result[i].Name = v
		}
		if v, ok := c.ParentID.Get(); ok {
			result[i].ParentID = v
		}
		if v, ok := c.ProjectID.Get(); ok {
			result[i].ProjectID = v
		}
		if v, ok := c.CreatedAt.Get(); ok {
			result[i].CreatedAt = v
		}
	}
	return result
}

// ideaFeatureFromAPI converts an API idea feature.
func ideaFeatureFromAPI(f api.IdeaFeature) *IdeaFeature {
	feature := &IdeaFeature{}
	if v, ok := f.ID.Get(); ok {
		feature.ID = v
	}
	if v, ok := f.ReferenceNum.Get(); ok {
		feature.ReferenceNum = v
	}
	if v, ok := f.Name.Get(); ok {
		feature.Name = v
	}
	if v, ok := f.URL.Get(); ok {
		feature.URL = v
	}
	if v, ok := f.Resource.Get(); ok {
		feature.Resource = v
	}
	if v, ok := f.ProductID.Get(); ok {
		feature.ProductID = v
	}
	if v, ok := f.CreatedAt.Get(); ok {
		feature.CreatedAt = v
	}
	return feature
}
