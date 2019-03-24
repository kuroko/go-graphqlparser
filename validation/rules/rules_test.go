package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// schemaDocument ...
	schemaDocument = []byte(`
		schema {
			query: Query
		}

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

		type Query {
			checkEnumValueUniqueness: String!
		}
	`)
)

// ruleTestCase ...
type ruleTestCase struct {
	msg    string
	query  string
	schema *types.Schema
	errs   *types.Errors
}

// queryRuleTester ...
func queryRuleTester(t *testing.T, tt []ruleTestCase, fn validation.VisitFunc) {
	schema, errs, err := graphql.BuildSchema(nil, schemaDocument)
	require.NoError(t, err, "failed to build schema")
	require.Equal(t, (*types.Errors)(nil), errs, "failed to validate schema")

	for _, tc := range tt {
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		if err != nil {
			require.NoError(t, err)
		}

		walker := validation.NewWalker([]validation.VisitFunc{fn})

		errs := validation.Validate(doc, schema, walker)
		assert.Equal(t, tc.errs, errs, tc.msg)
	}
}

// sdlRuleTester ...
func sdlRuleTester(t *testing.T, tt []ruleTestCase, fn validation.VisitFunc) {
	for _, tc := range tt {
		schemaParser := language.NewParser([]byte(tc.query))

		schemaDoc, err := schemaParser.Parse()
		require.NoError(t, err, "failed to parse schema document")

		walker := validation.NewWalker([]validation.VisitFunc{fn})

		errs := validation.ValidateSDL(schemaDoc, tc.schema, walker)
		assert.Equal(t, tc.errs, errs, tc.msg)
	}
}
