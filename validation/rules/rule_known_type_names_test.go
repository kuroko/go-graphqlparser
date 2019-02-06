package rules

import (
	"testing"
)

func TestKnownTypeNames(t *testing.T) {
	t.Run("schema definition language", func(t *testing.T) {
		// TODO
		//tt := []ruleTestCase{}
		//
		//sdlRuleTester(t, tt, knownTypeNames)
	})

	t.Run("query language", func(t *testing.T) {
		var tt []ruleTestCase

		queryRuleTester(t, tt, knownTypeNames)
	})
}
