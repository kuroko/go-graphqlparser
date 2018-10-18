package ast

func (w *Walker) Walk(doc Document) {
	w.OnDocumentEnter(doc)
	w.walkDefinitions(doc.Definitions)
	w.OnDocumentLeave(doc)
}

func (w *Walker) walkDefinitions(defs *Definitions) {
	w.OnDefinitionsEnter(defs)

	defs.ForEach(func(def Definition, i int) {
		w.walkDefinition(def)
	})

	w.OnDefinitionsLeave(defs)
}

func (w *Walker) walkDefinition(definition Definition) {
	w.OnDefinitionEnter(definition)

	// TODO...
	//w.walkDefinitions(doc.Definitions)

	w.OnDefinitionLeave(definition)
}
