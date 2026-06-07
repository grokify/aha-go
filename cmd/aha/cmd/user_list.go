package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long: `List all users in the Aha account.

Examples:
  aha user list
  aha user list --output json`,
	RunE: runListUsers,
}

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	userListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListUsers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	list, err := client.ListUsers(ctx)
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
	if len(list.Users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	fmt.Printf("Users (%d total):\n\n", list.Pagination.TotalRecords)
	for _, u := range list.Users {
		role := ""
		if u.Role != "" {
			role = fmt.Sprintf(" [%s]", u.Role)
		}
		fmt.Printf("  %s  %s%s\n", u.Email, u.Name(), role)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
