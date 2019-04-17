package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// DuplicateOperationNameError ...
func DuplicateOperationNameError(operationName string, line, col int) types.Error {
	return types.NewError("There can be only one operation named " + operationName + ".")
}

// UniqueOperationNames must each defined operation have.
func UniqueOperationNames(w *validation.Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *validation.Context, od *ast.OperationDefinition) {
		nameMatchDef, ok := ctx.OperationDefinitions[od.Name]
		if ok && nameMatchDef != od {
			ctx.AddError(DuplicateOperationNameError(od.Name, 0, 0))
		}
	})
}
