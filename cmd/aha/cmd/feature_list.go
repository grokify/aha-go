package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	featureQuery    string
	featureAssignee string
	featureTag      string
	featureRelease  string
)

var featureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List features",
	Long: `List features with optional filtering.

Examples:
  aha feature list
  aha feature list --query "login"
  aha feature list --assignee user@example.com
  aha feature list --tag important
  aha feature list --release REL-1
  aha feature list --output json`,
	RunE: runListFeatures,
}

func init() {
	featureCmd.AddCommand(featureListCmd)

	featureListCmd.Flags().StringVarP(&featureQuery, "query", "q", "", "Search query")
	featureListCmd.Flags().StringVarP(&featureAssignee, "assignee", "a", "", "Filter by assignee email")
	featureListCmd.Flags().StringVarP(&featureTag, "tag", "t", "", "Filter by tag")
	featureListCmd.Flags().StringVarP(&featureRelease, "release", "r", "", "Filter by release ID")
	featureListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	featureListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListFeatures(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var list *aha.FeatureList
	var err error

	if featureRelease != "" {
		// List features in a specific release
		list, err = client.ListReleaseFeatures(ctx, featureRelease)
	} else {
		// List all features with filters
		var opts []aha.ListFeaturesOption
		if featureQuery != "" {
			opts = append(opts, aha.WithFeatureQuery(featureQuery))
		}
		if featureAssignee != "" {
			opts = append(opts, aha.WithFeatureAssignee(featureAssignee))
		}
		if featureTag != "" {
			opts = append(opts, aha.WithFeatureTag(featureTag))
		}
		list, err = client.ListFeatures(ctx, opts...)
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
	if len(list.Features) == 0 {
		fmt.Println("No features found.")
		return nil
	}

	fmt.Printf("Features (%d total):\n\n", list.Pagination.TotalRecords)
	for _, f := range list.Features {
		fmt.Printf("  %s  %s\n", f.ReferenceNum, f.Name)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
