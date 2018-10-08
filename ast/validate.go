package ast

type ValidationContext struct {
	Schema Document
}

func (d *Document) validate(vctx *ValidationContext) *Errors {
	var errs *Errors

	d.Definitions.ForEach(func(definition Definition, _ int) {
		e := definition.validate(vctx)
		if e == nil {
			return
		}

		e.Join(errs)
		errs = e
	})

	return errs.Reverse()
}

func (d *Definition) validate(vctx *ValidationContext) *Errors {
	var errs *Errors

	switch d.Kind {
	case DefinitionKindExecutable:
		errs = d.ExecutableDefinition.validate(vctx)
	case DefinitionKindTypeSystem:
		errs = d.TypeSystemDefinition.validate(vctx)
	case DefinitionKindTypeSystemExtension:
		errs = d.TypeSystemExtension.validate(vctx)
	}

	return errs
}

func (d *ExecutableDefinition) validate(vctx *ValidationContext) *Errors {
	var errs *Errors

	switch d.Kind {
	case ExecutableDefinitionKindOperation:

		switch d.OperationType {
		case OperationTypeQuery:
		case OperationTypeMutation:
		case OperationTypeSubscription:
		}

	case ExecutableDefinitionKindFragment:

	}

	return errs
}

func (d *TypeSystemDefinition) validate(vctx *ValidationContext) *Errors {
	var errs *Errors
	return errs
}

func (d *TypeSystemExtension) validate(vctx *ValidationContext) *Errors {
	var errs *Errors
	return errs
}
