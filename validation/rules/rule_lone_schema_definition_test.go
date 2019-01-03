package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: when schema is implemented add tests for canNotDefineSchemaWithinExtensionMessage.
func TestLoneSchemaDefinition(t *testing.T) {
	tt := []struct {
		msg   string
		query string
		errs  *graphql.Errors
	}{
		{
			msg: "no schema",
			query: `
			type Foo { foo: String }
			`,
			errs: nil,
		},
		{
			msg: "one schema definition",
			query: `
			schema { query: Foo }
			type Foo { foo: String }
			`,
			errs: nil,
		},
		{
			msg: "multiple schema definitions",
			query: `
			schema { query: Foo }
			type Foo { foo: String }
			schema { mutation: Foo }
			schema { subscription: Foo }
			`,
			errs: (*graphql.Errors).
				Add(nil, schemaDefinitionNotAlone(0, 0)).
				Add(schemaDefinitionNotAlone(0, 0)),
		},
	}

	for _, tc := range tt {
		walker := validation.NewWalker([]validation.RuleFunc{loneSchemaDefinition})
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		if err != nil {
			require.NoError(t, err)
		}

		errs := validation.Validate(doc, walker)
		assert.Equal(t, tc.errs, errs)
	}
}
