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

	// TODO: Split to function v
	schemaDef := ctx.SDLContext.SchemaDefinition
	if schemaDef == nil {
		// TODO: What if there's no Query type defined? When do we add an error, given that this is
		// not validation, and we're trying to avoid adding errors here.
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
	// TODO: Split to function ^

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
