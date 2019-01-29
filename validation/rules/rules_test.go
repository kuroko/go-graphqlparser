package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ruleTestCase struct {
	msg   string
	query string
	errs  *graphql.Errors
}

func ruleTester(t *testing.T, tt []ruleTestCase, fn validation.VisitFunc) {
	for _, tc := range tt {
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		if err != nil {
			require.NoError(t, err)
		}

		walker := validation.NewWalker([]validation.VisitFunc{loneSchemaDefinition})

		errs := validation.Validate(doc, walker)
		assert.Equal(t, tc.errs, errs, tc.msg)
	}
}
