package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, walker *Walker) *graphql.Errors {
	ctx := NewContext(doc)
	walker.Walk(ctx, doc)

	return ctx.Errors
}
