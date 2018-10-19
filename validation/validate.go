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
type RuleFunc func(vctx *Context, walker *ast.Walker)

// Validate ...
func Validate(doc ast.Document, rules []RuleFunc) *graphql.Errors {
	vctx := &Context{}
	walker := ast.NewWalker()

	// Apply all specified validation rules to the walker so that the AST can be validated as it is
	// traversed by the walker, populating the Context along the way.
	for _, rule := range rules {
		rule(vctx, walker)
	}

	walker.Walk(doc)

	return vctx.Errors
}
