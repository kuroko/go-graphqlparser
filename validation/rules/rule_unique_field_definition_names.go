package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueFieldDefinitionNames ...
//
// TODO: Move to BuildSchema:
// NOTE: Maybe we can run this rule when we begin merging types with their extensions. We'll be
// going over this same data, stored in a similar format again, but maybe using less memory in fact.
func UniqueFieldDefinitionNames(w *validation.Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		typeName := def.Name

		prepareFieldDefinitionSymbolTables(ctx, typeName)

		// Objects and interfaces have fields.
		if (ast.IsObjectTypeDefinition(def) || ast.IsInterfaceTypeDefinition(def)) && def.FieldsDefinition.Len() > 0 {
			def.FieldsDefinition.ForEach(func(fd ast.FieldDefinition, i int) {
				checkFieldDefinitionNameUniqueness(ctx, typeName, fd.Name)
			})
		}

		// Input objects have input fields.
		if ast.IsInputObjectTypeDefinition(def) && def.InputFieldsDefinition.Len() > 0 {
			def.InputFieldsDefinition.ForEach(func(ivd ast.InputValueDefinition, i int) {
				checkFieldDefinitionNameUniqueness(ctx, typeName, ivd.Name)
			})
		}
	})

	w.AddTypeExtensionEnterEventHandler(func(ctx *validation.Context, ext *ast.TypeExtension) {
		typeName := ext.Name

		prepareFieldDefinitionSymbolTables(ctx, typeName)

		// Objects and interfaces have fields.
		if (ast.IsObjectTypeExtension(ext) || ast.IsInterfaceTypeExtension(ext)) && ext.FieldsDefinition.Len() > 0 {
			ext.FieldsDefinition.ForEach(func(fd ast.FieldDefinition, i int) {
				checkFieldDefinitionNameUniqueness(ctx, typeName, fd.Name)
			})
		}

		// Input objects have input fields.
		if ast.IsInputObjectTypeExtension(ext) && ext.InputFieldsDefinition.Len() > 0 {
			ext.InputFieldsDefinition.ForEach(func(ivd ast.InputValueDefinition, i int) {
				checkFieldDefinitionNameUniqueness(ctx, typeName, ivd.Name)
			})
		}
	})
}

// DuplicateFieldDefinitionNameMessage ...
func DuplicateFieldDefinitionNameError(typeName, fieldName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + typeName + "." + fieldName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// ExistedFieldDefinitionNameMessage ...
func ExistedFieldDefinitionNameError(typeName, fieldName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + typeName + "." + fieldName + "\" already exists in the schema. It cannot also be defined in this type extension.",
		// TODO: Location.
	)
}

// prepareFieldDefinitionSymbolTables ...
func prepareFieldDefinitionSymbolTables(ctx *validation.Context, typeName string) {
	if ctx.SDLContext.KnownFieldNames == nil {
		ctx.SDLContext.KnownFieldNames = make(map[string]map[string]struct{}, ctx.Document.TypeDefinitions)
	}

	if _, ok := ctx.SDLContext.KnownFieldNames[typeName]; !ok {
		var fieldCount int

		if typeDef, ok := ctx.SDLContext.TypeDefinitions[typeName]; ok {
			if typeDef.Kind == ast.TypeDefinitionKindInterface || typeDef.Kind == ast.TypeDefinitionKindObject {
				fieldCount = typeDef.FieldsDefinition.Len()
			} else if typeDef.Kind == ast.TypeDefinitionKindInputObject {
				fieldCount = typeDef.InputFieldsDefinition.Len()
			}
		}

		if typeExt, ok := ctx.SDLContext.TypeExtensions[typeName]; ok {
			if typeExt.Kind == ast.TypeExtensionKindInterface || typeExt.Kind == ast.TypeExtensionKindObject {
				fieldCount = typeExt.FieldsDefinition.Len()
			} else if typeExt.Kind == ast.TypeExtensionKindInputObject {
				fieldCount = typeExt.InputFieldsDefinition.Len()
			}
		}

		ctx.SDLContext.KnownFieldNames[typeName] = make(map[string]struct{}, fieldCount)
	}
}

// checkFieldDefinitionNameUniqueness ...
func checkFieldDefinitionNameUniqueness(ctx *validation.Context, typeName, fieldName string) {
	if hasField(ctx.Schema.Types[typeName], fieldName) {
		ctx.AddError(ExistedFieldDefinitionNameError(typeName, fieldName, 0, 0))
	} else if _, ok := ctx.SDLContext.KnownFieldNames[typeName][fieldName]; ok {
		ctx.AddError(DuplicateFieldDefinitionNameError(typeName, fieldName, 0, 0))
	} else {
		ctx.SDLContext.KnownFieldNames[typeName][fieldName] = struct{}{}
	}
}

// hasField ...
func hasField(def *ast.TypeDefinition, fieldName string) bool {
	if def == nil {
		return false
	}

	if ast.IsObjectTypeDefinition(def) || ast.IsInterfaceTypeDefinition(def) {
		var found bool

		// TODO: It'd be good to be able to break out of this loop early if we need to.
		def.FieldsDefinition.ForEach(func(fd ast.FieldDefinition, i int) {
			if fd.Name == fieldName {
				found = true
			}
		})

		return found
	}

	if ast.IsInputObjectTypeDefinition(def) {
		var found bool

		// TODO: It'd be good to be able to break out of this loop early if we need to.
		def.InputFieldsDefinition.ForEach(func(ivd ast.InputValueDefinition, i int) {
			if ivd.Name == fieldName {
				found = true
			}
		})

		return found
	}

	return false
}
