package validation

import "github.com/bucketd/go-graphqlparser/ast"

// Walk ...
func (w *Walker) Walk(ctx *Context, doc ast.Document) {
	w.OnDocumentEnter(ctx, doc)
	w.walkDefinitions(ctx, doc.Definitions)
	w.OnDocumentLeave(ctx, doc)
}

// walkDefinitions ...
func (w *Walker) walkDefinitions(ctx *Context, defs *ast.Definitions) {
	w.OnDefinitionsEnter(ctx, defs)
	defs.ForEach(func(def ast.Definition, i int) {
		w.walkDefinition(ctx, def)
	})
	w.OnDefinitionsLeave(ctx, defs)
}

// walkDefinition ...
func (w *Walker) walkDefinition(ctx *Context, def ast.Definition) {
	w.OnDefinitionEnter(ctx, def)
	switch def.Kind {
	case ast.DefinitionKindExecutable:
		w.walkExecutableDefinition(ctx, def.ExecutableDefinition)
	case ast.DefinitionKindTypeSystem:
		w.walkTypeSystemDefinition(ctx, def.TypeSystemDefinition)
	case ast.DefinitionKindTypeSystemExtension:
		w.walkTypeSystemExtension(ctx, def.TypeSystemExtension)
	}
	w.OnDefinitionLeave(ctx, def)
}

// walkExecutableDefinition ...
func (w *Walker) walkExecutableDefinition(ctx *Context, def *ast.ExecutableDefinition) {
	w.OnExecutableDefinitionEnter(ctx, def)
	w.OnExecutableDefinitionLeave(ctx, def)
}

// walkTypeSystemDefinition ...
func (w *Walker) walkTypeSystemDefinition(ctx *Context, def *ast.TypeSystemDefinition) {
	w.OnTypeSystemDefinitionEnter(ctx, def)
	w.OnTypeSystemDefinitionLeave(ctx, def)
}

// walkTypeSystemExtension ...
func (w *Walker) walkTypeSystemExtension(ctx *Context, def *ast.TypeSystemExtension) {
	w.OnTypeSystemExtensionEnter(ctx, def)
	w.OnTypeSystemExtensionLeave(ctx, def)
}
