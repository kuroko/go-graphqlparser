package rules

import (
	"testing"
)

func TestUniqueFieldDefinitionNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueFieldDefinitionNames)
}
