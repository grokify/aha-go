package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var userGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get a user by ID or email",
	Long: `Get details about a specific user.

Examples:
  aha user get user@example.com
  aha user get 12345678
  aha user get 12345678 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runGetUser,
}

func init() {
	userCmd.AddCommand(userGetCmd)
}

func runGetUser(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	userID := args[0]

	user, err := client.GetUser(ctx, userID)
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
		return output.NewPrinter(format).Print(user)
	}

	// Table output
	fmt.Printf("User: %s\n", user.Name())
	fmt.Printf("  ID:         %s\n", user.ID)
	fmt.Printf("  Email:      %s\n", user.Email)
	fmt.Printf("  First Name: %s\n", user.FirstName)
	fmt.Printf("  Last Name:  %s\n", user.LastName)
	if user.Role != "" {
		fmt.Printf("  Role:       %s\n", user.Role)
	}
	if user.CreatedAt != nil {
		fmt.Printf("  Created:    %s\n", user.CreatedAt.Format("2006-01-02"))
	}

	return nil
}
