// Command aha provides a CLI for interacting with Aha.io.
package main

import (
	"os"

	"github.com/grokify/aha-go/cmd/aha/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
