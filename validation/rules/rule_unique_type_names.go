package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueTypeNames ...
func UniqueTypeNames(w *validation.Walker) {
	w.AddDocumentEnterEventHandler(func(ctx *validation.Context, doc ast.Document) {
		knownTypeNames := make(map[string]struct{}, doc.TypeDefinitions)

		doc.Definitions.ForEach(func(def ast.Definition, i int) {
			if def.Kind != ast.DefinitionKindTypeSystem {
				return
			}

			if def.TypeSystemDefinition.Kind != ast.TypeSystemDefinitionKindType {
				return
			}

			tdef := def.TypeSystemDefinition.TypeDefinition

			if _, ok := ctx.Schema.Types[tdef.Name]; ok {
				ctx.AddError(ExistedTypeNameError(tdef.Name, 0, 0))
				return
			}

			if _, ok := knownTypeNames[tdef.Name]; ok {
				ctx.AddError(DuplicateTypeNameError(tdef.Name, 0, 0))
				return
			} else {
				knownTypeNames[tdef.Name] = struct{}{}
			}
		})
	})
}

// DuplicateTypeNameError ...
func DuplicateTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// ExistedTypeNameError ...
func ExistedTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}
