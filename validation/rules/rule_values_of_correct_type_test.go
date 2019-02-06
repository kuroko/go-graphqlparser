package rules

import (
	"testing"
)

func TestValuesOfCorrectType(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, valuesOfCorrectType)
}
