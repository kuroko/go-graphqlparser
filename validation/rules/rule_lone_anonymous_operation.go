package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// LoneAnonymousOperation: Lone anonymous operation
//
// A GraphQL document is only valid if when it contains an anonymous operation
// (the query short-hand) it contains only that one operation definition.
func LoneAnonymousOperation(w *validation.Walker) {
	w.AddDocumentEnterEventHandler(func(ctx *validation.Context, document ast.Document) {
		document.Definitions.ForEach(func(d ast.Definition, i int) {
			if d.Kind == ast.DefinitionKindExecutable {
				if d.ExecutableDefinition.Kind == ast.ExecutableDefinitionKindOperation {
					ctx.OperationsCount++
				}
			}
		})
	})

	w.AddOperationDefinitionEnterEventHandler(func(ctx *validation.Context, definition *ast.OperationDefinition) {
		if definition.Name == "" && ctx.OperationsCount > 1 {
			ctx.AddError(AnonOperationNotAloneError())
		}
	})
}

// AnonOperationNotAloneError ...
func AnonOperationNotAloneError() types.Error {
	return types.NewError(
		anonOperationNotAloneMessage(),
		// TODO: Location.
	)
}

// anonOperationNotAloneMessage ...
func anonOperationNotAloneMessage() string {
	return "This anonymous operation must be the only defined operation."
}
