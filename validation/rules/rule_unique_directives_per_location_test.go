package rules

import (
	"testing"
)

func TestUniqueDirectivesPerLocation(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueDirectivesPerLocation)
}
