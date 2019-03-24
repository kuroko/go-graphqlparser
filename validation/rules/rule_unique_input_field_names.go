package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueInputFieldNames ...
func UniqueInputFieldNames(w *validation.Walker) {
	w.AddObjectValueEnterEventHandler(func(ctx *validation.Context, val ast.Value) {
		// TODO: Maybe a better type for this, slice?
		knownNames := make(map[string]struct{}, len(val.ObjectValue))

		for _, field := range val.ObjectValue {
			fieldName := field.Name

			if _, ok := knownNames[fieldName]; ok {
				ctx.AddError(DuplicateInputFieldError(fieldName, 0, 0))
			} else {
				knownNames[fieldName] = struct{}{}
			}
		}
	})
}

// DuplicateInputFieldError ...
func DuplicateInputFieldError(fieldName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one input field named \"" + fieldName + "\".",
		// TODO: Location.
	)
}
