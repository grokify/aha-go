package cmd

import (
	"context"
	"fmt"

	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var releaseGetCmd = &cobra.Command{
	Use:   "get <release-id>",
	Short: "Get a release by ID or reference number",
	Long: `Get details about a specific release.

Examples:
  aha release get PROD-R-1
  aha release get 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runGetRelease,
}

func init() {
	releaseCmd.AddCommand(releaseGetCmd)
}

func runGetRelease(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	releaseID := args[0]

	release, err := client.GetRelease(ctx, releaseID)
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
		return output.NewPrinter(format).Print(release)
	}

	// Table output
	fmt.Printf("Release: %s\n", release.Name)
	fmt.Printf("  Reference:    %s\n", release.ReferenceNum)
	fmt.Printf("  ID:           %s\n", release.ID)

	if release.StartDate != nil {
		fmt.Printf("  Start Date:   %s\n", release.StartDate.Format("2006-01-02"))
	}
	if release.ReleaseDate != nil {
		fmt.Printf("  Release Date: %s\n", release.ReleaseDate.Format("2006-01-02"))
	}
	if release.ExternalReleaseDate != nil {
		fmt.Printf("  External:     %s\n", release.ExternalReleaseDate.Format("2006-01-02"))
	}

	if release.Released {
		fmt.Printf("  Status:       Released\n")
	} else if release.ParkingLot {
		fmt.Printf("  Status:       Parking Lot\n")
	} else {
		fmt.Printf("  Status:       Active\n")
	}

	fmt.Printf("  URL:          %s\n", release.URL)

	return nil
}
