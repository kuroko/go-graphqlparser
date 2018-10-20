package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// executableDefinitions ...
func executableDefinitions(walker *validation.Walker) {
	walker.AddDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.Definition) {
		if def.Kind != ast.DefinitionKindExecutable {
			// TODO...
			//ctx.Errors = ctx.Errors.Add(errors.New(
			//	nonExecutableDefinitionMessage(def),
			//))
		}
	})
}

// nonExecutableDefinitionMessage ...
func nonExecutableDefinitionMessage(def ast.Definition) string {
	var name string

	// TODO(elliot): We really need to move the name field to the top level of Definition...
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
