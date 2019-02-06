package rules

import (
	"testing"
)

func TestUniqueDirectiveNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueDirectiveNames)
}
