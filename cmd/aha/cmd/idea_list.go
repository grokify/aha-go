package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var (
	ideaQuery  string
	ideaStatus string
	ideaSort   string
	ideaTag    string
)

var ideaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ideas",
	Long: `List ideas with optional filtering.

Sort options: recent, trending, popular

Examples:
  aha idea list
  aha idea list --query "dashboard"
  aha idea list --status "Under consideration"
  aha idea list --sort popular
  aha idea list --tag important`,
	RunE: runListIdeas,
}

func init() {
	ideaCmd.AddCommand(ideaListCmd)

	ideaListCmd.Flags().StringVarP(&ideaQuery, "query", "q", "", "Search query")
	ideaListCmd.Flags().StringVarP(&ideaStatus, "status", "s", "", "Filter by workflow status")
	ideaListCmd.Flags().StringVar(&ideaSort, "sort", "", "Sort order: recent, trending, popular")
	ideaListCmd.Flags().StringVarP(&ideaTag, "tag", "t", "", "Filter by tag")
	ideaListCmd.Flags().IntVar(&page, "page", 1, "Page number")
	ideaListCmd.Flags().IntVar(&perPage, "per-page", 30, "Results per page")
}

func runListIdeas(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var opts []aha.ListIdeasOption
	if ideaQuery != "" {
		opts = append(opts, aha.WithIdeaQuery(ideaQuery))
	}
	if ideaStatus != "" {
		opts = append(opts, aha.WithIdeaStatus(ideaStatus))
	}
	if ideaSort != "" {
		opts = append(opts, aha.WithIdeaSort(ideaSort))
	}
	if ideaTag != "" {
		opts = append(opts, aha.WithIdeaTag(ideaTag))
	}

	list, err := client.ListIdeas(ctx, opts...)
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
	if len(list.Ideas) == 0 {
		fmt.Println("No ideas found.")
		return nil
	}

	fmt.Printf("Ideas (%d total):\n\n", list.Pagination.TotalRecords)
	for _, i := range list.Ideas {
		status := ""
		if i.WorkflowStatus != nil {
			status = fmt.Sprintf(" [%s]", i.WorkflowStatus.Name)
		}
		fmt.Printf("  %s  %s (%d votes)%s\n", i.ReferenceNum, i.Name, i.Votes, status)
	}

	if list.Pagination.TotalPages > 1 {
		fmt.Printf("\nPage %d of %d\n", list.Pagination.CurrentPage, list.Pagination.TotalPages)
	}

	return nil
}
