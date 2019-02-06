package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

var (
	// DefaultQueryWalker is the default query walker.
	DefaultQueryWalker = NewWalker(rules.Specified)
	// DefaultSDLWalker is the default SDL walker.
	DefaultSDLWalker = NewWalker(rules.SpecifiedSDL)
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *graphql.Schema) *graphql.Errors {
	ctx := NewQueryContext(doc, schema)

	DefaultQueryWalker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateWithWalker ...
func ValidateWithWalker(doc ast.Document, schema *graphql.Schema, walker *Walker) *graphql.Errors {
	ctx := NewQueryContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *graphql.Schema) *graphql.Errors {
	ctx := NewSDLContext(doc, schema)

	DefaultSDLWalker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateSDLWithWalker ...
func ValidateSDLWithWalker(doc ast.Document, schema *graphql.Schema, walker *Walker) *graphql.Errors {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}
