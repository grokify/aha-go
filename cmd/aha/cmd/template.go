package cmd

import (
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage strategic model templates (browser automation)",
	Long: `Manage strategic model templates via browser automation.

Since Aha does not provide APIs for creating strategic model templates,
this command uses browser automation (go-rod) to create them through the UI.

Requires browser credentials:
  - AHA_SUBDOMAIN: Your Aha account subdomain
  - AHA_EMAIL: Your Aha login email
  - AHA_PASSWORD: Your Aha login password

Examples:
  aha template list-predefined
  aha template create --name "Capability Stack"
  aha template create-predefined capability-stack`,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
