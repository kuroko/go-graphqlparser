package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// ExecutableDefinitions ...
func ExecutableDefinitions(w *validation.Walker) {
	w.AddDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.Definition) {
		if def.Kind != ast.DefinitionKindExecutable {
			ctx.AddError(validation.NonExecutableDefinitionError(getDefinitionName(def), 0, 0))
		}
	})
}

// getDefinitionName ...
func getDefinitionName(def ast.Definition) string {
	var name string

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

	return name
}
