package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	commentCreateFeatureID string
	commentCreateIdeaID    string
	commentCreateBody      string
)

var commentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new comment",
	Long: `Create a new comment on a feature or idea.

Examples:
  aha comment create --feature PROD-123 --body "Looks great!"
  aha comment create --idea IDEA-1 --body "Internal note about this idea"`,
	RunE: runCreateComment,
}

func init() {
	commentCmd.AddCommand(commentCreateCmd)

	commentCreateCmd.Flags().StringVar(&commentCreateFeatureID, "feature", "", "Feature ID to comment on")
	commentCreateCmd.Flags().StringVar(&commentCreateIdeaID, "idea", "", "Idea ID to comment on")
	commentCreateCmd.Flags().StringVarP(&commentCreateBody, "body", "b", "", "Comment body (required)")

	_ = commentCreateCmd.MarkFlagRequired("body")
}

func runCreateComment(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	opts := []aha.CreateCommentOption{
		aha.WithCommentBody(commentCreateBody),
	}

	var comment *aha.Comment
	var err error
	var resourceType string

	switch {
	case commentCreateFeatureID != "":
		comment, err = client.CreateFeatureComment(ctx, commentCreateFeatureID, opts...)
		resourceType = "feature " + commentCreateFeatureID
	case commentCreateIdeaID != "":
		comment, err = client.CreateIdeaComment(ctx, commentCreateIdeaID, opts...)
		resourceType = "idea " + commentCreateIdeaID
	default:
		return fmt.Errorf("specify one of: --feature, --idea")
	}

	if err != nil {
		return err
	}

	fmt.Printf("Created comment on %s\n", resourceType)
	fmt.Printf("  ID:      %s\n", comment.ID)
	fmt.Printf("  Created: %s\n", comment.CreatedAt.Format("2006-01-02 15:04"))

	return nil
}
