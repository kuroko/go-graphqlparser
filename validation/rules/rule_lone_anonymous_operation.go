package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// loneAnonymousOperation: Lone anonymous operation
//
// A GraphQL document is only valid if when it contains an anonymous operation
// (the query short-hand) it contains only that one operation definition.
func loneAnonymousOperation(ctx *validation.Context) validation.VisitFunc {
	var operations int

	return func(w *validation.Walker) {
		w.AddDocumentEnterEventHandler(func(document ast.Document) {
			document.Definitions.ForEach(func(d ast.Definition, i int) {
				if d.Kind == ast.DefinitionKindExecutable {
					if d.ExecutableDefinition.Kind == ast.ExecutableDefinitionKindOperation {
						operations++
					}
				}
			})
		})

		w.AddOperationDefinitionEnterEventHandler(func(definition *ast.OperationDefinition) {
			if definition.Name == "" && operations > 1 {
				ctx.Errors = ctx.Errors.Add(graphql.NewError(
					"This anonymous operation must be the only defined operation.",
					// TODO: Location.
				))
			}
		})

		w.AddDocumentLeaveEventHandler(func(document ast.Document) {
			operations = 0
		})
	}
}
