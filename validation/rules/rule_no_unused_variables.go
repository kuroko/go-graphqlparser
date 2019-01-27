package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// noUnusedVariables ...
func noUnusedVariables(w *validation.Walker) {
	w.AddExecutableDefinitionLeaveEventHandler(func(ctx *validation.Context, def *ast.ExecutableDefinition) {
		// TODO: Should we evaluate adding parent/child relationships between AST nodes to avoid
		// this kind of thing? The reason I've switched this to use ExecutableDefinition is because
		// it's more appropriate for us to store the variable usages by ExecutableDefinition, but
		// from the OperationDefinition we'd spend more time looking up the ExecutableDefinition for
		// the OperationDefinition. It'd be cheap to set this up in the parser, as long as memory
		// usage wasn't a problem.
		if def.Kind != ast.ExecutableDefinitionKindOperation {
			return
		}

		opDef := def.OperationDefinition

		// TODO: Should use recursive variable usages.
		variableUses := ctx.VariableUsages(def)

		opDef.VariableDefinitions.ForEach(func(varDef ast.VariableDefinition, _ int) {
			var used bool
			for _, u := range variableUses {
				if u == varDef.Name {
					used = true
				}
			}

			if !used {
				ctx.Errors = ctx.Errors.Add(unusedVariableError(varDef.Name, opDef.Name, 0, 0))
			}
		})
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
