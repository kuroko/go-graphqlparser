package validation

import "github.com/bucketd/go-graphqlparser/ast"

// Schema ...
type Schema struct {
	// NOTE: We can't really include this, because we may re-use a Schema instance to extend an
	// existing schema, e.g. if we were stitching together multiple schema documents. One
	// alternative would be to keep a slice / map of documents instead - if we needed them. We
	// should aim to store enough information on this type that we don't need the original
	// ast.Document for the schema.
	//Document ast.Document

	// NOTE: The following should all be ast.Type instances with the kind ast.TypeKindNamed.
	QueryType        *ast.Type
	MutationType     *ast.Type
	SubscriptionType *ast.Type

	// TODO: Probably not final type.
	// TODO: Not actually populated yet.
	Types map[string]struct{}
}
