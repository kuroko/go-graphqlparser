package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *graphql.Schema, walker *Walker) *Context {
	ctx := NewContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *graphql.Schema, walker *Walker) *Context {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx
}
