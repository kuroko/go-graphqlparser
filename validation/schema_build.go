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

	buildSchemaDefinition(ctx)
	buildDirectives(ctx)
	buildTypes(ctx)
}

// buildDirectives ...
func buildDirectives(ctx *Context) {
	if !ctx.SDLContext.IsExtending {
		// The map on SDLContext is created with the right size to also contain the built-in types
		// without needing to grow.
		ctx.Schema.Directives = ctx.SDLContext.DirectiveDefinitions

		for name, def := range graphql.SpecifiedDirectives() {
			if _, ok := ctx.Schema.Directives[name]; !ok {
				ctx.Schema.Directives[name] = def
			}
		}
	} else {
		directivesSize := len(ctx.Schema.Directives) + len(ctx.SDLContext.DirectiveDefinitions)
		directives := make(map[string]*ast.DirectiveDefinition, directivesSize)

		for name, def := range ctx.SDLContext.DirectiveDefinitions {
			directives[name] = def
		}

		ctx.Schema.Directives = directives
	}
}

// buildTypes ...
func buildTypes(ctx *Context) {
	if !ctx.SDLContext.IsExtending {
		// The map on SDLContext is created with the right size to also contain the built-in types
		// without needing to grow.
		ctx.Schema.Types = ctx.SDLContext.TypeDefinitions

		for name, def := range graphql.SpecifiedTypes() {
			if _, ok := ctx.Schema.Types[name]; !ok {
				ctx.Schema.Types[name] = def
			}
		}
	} else {
		typesSize := len(ctx.Schema.Types) + len(ctx.SDLContext.TypeDefinitions)
		types := make(map[string]*ast.TypeDefinition, typesSize)

		for name, def := range ctx.SDLContext.TypeDefinitions {
			types[name] = def
		}

		ctx.Schema.Types = types
	}
}

// buildSchemaDefinition ...
func buildSchemaDefinition(ctx *Context) {
	schemaDef := ctx.SDLContext.SchemaDefinition
	if schemaDef == nil {
		// This is validated by the next step, as it would be if we were using a definition found on
		// the SDLContext instance.
		schemaDef = &ast.SchemaDefinition{
			OperationTypeDefinitions: (*ast.OperationTypeDefinitions)(nil).
				Add(ast.OperationTypeDefinition{
					NamedType: ast.Type{
						NamedType: "Query",
						Kind:      ast.TypeKindNamed,
					},
					OperationType: ast.OperationDefinitionKindQuery,
				}),
		}
	}

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
	for typeName, typeExts := range ctx.SDLContext.TypeExtensions {
		typeDef, _ := ctx.TypeDefinition(typeName)
		if typeDef == nil {
			// Handled by rule PossibleTypeExtensions.
			return
		}

		for _, typeExt := range typeExts {
			typeDef.Directives.Join(typeExt.Directives)

			switch {
			case ast.IsObjectTypeExtension(typeExt):
				typeDef.FieldsDefinition.Join(typeExt.FieldsDefinition)
				typeDef.ImplementsInterface.Join(typeExt.ImplementsInterface)
			case ast.IsInterfaceTypeExtension(typeExt):
				typeDef.FieldsDefinition.Join(typeExt.FieldsDefinition)
			case ast.IsUnionTypeExtension(typeExt):
				typeDef.UnionMemberTypes.Join(typeExt.UnionMemberTypes)
			case ast.IsEnumTypeExtension(typeExt):
				typeDef.EnumValuesDefinition.Join(typeExt.EnumValuesDefinition)
			case ast.IsInputObjectTypeExtension(typeExt):
				typeDef.InputFieldsDefinition.Join(typeExt.InputFieldsDefinition)
			}
		}
	}
}
