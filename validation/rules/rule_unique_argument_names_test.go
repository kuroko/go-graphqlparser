package rules

import (
	"testing"
)

func TestUniqueArgumentNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueArgumentNames)
}
