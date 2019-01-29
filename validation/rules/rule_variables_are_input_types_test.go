package rules

import (
	"testing"
)

func TestVariablesAreInputTypes(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, variablesAreInputTypes)
}
