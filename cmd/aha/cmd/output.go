package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	aha "github.com/grokify/aha-go"
)

// printCanvasCreated prints the created canvas details.
func printCanvasCreated(typeName string, sm *aha.StrategicModel) {
	fmt.Printf("%s created successfully!\n\n", typeName)
	printCanvasDetails(sm)
}

// printCanvasDetails prints canvas details in tabular format.
func printCanvasDetails(sm *aha.StrategicModel) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "ID:\t%s\n", sm.ID)
	_, _ = fmt.Fprintf(w, "Reference:\t%s\n", sm.ReferenceNum)
	_, _ = fmt.Fprintf(w, "Name:\t%s\n", sm.Name)
	_, _ = fmt.Fprintf(w, "Kind:\t%s\n", sm.Kind)
	_, _ = fmt.Fprintf(w, "URL:\t%s\n", sm.URL)
	_ = w.Flush()

	if len(sm.Components) > 0 {
		fmt.Println()
		fmt.Println("Components (blocks):")
		for _, comp := range sm.Components {
			fmt.Printf("  - %s (ID: %s)\n", comp.Name, comp.ID)
		}
	}
}
