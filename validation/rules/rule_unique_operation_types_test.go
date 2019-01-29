package rules

import (
	"testing"
)

func TestUniqueOperationTypes(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueOperationTypes)
}
