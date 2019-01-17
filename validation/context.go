package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var contextDecoratorWalker = NewWalker([]VisitFunc{
	setFragment,
	setVariableUsages,
	setRecursiveVariableUsages,
	setRecursivelyReferencedFragments,
	setFragmentSpreads,
})

// NewContext instantiates a validation context struct, this involves the walker
// doing a preliminary pass of the document, gathering basic information for the
// more complicated validation walk to come.
func NewContext(doc ast.Document) *Context {
	ctx := &Context{}

	contextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Errors *graphql.Errors
	Schema *graphql.Schema

	// Used by validation rules.
	VariableDefs *ast.VariableDefinitions

	// Internal pre-cached with methods to access.
	fragments                      map[string]*ast.FragmentDefinition
	fragmentSpreads                map[*ast.Selections]map[string]bool
	recursivelyReferencedFragments map[string]map[string]bool
	variableUsages                 map[string]map[string]bool
	recursiveVariableUsages        map[string]map[string]bool

	fragmentSpreadsSelectionSet *ast.Selections
	name                        string
	nameToSelections            map[string]*ast.Selections
}

// Fragment returns a FragmentDefinition by name.
func (ctx *Context) Fragment(fragName string) *ast.FragmentDefinition {
	return ctx.fragments[fragName]
}

func setFragment(w *Walker) {
	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		if ctx.fragments == nil {
			ctx.fragments = make(map[string]*ast.FragmentDefinition)
		}

		ctx.fragments[fd.Name] = fd
	})
}

// FragmentSpreads returns all nested usages of fragment spreads in this Selections.
func (ctx *Context) FragmentSpreads(ss *ast.Selections) map[string]bool {
	return ctx.fragmentSpreads[ss]
}

func setFragmentSpreads(w *Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *Context, od *ast.OperationDefinition) {
		ctx.fragmentSpreadsSelectionSet = od.SelectionSet
	})

	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		ctx.fragmentSpreadsSelectionSet = fd.SelectionSet
	})

	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		if ctx.fragmentSpreads == nil {
			ctx.fragmentSpreads = make(map[*ast.Selections]map[string]bool)
		}

		if ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet] == nil {
			ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet] = make(map[string]bool)
		}

		ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet][s.Name] = true
	})
}

// RecursivelyReferencedFragments returns all the recursively referenced
// fragments used by an operation or fragment definition.
func (ctx *Context) RecursivelyReferencedFragments(exDefName string) map[string]bool {
	return ctx.recursivelyReferencedFragments[exDefName]
}

func setRecursivelyReferencedFragments(w *Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *Context, od *ast.OperationDefinition) {
		ctx.name = od.Name
	})

	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		ctx.name = fd.Name
	})

	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		if ctx.nameToSelections == nil {
			ctx.nameToSelections = make(map[string]*ast.Selections)
		}

		ctx.nameToSelections[ctx.name] = ctx.fragmentSpreadsSelectionSet
	})

	w.AddDocumentLeaveEventHandler(func(ctx *Context, d ast.Document) {
		for exDefName := range ctx.nameToSelections {

			if _, ok := ctx.recursivelyReferencedFragments[exDefName]; ok {
				continue
			}

			_ = set(ctx, exDefName, []string{exDefName})
		}
	})
}

func set(ctx *Context, name string, parents []string) []string {
	var children []string
	for frag := range ctx.FragmentSpreads(ctx.nameToSelections[name]) {
		if in(frag, parents) {
			continue
		}
		c := set(
			ctx,
			frag,
			append(parents, name),
		)
		children = append(children, c...)
	}

	ctx.recursivelyReferencedFragments[name] = mappend(ctx.FragmentSpreads(ctx.nameToSelections[name]))

	for _, child := range children {
		quzz := mappend(
			ctx.recursivelyReferencedFragments[name],
			ctx.FragmentSpreads(ctx.nameToSelections[child]),
		)
		ctx.recursivelyReferencedFragments[name] = quzz
	}

	return append(children, name)
}

func in(item string, list []string) bool {
	for _, str := range list {
		if str == item {
			return true
		}
	}
	return false
}

func mappend(maps ...map[string]bool) map[string]bool {
	if len(maps) < 1 {
		return make(map[string]bool)
	}

	if maps[0] == nil {
		maps[0] = make(map[string]bool)
	}

	if len(maps) < 2 {
		return maps[0]
	}

	for _, m := range maps[1:] {
		for k := range m {
			maps[0][k] = true
		}
	}
	return maps[0]
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(exDefName string) map[string]bool {
	return ctx.variableUsages[exDefName]
}

func setVariableUsages(w *Walker) {

}

// RecursiveVariableUsages returns all recursively referenced variable usages for an operation.
func (ctx *Context) RecursiveVariableUsages(opName string) map[string]bool {
	return ctx.recursiveVariableUsages[opName]
}

func setRecursiveVariableUsages(w *Walker) {

}
