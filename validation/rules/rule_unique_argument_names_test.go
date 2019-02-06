package rules

import (
	"testing"
)

func TestUniqueArgumentNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueArgumentNames)
}
