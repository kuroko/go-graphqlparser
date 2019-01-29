package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	for _, ruleName := range rules {
		snekName := camToSnek(ruleName)

		ruleBody := bytes.NewBuffer(nil)
		fmt.Fprintf(ruleBody, packageName)

		err := ruleBodyTmpl.Execute(ruleBody, map[string]string{
			"RuleNameLCF": lcfirst(ruleName),
			"RuleName":    ruleName,
		})
		if err != nil {
			log.Fatal(err)
		}

		outPath := fmt.Sprintf("%s/rule_%s%s.go", pathPrefix, snekName, "")
		err = ioutil.WriteFile(outPath, ruleBody.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}

		testBody := bytes.NewBuffer(nil)
		fmt.Fprintf(testBody, packageName)

		err = ruleTestTmpl.Execute(testBody, map[string]string{
			"RuleNameLCF": lcfirst(ruleName),
			"RuleName":    ruleName,
		})
		if err != nil {
			log.Fatal(err)
		}

		testOutPath := fmt.Sprintf("%s/rule_%s%s.go", pathPrefix, snekName, "_test")
		err = ioutil.WriteFile(testOutPath, testBody.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var pathPrefix = "/tmp"

var packageName = `
package rules
`

var ruleBodyTmpl = template.Must(template.New("ruleBodyTmpl").Parse(`
import (
	"github.com/bucketd/go-graphqlparser/validation"
)

// {{.RuleNameLCF}} ...
func {{.RuleNameLCF}}(w *validation.Walker) {}
`))

var ruleTestTmpl = template.Must(template.New("ruleTestTmpl").Parse(`
import (
	"testing"
)

func Test{{.RuleName}}(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, {{.RuleNameLCF}})
}
`))

func camToSnek(in string) string {
	snek := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(in, "${1}_${2}")
	snek = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snek, "${1}_${2}")
	return strings.ToLower(snek)
}

func lcfirst(in string) string {
	if len(in) == 0 {
		return in
	}

	if len(in) == 1 {
		return strings.ToLower(in)
	}

	fr, w := utf8.DecodeRuneInString(in)

	return strings.ToLower(string(fr)) + in[w:]
}

var rules = []string{
	"ExecutableDefinitions",
	"FieldsOnCorrectType",
	"FragmentsOnCompositeTypes",
	"KnownArgumentNames",
	"KnownDirectives",
	"KnownFragmentNames",
	"KnownTypeNames",
	"LoneAnonymousOperation",
	"LoneSchemaDefinition",
	"NoFragmentCycles",
	"NoUndefinedVariables",
	"NoUnusedFragments",
	"NoUnusedVariables",
	"OverlappingFieldsCanBeMerged",
	"PossibleFragmentSpreads",
	"PossibleTypeExtensions",
	"ProvidedRequiredArguments",
	"ScalarLeafs",
	"SingleFieldSubscriptions",
	"UniqueArgumentNames",
	"UniqueDirectiveNames",
	"UniqueDirectivesPerLocation",
	"UniqueEnumValueNames",
	"UniqueFieldDefinitionNames",
	"UniqueFragmentNames",
	"UniqueInputFieldNames",
	"UniqueOperationNames",
	"UniqueOperationTypes",
	"UniqueTypeNames",
	"UniqueVariableNames",
	"ValuesOfCorrectType",
	"VariablesAreInputTypes",
	"VariablesInAllowedPosition",
}
