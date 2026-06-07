package cmd

import (
	"fmt"

	"github.com/grokify/aha-go/browser"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/spf13/cobra"
)

var templateListPredefinedCmd = &cobra.Command{
	Use:   "list-predefined",
	Short: "List available predefined templates",
	Long: `List all predefined strategic model templates that can be created.

These templates are built-in configurations for common canvas types:
  - capability-stack: Layered capability model (4 layers)
  - maturity-model: Capability maturity assessment grid (5x5)
  - opportunity-patton: Jeff Patton's 10-block opportunity canvas
  - feature-canvas: Nikita Efimov's feature planning canvas

Examples:
  aha template list-predefined
  aha template list-predefined --output json`,
	RunE: runListPredefinedTemplates,
}

func init() {
	templateCmd.AddCommand(templateListPredefinedCmd)
}

type predefinedTemplateInfo struct {
	Name        string   `json:"name" yaml:"name"`
	DisplayName string   `json:"displayName" yaml:"displayName"`
	Description string   `json:"description" yaml:"description"`
	Rows        int      `json:"rows" yaml:"rows"`
	Columns     int      `json:"columns" yaml:"columns"`
	BlockCount  int      `json:"blockCount" yaml:"blockCount"`
	Blocks      []string `json:"blocks,omitempty" yaml:"blocks,omitempty"`
}

func runListPredefinedTemplates(cmd *cobra.Command, args []string) error {
	// Parse output format
	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	// Build template info list
	var templates []predefinedTemplateInfo
	for name, config := range browser.PredefinedTemplates {
		blockNames := make([]string, len(config.Blocks))
		for i, b := range config.Blocks {
			blockNames[i] = b.Name
		}
		templates = append(templates, predefinedTemplateInfo{
			Name:        name,
			DisplayName: config.Name,
			Description: config.Description,
			Rows:        config.Rows,
			Columns:     config.Columns,
			BlockCount:  len(config.Blocks),
			Blocks:      blockNames,
		})
	}

	// Handle structured output
	if format.IsStructured() {
		return output.NewPrinter(format).Print(templates)
	}

	// Table output
	fmt.Println("Available Predefined Templates:")
	fmt.Println()
	for _, t := range templates {
		fmt.Printf("  %s\n", t.Name)
		fmt.Printf("    Display Name: %s\n", t.DisplayName)
		fmt.Printf("    Description:  %s\n", t.Description)
		fmt.Printf("    Grid:         %dx%d (%d blocks)\n", t.Rows, t.Columns, t.BlockCount)
		fmt.Println()
	}

	return nil
}
