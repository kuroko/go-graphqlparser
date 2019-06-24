package validation

import (
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

const (
	scalar = ast.TypeDefinitionKindScalar
	object = ast.TypeDefinitionKindObject
	iface  = ast.TypeDefinitionKindInterface
	union  = ast.TypeDefinitionKindUnion
	enum   = ast.TypeDefinitionKindEnum
	inobj  = ast.TypeDefinitionKindInputObject
)

// IsAbstractType ...
func IsAbstractType(schema *graphql.Schema, t ast.Type) bool {
	return IsInterfaceType(schema, t) || IsUnionType(schema, t)
}

// IsInputType ...
func IsInputType(schema *graphql.Schema, t ast.Type) bool {
	if t.Kind == ast.TypeKindList {
		return IsInputType(schema, *t.ListType)
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		switch def.Kind {
		case scalar, enum, inobj:
			return true
		}
	}

	return false
}

// IsInterfaceType
func IsInterfaceType(schema *graphql.Schema, t ast.Type) bool {
	if def, ok := schema.Types[t.NamedType]; ok {
		return def.Kind == ast.TypeDefinitionKindInterface
	}

	return false
}

// IsObjectType ...
func IsObjectType(schema *graphql.Schema, t ast.Type) bool {
	if def, ok := schema.Types[t.NamedType]; ok {
		return def.Kind == ast.TypeDefinitionKindObject
	}

	return false
}

// IsUnionType ...
func IsUnionType(schema *graphql.Schema, t ast.Type) bool {
	if def, ok := schema.Types[t.NamedType]; ok {
		return def.Kind == ast.TypeDefinitionKindUnion
	}

	return false
}

// IsOutputType ...
func IsOutputType(schema *graphql.Schema, t ast.Type) bool {
	if t.Kind == ast.TypeKindList {
		return IsOutputType(schema, *t.ListType)
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		switch def.Kind {
		case scalar, object, iface, union, enum:
			return true
		}
	}

	return false
}

// IsTypeSubTypeOf ...
func IsTypeSubTypeOf(schema *graphql.Schema, maybeSubType, superType ast.Type) bool {
	if maybeSubType == superType {
		return true
	}

	if superType.NonNullable {
		if maybeSubType.NonNullable {
			maybeSubType.NonNullable = false
			superType.NonNullable = false
			return IsTypeSubTypeOf(schema, maybeSubType, superType)
		}

		return false
	}

	if maybeSubType.NonNullable {
		maybeSubType.NonNullable = false
		return IsTypeSubTypeOf(schema, maybeSubType, superType)
	}

	if superType.Kind == ast.TypeKindList {
		if maybeSubType.Kind == ast.TypeKindList {
			return IsTypeSubTypeOf(schema, *maybeSubType.ListType, *superType.ListType)
		}

		return false
	}

	if maybeSubType.Kind == ast.TypeKindList {
		return false
	}

	isAbstractSuperType := IsAbstractType(schema, superType)
	isObjectSubType := IsObjectType(schema, maybeSubType)

	if isAbstractSuperType && isObjectSubType && IsPossibleType(schema, superType, maybeSubType) {
		return true
	}

	return false
}

// IsPossibleType ...
func IsPossibleType(schema *graphql.Schema, abstractType, possibleType ast.Type) bool {
	var found bool

	possibleTypes := PossibleTypes(schema, abstractType)
	gen := possibleTypes.Generator()

	for t, i := gen.Next(); i < possibleTypes.Len(); t, i = gen.Next() {
		if t == possibleType {
			found = true
			break
		}
	}

	return found
}

// PossibleTypes ...
func PossibleTypes(schema *graphql.Schema, abstractType ast.Type) *ast.Types {
	if IsUnionType(schema, abstractType) {
		unionType := schema.Types[abstractType.NamedType]
		return unionType.UnionMemberTypes
	}

	// TODO: This is crazy inefficient. We should do this once per schema.
	implementations := InterfaceImplementations(schema)

	return implementations[abstractType.NamedType]
}

// InterfaceImplementations ...
func InterfaceImplementations(schema *graphql.Schema) map[string]*ast.Types {
	// TODO: Maybe do this when we load the data.
	var interfaces int
	for _, typeDef := range schema.Types {
		if ast.IsInterfaceTypeDefinition(typeDef) {
			interfaces++
		}
	}

	implementations := make(map[string]*ast.Types, interfaces)

	for typeName, typeDef := range schema.Types {
		if ast.IsObjectTypeDefinition(typeDef) {
			typeDef.ImplementsInterface.ForEach(func(iface ast.Type, i int) {
				if IsInterfaceType(schema, iface) {
					if _, ok := implementations[iface.NamedType]; ok {
						implementations[iface.NamedType] = implementations[iface.NamedType].
							Add(ast.Type{NamedType: typeName})
					} else {
						implementations[iface.NamedType] = (*ast.Types)(nil).
							Add(ast.Type{NamedType: typeName})
					}
				}
			})
		} else if ast.IsInterfaceTypeDefinition(typeDef) {
			if _, ok := implementations[typeName]; !ok {
				implementations[typeName] = (*ast.Types)(nil)
			}
		}
	}

	return implementations
}

// IsUnionMemberType ...
func IsUnionMemberType(schema *graphql.Schema, t ast.Type) bool {
	if t.NonNullable {
		return false
	}

	if t.Kind == ast.TypeKindList {
		return false
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		if def.Kind == object {
			return true
		}
	}

	return false
}

// OrList takes a string slice, and returns each item comma separated, up to the last item which
// will be separated with ' or '.
func OrList(items []string) string {
	itemCount := len(items)
	if itemCount == 0 {
		panic("validation: expected at least one item")
	}

	if itemCount == 1 {
		return items[0]
	}

	if itemCount == 2 {
		return items[0] + " or " + items[1]
	}

	if itemCount > 5 {
		items = items[:5]
		itemCount = 5
	}

	items, last := items[:itemCount-1], items[itemCount-1]

	return strings.Join(items, ", ") + ", or " + last
}

// QuotedOrList is the same as OrList, but each item will also be wrapped in double-quotes.
func QuotedOrList(items []string) string {
	for i := 0; i < len(items); i++ {
		items[i] = `"` + items[i] + `"`
	}

	return OrList(items)
}

// SuggestionList accepts an invalid input string and a list of valid options, and returns a
// filtered list of the valid options sorted based on their similarity with the input.
//
// Source: https://github.com/graphql/graphql-js/blob/v14.2.0/src/jsutils/suggestionList.js
func SuggestionList(input string, options []string) []string {
	//ol := len(options)
	it := len(input) / 2

	optionsByDistance := make(map[string]int)

	for _, option := range options {
		ot := len(option) / 2

		threshold := it
		if ot > it {
			threshold = ot
		}

		if threshold < 1 {
			threshold = 1
		}

		distance := lexicalDistance(input, option)

		if distance <= threshold {
			optionsByDistance[option] = distance
		}
	}

	results := make([]string, 0, len(optionsByDistance))
	for _, o := range options {
		if _, ok := optionsByDistance[o]; ok {
			results = append(results, o)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		// Sort by length if the distance is the same, so we get the most specific suggestion first.
		if optionsByDistance[results[i]] == optionsByDistance[results[j]] {
			return len(results[i]) > len(results[j])
		}

		return optionsByDistance[results[i]] < optionsByDistance[results[j]]
	})

	return results
}

// lexicalDistance computes the lexical distance between strings A and B.
//
// The "distance" between two strings is given by counting the minimum number of edits needed to
// transform string A into string B. An edit can be an insertion, deletion, or substitution of a
// single character, or a swap of two adjacent characters.
//
// This distance can be useful for detecting typos in input or sorting
//
// Source: https://github.com/graphql/graphql-js/blob/v14.2.0/src/jsutils/suggestionList.js
func lexicalDistance(a, b string) int {
	if a == b {
		return 0
	}

	a = strings.ToLower(a)
	b = strings.ToLower(b)

	// Any case change counts as a single edit
	if a == b {
		return 1
	}

	return levenshtein.ComputeDistance(a, b)
}
