package rules

import (
	"testing"
)

func TestUniqueTypeNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueTypeNames)
}
