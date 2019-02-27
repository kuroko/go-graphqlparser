package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *Schema, walker *Walker) *graphql.Errors {
	ctx := NewQueryContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *Schema, walker *Walker) *graphql.Errors {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}
