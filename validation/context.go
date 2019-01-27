package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var contextDecoratorWalker = NewWalker([]VisitFunc{
	setExecutableDefinition,
	setReferencedFragments,
	setVariableUsages,
})

// NewContext instantiates a validation context struct, this involves the walker
// doing a preliminary pass of the document, gathering basic information for the
// more complicated validation walk to come.
func NewContext(doc ast.Document) *Context {
	ctx := &Context{
		document: doc,
	}

	contextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Errors   *graphql.Errors
	Schema   *graphql.Schema
	document ast.Document

	// Used by validation rules.
	VariableDefs *ast.VariableDefinitions

	// Internal pre-cached with methods to access.
	variableUsages map[*ast.ExecutableDefinition]map[string]bool

	// Maybe this is all we need here? This would store directly referenced fragments, as a list of
	// definitions. By storing them as pointers, do we save memory by re-using the existing memory
	// we've allocated during parsing?
	referencedFragments map[*ast.ExecutableDefinition]*ast.Definitions

	// Maybe this doesn't need to exist?
	recursivelyReferencedFragments map[*ast.Definition]*ast.Definitions

	executableDefinition *ast.ExecutableDefinition
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(def *ast.ExecutableDefinition) map[string]bool {
	return ctx.variableUsages[def]
}

// NOTE: In practice, def would likely be an operation definition, but this isn't a requirement.
func (ctx *Context) RecursivelyReferencedFragments(def *ast.ExecutableDefinition) *ast.Definitions {
	if ctx.referencedFragments[def] == nil || ctx.referencedFragments[def].Len() == 0 {
		return nil
	}

	var result *ast.Definitions

	ctx.referencedFragments[def].ForEach(func(d ast.Definition, _ int) {
		// TODO: Is this necessary? We should only put in valid data to the map.
		if d.Kind != ast.DefinitionKindExecutable {
			return
		}

		result = result.Add(d)
		result.Join(ctx.RecursivelyReferencedFragments(d.ExecutableDefinition))
	})

	return result
}

func setExecutableDefinition(w *Walker) {
	w.AddExecutableDefinitionEnterEventHandler(func(ctx *Context, def *ast.ExecutableDefinition) {
		ctx.executableDefinition = def
	})
}

func setReferencedFragments(w *Walker) {
	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		ctx.document.Definitions.ForEach(func(d ast.Definition, i int) {
			if d.Kind != ast.DefinitionKindExecutable {
				return
			}

			if d.ExecutableDefinition.Kind != ast.ExecutableDefinitionKindFragment {
				return
			}

			if d.ExecutableDefinition.FragmentDefinition.Name == s.Name {
				if ctx.referencedFragments == nil {
					ctx.referencedFragments = make(map[*ast.ExecutableDefinition]*ast.Definitions)
				}

				ctx.referencedFragments[ctx.executableDefinition] =
					ctx.referencedFragments[ctx.executableDefinition].Add(d)
			}
		})
	})

	//w.AddDocumentLeaveEventHandler(func(ctx *Context, doc ast.Document) {
	//	spew.Dump(ctx.referencedFragments)
	//})
}

func setVariableUsages(w *Walker) {
	w.AddVariableValueEnterEventHandler(func(ctx *Context, v ast.Value) {
		if ctx.variableUsages == nil {
			ctx.variableUsages = make(map[*ast.ExecutableDefinition]map[string]bool)
		}

		if ctx.variableUsages[ctx.executableDefinition] == nil {
			ctx.variableUsages[ctx.executableDefinition] = make(map[string]bool)
		}

		ctx.variableUsages[ctx.executableDefinition][v.StringValue] = true
	})
}
