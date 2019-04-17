package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *types.Schema, walker *Walker) *types.Errors {
	ctx := NewContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx.Errors
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *types.Schema, walker *Walker) (*types.Schema, *types.Errors) {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	buildTypes(ctx)

	return ctx.Schema, ctx.Errors
}

// buildTypes ...
func buildTypes(ctx *Context) {
	if !ctx.SDLContext.IsExtending {
		for typeName, typeDef := range ctx.SDLContext.TypeDefinitions {
			if ast.IsObjectTypeDefinition(typeDef) || ast.IsInterfaceTypeDefinition(typeDef) {
				fldAGen := typeDef.FieldsDefinition.Generator()
				fldBGen := typeDef.FieldsDefinition.Generator()

				for fieldA, i := fldAGen.Emit(); i >= 0; fieldA, i = fldAGen.Emit() {
					for fieldB, j := fldBGen.Emit(); j >= 0; fieldB, j = fldBGen.Emit() {
						if i == j {
							continue
						}

						if fieldA.Name == fieldB.Name {
							ctx.AddError(DuplicateFieldDefinitionNameError(typeName, fieldA.Name, 0, 0))
						}
					}

					fldBGen.Reset()
				}
			}

			if ast.IsInputObjectTypeDefinition(typeDef) {
				fldAGen := typeDef.InputFieldsDefinition.Generator()
				fldBGen := typeDef.InputFieldsDefinition.Generator()

				for fieldA, i := fldAGen.Emit(); i >= 0; fieldA, i = fldAGen.Emit() {
					for fieldB, j := fldBGen.Emit(); j >= 0; fieldB, j = fldBGen.Emit() {
						if i == j {
							continue
						}

						if fieldA.Name == fieldB.Name {
							ctx.AddError(DuplicateFieldDefinitionNameError(typeName, fieldA.Name, 0, 0))
						}
					}

					fldBGen.Reset()
				}
			}
		}

		ctx.Schema.Types = ctx.SDLContext.TypeDefinitions
	}
}

// DuplicateFieldDefinitionNameError ...
func DuplicateFieldDefinitionNameError(typeName, fieldName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + typeName + "." + fieldName + "\" can only be defined once.",
		// TODO: Location.
	)
}
