package cmd

import (
	"context"
	"fmt"

	aha "github.com/grokify/aha-go"
	"github.com/spf13/cobra"
)

var (
	ideaUpdateName        string
	ideaUpdateDescription string
	ideaUpdateStatus      string
	ideaUpdateVisibility  string
)

var ideaUpdateCmd = &cobra.Command{
	Use:   "update <idea-id>",
	Short: "Update an existing idea",
	Long: `Update an existing idea's properties.

Examples:
  aha idea update IDEA-123 --name "Updated Name"
  aha idea update IDEA-123 --status "Under consideration"
  aha idea update IDEA-123 --description "New description"
  aha idea update IDEA-123 --visibility public`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdateIdea,
}

func init() {
	ideaCmd.AddCommand(ideaUpdateCmd)

	ideaUpdateCmd.Flags().StringVarP(&ideaUpdateName, "name", "n", "", "Idea name")
	ideaUpdateCmd.Flags().StringVarP(&ideaUpdateDescription, "description", "d", "", "Idea description")
	ideaUpdateCmd.Flags().StringVarP(&ideaUpdateStatus, "status", "s", "", "Workflow status")
	ideaUpdateCmd.Flags().StringVar(&ideaUpdateVisibility, "visibility", "", "Visibility (public, private)")
}

func runUpdateIdea(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ideaID := args[0]

	var opts []aha.UpdateIdeaOption

	if ideaUpdateName != "" {
		opts = append(opts, aha.WithUpdateIdeaName(ideaUpdateName))
	}
	if ideaUpdateDescription != "" {
		opts = append(opts, aha.WithUpdateIdeaDescription(ideaUpdateDescription))
	}
	if ideaUpdateStatus != "" {
		opts = append(opts, aha.WithUpdateIdeaStatus(ideaUpdateStatus))
	}
	if ideaUpdateVisibility != "" {
		opts = append(opts, aha.WithUpdateIdeaVisibility(ideaUpdateVisibility))
	}

	if len(opts) == 0 {
		return fmt.Errorf("no update options provided, use --name, --status, --description, etc")
	}

	idea, err := client.UpdateIdea(ctx, ideaID, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("Idea updated successfully!\n\n")
	fmt.Printf("  Reference:  %s\n", idea.ReferenceNum)
	fmt.Printf("  Name:       %s\n", idea.Name)
	if idea.WorkflowStatus != nil {
		fmt.Printf("  Status:     %s\n", idea.WorkflowStatus.Name)
	}
	fmt.Printf("  Votes:      %d\n", idea.Votes)

	return nil
}
