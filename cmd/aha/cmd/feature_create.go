package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	featureCreateRelease     string
	featureCreateName        string
	featureCreateDescription string
	featureCreateStatus      string
	featureCreateAssignee    string
	featureCreateTags        string
)

var featureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new feature",
	Long: `Create a new feature in a release.

Examples:
  aha feature create --release REL-1 --name "User Authentication"
  aha feature create --release REL-1 --name "OAuth Support" --description "Add OAuth2 login"
  aha feature create --release REL-1 --name "Dashboard" --assignee user@example.com`,
	RunE: runCreateFeature,
}

func init() {
	featureCmd.AddCommand(featureCreateCmd)

	featureCreateCmd.Flags().StringVarP(&featureCreateRelease, "release", "r", "", "Release ID (required)")
	featureCreateCmd.Flags().StringVarP(&featureCreateName, "name", "n", "", "Feature name (required)")
	featureCreateCmd.Flags().StringVarP(&featureCreateDescription, "description", "d", "", "Feature description")
	featureCreateCmd.Flags().StringVarP(&featureCreateStatus, "status", "s", "", "Workflow status")
	featureCreateCmd.Flags().StringVarP(&featureCreateAssignee, "assignee", "a", "", "Assignee email")
	featureCreateCmd.Flags().StringVarP(&featureCreateTags, "tags", "t", "", "Tags (comma-separated)")

	_ = featureCreateCmd.MarkFlagRequired("release")
	_ = featureCreateCmd.MarkFlagRequired("name")
}

func runCreateFeature(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var opts []aha.CreateFeatureOption
	opts = append(opts, aha.WithFeatureName(featureCreateName))

	if featureCreateDescription != "" {
		opts = append(opts, aha.WithFeatureDescription(featureCreateDescription))
	}
	if featureCreateStatus != "" {
		opts = append(opts, aha.WithFeatureStatus(featureCreateStatus))
	}
	if featureCreateAssignee != "" {
		opts = append(opts, aha.WithFeatureAssignedTo(featureCreateAssignee))
	}
	if featureCreateTags != "" {
		opts = append(opts, aha.WithFeatureTags(featureCreateTags))
	}

	feature, err := client.CreateFeature(ctx, featureCreateRelease, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("Feature created successfully!\n\n")
	fmt.Printf("  Reference:  %s\n", feature.ReferenceNum)
	fmt.Printf("  Name:       %s\n", feature.Name)
	fmt.Printf("  URL:        %s\n", feature.URL)

	return nil
}
