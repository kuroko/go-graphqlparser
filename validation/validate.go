package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// RuleFunc ...
type RuleFunc func(ctx *Context) ast.VisitFunc

// Validate ...
func Validate(doc ast.Document, rules []RuleFunc) *graphql.Errors {
	ctx := NewContext(doc)

	visitFns := make([]ast.VisitFunc, 0, len(rules))
	for _, rule := range rules {
		visitFns = append(visitFns, rule(ctx))
	}

	walker := ast.NewWalker(visitFns)
	walker.Walk(doc)

	return ctx.Errors
}
