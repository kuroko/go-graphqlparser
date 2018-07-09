package codegen

// QueryResolverImpl ...
type QueryResolverImpl struct {
}

// Search ...
func (r *QueryResolverImpl) Search(input string) []SearchResult {
	// Actually perform search, in somewhere like Elasticsearch.
	// Loop through the results.
	// Map each result to a SearchResult type, and return them:

	return []SearchResult{
		{Kind: SearchResultKindDroid, Droid: nil},
		{Kind: SearchResultKindHuman, Human: nil},
	}
}

// MutationResolverImpl ...
type MutationResolverImpl struct {
}

// CreateDroid ...
func (r *MutationResolverImpl) CreateDroid(droid DroidInput) DroidResolver {
	// Take input, map to internal type?
	// Call RPC client or something?
	// Map result.

	// Return resolver with created droid.
	return nil
}
