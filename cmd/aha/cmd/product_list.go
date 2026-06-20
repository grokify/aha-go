package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	genql "github.com/Khan/genqlient/graphql"
	aha "github.com/grokify/aha-go"
	"github.com/grokify/aha-go/cmd/aha/internal/output"
	"github.com/grokify/aha-go/graphql"
	"github.com/grokify/aha-go/graphql/generated"
	"github.com/spf13/cobra"
)

var (
	// Global pagination variables used by multiple list commands
	page    int
	perPage int

	// Product-specific filter variables
	productUpdatedSince   string
	productWithIdeaPortal bool
	productShowTree       bool
)

// ProductNode represents a product in the tree structure (JSON IR).
type ProductNode struct {
	ID              string         `json:"id"`
	ReferencePrefix string         `json:"reference_prefix"`
	Name            string         `json:"name"`
	ProductLine     bool           `json:"product_line"`
	WorkspaceType   string         `json:"workspace_type,omitempty"`
	Children        []*ProductNode `json:"children,omitempty"`
}

var productListCmd = &cobra.Command{
	Use:   "list",
	Short: "List products (workspaces)",
	Long: `List all products (workspaces) you have access to.

Products are displayed as a tree showing product lines and their children.
Output is sorted alphabetically by reference prefix.

Examples:
  aha product list
  aha product list --with-idea-portals
  aha product list --tree
  aha product list --output json`,
	RunE: runListProducts,
}

func init() {
	productCmd.AddCommand(productListCmd)

	productListCmd.Flags().StringVar(&productUpdatedSince, "updated-since", "", "Filter to products updated after this date (ISO8601 or YYYY-MM-DD)")
	productListCmd.Flags().BoolVar(&productWithIdeaPortal, "with-idea-portals", false, "Only list products with idea portals")
	productListCmd.Flags().BoolVar(&productShowTree, "tree", false, "Show as tree structure (default for table output)")
}

func runListProducts(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse output format early
	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	// Try GraphQL first for tree structure (has parent info)
	tree, err := fetchProductsGraphQL(ctx)
	if err != nil {
		// Fall back to REST API if GraphQL fails
		allProducts, restErr := fetchAllProducts(ctx)
		if restErr != nil {
			return fmt.Errorf("GraphQL: %w; REST: %w", err, restErr)
		}

		if len(allProducts) == 0 {
			fmt.Println("No products found.")
			return nil
		}

		// Build flat tree from REST (no parent info)
		tree = buildProductTreeFromREST(allProducts)
	}

	if len(tree) == 0 {
		fmt.Println("No products found.")
		return nil
	}

	// Handle structured output - output the tree as JSON/YAML
	if format.IsStructured() {
		return output.NewPrinter(format).Print(tree)
	}

	// Text output - render tree with numbering
	printProductTree(tree)

	return nil
}

// fetchProductsGraphQL fetches all products via GraphQL with parent info for tree building.
// It recursively fetches children for each product line to build a complete tree.
func fetchProductsGraphQL(ctx context.Context) ([]*ProductNode, error) {
	gqlClient := graphql.NewGenqlientClient(client.Subdomain(), client.APIKey())

	// Fetch root products (no parent)
	roots, err := fetchProjectsPage(ctx, gqlClient, nil)
	if err != nil {
		return nil, err
	}

	// Build tree nodes and recursively fetch children
	var tree []*ProductNode
	for _, p := range roots {
		node := projectToNode(p)
		// Recursively fetch children for product lines
		if err := fetchChildrenRecursive(ctx, gqlClient, node, p.Id); err != nil {
			return nil, err
		}
		tree = append(tree, node)
	}

	// Sort the tree
	sortNodes(tree)

	return tree, nil
}

// fetchProjectsPage fetches all projects for a given parent (nil for roots).
func fetchProjectsPage(ctx context.Context, gqlClient genql.Client, parentID *string) ([]generated.ListProjectsProjectsProjectPageNodesProject, error) {
	var allProjects []generated.ListProjectsProjectsProjectPageNodesProject
	currentPage := 1
	pageSize := 100

	for {
		resp, err := generated.ListProjects(ctx, gqlClient, &currentPage, &pageSize, nil, parentID)
		if err != nil {
			return nil, fmt.Errorf("GraphQL ListProjects: %w", err)
		}

		allProjects = append(allProjects, resp.Projects.Nodes...)

		if resp.Projects.IsLastPage || currentPage >= resp.Projects.TotalPages {
			break
		}
		currentPage++
	}

	return allProjects, nil
}

// fetchChildrenRecursive fetches children for a node and recursively builds the subtree.
func fetchChildrenRecursive(ctx context.Context, gqlClient genql.Client, node *ProductNode, parentID string) error {
	children, err := fetchProjectsPage(ctx, gqlClient, &parentID)
	if err != nil {
		return err
	}

	for _, child := range children {
		childNode := projectToNode(child)
		// Recursively fetch children of this child
		if err := fetchChildrenRecursive(ctx, gqlClient, childNode, child.Id); err != nil {
			return err
		}
		node.Children = append(node.Children, childNode)
	}

	return nil
}

// projectToNode converts a GraphQL project to a ProductNode.
func projectToNode(p generated.ListProjectsProjectsProjectPageNodesProject) *ProductNode {
	return &ProductNode{
		ID:              p.Id,
		ReferencePrefix: p.ReferencePrefix,
		Name:            p.Name,
		ProductLine:     p.WorkspaceType == "product_line" || p.WorkspaceType == "portfolio",
		WorkspaceType:   p.WorkspaceType,
	}
}

// fetchAllProducts fetches all products across all pages via REST API.
func fetchAllProducts(ctx context.Context) ([]aha.ProductMeta, error) {
	var allProducts []aha.ProductMeta
	currentPage := 1
	pageSize := 100 // Fetch in larger batches

	for {
		opts := []aha.ListProductsOption{
			aha.WithProductsPage(currentPage),
			aha.WithProductsPerPage(pageSize),
		}

		if productWithIdeaPortal {
			opts = append(opts, aha.WithIdeaPortals())
		}
		if productUpdatedSince != "" {
			t, err := parseTime(productUpdatedSince)
			if err != nil {
				return nil, fmt.Errorf("invalid --updated-since format: %w", err)
			}
			opts = append(opts, aha.WithUpdatedSince(t))
		}

		list, err := client.ListProducts(ctx, opts...)
		if err != nil {
			return nil, err
		}

		allProducts = append(allProducts, list.Products...)

		// Check if we've fetched all pages
		if currentPage >= int(list.Pagination.TotalPages) {
			break
		}
		currentPage++
	}

	return allProducts, nil
}

// parseTime parses a time string in ISO8601 or YYYY-MM-DD format.
func parseTime(s string) (time.Time, error) {
	// Try ISO8601 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try date only
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("expected format: YYYY-MM-DD or ISO8601")
}

// sortNodes recursively sorts nodes and their children alphabetically.
func sortNodes(nodes []*ProductNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return strings.ToLower(nodes[i].ReferencePrefix) < strings.ToLower(nodes[j].ReferencePrefix)
	})
	for _, node := range nodes {
		if len(node.Children) > 0 {
			sortNodes(node.Children)
		}
	}
}

// buildProductTreeFromREST builds a flat tree from REST API response (no parent info).
// Since ProductMeta doesn't have ParentID, we use ProductLine flag
// to identify top-level product lines. Products without a parent
// in the list are shown at root level.
func buildProductTreeFromREST(products []aha.ProductMeta) []*ProductNode {
	// First, sort all products alphabetically by reference prefix
	sortedProducts := make([]aha.ProductMeta, len(products))
	copy(sortedProducts, products)
	sort.Slice(sortedProducts, func(i, j int) bool {
		return strings.ToLower(sortedProducts[i].ReferencePrefix) < strings.ToLower(sortedProducts[j].ReferencePrefix)
	})

	// Convert to nodes - since we don't have parent info in ProductMeta,
	// we show product lines first, then regular products
	var productLines []*ProductNode
	var regularProducts []*ProductNode

	for _, p := range sortedProducts {
		node := &ProductNode{
			ID:              p.ID,
			ReferencePrefix: p.ReferencePrefix,
			Name:            p.Name,
			ProductLine:     p.ProductLine,
			WorkspaceType:   p.WorkspaceType,
		}

		if p.ProductLine {
			productLines = append(productLines, node)
		} else {
			regularProducts = append(regularProducts, node)
		}
	}

	// Combine: product lines first, then regular products
	result := append(productLines, regularProducts...)
	return result
}

// printProductTree prints the product tree with numbering.
func printProductTree(nodes []*ProductNode) {
	fmt.Printf("Products (%d total):\n\n", countNodes(nodes))

	for i, node := range nodes {
		printNode(node, i+1, "")
	}
}

// printNode prints a single node with proper formatting.
func printNode(node *ProductNode, num int, indent string) {
	typeIndicator := ""
	if node.ProductLine {
		typeIndicator = " [product line]"
	}

	workspaceType := ""
	if node.WorkspaceType != "" && node.WorkspaceType != "product_workspace" {
		workspaceType = fmt.Sprintf(" (%s)", node.WorkspaceType)
	}

	fmt.Printf("%s%2d. %-10s %s%s%s\n",
		indent,
		num,
		node.ReferencePrefix,
		node.Name,
		typeIndicator,
		workspaceType,
	)

	// Print children with increased indent
	for i, child := range node.Children {
		printNode(child, i+1, indent+"    ")
	}
}

// countNodes counts total nodes including children.
func countNodes(nodes []*ProductNode) int {
	count := len(nodes)
	for _, node := range nodes {
		count += countNodes(node.Children)
	}
	return count
}

// MarshalJSON implements custom JSON marshaling for tree output.
func (n *ProductNode) MarshalJSON() ([]byte, error) {
	type Alias ProductNode
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}
