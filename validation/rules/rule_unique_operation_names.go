package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueOperationNames must each defined operation have.
func UniqueOperationNames(w *validation.Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *validation.Context, od *ast.OperationDefinition) {
		// TODO: Surely this rule doesn't work? It's stored in a map by name, so if there were more
		// than one definition with the same name, it'd end up overwriting the existing map key?
		nameMatchDef, ok := ctx.OperationDefinitions[od.Name]
		if ok && nameMatchDef != od {
			ctx.AddError(validation.DuplicateOperationNameError(od.Name, 0, 0))
		}
	})
}
