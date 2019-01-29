package rules

import (
	"testing"
)

func TestNoUndefinedVariables(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, noUndefinedVariables)
}
