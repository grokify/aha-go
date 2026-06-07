package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var commentDeleteForce bool

var commentDeleteCmd = &cobra.Command{
	Use:   "delete <comment-id>",
	Short: "Delete a comment",
	Long: `Delete a comment by ID.

Examples:
  aha comment delete 12345678
  aha comment delete 12345678 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDeleteComment,
}

func init() {
	commentCmd.AddCommand(commentDeleteCmd)

	commentDeleteCmd.Flags().BoolVar(&commentDeleteForce, "force", false, "Skip confirmation prompt")
}

func runDeleteComment(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	commentID := args[0]

	if !commentDeleteForce {
		fmt.Printf("Are you sure you want to delete comment %s? [y/N]: ", commentID)
		var confirm string
		_, err := fmt.Scanln(&confirm)
		if err != nil || (confirm != "y" && confirm != "Y") {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := client.DeleteComment(ctx, commentID); err != nil {
		return err
	}

	fmt.Printf("Deleted comment: %s\n", commentID)
	return nil
}
