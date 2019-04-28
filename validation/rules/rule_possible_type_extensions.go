package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

var defKindToExtKind = map[ast.TypeDefinitionKind]ast.TypeExtensionKind{
	ast.TypeDefinitionKindScalar:      ast.TypeExtensionKindScalar,
	ast.TypeDefinitionKindObject:      ast.TypeExtensionKindObject,
	ast.TypeDefinitionKindInterface:   ast.TypeExtensionKindInterface,
	ast.TypeDefinitionKindUnion:       ast.TypeExtensionKindUnion,
	ast.TypeDefinitionKindEnum:        ast.TypeExtensionKindEnum,
	ast.TypeDefinitionKindInputObject: ast.TypeExtensionKindInputObject,
}

// PossibleTypeExtensions ...
func PossibleTypeExtensions(w *validation.Walker) {
	w.AddTypeExtensionEnterEventHandler(func(ctx *validation.Context, ext *ast.TypeExtension) {
		typeName := ext.Name
		typeDef, _ := ctx.TypeDefinition(typeName)

		if typeDef != nil {
			expectedKind := defKindToExtKind[typeDef.Kind]
			if ext.Kind != expectedKind {
				ctx.AddError(validation.ExtendingDifferentTypeKindError(typeName, typeDef.Kind.String(), 0, 0))
			}
		} else {
			ctx.AddError(validation.ExtendingUnknownTypeError(typeName, 0, 0))
		}
	})
}
