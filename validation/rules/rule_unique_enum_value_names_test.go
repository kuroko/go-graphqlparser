package rules

import (
	"testing"
)

func TestUniqueEnumValueNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueEnumValueNames)
}
