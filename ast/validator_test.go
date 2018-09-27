package ast

import (
	"reflect"
	"testing"
)

func TestApplyExecutableRules(t *testing.T) {
	v := NewValidator(nil, nil)
	errs := v.ApplyExecutableRules()
	errs.ForEach(func(err Error, _ int) {
		t.Error(err)
	})
}

func TestFoo_executableDefinitions(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.executableDefinitions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.executableDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_fieldsOnCorrectType(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.fieldsOnCorrectType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.fieldsOnCorrectType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_fragmentsOnCompositeTypes(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.fragmentsOnCompositeTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.fragmentsOnCompositeTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_knownArgumentNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.knownArgumentNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.knownArgumentNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_knownDirectives(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.knownDirectives(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.knownDirectives() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_knownFragmentNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.knownFragmentNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.knownFragmentNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_knownTypeNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.knownTypeNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.knownTypeNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_loneAnonymousOperation(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.loneAnonymousOperation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.loneAnonymousOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_loneSchemaDefinition(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.loneSchemaDefinition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.loneSchemaDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_noFragmentCycles(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.noFragmentCycles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.noFragmentCycles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_noUndefinedVariables(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.noUndefinedVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.noUndefinedVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_noUnusedFragments(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.noUnusedFragments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.noUnusedFragments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_noUnusedVariables(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.noUnusedVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.noUnusedVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_overlappingFieldsCanBeMerged(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.overlappingFieldsCanBeMerged(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.overlappingFieldsCanBeMerged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_possibleFragmentSpreads(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.possibleFragmentSpreads(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.possibleFragmentSpreads() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_providedRequiredArguments(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.providedRequiredArguments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.providedRequiredArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_scalarLeafs(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.scalarLeafs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.scalarLeafs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_singleFieldSubscriptions(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.singleFieldSubscriptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.singleFieldSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueArgumentNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueArgumentNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueArgumentNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueDirectivesPerLocation(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueDirectivesPerLocation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueDirectivesPerLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueFragmentNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueFragmentNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueFragmentNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueInputFieldNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueInputFieldNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueInputFieldNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueOperationNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueOperationNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueOperationNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_uniqueVariableNames(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.uniqueVariableNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.uniqueVariableNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_valuesOfCorrectType(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.valuesOfCorrectType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.valuesOfCorrectType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_variablesAreInputTypes(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.variablesAreInputTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.variablesAreInputTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFoo_variablesInAllowedPosition(t *testing.T) {
	tests := []struct {
		name      string
		validator *Validator
		want      *Errors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.validator.variablesInAllowedPosition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.variablesInAllowedPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
