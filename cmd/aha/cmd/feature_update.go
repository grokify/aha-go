package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	featureUpdateName        string
	featureUpdateDescription string
	featureUpdateStatus      string
	featureUpdateAssignee    string
	featureUpdateTags        string
	featureUpdateRelease     string
	featureUpdateInitiative  string
)

var featureUpdateCmd = &cobra.Command{
	Use:   "update <feature-id>",
	Short: "Update an existing feature",
	Long: `Update an existing feature's properties.

Examples:
  aha feature update PROD-123 --name "Updated Name"
  aha feature update PROD-123 --status "In Progress"
  aha feature update PROD-123 --assignee user@example.com
  aha feature update PROD-123 --release REL-2
  aha feature update PROD-123 --tags "priority,q4"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdateFeature,
}

func init() {
	featureCmd.AddCommand(featureUpdateCmd)

	featureUpdateCmd.Flags().StringVarP(&featureUpdateName, "name", "n", "", "Feature name")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateDescription, "description", "d", "", "Feature description")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateStatus, "status", "s", "", "Workflow status")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateAssignee, "assignee", "a", "", "Assignee email")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateTags, "tags", "t", "", "Tags (comma-separated)")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateRelease, "release", "r", "", "Move to release")
	featureUpdateCmd.Flags().StringVarP(&featureUpdateInitiative, "initiative", "i", "", "Initiative")
}

func runUpdateFeature(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	featureID := args[0]

	// Build update options based on provided flags
	var opts []aha.UpdateFeatureOption

	if featureUpdateName != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.Name = featureUpdateName })
	}
	if featureUpdateDescription != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.Description = featureUpdateDescription })
	}
	if featureUpdateStatus != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.WorkflowStatus = featureUpdateStatus })
	}
	if featureUpdateAssignee != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.AssignedToUser = featureUpdateAssignee })
	}
	if featureUpdateTags != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.Tags = featureUpdateTags })
	}
	if featureUpdateRelease != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.Release = featureUpdateRelease })
	}
	if featureUpdateInitiative != "" {
		opts = append(opts, func(o *aha.UpdateFeatureOptions) { o.Initiative = featureUpdateInitiative })
	}

	if len(opts) == 0 {
		return fmt.Errorf("no update options provided, use --name, --status, --assignee, etc")
	}

	feature, err := client.UpdateFeature(ctx, featureID, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("Feature updated successfully!\n\n")
	fmt.Printf("  Reference:  %s\n", feature.ReferenceNum)
	fmt.Printf("  Name:       %s\n", feature.Name)
	if feature.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", feature.WorkflowStatus.Name)
	}
	fmt.Printf("  URL:        %s\n", feature.URL)

	return nil
}
