package validator

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden record files?")

var schema = strings.TrimSpace(`
type Query {
  dog: Dog
}

enum DogCommand { 
  SIT
  DOWN
  HEEL
}

type Dog implements Pet {
  name: String!
  nickname: String
  barkVolume: Int
  doesKnowCommand(dogCommand: DogCommand!): Boolean!
  isHousetrained(atOtherHomes: Boolean): Boolean!
  owner: Human
}

interface Sentient {
  name: String!
}

interface Pet {
  name: String!
}

type Alien implements Sentient {
  name: String!
  homePlanet: String
}

type Human implements Sentient {
  name: String!
}

enum CatCommand {
  JUMP
}

type Cat implements Pet {
  name: String!
  nickname: String
  doesKnowCommand(catCommand: CatCommand!): Boolean!
  meowVolume: Int
}

union CatOrDog = Cat | Dog
union DogOrHuman = Dog | Human
union HumanOrAlien = Human | Alien
`)

func TestValidator_applyRulesGolden(t *testing.T) {
	tests := []struct {
		index    string
		rule     rule
		document string
		schema   string
	}{
		// {"00-00", (*Validator).executableDefinitions, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		// {"00-01", (*Validator).executableDefinitions, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		{"00-00", (*Validator).executableDefinitions, "{ foo }", schema},
		{"00-01", (*Validator).executableDefinitions, "extend type foo", schema},

		//{"01-00", (*Validator).fieldsOnCorrectType, nil},
		//{"01-01", (*Validator).fieldsOnCorrectType, nil},

		//{"02-00", (*Validator).fragmentsOnCompositeTypes, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"02-01", (*Validator).fragmentsOnCompositeTypes, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		//{"03-00", (*Validator).knownArgumentNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"03-01", (*Validator).knownArgumentNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		//{"04-00", (*Validator).knownArgumentNamesOnDirectives, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"04-01", (*Validator).knownArgumentNamesOnDirectives, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		//{"05-00", (*Validator).knownDirectives, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"05-01", (*Validator).knownDirectives, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		//{"06-00", (*Validator).knownFragmentNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"06-01", (*Validator).knownFragmentNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		//{"07-00", (*Validator).knownTypeNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindExecutable}, nil}}, &Document{}}},
		//{"07-01", (*Validator).knownTypeNames, &Validator{&Document{&Definitions{Definition{Kind: DefinitionKindTypeSystem}, nil}}, &Document{}}},

		// ...
		// (*Validator).loneAnonymousOperation
		// (*Validator).loneSchemaDefinition
		// (*Validator).noFragmentCycles
		// (*Validator).noUndefinedVariables
		// (*Validator).noUnusedFragments
		// (*Validator).noUnusedVariables
		// (*Validator).overlappingFieldsCanBeMerged
		// (*Validator).possibleFragmentSpreads
		// (*Validator).providedRequiredArguments
		// (*Validator).providedRequiredArgumentsOnDirectives
		// (*Validator).scalarLeafs
		// (*Validator).singleFieldSubscriptions
		// (*Validator).uniqueArgumentNames
		// (*Validator).uniqueDirectivesPerLocation
		// (*Validator).uniqueFragmentNames
		// (*Validator).uniqueInputFieldNames
		// (*Validator).uniqueOperationNames
		// (*Validator).uniqueVariableNames
		// (*Validator).valuesOfCorrectType
		// (*Validator).variablesAreInputTypes
		// (*Validator).variablesInAllowedPosition

	}
	psr := parser.New([]byte(schema))

	schemaAST, err := psr.Parse()
	if err != nil {
		t.Fatal(err)
	}

	type record struct {
		InputAST string
		Schema   string
		Errors   *ast.Errors
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s", tc.index), func(t *testing.T) {

			psr := parser.New([]byte(tc.document))

			doc, err := psr.Parse()
			if err != nil {
				t.Fatal(err)
			}

			validator := &Validator{
				DocumentAST:   &doc,
				GraphQLSchema: &schemaAST,
			}

			errs := validator.applyRules([]rule{tc.rule})

			actual := record{
				InputAST: ast.Sdump(*validator.DocumentAST),
				Schema:   ast.Sdump(*validator.GraphQLSchema),
				Errors:   errs,
			}

			goldenFileName := fmt.Sprintf("testdata/applyRulesGolden.%s.json", tc.index)

			if *update {
				bs, err := json.MarshalIndent(actual, "", "  ")
				if err != nil {
					t.Error(err)
				}

				err = ioutil.WriteFile(goldenFileName, bs, 0666)
				if err != nil {
					t.Error(err)
				}

				return
			}

			goldenBs, err := ioutil.ReadFile(goldenFileName)
			if err != nil {
				t.Error(err)
			}

			expected := record{}

			err = json.Unmarshal(goldenBs, &expected)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, expected, actual)
		})
	}
}
