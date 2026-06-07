package cmd

import (
	"context"
	"fmt"
	"strings"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	commentFeatureID    string
	commentIdeaID       string
	commentReleaseID    string
	commentProductID    string
	commentInitiativeID string
	commentEpicID       string
	commentGoalID       string
)

var commentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments",
	Long: `List comments for a feature, idea, release, or other resource.

You must specify one of the resource flags.

Examples:
  aha comment list --feature PROD-123
  aha comment list --idea IDEA-1
  aha comment list --release PROD-R-1
  aha comment list --product PROD
  aha comment list --output json`,
	RunE: runListComments,
}

func init() {
	commentCmd.AddCommand(commentListCmd)

	commentListCmd.Flags().StringVar(&commentFeatureID, "feature", "", "Feature ID")
	commentListCmd.Flags().StringVar(&commentIdeaID, "idea", "", "Idea ID")
	commentListCmd.Flags().StringVar(&commentReleaseID, "release", "", "Release ID")
	commentListCmd.Flags().StringVar(&commentProductID, "product", "", "Product ID")
	commentListCmd.Flags().StringVar(&commentInitiativeID, "initiative", "", "Initiative ID")
	commentListCmd.Flags().StringVar(&commentEpicID, "epic", "", "Epic ID")
	commentListCmd.Flags().StringVar(&commentGoalID, "goal", "", "Goal ID")
	commentListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	commentListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListComments(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var list *aha.CommentList
	var err error
	var resourceType string

	switch {
	case commentFeatureID != "":
		list, err = client.ListFeatureComments(ctx, commentFeatureID)
		resourceType = "feature " + commentFeatureID
	case commentIdeaID != "":
		list, err = client.ListIdeaComments(ctx, commentIdeaID)
		resourceType = "idea " + commentIdeaID
	case commentReleaseID != "":
		list, err = client.ListReleaseComments(ctx, commentReleaseID)
		resourceType = "release " + commentReleaseID
	case commentProductID != "":
		list, err = client.ListProductComments(ctx, commentProductID)
		resourceType = "product " + commentProductID
	case commentInitiativeID != "":
		list, err = client.ListInitiativeComments(ctx, commentInitiativeID)
		resourceType = "initiative " + commentInitiativeID
	case commentEpicID != "":
		list, err = client.ListEpicComments(ctx, commentEpicID)
		resourceType = "epic " + commentEpicID
	case commentGoalID != "":
		list, err = client.ListGoalComments(ctx, commentGoalID)
		resourceType = "goal " + commentGoalID
	default:
		return fmt.Errorf("specify one of: --feature, --idea, --release, --product, --initiative, --epic, --goal")
	}

	if err != nil {
		return err
	}

	// Parse output format
	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	// Handle structured output
	if format.IsStructured() {
		return output.NewPrinter(format).Print(list)
	}

	// Table output
	if len(list.Comments) == 0 {
		fmt.Printf("No comments found for %s.\n", resourceType)
		return nil
	}

	fmt.Printf("Comments for %s (%d total):\n\n", resourceType, list.Pagination.TotalRecords)
	for _, c := range list.Comments {
		user := "Unknown"
		if c.User != nil {
			user = c.User.Name()
		}
		body := c.Body
		if len(body) > 60 {
			body = body[:57] + "..."
		}
		body = strings.ReplaceAll(body, "\n", " ")
		fmt.Printf("  %s  %s: %s\n", c.ID, user, body)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
