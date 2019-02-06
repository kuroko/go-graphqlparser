package rules

import (
	"testing"
)

func TestUniqueInputFieldNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueInputFieldNames)
}
