package example

// SearchDocumentsResponse is the response from the searchDocuments query.
type SearchDocumentsResponse struct {
	SearchDocuments SearchResults `json:"searchDocuments"`
}

// SearchResults contains paginated search results.
type SearchResults struct {
	Nodes       []DocumentNode `json:"nodes"`
	CurrentPage int            `json:"currentPage"`
	TotalCount  int            `json:"totalCount"`
	TotalPages  int            `json:"totalPages"`
	IsLastPage  bool           `json:"isLastPage"`
}

// DocumentNode represents a search result document.
type DocumentNode struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	SearchableID   string `json:"searchableId"`
	SearchableType string `json:"searchableType"`
}

// PageResponse is the response from the page query.
type PageResponse struct {
	Page *Page `json:"page"`
}

// Page represents an Aha page/note.
type Page struct {
	Name        string       `json:"name"`
	Description *Description `json:"description"`
	Children    []PageRef    `json:"children"`
	Parent      *PageRef     `json:"parent"`
}

// PageRef is a reference to a page.
type PageRef struct {
	Name         string `json:"name"`
	ReferenceNum string `json:"referenceNum"`
}

// Description contains formatted content.
type Description struct {
	MarkdownBody string `json:"markdownBody"`
}

// FeatureResponse is the response from the feature query.
type FeatureResponse struct {
	Feature *Feature `json:"feature"`
}

// Feature represents an Aha feature from GraphQL.
type Feature struct {
	Name        string       `json:"name"`
	Description *Description `json:"description"`
}

// RequirementResponse is the response from the requirement query.
type RequirementResponse struct {
	Requirement *Requirement `json:"requirement"`
}

// Requirement represents an Aha requirement from GraphQL.
type Requirement struct {
	Name        string       `json:"name"`
	Description *Description `json:"description"`
}
