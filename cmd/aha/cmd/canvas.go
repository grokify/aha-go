package cmd

import (
	"github.com/spf13/cobra"
)

// canvasCmd represents the canvas command.
var canvasCmd = &cobra.Command{
	Use:   "canvas",
	Short: "Manage Aha strategic canvases",
	Long: `Manage strategic canvases (strategic models) in Aha.io.

Strategic canvases help teams align on product strategy by capturing
key aspects like problems, users, solutions, and metrics in a
structured format.

Supported canvas types:
  - opportunity   Opportunity Canvas (Jeff Patton's 10-block structure)
  - leanux        Lean UX Canvas (Jeff Gothelf's 8-block structure)
  - bmc           Business Model Canvas (Osterwalder's 9-block structure)`,
}

func init() {
	rootCmd.AddCommand(canvasCmd)
}
