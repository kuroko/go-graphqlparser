package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var (
	// queryContextDecoratorWalker ...
	queryContextDecoratorWalker = NewWalker([]VisitFunc{
		setFragments,
		setExecutableDefinition,
		setReferencedFragments,
		setVariableUsages,
	})
)

// NewContext instantiates a validation context struct, this involves the walker doing a
// preliminary pass of a query document, gathering basic information for the more complicated
// validation walk to come.
func NewContext(doc ast.Document, schema *graphql.Schema) *Context {
	if schema == nil {
		schema = &graphql.Schema{}
	}

	ctx := &Context{
		Document: doc,
		Schema:   schema,
	}

	queryContextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// NewSDLContext ...
func NewSDLContext(doc ast.Document, schema *graphql.Schema) *Context {
	isExtending := schema != nil
	if !isExtending {
		schema = &graphql.Schema{}
	}

	ctx := &Context{
		Document: doc,
		Schema:   schema,
	}

	// Construct SDL specific structures.
	ctx.SDLContext = &SDLContext{
		IsExtending: isExtending,
	}

	// Perform some initial validation, set up some data structures for further validation.
	PrepareContextSDL(ctx)

	return ctx
}

// Context ...
type Context struct {
	Document ast.Document
	Errors   *graphql.Errors
	Schema   *graphql.Schema

	// Used if we're validating an SDL file. This contains state for validating SDL documents, along
	// with symbol tables for definitions that can only be used in SDL documents.
	SDLContext *SDLContext

	// OperationDefinitions contains all operation definitions found in the input document, by name,
	// except for anonymous operations (i.e. shorthand queries), which will be stored under the key
	// `$$query` - which is an otherwise invalid name.
	OperationDefinitions map[string]*ast.OperationDefinition

	// FragmentDefinitions contains all fragment definitions found in the input document, by name.
	FragmentDefinitions map[string]*ast.FragmentDefinition

	// referencedFragments stores the fragment definitions referenced directly by an executable
	// definition, i.e. this is not recursively referenced FragmentDefinitions.
	referencedFragments map[*ast.ExecutableDefinition][]ast.Definition

	// variableUsages stores the variable usages referenced directly by an executable definition,
	// i.e. this is not recursive variable usages.
	variableUsages map[*ast.ExecutableDefinition][]string

	// executableDefinition is the current executable definition being walked over.
	executableDefinition *ast.ExecutableDefinition
}

// AddError adds an error to the linked list of errors on this Context.
func (ctx *Context) AddError(err graphql.Error) {
	ctx.Errors = ctx.Errors.Add(err)
}

// DirectiveDefinition ...
func (ctx *Context) DirectiveDefinition(name string) (*ast.DirectiveDefinition, bool) {
	var dirDef *ast.DirectiveDefinition
	var isInSchema bool

	if ctx.SDLContext.IsExtending {
		// If we're extending a schema, we might be extending a type that's defined there.
		dirDef, isInSchema = ctx.Schema.Directives[name]
	}

	if !isInSchema {
		// If we couldn't find it in the schema, or it's not possible for it to exist there, then
		// check in the current document.
		dirDef = ctx.SDLContext.DirectiveDefinitions[name]
	}

	if dirDef == nil {
		dirDef = graphql.SpecifiedDirectives()[name]
	}

	return dirDef, isInSchema
}

// Fragment ...
func (ctx *Context) Fragment(name string) *ast.FragmentDefinition {
	return ctx.FragmentDefinitions[name]
}

// TypeDefinition ...
func (ctx *Context) TypeDefinition(name string) (*ast.TypeDefinition, bool) {
	var typeDef *ast.TypeDefinition
	var isInSchema bool

	if ctx.SDLContext.IsExtending {
		// If we're extending a schema, we might be extending a type that's defined there.
		typeDef, isInSchema = ctx.Schema.Types[name]
	}

	if !isInSchema {
		// If we couldn't find it in the schema, or it's not possible for it to exist there, then
		// check in the current document.
		typeDef = ctx.SDLContext.TypeDefinitions[name]
	}

	return typeDef, isInSchema
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(def *ast.ExecutableDefinition) []string {
	return ctx.variableUsages[def]
}

// RecursiveVariableUsages ...
func (ctx *Context) RecursiveVariableUsages(def *ast.ExecutableDefinition) map[string]struct{} {
	// TODO: Maybe we could make this a slice too?.. and maybe we should, given VariableUsages
	// returns a string slice too, this is pretty inconsistent.
	result := make(map[string]struct{})

	ctx.recursiveVariableUsagesIter(def, result, make(map[*ast.ExecutableDefinition]struct{}))

	return result
}

// recursiveVariableUsagesIter ...
func (ctx *Context) recursiveVariableUsagesIter(def *ast.ExecutableDefinition, agg map[string]struct{}, seen map[*ast.ExecutableDefinition]struct{}) {
	for _, vu := range ctx.variableUsages[def] {
		agg[vu] = struct{}{}
	}

	// TODO: Can this be swapped to use a cached version of recursively referenced FragmentDefinitions.
	for _, rd := range ctx.referencedFragments[def] {
		// We only want to recurse deeper if we've never seen the fragment before.
		if _, ok := seen[rd.ExecutableDefinition]; !ok {
			seen[rd.ExecutableDefinition] = struct{}{}
			ctx.recursiveVariableUsagesIter(rd.ExecutableDefinition, agg, seen)
		}
	}
}

// ReferencedFragments returns the FragmentDefinitions directly referenced by the given executable definition.
func (ctx *Context) ReferencedFragments(def *ast.ExecutableDefinition) []ast.Definition {
	return ctx.referencedFragments[def]
}

// RecursivelyReferencedFragments ...
// NOTE: In practice, def would likely be an operation definition, but this isn't a requirement.
func (ctx *Context) RecursivelyReferencedFragments(def *ast.ExecutableDefinition) map[ast.Definition]struct{} {
	if ctx.referencedFragments[def] == nil || len(ctx.referencedFragments[def]) == 0 {
		return nil
	}

	// Maybe we could make this a slice too?
	result := make(map[ast.Definition]struct{})

	ctx.recursivelyReferencedFragmentsIter(def, result, make(map[*ast.ExecutableDefinition]struct{}))

	return result
}

// recursivelyReferencedFragmentsIter is the inner iteration method for finding recursively
// referenced FragmentDefinitions for a given executable definition. It modifies the given aggregate map of
// results.
func (ctx *Context) recursivelyReferencedFragmentsIter(def *ast.ExecutableDefinition, agg map[ast.Definition]struct{}, seen map[*ast.ExecutableDefinition]struct{}) {
	// For each referenced fragment in the current executable definition...
	for _, rd := range ctx.referencedFragments[def] {
		agg[rd] = struct{}{}

		// We only want to recurse deeper if we've never seen the fragment before.
		if _, ok := seen[rd.ExecutableDefinition]; !ok {
			seen[rd.ExecutableDefinition] = struct{}{}
			ctx.recursivelyReferencedFragmentsIter(rd.ExecutableDefinition, agg, seen)
		}
	}
}

// SDLContext ...
type SDLContext struct {
	DirectiveDefinitions map[string]*ast.DirectiveDefinition
	TypeDefinitions      map[string]*ast.TypeDefinition
	TypeExtensions       map[string][]*ast.TypeExtension
	SchemaDefinition     *ast.SchemaDefinition
	SchemaExtensions     []*ast.SchemaExtension

	KnownEnumValueNames map[string]map[string]struct{}
	KnownFieldNames     map[string]map[string]struct{}

	QueryTypeDefined        bool
	MutationTypeDefined     bool
	SubscriptionTypeDefined bool

	// IsExtending is true if this context was created with an existing Schema, and it's being
	// extended by another SDL file.
	IsExtending bool
}

// setExecutableDefinition ...
func setExecutableDefinition(w *Walker) {
	w.AddExecutableDefinitionEnterEventHandler(func(ctx *Context, def *ast.ExecutableDefinition) {
		ctx.executableDefinition = def
	})
}

// setFragments ...
func setFragments(w *Walker) {
	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, def *ast.FragmentDefinition) {
		if ctx.FragmentDefinitions == nil {
			ctx.FragmentDefinitions = make(map[string]*ast.FragmentDefinition, ctx.Document.FragmentDefinitions)
		}

		ctx.FragmentDefinitions[def.Name] = def
	})
}

// setReferencedFragments ...
func setReferencedFragments(w *Walker) {
	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		ctx.Document.Definitions.ForEach(func(d ast.Definition, i int) {
			if d.Kind != ast.DefinitionKindExecutable {
				return
			}

			if d.ExecutableDefinition.Kind != ast.ExecutableDefinitionKindFragment {
				return
			}

			if d.ExecutableDefinition.FragmentDefinition.Name == s.Name {
				if ctx.referencedFragments == nil {
					ctx.referencedFragments = make(map[*ast.ExecutableDefinition][]ast.Definition)
				}

				for _, v := range ctx.referencedFragments[ctx.executableDefinition] {
					if v == d {
						return
					}
				}

				ctx.referencedFragments[ctx.executableDefinition] =
					append(ctx.referencedFragments[ctx.executableDefinition], d)
			}
		})
	})
}

// setVariableUsages ...
func setVariableUsages(w *Walker) {
	w.AddVariableValueEnterEventHandler(func(ctx *Context, v ast.Value) {
		if ctx.variableUsages == nil {
			ctx.variableUsages = make(map[*ast.ExecutableDefinition][]string)
		}

		for _, u := range ctx.variableUsages[ctx.executableDefinition] {
			if u == v.StringValue {
				return
			}
		}

		ctx.variableUsages[ctx.executableDefinition] = append(ctx.variableUsages[ctx.executableDefinition], v.StringValue)
	})
}

// PrepareContextSDL is a function that populates information needed prior to walking the AST. It
// also performs some validation that would be far more inefficient if we were to walk (e.g. we'd
// either need to allocate another map and duplicate data, or we'd need to walk the whole AST twice.
// Instead, we manually touch the specific portions of the AST we need in one go. Currently we only
// want to use quite shallow data in from the AST at this phase, so using the Walker would be quite
// inefficient as it would hit many leaf nodes at a depth we simply don't need here).
func PrepareContextSDL(ctx *Context) {
	// Prepare maps, we always need these, at least for built-in types, etc.
	if ctx.OperationDefinitions == nil {
		ctx.OperationDefinitions = make(map[string]*ast.OperationDefinition, ctx.Document.OperationDefinitions)
	}

	if ctx.FragmentDefinitions == nil {
		ctx.FragmentDefinitions = make(map[string]*ast.FragmentDefinition, ctx.Document.FragmentDefinitions)
	}

	if ctx.SDLContext.DirectiveDefinitions == nil {
		size := int(ctx.Document.DirectiveDefinitions) + len(graphql.SpecifiedDirectives())
		ctx.SDLContext.DirectiveDefinitions = make(map[string]*ast.DirectiveDefinition, size)
	}

	if ctx.SDLContext.TypeDefinitions == nil {
		size := int(ctx.Document.TypeDefinitions) + len(graphql.SpecifiedTypes())
		ctx.SDLContext.TypeDefinitions = make(map[string]*ast.TypeDefinition, size)
	}

	if ctx.SDLContext.SchemaExtensions == nil {
		ctx.SDLContext.SchemaExtensions = make([]*ast.SchemaExtension, 0, ctx.Document.SchemaExtensions)
	}

	if ctx.SDLContext.TypeExtensions == nil {
		ctx.SDLContext.TypeExtensions = make(map[string][]*ast.TypeExtension, ctx.Document.TypeExtensions)
	}

	ctx.Document.Definitions.ForEach(func(def ast.Definition, i int) {
		switch def.Kind {
		case ast.DefinitionKindExecutable:
			switch def.ExecutableDefinition.Kind {
			case ast.ExecutableDefinitionKindOperation:
				odef := def.ExecutableDefinition.OperationDefinition
				if odef.Name == "" {
					odef.Name = "$$query"
				}

				ctx.OperationDefinitions[odef.Name] = odef

			case ast.ExecutableDefinitionKindFragment:
				fdef := def.ExecutableDefinition.FragmentDefinition

				ctx.FragmentDefinitions[fdef.Name] = fdef
			}

		case ast.DefinitionKindTypeSystem:
			switch def.TypeSystemDefinition.Kind {
			// UniqueDirectiveNames:
			case ast.TypeSystemDefinitionKindDirective:
				ddef := def.TypeSystemDefinition.DirectiveDefinition

				if _, ok := ctx.Schema.Directives[ddef.Name]; ok {
					ctx.AddError(ExistedDirectiveNameError(ddef.Name, 0, 0))
					return
				}

				if _, ok := ctx.SDLContext.DirectiveDefinitions[ddef.Name]; ok {
					ctx.AddError(DuplicateDirectiveNameError(ddef.Name, 0, 0))
				} else {
					ctx.SDLContext.DirectiveDefinitions[ddef.Name] = ddef
				}

			// LoneSchemaDefinition:
			case ast.TypeSystemDefinitionKindSchema:
				if ctx.SDLContext.IsExtending {
					// There cannot be any schema definitions in schema extensions, as either one will have
					// already been defined, or one will have been automatically created.
					ctx.AddError(CanNotDefineSchemaWithinExtensionError(0, 0))
					return
				} else if !ctx.SDLContext.IsExtending && ctx.SDLContext.SchemaDefinition != nil {
					// There should only be one schema definition in schema document, when not extending.
					ctx.AddError(SchemaDefinitionNotAloneError(0, 0))
					return
				}

				sdef := def.TypeSystemDefinition.SchemaDefinition

				// TODO: Validate here.
				ctx.SDLContext.SchemaDefinition = sdef

			// UniqueTypeNames:
			case ast.TypeSystemDefinitionKindType:
				tdef := def.TypeSystemDefinition.TypeDefinition

				if _, ok := ctx.Schema.Types[tdef.Name]; ok {
					ctx.AddError(ExistedTypeNameError(tdef.Name, 0, 0))
					return
				}

				if _, ok := ctx.SDLContext.TypeDefinitions[tdef.Name]; ok {
					ctx.AddError(DuplicateTypeNameError(tdef.Name, 0, 0))
				} else {
					ctx.SDLContext.TypeDefinitions[tdef.Name] = tdef
				}
			}

		case ast.DefinitionKindTypeSystemExtension:
			switch def.TypeSystemExtension.Kind {
			case ast.TypeSystemExtensionKindSchema:
				ext := def.TypeSystemExtension.SchemaExtension

				ctx.SDLContext.SchemaExtensions = append(ctx.SDLContext.SchemaExtensions, ext)

			case ast.TypeSystemExtensionKindType:
				ext := def.TypeSystemExtension.TypeExtension

				ctx.SDLContext.TypeExtensions[ext.Name] = append(ctx.SDLContext.TypeExtensions[ext.Name], ext)
			}
		}
	})
}
