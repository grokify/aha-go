package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var initiativeGetCmd = &cobra.Command{
	Use:   "get <initiative-id>",
	Short: "Get an initiative by ID or reference number",
	Long: `Get details about a specific initiative.

Examples:
  aha initiative get PROD-I-1
  aha initiative get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetInitiative,
}

func init() {
	initiativeCmd.AddCommand(initiativeGetCmd)
}

func runGetInitiative(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	initiativeID := args[0]

	initiative, err := client.GetInitiative(ctx, initiativeID)
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
		return output.NewPrinter(format).Print(initiative)
	}

	// Table output
	fmt.Printf("Initiative: %s\n", initiative.Name)
	fmt.Printf("  Reference:  %s\n", initiative.ReferenceNum)
	fmt.Printf("  ID:         %s\n", initiative.ID)

	if initiative.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", initiative.WorkflowStatus.Name)
	}
	fmt.Printf("  Progress:   %.0f%%\n", initiative.Progress*100)

	if initiative.Value > 0 {
		fmt.Printf("  Value:      %.0f\n", initiative.Value)
	}
	if initiative.Effort > 0 {
		fmt.Printf("  Effort:     %.0f\n", initiative.Effort)
	}

	if initiative.StartDate != nil {
		fmt.Printf("  Start:      %s\n", initiative.StartDate.Format("2006-01-02"))
	}
	if initiative.EndDate != nil {
		fmt.Printf("  End:        %s\n", initiative.EndDate.Format("2006-01-02"))
	}

	if initiative.Epic != nil {
		fmt.Printf("  Epic:       %s (%s)\n", initiative.Epic.Name, initiative.Epic.ReferenceNum)
	}

	if len(initiative.Features) > 0 {
		fmt.Printf("  Features:   %d\n", len(initiative.Features))
	}

	fmt.Printf("  URL:        %s\n", initiative.URL)
	fmt.Printf("  Created:    %s\n", initiative.CreatedAt.Format("2006-01-02"))

	if initiative.Description != "" && verbose {
		fmt.Printf("\nDescription:\n%s\n", initiative.Description)
	}

	return nil
}
