package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandsRegistered(t *testing.T) {
	// Verify root command exists
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}

	// Expected subcommands
	expectedCommands := []string{
		"product",
		"feature",
		"release",
		"idea",
		"goal",
		"initiative",
		"epic",
		"requirement",
		"canvas",
		"completion",
	}

	commands := rootCmd.Commands()
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("Expected command %q not found in root command", expected)
		}
	}
}

func TestProductSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, productCmd, "product", expectedSubcommands)
}

func TestFeatureSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get", "create", "update"}
	verifySubcommands(t, featureCmd, "feature", expectedSubcommands)
}

func TestReleaseSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, releaseCmd, "release", expectedSubcommands)
}

func TestIdeaSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, ideaCmd, "idea", expectedSubcommands)
}

func TestGoalSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, goalCmd, "goal", expectedSubcommands)
}

func TestInitiativeSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, initiativeCmd, "initiative", expectedSubcommands)
}

func TestEpicSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get"}
	verifySubcommands(t, epicCmd, "epic", expectedSubcommands)
}

func TestRequirementSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list", "get", "create", "update", "delete"}
	verifySubcommands(t, requirementCmd, "requirement", expectedSubcommands)
}

func verifySubcommands(t *testing.T, parent *cobra.Command, parentName string, expected []string) {
	t.Helper()

	if parent == nil {
		t.Fatalf("%s command is nil", parentName)
	}

	commands := parent.Commands()
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, exp := range expected {
		if !commandMap[exp] {
			t.Errorf("Expected subcommand %q not found in %s command", exp, parentName)
		}
	}
}

func TestRootCommandFlags(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	tests := []struct {
		name     string
		flagName string
	}{
		{"api-key flag", "api-key"},
		{"subdomain flag", "subdomain"},
		{"verbose flag", "verbose"},
		{"output flag", "output"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %q not found in root command", tt.flagName)
			}
		})
	}
}

func TestCompletionValidArgs(t *testing.T) {
	if completionCmd == nil {
		t.Fatal("completionCmd is nil")
	}

	expectedArgs := []string{"bash", "zsh", "fish", "powershell"}

	for _, expected := range expectedArgs {
		found := false
		for _, arg := range completionCmd.ValidArgs {
			if arg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected valid arg %q not found in completion command", expected)
		}
	}
}
