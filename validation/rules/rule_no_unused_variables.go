package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// NoUnusedVariables ...
func NoUnusedVariables(w *validation.Walker) {
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
