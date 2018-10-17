package ast

func (w *Walker) Walk(doc Document) {
	w.OnDocumentEvent(WalkerEnter, doc)
	w.walkDefinitions(doc.Definitions)
	w.OnDocumentEvent(WalkerLeave, doc)
}

func (w *Walker) walkDefinitions(defs *Definitions) {
	w.OnDefinitionsEvent(WalkerEnter, defs)

	defs.ForEach(func(def Definition, i int) {
		w.walkDefinition(def)
	})

	w.OnDefinitionsEvent(WalkerLeave, defs)
}

func (w *Walker) walkDefinition(definition Definition) {
	w.OnDefinitionEvent(WalkerEnter, definition)

	// TODO...
	//w.walkDefinitions(doc.Definitions)

	w.OnDefinitionEvent(WalkerLeave, definition)
}
