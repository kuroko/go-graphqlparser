package main

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
)

func main() {
	parser := language.NewParser([]byte(`
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
			...FragB
			...FragC
			...FragB
			...FragC
			...FragB
			...FragC
			...FragB
			...FragC
			...FragB
			...FragC
			...FragD
		  }
		}
		query Foo($a: String, $b: String, $c: String) {
		  ...FragA
		}
		fragment FragC on Type {
		  field(c: $c)
		}
		fragment FragD on Type {
		  field(c: $c)
		}
	`))

	doc, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	ctx := validation.NewContext(doc)
	_ = ctx

	doc.Definitions.ForEach(func(d ast.Definition, i int) {
		// We want to see the fragments referenced on the query...
		if d.Kind != ast.DefinitionKindExecutable || d.ExecutableDefinition.Kind != ast.ExecutableDefinitionKindOperation {
			return
		}

		defs := ctx.RecursivelyReferencedFragments(d.ExecutableDefinition)
		_ = defs

		for k := range defs {
			fmt.Println(k.ExecutableDefinition.FragmentDefinition.Name)
		}

		//defs.ForEach(func(d ast.Definition, i int) {
		//	fmt.Println(d.ExecutableDefinition.FragmentDefinition.Name)
		//})
	})
}
