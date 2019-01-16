package validation

import (
	"sort"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newContext() *Context {
	return &Context{
		fragment:                       make(map[string]*ast.FragmentDefinition),
		fragmentSpreads:                make(map[*ast.Selections]map[string]bool),
		recursivelyReferencedFragments: make(map[string]map[string]bool),
		variableUsages:                 make(map[string]map[string]bool),
		recursiveVariableUsages:        make(map[string]map[string]bool),
	}
}

func TestSetFragment(t *testing.T) {
	ctx := newContext()
	visitFns := []VisitFunc{setFragment(ctx)}
	walker := NewWalker(visitFns)

	query := `
	query Foo($a: String, $b: String) {
		...FragA
	}
	
	fragment FragA on Type {
		field(a: $a) {
			...FragB
		}
	}
	
	fragment FragB on Type {
		field(b: $b)
	}
	`
	parser := language.NewParser([]byte(query))

	doc, err := parser.Parse()
	if err != nil {
		require.NoError(t, err)
	}

	walker.Walk(ctx, doc)

	seeking := "FragA"
	found := ctx.Fragment(seeking)
	assert.Equal(t, seeking, found.Name)
}

func TestSetFragmentSpreads(t *testing.T) {
	ctx := newContext()
	visitFns := []VisitFunc{setFragmentSpreads(ctx)}
	walker := NewWalker(visitFns)

	query := `
	query Foo($a: String, $b: String) {
		...Frag11
		field1 {
			...Frag21
			field2 {
				...Frag31
			}
		}
		...Frag12
	}
	`
	parser := language.NewParser([]byte(query))

	doc, err := parser.Parse()
	if err != nil {
		require.NoError(t, err)
	}

	walker.Walk(ctx, doc)

	var s1, s2, s3 *ast.Selections
	for ss, v := range ctx.fragmentSpreads {
		switch {
		case len(v) == 4:
			s1 = ss
		case len(v) == 2:
			s2 = ss
		case len(v) == 1:
			s3 = ss
		default:
			t.Fatal("Found unexpected selection")
		}
	}

	frags1, seen1 := ctx.FragmentSpreads(s1), make([]string, 0)
	for k := range frags1 {
		seen1 = append(seen1, k)
	}
	sort.Strings(seen1)

	assert.Equal(t, []string{
		"Frag11",
		"Frag12",
		"Frag21",
		"Frag31",
	}, seen1)

	assert.False(t, ctx.FragmentSpreads(s2)["Frag12"])
	assert.False(t, ctx.FragmentSpreads(s3)["Frag12"])
}

func TestSetVariableUsages(t *testing.T) {}

func TestSetRecursiveVariableUsages(t *testing.T) {}

func TestSetRecursivelyReferencedFragments(t *testing.T) {}
