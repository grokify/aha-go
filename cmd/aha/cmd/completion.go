package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for aha CLI.

To load completions:

Bash:
  $ source <(aha completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ aha completion bash > /etc/bash_completion.d/aha
  # macOS:
  $ aha completion bash > $(brew --prefix)/etc/bash_completion.d/aha

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ aha completion zsh > "${fpath[1]}/_aha"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ aha completion fish | source

  # To load completions for each session, execute once:
  $ aha completion fish > ~/.config/fish/completions/aha.fish

PowerShell:
  PS> aha completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> aha completion powershell > aha.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
