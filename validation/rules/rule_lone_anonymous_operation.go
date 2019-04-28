package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// LoneAnonymousOperation: Lone anonymous operation
//
// A GraphQL document is only valid if when it contains an anonymous operation
// (the query short-hand) it contains only that one operation definition.
func LoneAnonymousOperation(w *validation.Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *validation.Context, definition *ast.OperationDefinition) {
		if definition.Name == "" && ctx.Document.OperationDefinitions > 1 {
			ctx.AddError(validation.AnonOperationNotAloneError(0, 0))
		}
	})
}
