package validation

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWalker(t *testing.T) {
	t.Run("should not return nil", func(t *testing.T) {
		assert.NotNil(t, NewWalker(nil))
	})
}

func TestWalker_Walk(t *testing.T) {
	var pass bool

	tt := []struct {
		name  string
		query []byte
		rules []RuleFunc
	}{
		{
			name:  "simple selection",
			query: []byte(`{ hello }`),
			rules: []RuleFunc{func(w *Walker) {
				w.AddSelectionEnterEventHandler(func(context *Context, selection ast.Selection) {
					pass = selection.Name == "hello"
				})
			}},
		},
		{
			name: "simple type system",
			query: []byte(`
				type Climbing implements
						& Ropes
						& Rocks
						& Chalk {
					name: String
					weight: Int
				}
			`),
			rules: []RuleFunc{func(w *Walker) {
				var count int
				w.AddTypeEnterEventHandler(func(context *Context, gt ast.Type) {
					count++
					pass = count == 5
				})
			}},
		},
		{
			name: "deeply nested directive",
			query: []byte(`
				{ foo { bar { baz { qux { quux { corge { uier { grault @garply } } } } } } } }
			`),
			rules: []RuleFunc{func(w *Walker) {
				w.AddDirectiveEnterEventHandler(func(context *Context, directive ast.Directive) {
					pass = directive.Name == "garply"
				})
			}},
		},
		{
			name: "item within list value",
			query: []byte(`
				{
					foo(list: ["bar"])
				}
			`),
			rules: []RuleFunc{func(w *Walker) {
				w.AddValueEnterEventHandler(func(context *Context, value ast.Value) {
					t.Log(value)
				})

				w.AddStringValueEnterEventHandler(func(context *Context, value ast.Value) {
					pass = value.StringValue == "bar"
					t.Log(value)
				})
			}},
		},
	}

	for _, tc := range tt {
		pass = false

		parser := language.NewParser(tc.query)

		doc, err := parser.Parse()
		require.NoError(t, err)

		ctx := &Context{}

		walker := NewWalker(tc.rules)
		walker.Walk(ctx, doc)

		assert.True(t, pass, "test case %q failed", tc.name)
	}
}
