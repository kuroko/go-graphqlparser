package rules

import (
	"testing"
)

func TestUniqueOperationNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueOperationNames)
}
