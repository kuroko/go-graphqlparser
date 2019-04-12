package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

// schemaWithSDLDirectives ...
var schemaWithSDLDirectives = graphql.MustBuildSchema(nil, []byte(`
	directive @onSchema on SCHEMA
	directive @onScalar on SCALAR
	directive @onObject on OBJECT
	directive @onFieldDefinition on FIELD_DEFINITION
	directive @onArgumentDefinition on ARGUMENT_DEFINITION
	directive @onInterface on INTERFACE
	directive @onUnion on UNION
	directive @onEnum on ENUM
	directive @onEnumValue on ENUM_VALUE
	directive @onInputObject on INPUT_OBJECT
	directive @onInputFieldDefinition on INPUT_FIELD_DEFINITION
`))

func TestKnownDirectives(t *testing.T) {
	t.Run("query document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "with no directives",
				query: `
					query Foo {
						name
						...Frag
					}
					fragment Frag on Dog {
						name
					}
				`,
			},
			{
				msg: "with known directives",
				query: `
					{
						dog @include(if: true) {
							name
						}
						human @skip(if: false) {
							name
						}
					}
				`,
			},
			{
				msg: "with unknown directive",
				query: `
					{
						dog @unknown(directive: "value") {
							name
						}
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownDirectiveError("unknown", 0, 0)),
			},
			{
				msg: "with many unknown directives",
				query: `
					{
						dog @unknown(directive: "value") {
							name
						}
						human @unknown(directive: "value") {
							name
							pets @unknown(directive: "value") {
								name
							}
						}
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownDirectiveError("unknown", 0, 0)).
					Add(rules.UnknownDirectiveError("unknown", 0, 0)).
					Add(rules.UnknownDirectiveError("unknown", 0, 0)),
			},
			{
				msg: "with well placed directives",
				query: `
					query Foo($var: Boolean) @onQuery {
						name @include(if: $var)
						...Frag @include(if: true)
						skippedField @skip(if: true)
						...SkippedFrag @skip(if: true)
					}

					mutation Bar @onMutation {
						someField
					}
				`,
			},
			// TODO: This is in the working draft, not the 2018 spec.
			//{
			//	msg: "with well placed variable definition directive",
			//	query: `
			//		query Foo($var: Boolean @onVariableDefinition) {
			//			name
			//		}
			//	`,
			//},
			{
				msg: "with misplaced directives",
				query: `
					query Foo($var: Boolean) @include(if: true) {
						name @onQuery @include(if: $var)
						...Frag @onQuery
					}

					mutation Bar @onQuery {
						someField
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.MisplacedDirectiveError("include", ast.DirectiveLocationKindQuery, 0, 0)).
					Add(rules.MisplacedDirectiveError("onQuery", ast.DirectiveLocationKindField, 0, 0)).
					Add(rules.MisplacedDirectiveError("onQuery", ast.DirectiveLocationKindFragmentSpread, 0, 0)).
					Add(rules.MisplacedDirectiveError("onQuery", ast.DirectiveLocationKindMutation, 0, 0)),
			},
			// TODO: This is in the working draft, not the 2018 spec.
			//{
			//	msg: "with misplaced variable definition directive",
			//	query: `
			//		query Foo($var: Boolean @onField) {
			//			name
			//		}
			//	`,
			//	errs: (*types.Errors)(nil).
			//		Add(rules.MisplacedDirectiveError("onField", ast.DirectiveLocationKindVariableDefinition, 0, 0)),
			//},
		}

		queryRuleTester(t, tt, rules.KnownDirectives)
	})

	t.Run("sdl document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "with directive defined inside SDL",
				query: `
					type Query {
						foo: String @test
					}

					directive @test on FIELD_DEFINITION
				`,
			},
			{
				msg: "with standard directive",
				query: `
					type Query {
						foo: String @deprecated
					}
				`,
			},
			{
				msg: "with overridden standard directive",
				query: `
					schema @deprecated {
						query: Query
					}

					directive @deprecated on SCHEMA
				`,
			},
			{
				msg: "with directive defined in schema extension",
				schema: graphql.MustBuildSchema(nil, []byte(`
					type Query {
						foo: String
					}
				`)),
				query: `
					directive @test on OBJECT

					extend type Query @test
				`,
			},
			{
				msg: "with unknown directive in schema extension",
				schema: graphql.MustBuildSchema(nil, []byte(`
					type Query {
						foo: String
					}
				`)),
				query: `
					extend type Query @unknown
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownDirectiveError("unknown", 0, 0)),
			},
			{
				msg: "well placed on schema",
				query: `
					directive @onSchema on SCHEMA

					schema @onSchema {
						query: MyQuery
					}
				`,
			},
			{
				msg:    "with well placed directives",
				schema: schemaWithSDLDirectives,
				query: `
					type MyObj implements MyInterface @onObject {
						myField(myArg: Int @onArgumentDefinition): String @onFieldDefinition
					}

					extend type MyObj @onObject

					scalar MyScalar @onScalar

					extend scalar MyScalar @onScalar

					interface MyInterface @onInterface {
						myField(myArg: Int @onArgumentDefinition): String @onFieldDefinition
					}

					extend interface MyInterface @onInterface

					union MyUnion @onUnion = MyObj | Other

					extend union MyUnion @onUnion

					enum MyEnum @onEnum {
						MY_VALUE @onEnumValue
					}

					extend enum MyEnum @onEnum

					input MyInput @onInputObject {
						myField: Int @onInputFieldDefinition
					}

					extend input MyInput @onInputObject

					extend schema @onSchema
				`,
			},
			{
				msg:    "with misplaced directives",
				schema: schemaWithSDLDirectives,
				query: `
					type MyObj implements MyInterface @onInterface {
						myField(myArg: Int @onInputFieldDefinition): String @onInputFieldDefinition
					}

					scalar MyScalar @onEnum

					interface MyInterface @onObject {
						myField(myArg: Int @onInputFieldDefinition): String @onInputFieldDefinition
					}

					union MyUnion @onEnumValue = MyObj | Other

					enum MyEnum @onScalar {
						MY_VALUE @onUnion
					}

					input MyInput @onEnum {
						myField: Int @onArgumentDefinition
					}

					enum NotSchema @onObject {
						NOT_SCHEMA
					}

					extend schema @onObject
				`,
				errs: (*types.Errors)(nil).
					Add(rules.MisplacedDirectiveError("onInterface", ast.DirectiveLocationKindObject, 0, 0)).
					Add(rules.MisplacedDirectiveError("onInputFieldDefinition", ast.DirectiveLocationKindArgumentDefinition, 0, 0)).
					Add(rules.MisplacedDirectiveError("onInputFieldDefinition", ast.DirectiveLocationKindFieldDefinition, 0, 0)).
					Add(rules.MisplacedDirectiveError("onEnum", ast.DirectiveLocationKindScalar, 0, 0)).
					Add(rules.MisplacedDirectiveError("onObject", ast.DirectiveLocationKindInterface, 0, 0)).
					Add(rules.MisplacedDirectiveError("onInputFieldDefinition", ast.DirectiveLocationKindArgumentDefinition, 0, 0)).
					Add(rules.MisplacedDirectiveError("onInputFieldDefinition", ast.DirectiveLocationKindFieldDefinition, 0, 0)).
					Add(rules.MisplacedDirectiveError("onEnumValue", ast.DirectiveLocationKindUnion, 0, 0)).
					Add(rules.MisplacedDirectiveError("onScalar", ast.DirectiveLocationKindEnum, 0, 0)).
					Add(rules.MisplacedDirectiveError("onUnion", ast.DirectiveLocationKindEnumValue, 0, 0)).
					Add(rules.MisplacedDirectiveError("onEnum", ast.DirectiveLocationKindInputObject, 0, 0)).
					Add(rules.MisplacedDirectiveError("onArgumentDefinition", ast.DirectiveLocationKindInputFieldDefinition, 0, 0)).
					Add(rules.MisplacedDirectiveError("onObject", ast.DirectiveLocationKindEnum, 0, 0)).
					Add(rules.MisplacedDirectiveError("onObject", ast.DirectiveLocationKindSchema, 0, 0)),
			},
		}

		sdlRuleTester(t, tt, rules.KnownDirectives)
	})
}
