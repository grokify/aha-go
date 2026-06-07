package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// Comment represents an Aha comment.
type Comment struct {
	ID          string
	Body        string
	URL         string
	Resource    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        *User
	Attachments []Attachment
	Commentable *Commentable
}

// CommentList represents a paginated list of comments.
type CommentList struct {
	Comments   []CommentMeta
	Pagination Pagination
}

// CommentMeta represents comment metadata in list responses.
type CommentMeta struct {
	ID        string
	Body      string
	URL       string
	Resource  string
	CreatedAt time.Time
	User      *User
}

// Commentable represents the object that was commented on.
type Commentable struct {
	Type      string
	ID        string
	ProductID string
	URL       string
	Resource  string
}

// Attachment represents a file attachment on a comment.
type Attachment struct {
	ID          string
	DownloadURL string
	FileName    string
	FileSize    int64
	ContentType string
}

// GetComment retrieves a comment by ID.
func (c *Client) GetComment(ctx context.Context, id string) (*Comment, error) {
	resp, err := c.apiClient.GetComment(ctx, api.GetCommentParams{
		CommentID: id,
	})
	if err != nil {
		return nil, wrapError("GetComment", err)
	}

	// Handle different response types
	switch r := resp.(type) {
	case *api.CommentResponse:
		if comment, ok := r.Comment.Get(); ok {
			return commentFromAPI(comment), nil
		}
		return nil, &APIError{StatusCode: 404, Message: "comment not found"}
	default:
		return nil, &APIError{StatusCode: 404, Message: "comment not found"}
	}
}

// ListProductComments lists comments in a product.
func (c *Client) ListProductComments(ctx context.Context, productID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListProductCommentsParams{
		ProductID: productID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListProductComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListFeatureComments lists comments on a feature.
func (c *Client) ListFeatureComments(ctx context.Context, featureID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListFeatureCommentsParams{
		FeatureID: featureID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListFeatureComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListFeatureComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListIdeaComments lists comments on an idea.
func (c *Client) ListIdeaComments(ctx context.Context, ideaID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListIdeaCommentsParams{
		IdeaID: ideaID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListIdeaComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListIdeaComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListReleaseComments lists comments on a release.
func (c *Client) ListReleaseComments(ctx context.Context, releaseID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListReleaseCommentsParams{
		ReleaseID: releaseID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListReleaseComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListReleaseComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListInitiativeComments lists comments on an initiative.
func (c *Client) ListInitiativeComments(ctx context.Context, initiativeID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListInitiativeCommentsParams{
		InitiativeID: initiativeID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListInitiativeComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListInitiativeComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListEpicComments lists comments on an epic.
func (c *Client) ListEpicComments(ctx context.Context, epicID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListEpicCommentsParams{
		EpicID: epicID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListEpicComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListEpicComments", err)
	}

	return commentListFromAPI(resp), nil
}

// ListGoalComments lists comments on a goal.
func (c *Client) ListGoalComments(ctx context.Context, goalID string, opts ...ListOption) (*CommentList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListGoalCommentsParams{
		GoalID: goalID,
	}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page))
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage))
	}

	resp, err := c.apiClient.ListGoalComments(ctx, params)
	if err != nil {
		return nil, wrapError("ListGoalComments", err)
	}

	return commentListFromAPI(resp), nil
}

// CreateCommentOptions configures CreateComment calls.
type CreateCommentOptions struct {
	Body string
}

// CreateCommentOption configures a CreateComment call.
type CreateCommentOption func(*CreateCommentOptions)

// WithCommentBody sets the comment body.
func WithCommentBody(body string) CreateCommentOption {
	return func(o *CreateCommentOptions) { o.Body = body }
}

// CreateFeatureComment creates a comment on a feature.
func (c *Client) CreateFeatureComment(ctx context.Context, featureID string, opts ...CreateCommentOption) (*Comment, error) {
	cfg := &CreateCommentOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Body == "" {
		return nil, &APIError{StatusCode: 400, Message: "comment body is required"}
	}

	req := &api.CommentCreateRequest{
		Comment: api.CommentCreate{
			Body: cfg.Body,
		},
	}

	resp, err := c.apiClient.CreateFeatureComment(ctx, req, api.CreateFeatureCommentParams{
		FeatureID: featureID,
	})
	if err != nil {
		return nil, wrapError("CreateFeatureComment", err)
	}

	if comment, ok := resp.Comment.Get(); ok {
		return commentFromAPI(comment), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: comment not returned"}
}

// CreateIdeaComment creates an internal comment on an idea.
func (c *Client) CreateIdeaComment(ctx context.Context, ideaID string, opts ...CreateCommentOption) (*Comment, error) {
	cfg := &CreateCommentOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Body == "" {
		return nil, &APIError{StatusCode: 400, Message: "comment body is required"}
	}

	req := &api.CommentCreateRequest{
		Comment: api.CommentCreate{
			Body: cfg.Body,
		},
	}

	resp, err := c.apiClient.CreateIdeaComment(ctx, req, api.CreateIdeaCommentParams{
		IdeaID: ideaID,
	})
	if err != nil {
		return nil, wrapError("CreateIdeaComment", err)
	}

	if comment, ok := resp.Comment.Get(); ok {
		return commentFromAPI(comment), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: comment not returned"}
}

// UpdateCommentOptions configures UpdateComment calls.
type UpdateCommentOptions struct {
	Body string
}

// UpdateCommentOption configures an UpdateComment call.
type UpdateCommentOption func(*UpdateCommentOptions)

// WithUpdateCommentBody sets the comment body.
func WithUpdateCommentBody(body string) UpdateCommentOption {
	return func(o *UpdateCommentOptions) { o.Body = body }
}

// UpdateComment updates an existing comment.
func (c *Client) UpdateComment(ctx context.Context, id string, opts ...UpdateCommentOption) (*Comment, error) {
	cfg := &UpdateCommentOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	req := api.CommentUpdate{}
	if cfg.Body != "" {
		req.Body = api.NewOptString(cfg.Body)
	}

	reqBody := &api.CommentUpdateRequest{
		Comment: req,
	}

	resp, err := c.apiClient.UpdateComment(ctx, reqBody, api.UpdateCommentParams{
		CommentID: id,
	})
	if err != nil {
		return nil, wrapError("UpdateComment", err)
	}

	if comment, ok := resp.Comment.Get(); ok {
		return commentFromAPI(comment), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: comment not returned"}
}

// DeleteComment deletes a comment.
func (c *Client) DeleteComment(ctx context.Context, id string) error {
	err := c.apiClient.DeleteComment(ctx, api.DeleteCommentParams{
		CommentID: id,
	})
	if err != nil {
		return wrapError("DeleteComment", err)
	}
	return nil
}

// commentFromAPI converts an API comment to a domain comment.
func commentFromAPI(c api.Comment) *Comment {
	comment := &Comment{
		ID:        c.ID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if v, ok := c.URL.Get(); ok {
		comment.URL = v
	}
	if v, ok := c.Resource.Get(); ok {
		comment.Resource = v
	}
	if v, ok := c.User.Get(); ok {
		comment.User = userFromAPI(v)
	}
	if v, ok := c.Commentable.Get(); ok {
		comment.Commentable = commentableFromAPI(v)
	}

	if attachments := c.Attachments; len(attachments) > 0 {
		comment.Attachments = make([]Attachment, len(attachments))
		for i, a := range attachments {
			comment.Attachments[i] = attachmentFromAPI(a)
		}
	}

	return comment
}

// commentListFromAPI converts an API comments response to a domain comment list.
func commentListFromAPI(resp *api.CommentsResponse) *CommentList {
	list := &CommentList{}

	list.Comments = make([]CommentMeta, len(resp.Comments))
	for i, c := range resp.Comments {
		list.Comments[i] = commentMetaFromAPI(c)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}

// commentMetaFromAPI converts an API comment meta to a domain comment meta.
func commentMetaFromAPI(c api.CommentMeta) CommentMeta {
	meta := CommentMeta{}
	if v, ok := c.ID.Get(); ok {
		meta.ID = v
	}
	if v, ok := c.Body.Get(); ok {
		meta.Body = v
	}
	if v, ok := c.URL.Get(); ok {
		meta.URL = v
	}
	if v, ok := c.Resource.Get(); ok {
		meta.Resource = v
	}
	if v, ok := c.CreatedAt.Get(); ok {
		meta.CreatedAt = v
	}
	if v, ok := c.User.Get(); ok {
		meta.User = userFromAPI(v)
	}
	return meta
}

// commentableFromAPI converts an API commentable to a domain commentable.
func commentableFromAPI(c api.Commentable) *Commentable {
	commentable := &Commentable{}
	if v, ok := c.Type.Get(); ok {
		commentable.Type = v
	}
	if v, ok := c.ID.Get(); ok {
		commentable.ID = v
	}
	if v, ok := c.ProductID.Get(); ok {
		commentable.ProductID = v
	}
	if v, ok := c.URL.Get(); ok {
		commentable.URL = v
	}
	if v, ok := c.Resource.Get(); ok {
		commentable.Resource = v
	}
	return commentable
}

// attachmentFromAPI converts an API attachment to a domain attachment.
func attachmentFromAPI(a api.Attachment) Attachment {
	attachment := Attachment{}
	if v, ok := a.ID.Get(); ok {
		attachment.ID = v
	}
	if v, ok := a.DownloadURL.Get(); ok {
		attachment.DownloadURL = v
	}
	if v, ok := a.FileName.Get(); ok {
		attachment.FileName = v
	}
	if v, ok := a.FileSize.Get(); ok {
		attachment.FileSize = v
	}
	if v, ok := a.ContentType.Get(); ok {
		attachment.ContentType = v
	}
	return attachment
}
