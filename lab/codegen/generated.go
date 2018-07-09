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
	if value.Kind != ast.ValueKindObjectValue {
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
// the exception of other resolvers:

// NullString ...
type NullString struct {
	String string
	IsNull bool
}

// For resolvers, the generated server will know which resolver are allowed to be nullable. If a
// resolver is marked as non-nullable, and nil is returned, then we know we can error. Just like
// resolvers, arrays will also be checked by the generated server.

// TODO(seeruk): Can you pass any type as input to something?
// - No.
