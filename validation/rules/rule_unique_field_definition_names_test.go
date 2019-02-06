package rules

import (
	"testing"
)

func TestUniqueFieldDefinitionNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueFieldDefinitionNames)
}
