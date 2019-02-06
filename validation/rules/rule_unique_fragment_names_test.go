package rules

import (
	"testing"
)

func TestUniqueFragmentNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueFragmentNames)
}
