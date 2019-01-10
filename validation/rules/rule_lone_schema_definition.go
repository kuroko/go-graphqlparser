package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// Lone Schema definition
//
// A GraphQL document is only valid if it contains only one schema definition.
func loneSchemaDefinition(ctx *validation.Context) ast.VisitFunc {
	var schemaDefinitions int

	return func(w *ast.Walker) {
		w.AddSchemaDefinitionEnterEventHandler(func(_ *ast.SchemaDefinition) {
			// TODO: implement logic once schema is implemented.
			_ = ctx.Schema
			if false {
				ctx.Errors = ctx.Errors.Add(canNotDefineSchemaWithinExtensionError(0, 0))
			}

			if schemaDefinitions > 0 {
				ctx.Errors = ctx.Errors.Add(schemaDefinitionNotAloneError(0, 0))
			}
			schemaDefinitions++
		})
	}
}

func schemaDefinitionNotAloneError(line, col int) graphql.Error {
	return graphql.NewError(
		"Must provide only one schema definition.",
		// TODO: Location.
	)
}

func canNotDefineSchemaWithinExtensionError(line, col int) graphql.Error {
	return graphql.NewError(
		"Cannot define a new schema within a schema extension.",
		// TODO: Location.
	)
}
