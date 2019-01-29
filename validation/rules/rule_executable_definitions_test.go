package rules

import (
	"testing"
)

func TestExecutableDefinitions(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, executableDefinitions)
}
