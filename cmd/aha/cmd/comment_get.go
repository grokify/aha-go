package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var commentGetCmd = &cobra.Command{
	Use:   "get <comment-id>",
	Short: "Get a comment by ID",
	Long: `Get details about a specific comment.

Examples:
  aha comment get 12345678
  aha comment get 12345678 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runGetComment,
}

func init() {
	commentCmd.AddCommand(commentGetCmd)
}

func runGetComment(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	commentID := args[0]

	comment, err := client.GetComment(ctx, commentID)
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
		return output.NewPrinter(format).Print(comment)
	}

	// Table output
	fmt.Printf("Comment: %s\n", comment.ID)
	if comment.User != nil {
		fmt.Printf("  Author:     %s\n", comment.User.Name())
	}
	fmt.Printf("  Created:    %s\n", comment.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("  Updated:    %s\n", comment.UpdatedAt.Format("2006-01-02 15:04"))

	if comment.Commentable != nil {
		fmt.Printf("  On:         %s (%s)\n", comment.Commentable.Type, comment.Commentable.ID)
	}

	if len(comment.Attachments) > 0 {
		fmt.Printf("  Attachments: %d\n", len(comment.Attachments))
	}

	fmt.Printf("\nBody:\n%s\n", comment.Body)

	return nil
}
