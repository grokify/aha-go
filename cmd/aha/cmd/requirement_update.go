package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	requirementUpdateName        string
	requirementUpdateDescription string
	requirementUpdateStatus      string
	requirementUpdateAssignee    string
	requirementUpdateWorkDone    float64
)

var requirementUpdateCmd = &cobra.Command{
	Use:   "update <requirement-id>",
	Short: "Update an existing requirement",
	Long: `Update a requirement by ID or reference number.

Examples:
  aha requirement update PROD-1-1 --name "Updated requirement name"
  aha requirement update PROD-1-1 --status "In progress"
  aha requirement update PROD-1-1 --work-done 4`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdateRequirement,
}

func init() {
	requirementCmd.AddCommand(requirementUpdateCmd)

	requirementUpdateCmd.Flags().StringVarP(&requirementUpdateName, "name", "n", "", "Requirement name")
	requirementUpdateCmd.Flags().StringVarP(&requirementUpdateDescription, "description", "d", "", "Requirement description")
	requirementUpdateCmd.Flags().StringVarP(&requirementUpdateStatus, "status", "s", "", "Workflow status")
	requirementUpdateCmd.Flags().StringVarP(&requirementUpdateAssignee, "assignee", "a", "", "Assigned user email")
	requirementUpdateCmd.Flags().Float64Var(&requirementUpdateWorkDone, "work-done", 0, "Work done (hours/points)")
}

func runUpdateRequirement(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	requirementID := args[0]

	var opts []aha.UpdateRequirementOption

	if requirementUpdateName != "" {
		opts = append(opts, aha.WithUpdateRequirementName(requirementUpdateName))
	}
	if requirementUpdateDescription != "" {
		opts = append(opts, aha.WithUpdateRequirementDescription(requirementUpdateDescription))
	}
	if requirementUpdateStatus != "" {
		opts = append(opts, aha.WithUpdateRequirementStatus(requirementUpdateStatus))
	}
	if requirementUpdateWorkDone > 0 {
		opts = append(opts, aha.WithUpdateRequirementWorkDone(requirementUpdateWorkDone))
	}

	if len(opts) == 0 {
		return fmt.Errorf("at least one field to update must be specified")
	}

	req, err := client.UpdateRequirement(ctx, requirementID, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("Updated requirement: %s\n", req.ReferenceNum)
	fmt.Printf("  Name:   %s\n", req.Name)
	if req.WorkflowStatus != nil {
		fmt.Printf("  Status: %s\n", req.WorkflowStatus.Name)
	}
	fmt.Printf("  URL:    %s\n", req.URL)

	return nil
}
