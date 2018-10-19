package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// Context ...
type Context struct {
	Errors *graphql.Errors
	Schema *graphql.Schema
}

// RuleFunc ...
type RuleFunc func(walker *Walker)

// Validate ...
func Validate(doc ast.Document, walker *Walker) *graphql.Errors {
	ctx := &Context{}

	walker.Walk(ctx, doc)

	return ctx.Errors
}
