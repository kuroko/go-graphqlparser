package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// NoUnusedVariables ...
func NoUnusedVariables(w *validation.Walker) {
	w.AddExecutableDefinitionLeaveEventHandler(func(ctx *validation.Context, def *ast.ExecutableDefinition) {
		if def.Kind != ast.ExecutableDefinitionKindOperation {
			return
		}

		opDef := def.OperationDefinition
		variableUses := ctx.RecursiveVariableUsages(def)

		opDef.VariableDefinitions.ForEach(func(varDef ast.VariableDefinition, _ int) {
			var used bool
			for k := range variableUses {
				if k == varDef.Name {
					used = true
				}
			}

			if !used {
				ctx.Errors = ctx.Errors.Add(validation.UnusedVariableError(varDef.Name, opDef.Name, 0, 0))
			}
		})
	})
}
