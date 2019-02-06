package rules

import (
	"testing"
)

func TestVariablesInAllowedPosition(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, variablesInAllowedPosition)
}
