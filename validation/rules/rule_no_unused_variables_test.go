package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoUnusedVariables(t *testing.T) {
	tt := []struct {
		msg   string
		query string
		errs  *graphql.Errors
	}{
		{
			msg: "uses all variables",
			query: `
			query ($a: String, $b: String, $c: String) {
			  field(a: $a, b: $b, c: $c)
			}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables deeply",
			query: `
			query Foo($a: String, $b: String, $c: String) {
			  field(a: $a) {
					field(b: $b) {
				  	field(c: $c)
					}
			  }
			}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables deeply in inline fragments",
			query: `
			query Foo($a: String, $b: String, $c: String) {
			  	... on Type {
						field(a: $a) {
					  	field(b: $b) {
								... on Type {
							  	field(c: $c)
							}
					  }
					}
			  }
			}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables in fragments",
			query: ` 
			query Foo($a: String, $b: String, $c: String) {
			  ...FragA
			}
			fragment FragA on Type {
			  field(a: $a) {
					...FragB
			  }
			}
			fragment FragB on Type {
			  field(b: $b) {
					...FragC
			  }
			}
			fragment FragC on Type {
			  field(c: $c)
			}
			`,
			errs: nil,
		},
		{
			msg: "variable used by fragment in multiple operations",
			query: `
			query Foo($a: String) {
			  ...FragA
			}
			query Bar($b: String) {
			  ...FragB
			}
			fragment FragA on Type {
			  field(a: $a)
			}
			fragment FragB on Type {
			  field(b: $b)
			}
			`,
			errs: nil,
		},
		{
			msg: "variable used by recursive fragment",
			query: `
			query Foo($a: String) {
			  ...FragA
			}
			fragment FragA on Type {
			  field(a: $a) {
				...FragA
			  }
			}
			`,
			errs: nil,
		},
		{
			msg: "variable not used",
			query: `
			query ($a: String, $b: String, $c: String) {
			  field(a: $a, b: $b)
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("c", "", 0, 0)),
		},
		{
			msg: "multiple variables not used",
			query: `
			query Foo($a: String, $b: String, $c: String) {
			  field(b: $b)
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("a", "Foo", 0, 0)).
				Add(unusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "variable not used in fragments",
			query: `
			query Foo($a: String, $b: String, $c: String) {
			  ...FragA
			}
			fragment FragA on Type {
			  field(a: $a) {
				...FragB
			  }
			}
			fragment FragB on Type {
			  field(b: $b) {
				...FragC
			  }
			}
			fragment FragC on Type {
			  field
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "multiple variables not used in fragments",
			query: `
			query Foo($a: String, $b: String, $c: String) {
			  ...FragA
			}
			fragment FragA on Type {
			  field {
				...FragB
			  }
			}
			fragment FragB on Type {
			  field(b: $b) {
				...FragC
			  }
			}
			fragment FragC on Type {
			  field
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("a", "Foo", 0, 0)).
				Add(unusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "variable not used by unreferenced fragment",
			query: `
			query Foo($b: String) {
			  ...FragA
			}
			fragment FragA on Type {
			  field(a: $a)
			}
			fragment FragB on Type {
			  field(b: $b)
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("b", "Foo", 0, 0)),
		},
		{
			msg: "variable not used by fragment used by other operation",
			query: `
			query Foo($b: String) {
			  ...FragA
			}
			query Bar($a: String) {
			  ...FragB
			}
			fragment FragA on Type {
			  field(a: $a)
			}
			fragment FragB on Type {
			  field(b: $b)
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, unusedVariableError("b", "Foo", 0, 0)).
				Add(unusedVariableError("a", "Bar", 0, 0)),
		},
	}

	for _, tc := range tt {
		parser := language.NewParser([]byte(tc.query))

		doc, err := parser.Parse()
		if err != nil {
			require.NoError(t, err)
		}

		walker := validation.NewWalker([]validation.VisitFunc{noUnusedVariables})

		errs := validation.Validate(doc, walker)
		assert.Equal(t, tc.errs, errs, tc.msg)
	}
}
