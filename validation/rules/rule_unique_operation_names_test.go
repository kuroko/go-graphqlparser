package rules

import (
	"testing"
)

func TestUniqueOperationNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueOperationNames)
}
