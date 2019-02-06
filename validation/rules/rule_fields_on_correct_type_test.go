package rules

import (
	"testing"
)

func TestFieldsOnCorrectType(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, fieldsOnCorrectType)
}
