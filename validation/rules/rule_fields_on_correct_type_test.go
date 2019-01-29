package rules

import (
	"testing"
)

func TestFieldsOnCorrectType(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, fieldsOnCorrectType)
}
