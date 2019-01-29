package rules

import (
	"testing"
)

func TestValuesOfCorrectType(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, valuesOfCorrectType)
}
