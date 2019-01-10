package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// noUnusedVariables ...
func noUnusedVariables(ctx *validation.Context) ast.VisitFunc {
	var variableDefs *ast.VariableDefinitions

	return func(w *ast.Walker) {
		w.AddOperationDefinitionEnterEventHandler(func(definition *ast.OperationDefinition) {
			variableDefs = &ast.VariableDefinitions{}
		})

		w.AddOperationDefinitionLeaveEventHandler(func(definition *ast.OperationDefinition) {
			// TODO: Magic: https://github.com/graphql/graphql-js/blob/master/src/validation/rules/NoUnusedVariables.js#L37
		})

		w.AddVariableDefinitionEnterEventHandler(func(def ast.VariableDefinition) {
			variableDefs.Add(def)
		})
	}
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
