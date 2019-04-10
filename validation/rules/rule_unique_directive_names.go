package rules

// UniqueDirectiveNames ...
//func UniqueDirectiveNames(w *validation.Walker) {
//	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
//		if ctx.SDLContext.DirectiveDefinitions == nil {
//			ctx.SDLContext.DirectiveDefinitions = make(map[string]*ast.DirectiveDefinition, ctx.Document.DirectiveDefinitions)
//		}
//
//		directiveName := def.Name
//
//		if _, ok := ctx.Schema.Directives[directiveName]; ok {
//			ctx.AddError(ExistedDirectiveNameError(directiveName, 0, 0))
//			return
//		}
//
//		if _, ok := ctx.SDLContext.DirectiveDefinitions[directiveName]; ok {
//			ctx.AddError(DuplicateDirectiveNameError(directiveName, 0, 0))
//		} else {
//			ctx.SDLContext.DirectiveDefinitions[directiveName] = def
//		}
//	})
//}
//
//// DuplicateDirectiveNameError ...
//func DuplicateDirectiveNameError(directiveName string, line, col int) types.Error {
//	return types.NewError(
//		"There can be only one directive named \"" + directiveName + "\".",
//		// TODO: Location.
//	)
//}
//
//// ExistedDirectiveNameError ...
//func ExistedDirectiveNameError(directiveName string, line, col int) types.Error {
//	return types.NewError(
//		"Directive \"" + directiveName + "\" already exists in the schema. It cannot be redefined.",
//		// TODO: Location.
//	)
//}
