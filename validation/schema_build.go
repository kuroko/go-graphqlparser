package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// MustBuildSchema ...
func MustBuildSchema(ctx *Context) *graphql.Schema {
	schema, err := BuildSchema(ctx)
	if err != nil {
		panic(err)
	}

	return schema
}

// BuildSchema ...
func BuildSchema(ctx *Context) (*graphql.Schema, error) {
	buildSchema(ctx)

	return ctx.Schema, nil
}

// buildSchema ...
func buildSchema(ctx *Context) {
	mergeSchemaExtensions(ctx)
	mergeTypeExtensions(ctx)

	schemaDef := ctx.SDLContext.SchemaDefinition

	// TODO: Split to function.
	if schemaDef != nil {
		if operationDefs := schemaDef.OperationTypeDefinitions; operationDefs != nil {
			operationDefs.ForEach(func(otd ast.OperationTypeDefinition, i int) {
				switch otd.OperationType {
				case ast.OperationDefinitionKindQuery:
					ctx.Schema.QueryType = &otd.NamedType
				case ast.OperationDefinitionKindMutation:
					ctx.Schema.MutationType = &otd.NamedType
				case ast.OperationDefinitionKindSubscription:
					ctx.Schema.SubscriptionType = &otd.NamedType
				}
			})
		}

		ctx.Schema.Definition = schemaDef
	} else {
		// TODO: Include default schema definition, with query option set up.
		// TODO: Ensure there's some validation around the Query type being defined in this case.
	}

	if !ctx.SDLContext.IsExtending {
		// The map on SDLContext is created with the right size to also contain the built-in
		// directives without needing to grow.
		ctx.Schema.Directives = ctx.SDLContext.DirectiveDefinitions

		for name, def := range graphql.SpecifiedDirectives() {
			// TODO: Should we allow overriding built-in directives?
			// TODO: ^ Yes, maybe?
			ctx.Schema.Directives[name] = def
		}

		// The map on SDLContext is created with the right size to also contain the built-in types
		// without needing to grow.
		ctx.Schema.Types = ctx.SDLContext.TypeDefinitions

		for name, def := range graphql.SpecifiedTypes() {
			// TODO: Should we allow overriding built-in types?
			ctx.Schema.Types[name] = def
		}
	} else {
		// TODO: Handle adding the new types / directives to an existing schema.
	}
}

// mergeSchemaExtensions ...
func mergeSchemaExtensions(ctx *Context) {
	if ctx.SDLContext.SchemaDefinition == nil || len(ctx.SDLContext.SchemaExtensions) == 0 {
		return
	}

	for _, schemaExt := range ctx.SDLContext.SchemaExtensions {
		ctx.SDLContext.SchemaDefinition.Directives.Join(schemaExt.Directives)
		ctx.SDLContext.SchemaDefinition.OperationTypeDefinitions.
			Join(schemaExt.OperationTypeDefinitions)
	}
}

// mergeTypeExtensions ...
func mergeTypeExtensions(ctx *Context) {
	for _, typeExts := range ctx.SDLContext.TypeExtensions {
		for _, typeExt := range typeExts {
			switch {
			case ast.IsObjectTypeExtension(typeExt) || ast.IsInterfaceTypeExtension(typeExt):
				mergeTypeExtensionFieldDefinitions(ctx, typeExt)
			case ast.IsInputObjectTypeExtension(typeExt):
				mergeTypeExtensionInputFieldDefinitions(ctx, typeExt)
			case ast.IsEnumTypeExtension(typeExt):
				mergeTypeExtensionEnumValues(ctx, typeExt)
			}
		}
	}
}

// mergeTypeExtensionFieldDefinitions ...
func mergeTypeExtensionFieldDefinitions(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// Handled by rule PossibleTypeExtensions.
		return
	}

	typeDef.FieldsDefinition.Join(typeExt.FieldsDefinition)
}

// mergeTypeExtensionInputFieldDefinitions ...
func mergeTypeExtensionInputFieldDefinitions(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// Handled by rule PossibleTypeExtensions.
		return
	}

	typeDef.InputFieldsDefinition.Join(typeExt.InputFieldsDefinition)
}

// mergeTypeExtensionEnumValues ...
func mergeTypeExtensionEnumValues(ctx *Context, typeExt *ast.TypeExtension) {
	typeDef, _ := ctx.TypeDefinition(typeExt.Name)
	if typeDef == nil {
		// Handled by rule PossibleTypeExtensions.
		return
	}

	typeDef.EnumValuesDefinition.Join(typeExt.EnumValuesDefinition)
}
