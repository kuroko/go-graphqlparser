package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// VisitFunc ...
type VisitFunc func(w *Walker)

// Validate ...
func Validate(doc ast.Document, schema *graphql.Schema, walker *Walker) *Context {
	ctx := NewContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document, schema *graphql.Schema, walker *Walker) *Context {
	ctx := NewSDLContext(doc, schema)

	walker.Walk(ctx, doc)

	return ctx
}

// validateTypeDefinitions ...
func validateTypeDefinitions(ctx *Context) {
	// Check if field definitions in current document are unique.
	for _, typeDef := range ctx.SDLContext.TypeDefinitions {
		switch {
		case ast.IsObjectTypeDefinition(typeDef) || ast.IsInterfaceTypeDefinition(typeDef):
			validateTypeDefinitionFieldDefinitions(ctx, typeDef)
		case ast.IsInputObjectTypeDefinition(typeDef):
			validateTypeDefinitionInputFieldDefinitions(ctx, typeDef)
		case ast.IsEnumTypeDefinition(typeDef):
			validateTypeDefinitionEnumValues(ctx, typeDef)
		}
	}
}

// validateAndMergeTypeExtensions ...
func validateAndMergeTypeExtensions(ctx *Context) {
	for _, typeExts := range ctx.SDLContext.TypeExtensions {
		for _, typeExt := range typeExts {
			switch {
			case ast.IsObjectTypeExtension(typeExt) || ast.IsInterfaceTypeExtension(typeExt):
				validateTypeExtensionFieldDefinitions(ctx, typeExt)
			case ast.IsInputObjectTypeExtension(typeExt):
				validateTypeExtensionInputFieldDefinitions(ctx, typeExt)
			case ast.IsEnumTypeExtension(typeExt):
				validateTypeExtensionEnumValues(ctx, typeExt)
			}
		}
	}
}

// validateTypeDefinitionFieldDefinitions ...
func validateTypeDefinitionFieldDefinitions(ctx *Context, typeDef *ast.TypeDefinition) {
	// Check each field against every other field on this type definition, and ensure that it has a
	// unique name (i.e. part of UniqueFieldDefinitionNames).
	fldAGen := typeDef.FieldsDefinition.Generator()
	fldBGen := typeDef.FieldsDefinition.Generator()

	for fieldA, i := fldAGen.Next(); i >= 0; fieldA, i = fldAGen.Next() {
		for fieldB, j := fldBGen.Next(); j >= 0; fieldB, j = fldBGen.Next() {
			if i == j {
				continue
			}

			if fieldA.Name == fieldB.Name {
				ctx.AddError(DuplicateFieldDefinitionNameError(typeDef.Name, fieldA.Name, 0, 0))
			}
		}

		fldBGen.Reset()
	}
}

// validateTypeDefinitionInputFieldDefinitions ...
func validateTypeDefinitionInputFieldDefinitions(ctx *Context, typeDef *ast.TypeDefinition) {
	// Check each input field against every other field on this type definition, and ensure that it
	// has a unique name (i.e. part of UniqueFieldDefinitionNames).
	fldAGen := typeDef.InputFieldsDefinition.Generator()
	fldBGen := typeDef.InputFieldsDefinition.Generator()

	for fieldA, i := fldAGen.Next(); i >= 0; fieldA, i = fldAGen.Next() {
		for fieldB, j := fldBGen.Next(); j >= 0; fieldB, j = fldBGen.Next() {
			if i == j {
				continue
			}

			if fieldA.Name == fieldB.Name {
				ctx.AddError(DuplicateFieldDefinitionNameError(typeDef.Name, fieldA.Name, 0, 0))
			}
		}

		fldBGen.Reset()
	}
}

// validateTypeDefinitionEnumValues ...
func validateTypeDefinitionEnumValues(ctx *Context, typeDef *ast.TypeDefinition) {
	// Check each enum value against every other enum value on this type definition, and ensure that
	// it has a unique name (i.e. part of UniqueEnumValueNames).
	enumAGen := typeDef.EnumValuesDefinition.Generator()
	enumBGen := typeDef.EnumValuesDefinition.Generator()

	for valA, i := enumAGen.Next(); i >= 0; valA, i = enumAGen.Next() {
		for valB, j := enumBGen.Next(); j >= 0; valB, j = enumBGen.Next() {
			if i == j {
				continue
			}

			if valA.EnumValue == valB.EnumValue {
				ctx.AddError(DuplicateEnumValueNameError(typeDef.Name, valA.EnumValue, 0, 0))
			}
		}

		enumBGen.Reset()
	}
}

// validateTypeExtensionFieldDefinitions ...
func validateTypeExtensionFieldDefinitions(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// TODO: This error needs to be defined elsewhere...
		//ctx.AddError(rules.UnknownTypeError())
		return
	}

	extFldGen := typeExt.FieldsDefinition.Generator()

	for extFld, i := extFldGen.Next(); i >= 0; extFld, i = extFldGen.Next() {
		defFldGen := typeDef.FieldsDefinition.Generator()

		// If there are no fields, it can't collide with anything, so just add the field.
		if typeDef.FieldsDefinition.Len() == 0 {
			typeDef.FieldsDefinition = typeDef.FieldsDefinition.Add(extFld)
			continue
		}

		for defFld, i := defFldGen.Next(); i >= 0; defFld, i = defFldGen.Next() {
			if extFld.Name == defFld.Name {
				ctx.AddError(DuplicateFieldDefinitionNameError(typeDef.Name, extFld.Name, 0, 0))
				continue
			}

			typeDef.FieldsDefinition = typeDef.FieldsDefinition.Add(extFld)
		}
	}
}

// validateTypeExtensionInputFieldDefinitions ...
func validateTypeExtensionInputFieldDefinitions(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// TODO: This error needs to be defined elsewhere...
		//ctx.AddError(rules.UnknownTypeError())
		return
	}

	extFldGen := typeExt.InputFieldsDefinition.Generator()

	for extFld, i := extFldGen.Next(); i >= 0; extFld, i = extFldGen.Next() {
		defFldGen := typeDef.InputFieldsDefinition.Generator()

		// If there are no fields, it can't collide with anything, so just add the field.
		if typeDef.InputFieldsDefinition.Len() == 0 {
			typeDef.InputFieldsDefinition = typeDef.InputFieldsDefinition.Add(extFld)
			continue
		}

		for defFld, i := defFldGen.Next(); i >= 0; defFld, i = defFldGen.Next() {
			if extFld.Name == defFld.Name {
				ctx.AddError(DuplicateFieldDefinitionNameError(typeDef.Name, extFld.Name, 0, 0))
				continue
			}

			typeDef.InputFieldsDefinition = typeDef.InputFieldsDefinition.Add(extFld)
		}
	}
}

// validateTypeExtensionInputFieldDefinitions ...
func validateTypeExtensionEnumValues(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// TODO: This error needs to be defined elsewhere...
		//ctx.AddError(rules.UnknownTypeError())
		return
	}

	extEnumGen := typeExt.EnumValuesDefinition.Generator()

	for extVal, i := extEnumGen.Next(); i >= 0; extVal, i = extEnumGen.Next() {
		defEnumGen := typeDef.EnumValuesDefinition.Generator()

		// If there are no fields, it can't collide with anything, so just add the field.
		if typeDef.EnumValuesDefinition.Len() == 0 {
			typeDef.EnumValuesDefinition = typeDef.EnumValuesDefinition.Add(extVal)
			continue
		}

		for defVal, i := defEnumGen.Next(); i >= 0; defVal, i = defEnumGen.Next() {
			if extVal.EnumValue == defVal.EnumValue {
				ctx.AddError(DuplicateEnumValueNameError(typeDef.Name, extVal.EnumValue, 0, 0))
				continue
			}

			typeDef.EnumValuesDefinition = typeDef.EnumValuesDefinition.Add(extVal)
		}
	}
}
