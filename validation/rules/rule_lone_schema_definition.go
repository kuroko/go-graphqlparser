package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// Lone Schema definition
//
// A GraphQL document is only valid if it contains only one schema definition.
func loneSchemaDefinition(walker *validation.Walker) {
	var schemaDefinitions int

	walker.AddSchemaDefinitionEnterEventHandler(func(context *validation.Context, _ *ast.SchemaDefinition) {
		// TODO: implement logic once schema is implemented.
		_ = context.Schema
		if false {
			context.Errors = context.Errors.Add(canNotDefineSchemaWithinExtensionMessage(0, 0))
		}

		if schemaDefinitions > 0 {
			context.Errors = context.Errors.Add(schemaDefinitionNotAlone(0, 0))
		}
		schemaDefinitions++
	})
}

func schemaDefinitionNotAlone(line, col int) graphql.Error {
	return graphql.NewError(
		"Must provide only one schema definition.",
		// TODO(seeruk): Location.
	)
}

func canNotDefineSchemaWithinExtensionMessage(line, col int) graphql.Error {
	return graphql.NewError(
		"Cannot define a new schema within a schema extension.",
		// TODO(seeruk): Location.
	)
}
