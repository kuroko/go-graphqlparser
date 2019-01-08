package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

func noUnusedVariables(walker *validation.Walker) {}

func unusedVariableMessage(varName, opName string, line, col int) graphql.Error {
	msg := fmt.Sprintf("Variable %s is never used", varName)

	if len(opName) > 0 {
		msg += fmt.Sprintf(" in operation %s", opName)
	}

	return graphql.NewError(
		msg + ".",
		// TODO(seeruk): Location.
	)
}
