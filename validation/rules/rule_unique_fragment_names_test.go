package rules

import (
	"testing"
)

func TestUniqueFragmentNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueFragmentNames)
}
