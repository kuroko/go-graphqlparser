package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// fragmentsOnCompositeTypes ...
func fragmentsOnCompositeTypes(w *validation.Walker) {}

func fragmentsOnCompositeTypesError(fragName, typeName string, line, col int) graphql.Error {
	message := "Fragment "
	if len(fragName) != 0 {
		message += fragName + " "
	}
	message += fmt.Sprintf(`cannot condition on non composite type "%s"`, typeName)

	return graphql.NewError(
		message,
		// TODO: Location.
	)
}
