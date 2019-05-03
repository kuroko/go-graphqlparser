package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UnionHasMembers ...
func UnionHasMembers(w *validation.Walker) {
	w.AddUnionTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		if def.UnionMemberTypes.Len() == 0 {
			ctx.AddError(validation.UnionHasNoMembersError(def.Name, 0, 0))
		}
	})
}
