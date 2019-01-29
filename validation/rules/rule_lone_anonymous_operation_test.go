package rules

import (
	"testing"
)

func TestLoneAnonymousOperation(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, loneAnonymousOperation)
}
