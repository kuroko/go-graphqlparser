package rules

import (
	"testing"
)

func TestUniqueEnumValueNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueEnumValueNames)
}
