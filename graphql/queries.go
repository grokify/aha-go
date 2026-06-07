package graphql

// SearchDocumentsQuery is the GraphQL query for searching documents.
const SearchDocumentsQuery = `
query SearchDocuments($query: String!, $searchableType: [String!]!) {
  searchDocuments(filters: {query: $query, searchableType: $searchableType}) {
    nodes {
      name
      url
      searchableId
      searchableType
    }
    currentPage
    totalCount
    totalPages
    isLastPage
  }
}
`

// GetPageQuery is the GraphQL query for getting a page by reference.
const GetPageQuery = `
query GetPage($id: ID!, $includeParent: Boolean!) {
  page(id: $id) {
    name
    description {
      markdownBody
    }
    children {
      name
      referenceNum
    }
    parent @include(if: $includeParent) {
      name
      referenceNum
    }
  }
}
`

// GetFeatureQuery is the GraphQL query for getting a feature by reference.
const GetFeatureQuery = `
query GetFeature($id: ID!) {
  feature(id: $id) {
    name
    description {
      markdownBody
    }
  }
}
`

// GetRequirementQuery is the GraphQL query for getting a requirement by reference.
const GetRequirementQuery = `
query GetRequirement($id: ID!) {
  requirement(id: $id) {
    name
    description {
      markdownBody
    }
  }
}
`
