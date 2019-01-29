package rules

import (
	"testing"
)

func TestUniqueInputFieldNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueInputFieldNames)
}
