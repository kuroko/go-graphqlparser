package rules

import (
	"testing"
)

func TestUniqueTypeNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueTypeNames)
}
