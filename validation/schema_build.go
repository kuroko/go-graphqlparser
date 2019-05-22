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
	validateAndMergeTypeExtensions(ctx)

	// TODO: At this point, do we stop if we have errors? Or what?

	if schemaDef := ctx.SDLContext.SchemaDefinition; schemaDef != nil {
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
	}

	if !ctx.SDLContext.IsExtending {
		// The map on SDLContext is created with the right size to also contain the built-in
		// directives without needing to grow.
		ctx.Schema.Directives = ctx.SDLContext.DirectiveDefinitions

		for name, def := range graphql.SpecifiedDirectives() {
			// TODO: Should we allow overriding built-in directives?
			ctx.Schema.Directives[name] = def
		}

		// The map on SDLContext is created with the right size to also contain the built-in types
		// without needing to grow.
		ctx.Schema.Types = ctx.SDLContext.TypeDefinitions

		for name, def := range graphql.SpecifiedTypes() {
			// TODO: Should we allow overriding built-in directives?
			ctx.Schema.Types[name] = def
		}
	} else {
		// TODO: Handle adding the new types / directives to an existing schema.
	}
}
