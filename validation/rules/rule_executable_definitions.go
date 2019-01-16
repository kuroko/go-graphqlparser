package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// executableDefinitions ...
func executableDefinitions(w *validation.Walker) {
	w.AddDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.Definition) {
		if def.Kind != ast.DefinitionKindExecutable {
			err := graphql.NewError(nonExecutableDefinitionMessage(def))
			// TODO: Add locations to AST.
			err.Locations = err.Locations.Add(graphql.Location{Line: 6, Column: 7})
			// NOTE(elliot): err.Path is unused in validation.

			ctx.Errors = ctx.Errors.Add(err)
		}
	})
}

// nonExecutableDefinitionMessage ...
func nonExecutableDefinitionMessage(def ast.Definition) string {
	var name string

	// TODO: We really need to move the name field to the top level of Definition...
	switch def.Kind {
	case ast.DefinitionKindTypeSystem:
		tsDef := def.TypeSystemDefinition

		switch tsDef.Kind {
		case ast.TypeSystemDefinitionKindSchema:
			name = "schema"
		case ast.TypeSystemDefinitionKindType:
			name = tsDef.TypeDefinition.Name
		case ast.TypeSystemDefinitionKindDirective:
			name = tsDef.DirectiveDefinition.Name
		}
	case ast.DefinitionKindTypeSystemExtension:
		tseDef := def.TypeSystemExtension

		switch tseDef.Kind {
		case ast.TypeSystemExtensionKindSchema:
			name = "schema"
		case ast.TypeSystemExtensionKindType:
			name = tseDef.TypeExtension.Name
		}
	}

	return fmt.Sprintf("The %s definition is not executable.", name)
}
