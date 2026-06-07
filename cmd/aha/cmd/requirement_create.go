package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	requirementCreateFeatureID   string
	requirementCreateName        string
	requirementCreateDescription string
	requirementCreateStatus      string
	requirementCreateAssignee    string
	requirementCreateEstimate    float64
)

var requirementCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new requirement",
	Long: `Create a new requirement for a feature.

Examples:
  aha requirement create --feature PROD-1 --name "Add input validation"
  aha requirement create -f PROD-1 -n "Add validation" --description "Validate all user inputs"
  aha requirement create -f PROD-1 -n "Add tests" --status "Not started" --estimate 5`,
	RunE: runCreateRequirement,
}

func init() {
	requirementCmd.AddCommand(requirementCreateCmd)

	requirementCreateCmd.Flags().StringVarP(&requirementCreateFeatureID, "feature", "f", "", "Feature ID or reference number (required)")
	requirementCreateCmd.Flags().StringVarP(&requirementCreateName, "name", "n", "", "Requirement name (required)")
	requirementCreateCmd.Flags().StringVarP(&requirementCreateDescription, "description", "d", "", "Requirement description")
	requirementCreateCmd.Flags().StringVarP(&requirementCreateStatus, "status", "s", "", "Workflow status")
	requirementCreateCmd.Flags().StringVarP(&requirementCreateAssignee, "assignee", "a", "", "Assigned user email")
	requirementCreateCmd.Flags().Float64Var(&requirementCreateEstimate, "estimate", 0, "Original estimate (hours/points)")

	_ = requirementCreateCmd.MarkFlagRequired("feature")
	_ = requirementCreateCmd.MarkFlagRequired("name")
}

func runCreateRequirement(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var opts []aha.CreateRequirementOption

	opts = append(opts, aha.WithRequirementName(requirementCreateName))

	if requirementCreateDescription != "" {
		opts = append(opts, aha.WithRequirementDescription(requirementCreateDescription))
	}
	if requirementCreateStatus != "" {
		opts = append(opts, aha.WithRequirementStatus(requirementCreateStatus))
	}
	if requirementCreateAssignee != "" {
		opts = append(opts, aha.WithRequirementAssignedTo(requirementCreateAssignee))
	}
	if requirementCreateEstimate > 0 {
		opts = append(opts, aha.WithRequirementEstimate(requirementCreateEstimate))
	}

	req, err := client.CreateRequirement(ctx, requirementCreateFeatureID, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("Created requirement: %s\n", req.ReferenceNum)
	fmt.Printf("  Name:    %s\n", req.Name)
	fmt.Printf("  ID:      %s\n", req.ID)
	fmt.Printf("  Feature: %s\n", requirementCreateFeatureID)
	fmt.Printf("  URL:     %s\n", req.URL)

	return nil
}
