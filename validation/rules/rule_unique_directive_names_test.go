package rules

import (
	"testing"
)

func TestUniqueDirectiveNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueDirectiveNames)
}
