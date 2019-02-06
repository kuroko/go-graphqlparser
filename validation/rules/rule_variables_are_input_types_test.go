package rules

import (
	"testing"
)

func TestVariablesAreInputTypes(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, variablesAreInputTypes)
}
