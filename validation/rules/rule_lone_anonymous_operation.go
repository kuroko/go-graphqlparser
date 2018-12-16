package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// loneAnonymousOperation: Lone anonymous operation
//
// A GraphQL document is only valid if when it contains an anonymous operation
// (the query short-hand) that it contains only that one operation definition.
func loneAnonymousOperation(walker *validation.Walker) {
	var operations int

	walker.AddDocumentEnterEventHandler(func(context *validation.Context, document ast.Document) {
		document.Definitions.ForEach(func(d ast.Definition, i int) {
			if d.Kind == ast.DefinitionKindExecutable {
				if d.ExecutableDefinition.Kind == ast.ExecutableDefinitionKindOperation {
					operations++
				}
			}
		})
	})

	walker.AddOperationDefinitionEnterEventHandler(func(context *validation.Context, definition *ast.OperationDefinition) {
		if definition.Name == "" && operations > 1 {
			context.Errors = context.Errors.Add(graphql.NewError(
				"This anonymous operation must be the only defined operation.",
				// TODO(seeruk): Location.
			))
		}
	})
}
