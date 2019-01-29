package rules

import (
	"testing"
)

func TestVariablesInAllowedPosition(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, variablesInAllowedPosition)
}
