package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// schemaDocument is a predefined schema that is used when testing query documents.
	schemaDocument = []byte(`
		interface Being {
			name(surname: Boolean): String
		}
		
		interface Pet {
			name(surname: Boolean): String
		}
		
		interface Canine {
			name(surname: Boolean): String
		}
		
		enum DogCommand {
			SIT
			HEEL
			DOWN
		}
		
		type Dog implements Being & Pet & Canine {
			name(surname: Boolean): String
			nickname: String
			barkVolume: Int
			barks: Boolean
			doesKnownCommand(dogCommand: DogCommand): Boolean
			isHousetrained(atOtherHomes: Boolean = true): Boolean
			isAtLocation(x: Int, y: Int): Boolean
		}
		
		type Cat implements Being & Pet {
			name(surname: Boolean): String
			nickname: String
			meows: Boolean
			meowVolume: Int
			furColor: FurColor
		}
		
		union CatOrDog = Dog | Cat
		
		interface Intelligent {
			iq: Int
		}
		
		type Human implements Being & Intelligent {
			name(surname: Boolean): String
			pets: [Pet]
			relatives: [Human]
			iq: Int
		}
		
		type Alien implements Being & Intelligent {
			iq: Int
			name(surname: Boolean): String
			numEyes: Int
		}
		
		union DogOrHuman = Dog | Human
		
		union HumanOrAlien = Human | Alien
		
		enum FurColor {
			BROWN
			BLACK
			TAN
			SPOTTED
			NO_FUR
			UNKNOWN
		}
		
		input ComplexInput {
			requiredField: Boolean!
			nonNullField: Boolean! = false
			intField: Int
			stringField: String
			booleanField: Boolean
			stringListField: [String]
		}
		
		type ComplicatedArgs {
			intArgField(intArg: Int): String
			nonNullIntArgField(nonNullIntArg: Int!): String
			stringArgField(stringArg: String): String
			booleanArgField(booleanArg: Boolean): String
			enumArgField(enumArg: FurColor): String
			floatArgField(floatArg: Float): String
			idArgField(idArg: ID): String
			stringListArgField(stringListArg: [String]): String
			stringListNonNullArgField(stringListNonNullArg: [String!]): String
			complexArgField(complexArg: ComplexInput): String
			multipleReqs(req1: Int!, req2: Int!): String
			nonNullFieldWithDefault(arg: Int! = 0): String
			multipleOpts(opt1: Int = 0, opt2: Int = 0): String
			multipleOptAndReq(req1: Int!, req2: Int!, opt1: Int = 0, opt2: Int = 0): String
		}
		
		scalar Invalid
		
		scalar Any
		
		type QueryRoot {
			human(id: ID): Human
			alien: Alien
			dog: Dog
			cat: Cat
			pet: Pet
			catOrDog: CatOrDog
			dogOrHuman: DogOrHuman
			humanOrAlien: HumanOrAlien
			complicatedArgs: ComplicatedArgs
			invalidArg(arg: Invalid): String
			anyArg(arg: Any): String
		}
		
		schema {
			query: QueryRoot
		}
		
		directive @onQuery on QUERY
		directive @onMutation on MUTATION
		directive @onSubscription on SUBSCRIPTION
		directive @onField on FIELD
		directive @onFragmentDefinition on FRAGMENT_DEFINITION
		directive @onFragmentSpread on FRAGMENT_SPREAD
		directive @onInlineFragment on INLINE_FRAGMENT
		directive @onVariableDefinition on VARIABLE_DEFINITION
	`)
)

// ruleTestCase ...
type ruleTestCase struct {
	msg    string
	query  string
	schema *graphql.Schema
	errs   *graphql.Errors
}

// queryRuleBencher ...
func queryRuleBencher(b *testing.B, t ruleTestCase, fn validation.VisitFunc) {
	schema, errs, err := buildSchema(nil, schemaDocument)
	require.NoError(b, err, "failed to build schema")
	require.Equal(b, (*graphql.Errors)(nil), errs, "failed to validate schema")

	parser := language.NewParser([]byte(t.query))
	doc, err := parser.Parse()
	require.NoError(b, err)

	walker := validation.NewWalker([]validation.VisitFunc{fn})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		walker.Walk(validation.NewContext(doc, schema), doc)
	}
}

// queryRuleTester ...
func queryRuleTester(t *testing.T, tt []ruleTestCase, fn validation.VisitFunc) {
	schema, errs, err := buildSchema(nil, schemaDocument)
	require.NoError(t, err, "failed to build schema")
	require.Equal(t, (*graphql.Errors)(nil), errs, "failed to validate schema")

	for _, tc := range tt {
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		require.NoError(t, err)

		walker := validation.NewWalker([]validation.VisitFunc{fn})

		ctx := validation.Validate(doc, schema, walker)

		// We need to sort errors, because we use maps in some places, and it leads to unpredictable
		// result error ordering.
		testErrs := graphql.SortErrors(tc.errs)
		rsltErrs := graphql.SortErrors(ctx.Errors)

		assert.Equal(t, testErrs, rsltErrs, tc.msg)
	}
}

// sdlRuleTester ...
func sdlRuleTester(t *testing.T, tt []ruleTestCase, fn validation.VisitFunc) {
	for _, tc := range tt {
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		require.NoError(t, err, "failed to parse schema document")

		walker := validation.NewWalker([]validation.VisitFunc{fn})

		ctx := validation.ValidateSDL(doc, tc.schema, walker)

		// We need to sort errors, because we use maps in some places, and it leads to unpredictable
		// result error ordering.
		testErrs := graphql.SortErrors(tc.errs)
		rsltErrs := graphql.SortErrors(ctx.Errors)

		assert.Equal(t, testErrs, rsltErrs, tc.msg)
	}
}
