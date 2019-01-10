package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// RuleFunc ...
type RuleFunc func(ctx *Context) VisitFunc

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, rules []RuleFunc) *graphql.Errors {
	ctx := NewContext(doc)

	visitFns := make([]VisitFunc, 0, len(rules))
	for _, rule := range rules {
		visitFns = append(visitFns, rule(ctx))
	}

	walker := NewWalker(visitFns)
	walker.Walk(doc)

	return ctx.Errors
}
