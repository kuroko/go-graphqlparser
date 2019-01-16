package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// NewContext instantiates a validation context struct, this involves the walker
// doing a preliminary pass of the document, gathing basic information for the
// more complicated validation walk to come.
func NewContext(doc ast.Document) *Context {
	ctx := &Context{
		fragment:                       make(map[string]*ast.FragmentDefinition),
		fragmentSpreads:                make(map[*ast.Selections]map[string]bool),
		recursivelyReferencedFragments: make(map[string]map[string]bool),
		variableUsages:                 make(map[string]map[string]bool),
		recursiveVariableUsages:        make(map[string]map[string]bool),
	}

	visitFns := []VisitFunc{
		setFragment(ctx),
		setVariableUsages(ctx),
		setRecursiveVariableUsages(ctx),
		setRecursivelyReferencedFragments(ctx),
		setFragmentSpreads(ctx),
	}

	walker := NewWalker(visitFns)
	walker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Errors *graphql.Errors
	Schema *graphql.Schema

	fragment                       map[string]*ast.FragmentDefinition
	fragmentSpreads                map[*ast.Selections]map[string]bool
	recursivelyReferencedFragments map[string]map[string]bool
	variableUsages                 map[string]map[string]bool
	recursiveVariableUsages        map[string]map[string]bool
}

// Fragment returns a FragmentDefinition by name.
func (ctx *Context) Fragment(fragName string) *ast.FragmentDefinition {
	return ctx.fragment[fragName]
}

func setFragment(ctx *Context) VisitFunc {
	return func(w *Walker) {
		w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
			ctx.fragment[fd.Name] = fd
		})
	}
}

// FragmentSpreads returns all nested usages of fragment spreads in this Selections.
func (ctx *Context) FragmentSpreads(ss *ast.Selections) map[string]bool {
	return ctx.fragmentSpreads[ss]
}

func setFragmentSpreads(ctx *Context) VisitFunc {
	return func(w *Walker) {
		var parents []*ast.Selections

		w.AddSelectionsEnterEventHandler(func(ctx *Context, ss *ast.Selections) {
			ctx.fragmentSpreads[ss] = make(map[string]bool)
			parents = append(parents, ss)
		})

		w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
			for _, p := range parents {
				ctx.fragmentSpreads[p][s.Name] = true
			}
		})

		w.AddSelectionsLeaveEventHandler(func(ctx *Context, ss *ast.Selections) {
			parents = parents[:len(parents)-1]
		})
	}
}

// RecursivelyReferencedFragments returns all the recursively referenced
// fragments used by an operation or fragment definition.
func (ctx *Context) RecursivelyReferencedFragments(exDefName string) map[string]bool {
	return ctx.recursivelyReferencedFragments[exDefName]
}

func setRecursivelyReferencedFragments(ctx *Context) VisitFunc {
	return func(w *Walker) {}
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(exDefName string) map[string]bool {
	return ctx.variableUsages[exDefName]
}

func setVariableUsages(ctx *Context) VisitFunc {
	return func(w *Walker) {}
}

// RecursiveVariableUsages returns all recursively referenced variable usages for an operation.
func (ctx *Context) RecursiveVariableUsages(opName string) map[string]bool {
	return ctx.recursiveVariableUsages[opName]
}

func setRecursiveVariableUsages(ctx *Context) VisitFunc {
	return func(w *Walker) {}
}
