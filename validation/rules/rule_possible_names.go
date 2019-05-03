package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

const us = '_'

// PossibleNames ...
func PossibleNames(w *validation.Walker) {
	// Arguments...
	// InputFields...
	w.AddInputValueDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.InputValueDefinition) {
		if !isValidName(def.Name) {
			ctx.AddError(validation.NameStartsWithTwoUnderscoresError(def.Name, 0, 0))
		}
	})

	// Directives...
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
		if !isValidName(def.Name) {
			ctx.AddError(validation.NameStartsWithTwoUnderscoresError(def.Name, 0, 0))
		}
	})

	// Fields...
	w.AddFieldDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.FieldDefinition) {
		if !isValidName(def.Name) {
			ctx.AddError(validation.NameStartsWithTwoUnderscoresError(def.Name, 0, 0))
		}
	})

	// Types...
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		if !isValidName(def.Name) {
			ctx.AddError(validation.NameStartsWithTwoUnderscoresError(def.Name, 0, 0))
		}
	})
}

// isValidName ...
func isValidName(name string) bool {
	if len(name) < 2 {
		return true
	}

	return name[0] != us && name[1] != us
}
