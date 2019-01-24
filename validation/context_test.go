package validation

import (
	"fmt"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func BenchmarkNewContext(b *testing.B) {
	parser := language.NewParser([]byte(`
		query Foo($a: String, $b: String, $c: String) {
		  ...FragA
		}
		fragment FragA on Type {
		  field(a: $a) {
			foo {
				bar {
					baz {
						...FragB
						...FragC
					}
				}
			}
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
	`))

	doc, err := parser.Parse()
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := NewContext(doc)
		_ = ctx
	}
}

func TestDump(t *testing.T) {
	parser := language.NewParser([]byte(`
		query Foo($a: String, $b: String, $c: String) {
		  ...FragA
		}
		fragment FragA on Type {
		  field(a: $a) {
			foo {
				bar {
					baz {
						...FragB
						...FragC
					}
				}
			}
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
	`))

	doc, err := parser.Parse()
	require.NoError(t, err)

	ctx := NewContext(doc)
	_ = ctx

	doc.Definitions.ForEach(func(d ast.Definition, i int) {
		// We want to see the fragments referenced on the query...
		if d.Kind != ast.DefinitionKindExecutable || d.ExecutableDefinition.Kind != ast.ExecutableDefinitionKindOperation {
			return
		}

		spew.Dump(ctx)

		defs := ctx.RecursivelyReferencedFragments(d.ExecutableDefinition)

		spew.Dump(d.ExecutableDefinition)

		defs.ForEach(func(d ast.Definition, i int) {
			fmt.Println(d.ExecutableDefinition.FragmentDefinition.Name)
		})
	})

	t.Fail()
}

//func TestSetFragment(t *testing.T) {
//	ctx := &Context{}
//	visitFns := []VisitFunc{setFragment}
//	walker := NewWalker(visitFns)
//
//	query := `
//	query Foo($a: String, $b: String) {
//		...FragA
//	}
//
//	fragment FragA on Type {
//		field(a: $a) {
//			...FragB
//		}
//	}
//
//	fragment FragB on Type {
//		field(b: $b)
//	}
//	`
//	parser := language.NewParser([]byte(query))
//
//	doc, err := parser.Parse()
//	if err != nil {
//		require.NoError(t, err)
//	}
//
//	walker.Walk(ctx, doc)
//
//	seeking := "FragA"
//	found := ctx.Fragment(seeking)
//	assert.Equal(t, seeking, found.Name)
//}
//
//func TestSetFragmentSpreads(t *testing.T) {
//	ctx := &Context{}
//	visitFns := []VisitFunc{setFragmentSpreads}
//	walker := NewWalker(visitFns)
//
//	query := `
//	query Foo($a: String, $b: String) {
//		...Frag11
//		field1 {
//			...Frag21
//			field2 {
//				...Frag31
//			}
//		}
//		...Frag12
//	}
//	`
//	parser := language.NewParser([]byte(query))
//
//	doc, err := parser.Parse()
//	if err != nil {
//		require.NoError(t, err)
//	}
//
//	walker.Walk(ctx, doc)
//
//	var s *ast.Selections
//	for ss, v := range ctx.fragmentSpreads {
//		if len(v) == 4 {
//			s = ss
//		} else {
//			t.Fatal("Found unexpected selection")
//		}
//	}
//
//	frags, seen := ctx.FragmentSpreads(s), make([]string, 0)
//	for k := range frags {
//		seen = append(seen, k)
//	}
//	sort.Strings(seen)
//
//	assert.Equal(t, []string{
//		"Frag11",
//		"Frag12",
//		"Frag21",
//		"Frag31",
//	}, seen)
//}
//
//func TestSetRecursivelyReferencedFragments(t *testing.T) {
//	ctx := &Context{}
//	visitFns := []VisitFunc{setDefinition, setFragmentSpreads, setRecursivelyReferencedFragments}
//	walker := NewWalker(visitFns)
//
//	query := `
//	fragment FragB on Type {
//	  field(b: $b) {
//	    ...FragC
//	  }
//	}
//	query Foo($a: String, $b: String, $c: String) {
//	  ...FragA
//	}
//	fragment FragC on Type {
//	  field(c: $c)
//	}
//	fragment FragA on Type {
//	  field(a: $a) {
//	    ...FragB
//	  }
//	}
//	`
//	parser := language.NewParser([]byte(query))
//
//	doc, err := parser.Parse()
//	if err != nil {
//		require.NoError(t, err)
//	}
//
//	walker.Walk(ctx, doc)
//
//	found := ctx.RecursivelyReferencedFragments("Foo")
//	assert.True(t, found["FragC"])
//}
//
//func TestSetVariableUsages(t *testing.T) {
//	ctx := &Context{}
//	visitFns := []VisitFunc{setDefinition, setVariableUsages}
//	walker := NewWalker(visitFns)
//
//	query := `
//	query Foo($a: String, $b: String, $c: String) {
//		... on Type {
//			field(a: $a) {
//				field(b: $b) {
//					... on Type {
//						field(c: $c)
//  				}
//  			}
//  		}
//  	}
//  }
//	`
//	parser := language.NewParser([]byte(query))
//
//	doc, err := parser.Parse()
//	if err != nil {
//		require.NoError(t, err)
//	}
//
//	walker.Walk(ctx, doc)
//
//	usages := ctx.VariableUsages("Foo")
//	_, found := usages["c"]
//
//	assert.True(t, found)
//}

func TestSetRecursiveVariableUsages(t *testing.T) {}
