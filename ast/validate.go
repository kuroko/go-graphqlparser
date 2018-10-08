package ast

import "fmt"

type ValidationContext struct {
	Schema Document
}

type Vdata struct {
	operationNames     map[string]bool
	anonOperationCount int
}

func (d *Document) validate(vctx *ValidationContext, vd *Vdata) *Errors {
	var errs *Errors

	d.Definitions.ForEach(func(definition Definition, _ int) {
		e := definition.validate(vctx, vd)
		if e == nil {
			return
		}

		e.Join(errs)
		errs = e
	})

	if vd.anonOperationCount > 1 {
		err := fmt.Errorf("seen %d shorthand queries", vd.anonOperationCount)
		errs = errs.Add(err)
	}

	return errs.Reverse()
}

func (d *Definition) validate(vctx *ValidationContext, vd *Vdata) *Errors {
	var errs *Errors

	switch d.Kind {
	case DefinitionKindExecutable:
		errs = d.ExecutableDefinition.validate(vctx, vd)

	case DefinitionKindTypeSystem:
		err := fmt.Errorf("definition is of type TypeSystemDefinition, must be executable")
		errs = errs.Add(err)

		errs = d.TypeSystemDefinition.validate(vctx, vd)

	case DefinitionKindTypeSystemExtension:
		err := fmt.Errorf("definition is of type TypeSystemExtension, must be executable")
		errs = errs.Add(err)

		errs = d.TypeSystemExtension.validate(vctx, vd)
	}

	return errs
}

func (d *ExecutableDefinition) validate(vctx *ValidationContext, vd *Vdata) *Errors {
	var errs *Errors

	if vd.operationNames[d.Name] {
		err := fmt.Errorf("operation definitions should be unique, seen %s more than once", d.Name)
		errs = errs.Add(err)
	}
	vd.operationNames[d.Name] = true

	switch d.Kind {
	case ExecutableDefinitionKindOperation:

		switch d.OperationType {
		case OperationTypeQuery:
			if d.ShorthandQuery {
				vd.anonOperationCount++
			}
		case OperationTypeMutation:
		case OperationTypeSubscription:
		}

	case ExecutableDefinitionKindFragment:

	}

	return errs
}

func (d *TypeSystemDefinition) validate(vctx *ValidationContext, vd *Vdata) *Errors {
	var errs *Errors

	return errs
}

func (d *TypeSystemExtension) validate(vctx *ValidationContext, vd *Vdata) *Errors {
	var errs *Errors
	return errs
}
