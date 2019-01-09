package codegen

import (
	"errors"

	"github.com/bucketd/go-graphqlparser/ast"
)

// Enum values for: enum Gender
const (
	GenderMALE Gender = iota
	GenderFEMALE
)

// Gender represents enum values for: enum Gender
type Gender int

// TODO: Gender would also need some kind of Marshal method, to take an ast.Value or similar
// and create a Gender out of it.

// Union kinds for: union SearchResult
const (
	SearchResultKindDroid SearchResultKind = iota
	SearchResultKindHuman
)

// SearchResultKind represents union kinds for: union SearchResult
type SearchResultKind int

// SearchResult represents: union Searchresult
type SearchResult struct {
	Kind  SearchResultKind
	Droid DroidResolver
	Human HumanResolver
}

// QueryResolver represents: type Query
type QueryResolver interface {
	Character(id string) CharacterResolver
	Search(input string) []SearchResult
}

// MutationResolver represents: type Mutation
type MutationResolver interface {
	CreateDroid(droid DroidInput) DroidResolver
}

// CharacterResolver represents: interface Character
type CharacterResolver interface {
	ID() string
	Name() string
	JobTitle() NullString
	BestFriend(gender Gender) CharacterResolver
}

// DroidInput represents: input DroidInput
type DroidInput struct {
	Name             string
	JobTitle         NullString
	MaleBestFriend   NullString
	FemaleBestFriend NullString
	PrimaryFunction  string
}

// MarshalGraphQL ...
func (di *DroidInput) MarshalGraphQL(value ast.Value) error {
	if value.Kind != ast.ValueKindObject {
		return errors.New("TODO")
	}

	for _, f := range value.ObjectValue {
		switch f.Name {
		case "name":
			di.Name = f.Value.StringValue
		}
	}

	return nil
}

// DroidResolver represents: type Droid implements Character
type DroidResolver interface {
	CharacterResolver
	PrimaryFunction() string
}

// HumanResolver represents: type Human implements Character
type HumanResolver interface {
	CharacterResolver
	Credits() int
}

// Somewhere else we define a NullXyz type for every type that can be returned from a resolver. With
// the exception of other resolvers, and slices:

// NullString ...
type NullString struct {
	String string
	IsNull bool
}

// NOTE(seeruk): Nullability:
// * For resolvers, the generated server will know which resolvers are allowed to be nullable. If a
//   resolver is marked as non-nullable, and nil is returned, then we know we can error.
// * Just like resolvers, arrays will also be checked by the generated server. Luckily, Go slices
//   can be nil, so handling this case is easy.
// * Nullable scalars can be represented by custom types, just like sql.NullString and co.
//
// Generated code should only be able to return values that we can return as JSON. We control what
// gets generated in the end anyway.

// NOTE(seeruk): You cannot pass a GraphQL `type` as an argument, it must be an `input` type, and // those have restrictions on what fields can be defined as.

// NOTE(seeruk): Extending code generation.
// * Uses directives.
// * Could use directives on fields for context, errors, and even things like specifying bit size,
//   or whether something is signed or not, for number types. Plenty to be able to do with this, but
//   we'll start off very small and simple.
