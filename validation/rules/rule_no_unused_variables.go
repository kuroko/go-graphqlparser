package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// noUnusedVariables ...
func noUnusedVariables(w *validation.Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *validation.Context, opDef *ast.OperationDefinition) {
		ctx.VariableDefs = &ast.VariableDefinitions{}
	})

	w.AddOperationDefinitionLeaveEventHandler(func(ctx *validation.Context, opDef *ast.OperationDefinition) {
		variableUses := ctx.RecursiveVariableUsages(opDef.Name)

		ctx.VariableDefs.ForEach(func(varDef ast.VariableDefinition, _ int) {
			if !variableUses[varDef.Name] {
				ctx.Errors = ctx.Errors.Add(unusedVariableError(varDef.Name, opDef.Name, 0, 0))
			}
		})
	})

	w.AddVariableDefinitionEnterEventHandler(func(ctx *validation.Context, varDef ast.VariableDefinition) {
		ctx.VariableDefs.Add(varDef)
	})
}

// unusedVariableError ...
func unusedVariableError(varName, opName string, line, col int) graphql.Error {
	return graphql.NewError(
		unusedVariableMessage(varName, opName),
		// TODO: Location.
	)
}

// unusedVariableMessage ...
func unusedVariableMessage(varName, opName string) string {
	if len(opName) > 0 {
		return fmt.Sprintf("Variable %s is never used in operation %s", varName, opName)
	}

	return fmt.Sprintf("Variable %s is never used", varName)
}
