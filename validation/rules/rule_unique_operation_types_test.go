package rules

import (
	"testing"
)

func TestUniqueOperationTypes(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueOperationTypes)
}
