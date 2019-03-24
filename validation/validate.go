package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *types.Schema, walker *Walker) *types.Errors {
	ctx := NewContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *types.Schema, walker *Walker) *types.Errors {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}
