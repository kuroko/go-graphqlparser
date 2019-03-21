package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueEnumValueNames ...
func uniqueEnumValueNames(w *validation.Walker) {
	w.AddEnumTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		checkEnumValueUniqueness(ctx, def.Name, def.EnumValuesDefinition)
	})

	w.AddEnumTypeExtensionEnterEventHandler(func(ctx *validation.Context, ext *ast.TypeExtension) {
		checkEnumValueUniqueness(ctx, ext.Name, ext.EnumValuesDefinition)
	})
}

// checkEnumValueUniqueness ...
func checkEnumValueUniqueness(ctx *validation.Context, typeName string, evds *ast.EnumValueDefinitions) {
	if ctx.SDLContext.KnownEnumValueNames == nil {
		ctx.SDLContext.KnownEnumValueNames = make(map[string]map[string]struct{})
	}

	if _, ok := ctx.SDLContext.KnownEnumValueNames[typeName]; !ok {
		ctx.SDLContext.KnownEnumValueNames[typeName] = make(map[string]struct{})
	}

	if evds.Len() > 0 {
		valueNames := ctx.SDLContext.KnownEnumValueNames[typeName]

		evds.ForEach(func(evd ast.EnumValueDefinition, i int) {
			var isExistingEnumValue bool

			valueName := evd.EnumValue

			existingType, exists := ctx.Schema.Types[typeName]
			if exists {
				if _, ok := existingType.GetEnumValueDefinition(valueName); ok {
					isExistingEnumValue = true
				}
			}

			if isExistingEnumValue {
				ctx.AddError(existedEnumValueNameMessage(typeName, valueName))
			} else if _, ok := valueNames[valueName]; ok {
				ctx.AddError(duplicateEnumValueNameMessage(typeName, valueName))
			} else {
				valueNames[valueName] = struct{}{}
			}
		})
	}
}

// duplicateEnumValueNameMessage ...
func duplicateEnumValueNameMessage(typeName string, valueName string) graphql.Error {
	return graphql.NewError(
		"Enum value \"" + typeName + "." + valueName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// existedEnumValueNameMessage ...
func existedEnumValueNameMessage(typeName string, valueName string) graphql.Error {
	return graphql.NewError(
		"Enum value \"" + typeName + "." + valueName + "\" already exists in the schema. It cannot also be defined in this type extension.",
		// TODO: Location.
	)
}
