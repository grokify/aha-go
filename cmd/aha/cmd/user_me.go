package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get the current authenticated user",
	Long: `Get details about the currently authenticated user.

Examples:
  aha user me
  aha user me --output json`,
	RunE: runGetCurrentUser,
}

func init() {
	userCmd.AddCommand(userMeCmd)
}

func runGetCurrentUser(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	user, err := client.GetCurrentUser(ctx)
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
	fmt.Printf("Current User: %s\n", user.Name())
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
