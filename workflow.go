package aha

import (
	"context"

	"github.com/grokify/aha-go/internal/api"
)

// Workflow represents a workflow with its statuses in Aha.
// WorkflowStatus is defined in feature.go.
type Workflow struct {
	ID       string
	Name     string
	Statuses []WorkflowStatus
}

// WorkflowList represents a list of workflows.
type WorkflowList struct {
	Workflows []Workflow
}

// ListProductWorkflows lists all workflows and their statuses for a product.
func (c *Client) ListProductWorkflows(ctx context.Context, productID string) (*WorkflowList, error) {
	params := api.ListProductWorkflowsParams{
		ProductID: productID,
	}

	resp, err := c.apiClient.ListProductWorkflows(ctx, params)
	if err != nil {
		return nil, wrapError("ListProductWorkflows", err)
	}

	return workflowListFromAPI(resp), nil
}

// workflowListFromAPI converts an API workflows response to a WorkflowList.
func workflowListFromAPI(resp *api.WorkflowsResponse) *WorkflowList {
	list := &WorkflowList{
		Workflows: make([]Workflow, len(resp.Workflows)),
	}

	for i, w := range resp.Workflows {
		list.Workflows[i] = workflowFromAPI(w)
	}

	return list
}

// workflowFromAPI converts an API workflow to a Workflow.
func workflowFromAPI(w api.Workflow) Workflow {
	workflow := Workflow{
		ID:   w.ID,
		Name: w.Name,
	}

	if len(w.WorkflowStatuses) > 0 {
		workflow.Statuses = make([]WorkflowStatus, len(w.WorkflowStatuses))
		for i, s := range w.WorkflowStatuses {
			workflow.Statuses[i] = workflowStatusFromAPIStatus(s)
		}
	}

	return workflow
}

// workflowStatusFromAPIStatus converts an API WorkflowStatus to a WorkflowStatus.
func workflowStatusFromAPIStatus(s api.WorkflowStatus) WorkflowStatus {
	status := WorkflowStatus{}

	if v, ok := s.ID.Get(); ok {
		status.ID = v
	}
	if v, ok := s.Name.Get(); ok {
		status.Name = v
	}
	if v, ok := s.Position.Get(); ok {
		status.Position = v
	}
	if v, ok := s.Complete.Get(); ok {
		status.Complete = v
	}
	if v, ok := s.Color.Get(); ok {
		status.Color = v
	}

	return status
}
