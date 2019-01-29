package rules

import (
	"testing"
)

func TestProvidedRequiredArguments(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, providedRequiredArguments)
}
