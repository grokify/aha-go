package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/aha-go/browser"
	"github.com/spf13/cobra"
)

var templateHeadless bool

var templateCreateCmd = &cobra.Command{
	Use:   "create-predefined <template-name>",
	Short: "Create a predefined strategic model template",
	Long: `Create a predefined strategic model template in Aha via browser automation.

Available predefined templates:
  - capability-stack: Layered capability model (4 layers)
  - maturity-model: Capability maturity assessment grid (5x5)
  - opportunity-patton: Jeff Patton's 10-block opportunity canvas
  - feature-canvas: Nikita Efimov's feature planning canvas

Requires browser credentials:
  - AHA_SUBDOMAIN: Your Aha account subdomain
  - AHA_EMAIL: Your Aha login email
  - AHA_PASSWORD: Your Aha login password

Examples:
  aha template create-predefined capability-stack
  aha template create-predefined feature-canvas --no-headless`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: browser.ListPredefinedTemplates(),
	RunE:      runCreatePredefinedTemplate,
}

func init() {
	templateCmd.AddCommand(templateCreateCmd)
	templateCreateCmd.Flags().BoolVar(&templateHeadless, "headless", true, "Run browser in headless mode")
}

func runCreatePredefinedTemplate(cmd *cobra.Command, args []string) error {
	templateKey := strings.ToLower(args[0])

	// Validate template exists
	templates := browser.ListPredefinedTemplates()
	found := false
	for _, t := range templates {
		if t == templateKey {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown template %q, available: %s", templateKey, strings.Join(templates, ", "))
	}

	// Create browser client
	browserClient, err := browser.NewClient(
		browser.WithHeadless(templateHeadless),
	)
	if err != nil {
		return fmt.Errorf("failed to create browser client: %w", err)
	}

	ctx := context.Background()

	// Connect browser
	fmt.Println("Launching browser...")
	if err := browserClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect browser: %w", err)
	}
	defer func() { _ = browserClient.Close() }()

	// Login
	fmt.Println("Logging in to Aha...")
	if err := browserClient.Login(ctx); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// Create template
	config := browser.PredefinedTemplates[templateKey]
	fmt.Printf("Creating template %q...\n", config.Name)
	if err := browserClient.CreateTemplate(ctx, config); err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	fmt.Printf("Successfully created template %q\n", config.Name)
	fmt.Println()
	fmt.Println("You can now create strategic models using this template via the API:")
	fmt.Printf("  aha canvas create --product PROD --name \"My Canvas\" --kind %q\n", config.Name)

	return nil
}
