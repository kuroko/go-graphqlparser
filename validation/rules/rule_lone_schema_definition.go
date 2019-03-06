package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// Lone Schema definition
//
// A GraphQL document is only valid if it contains only one schema definition.
func loneSchemaDefinition(w *validation.Walker) {
	w.AddSchemaDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.SchemaDefinition) {
		if ctx.Schema.IsExtending {
			ctx.AddError(canNotDefineSchemaWithinExtensionError(0, 0))
			return
		}

		if ctx.HasSeenSchemaDefinition {
			ctx.AddError(schemaDefinitionNotAloneError(0, 0))
			return
		}

		ctx.HasSeenSchemaDefinition = true
	})
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
