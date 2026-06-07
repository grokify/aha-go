package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/aha-go/graphql"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search Aha documents using GraphQL",
	Long:  `Search for pages, features, and other documents in Aha.io using the GraphQL API.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

var (
	searchType   string
	searchFormat string
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "Page", "Document type to search (Page, Feature, Idea, etc.)")
	searchCmd.Flags().StringVarP(&searchFormat, "format", "f", "table", "Output format (table, json)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Use global subdomain and apiKey from root.go (populated by PersistentPreRunE)
	gqlClient := graphql.NewClient(subdomain, apiKey)

	// Debug: print the endpoint
	fmt.Fprintf(os.Stderr, "GraphQL endpoint: %s\n", gqlClient.Endpoint())
	fmt.Fprintf(os.Stderr, "Searching for: %q (type: %s)\n", query, searchType)

	variables := map[string]any{
		"query":          query,
		"searchableType": []string{searchType},
	}

	var result graphql.SearchDocumentsResponse
	err := gqlClient.Query(context.Background(), graphql.SearchDocumentsQuery, variables, &result)
	if err != nil {
		return fmt.Errorf("GraphQL query failed: %w", err)
	}

	if searchFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Table format
	fmt.Printf("\nResults: %d (page %d of %d)\n\n", result.SearchDocuments.TotalCount, result.SearchDocuments.CurrentPage, result.SearchDocuments.TotalPages)

	if len(result.SearchDocuments.Nodes) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("%-15s %-40s %s\n", "TYPE", "NAME", "URL")
	fmt.Printf("%-15s %-40s %s\n", "----", "----", "---")
	for _, node := range result.SearchDocuments.Nodes {
		name := node.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		fmt.Printf("%-15s %-40s %s\n", node.SearchableType, name, node.URL)
	}

	return nil
}
