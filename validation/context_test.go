package validation

import (
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

	walker.Walk(doc)

	seeking := "FragA"
	found := ctx.Fragment(seeking)
	assert.Equal(t, seeking, found.Name)
}

func TestSetVariableUsages(t *testing.T) {}

func TestSetRecursiveVariableUsages(t *testing.T) {}

func TestSetRecursivelyReferencedFragments(t *testing.T) {}

func TestSetFragmentSpreads(t *testing.T) {}
