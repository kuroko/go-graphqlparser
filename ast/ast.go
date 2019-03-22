package ast

import "github.com/bucketd/go-graphqlparser/graphql"

// 2.2 Document
// http://facebook.github.io/graphql/June2018/#sec-Language.Document

type Document struct {
	Definitions *Definitions
}

const (
	// @wg:field ExecutableDefinition
	DefinitionKindExecutable DefinitionKind = iota
	// @wg:field TypeSystemDefinition
	DefinitionKindTypeSystem
	// @wg:field TypeSystemExtension
	DefinitionKindTypeSystemExtension
)

type DefinitionKind int8

func (k DefinitionKind) String() string {
	switch k {
	case DefinitionKindExecutable:
		return "executable"
	case DefinitionKindTypeSystem:
		return "type system"
	case DefinitionKindTypeSystemExtension:
		return "type system extension"
	}

	return "invalid"
}

type Definition struct {
	Location             graphql.Location
	ExecutableDefinition *ExecutableDefinition
	TypeSystemDefinition *TypeSystemDefinition
	TypeSystemExtension  *TypeSystemExtension
	Kind                 DefinitionKind
}

// 2.3 Operations
// http://facebook.github.io/graphql/June2018/#sec-Language.Operations
// 2.8 Fragments
// http://facebook.github.io/graphql/June2018/#sec-Language.Fragments

const (
	// @wg:field OperationDefinition
	ExecutableDefinitionKindOperation ExecutableDefinitionKind = iota
	// @wg:field FragmentDefinition
	ExecutableDefinitionKindFragment
)

type ExecutableDefinitionKind int8

func (k ExecutableDefinitionKind) String() string {
	switch k {
	case ExecutableDefinitionKindOperation:
		return "operation"
	case ExecutableDefinitionKindFragment:
		return "fragment"
	}

	return "invalid"
}

type ExecutableDefinition struct {
	FragmentDefinition  *FragmentDefinition
	OperationDefinition *OperationDefinition
	Kind                ExecutableDefinitionKind
}

func (def ExecutableDefinition) String() string {
	switch def.Kind {
	case ExecutableDefinitionKindFragment:
		return def.FragmentDefinition.Name
	case ExecutableDefinitionKindOperation:
		return def.OperationDefinition.Name
	}

	return "UnnamedExecutableDefinition"
}

type FragmentDefinition struct {
	Name          string
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  *Selections
}

// @wg:field self
const (
	OperationDefinitionKindQuery OperationDefinitionKind = iota
	OperationDefinitionKindMutation
	OperationDefinitionKindSubscription
)

type OperationDefinitionKind int8

func (t OperationDefinitionKind) String() string {
	switch t {
	case OperationDefinitionKindQuery:
		return "query"
	case OperationDefinitionKindMutation:
		return "mutation"
	case OperationDefinitionKindSubscription:
		return "subscription"
	}

	return "invalid"
}

type OperationDefinition struct {
	Name                string
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        *Selections
	Kind                OperationDefinitionKind
}

// 2.4 Selection Sets
// http://facebook.github.io/graphql/June2018/#sec-Selection-Sets

// @wg:field self
const (
	SelectionKindField SelectionKind = iota
	SelectionKindFragmentSpread
	SelectionKindInlineFragment
)

type SelectionKind int8

func (k SelectionKind) String() string {
	switch k {
	case SelectionKindField:
		return "field"
	case SelectionKindFragmentSpread:
		return "fragment spread"
	case SelectionKindInlineFragment:
		return "inline fragment"
	}

	return "unknown"
}

type Selection struct {
	Name  string // but not "on"
	Alias string
	// @wg:on_kinds FragmentSpreadSelection
	TypeCondition *TypeCondition
	Arguments     *Arguments
	Directives    *Directives
	SelectionSet  *Selections
	Path          *graphql.PathNodes
	Kind          SelectionKind
}

// 2.6 Arguments
// http://facebook.github.io/graphql/June2018/#sec-Language.Arguments

// Argument :
type Argument struct {
	Name  string
	Value Value
}

// 2.8 Fragments
// http://facebook.github.io/graphql/June2018/#sec-Language.Fragments

type TypeCondition struct {
	NamedType Type // Only allow "TypeKindNamed" kind NamedType.
}

// 2.9 Input Values
// http://facebook.github.io/graphql/June2018/#sec-Input-Values

// Value :
// @wg:field self
const (
	ValueKindVariable ValueKind = iota
	ValueKindInt
	ValueKindFloat
	ValueKindString
	ValueKindBoolean
	ValueKindNull
	ValueKindEnum
	ValueKindList
	ValueKindObject
)

type ValueKind int8

func (k ValueKind) String() string {
	switch k {
	case ValueKindVariable:
		return "variable value"
	case ValueKindInt:
		return "int8 value"
	case ValueKindFloat:
		return "float value"
	case ValueKindString:
		return "string value"
	case ValueKindBoolean:
		return "boolean value"
	case ValueKindNull:
		return "null value"
	case ValueKindEnum:
		return "enum value"
	case ValueKindList:
		return "list value"
	case ValueKindObject:
		return "object value"
	}

	return "unknown value"
}

type Value struct {
	IntValue   int
	FloatValue float64
	// StringValue covers variables and enums, enums are names, but not `true`, `false`, or `null`.
	StringValue string
	// @wg:on_kinds ListValue
	ListValue []Value
	// @wg:on_kinds ObjectValue
	ObjectValue  []ObjectField
	BooleanValue bool
	Kind         ValueKind
}

type ObjectField struct {
	Name  string
	Value Value
}

// 2.10 Variables
// http://facebook.github.io/graphql/June2018/#sec-Language.Variables

type VariableDefinition struct {
	Name         string
	Type         Type
	DefaultValue *Value
}

// 2.11 Type References
// http://facebook.github.io/graphql/June2018/#sec-Type-References

// @wg:field self
const (
	TypeKindNamed TypeKind = iota
	TypeKindList
)

type TypeKind int8

func (k TypeKind) String() string {
	switch k {
	case TypeKindNamed:
		return "NamedType"
	case TypeKindList:
		return "ListType"
	}

	return "InvalidType"
}

type Type struct {
	NamedType   string
	ListType    *Type
	NonNullable bool
	Kind        TypeKind
}

// 2.12 Directives
// http://facebook.github.io/graphql/June2018/#sec-Language.Directives

// Directive :
type Directive struct {
	Name      string
	Arguments *Arguments
}

// 3.0 NamedType System
// http://facebook.github.io/graphql/June2018/#TypeSystemDefinition

const (
	// @wg:field SchemaDefinition
	TypeSystemDefinitionKindSchema TypeSystemDefinitionKind = iota
	// @wg:field TypeDefinition
	TypeSystemDefinitionKindType
	// @wg:field DirectiveDefinition
	TypeSystemDefinitionKindDirective
)

type TypeSystemDefinitionKind int8

func (k TypeSystemDefinitionKind) String() string {
	switch k {
	case TypeSystemDefinitionKindSchema:
		return "schema"
	case TypeSystemDefinitionKindType:
		return "type"
	case TypeSystemDefinitionKindDirective:
		return "directive"
	}

	return "invalid"
}

type TypeSystemDefinition struct {
	SchemaDefinition    *SchemaDefinition
	TypeDefinition      *TypeDefinition
	DirectiveDefinition *DirectiveDefinition
	Kind                TypeSystemDefinitionKind
}

// 3.1 Type System Extensions
const (
	// @wg:field SchemaExtension
	TypeSystemExtensionKindSchema TypeSystemExtensionKind = iota
	// @wg:field TypeExtension
	TypeSystemExtensionKindType
)

type TypeSystemExtensionKind uint8

type TypeSystemExtension struct {
	SchemaExtension *SchemaExtension
	TypeExtension   *TypeExtension
	Kind            TypeSystemExtensionKind
}

// 3.2 Schema
// http://facebook.github.io/graphql/June2018/#sec-Schema

type SchemaDefinition struct {
	Directives                   *Directives
	RootOperationTypeDefinitions *RootOperationTypeDefinitions
}

// TODO: Should NamedType be a string?
type RootOperationTypeDefinition struct {
	NamedType     Type
	OperationType OperationDefinitionKind
}

// 3.2.2 Schema Extension
type SchemaExtension struct {
	Directives               *Directives
	OperationTypeDefinitions *OperationTypeDefinitions
}

// TODO: Should NamedType be a string?
type OperationTypeDefinition struct {
	NamedType     Type
	OperationType OperationDefinitionKind
}

// 3.4 Types
// http://facebook.github.io/graphql/June2018/#sec-Types

// @wg:field self
const (
	TypeDefinitionKindScalar TypeDefinitionKind = iota
	TypeDefinitionKindObject
	TypeDefinitionKindInterface
	TypeDefinitionKindUnion
	TypeDefinitionKindEnum
	TypeDefinitionKindInputObject
)

type TypeDefinitionKind int8

type TypeDefinition struct {
	Description           string
	Directives            *Directives
	ImplementsInterface   *Types // Only allow "TypeKindNamed" kind NamedType.
	FieldsDefinition      *FieldDefinitions
	UnionMemberTypes      *Types // Only allow "TypeKindNamed" kind NamedType.
	EnumValuesDefinition  *EnumValueDefinitions
	InputFieldsDefinition *InputValueDefinitions
	Name                  string
	Kind                  TypeDefinitionKind
}

// GetEnumValueDefinition ...
func (d TypeDefinition) GetEnumValueDefinition(valueName string) (EnumValueDefinition, bool) {
	var evd EnumValueDefinition
	var found bool

	if IsEnumTypeDefinition(&d) {
		d.EnumValuesDefinition.ForEach(func(ievd EnumValueDefinition, i int) {
			if ievd.EnumValue == valueName {
				evd = ievd
				found = true
			}
		})
	}

	return evd, found
}

// IsScalarTypeDefinition ...
func IsScalarTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindScalar
}

// IsObjectTypeDefinition ...
func IsObjectTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindObject
}

// IsInterfaceTypeDefinition ...
func IsInterfaceTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindInterface
}

// IsUnionTypeDefinition ...
func IsUnionTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindUnion
}

// IsEnumTypeDefinition ...
func IsEnumTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindEnum
}

// IsInputObjectTypeDefinition ...
func IsInputObjectTypeDefinition(def *TypeDefinition) bool {
	return def.Kind == TypeDefinitionKindInputObject
}

type FieldDefinition struct {
	Description         string
	Name                string
	Directives          *Directives
	Type                Type
	ArgumentsDefinition *InputValueDefinitions
}

type EnumValueDefinition struct {
	Description string
	Directives  *Directives
	EnumValue   string // Name but not true or false or null.
}

// 3.4.3 Type Extensions
// @wg:field self
const (
	TypeExtensionKindScalar TypeExtensionKind = iota
	TypeExtensionKindObject
	TypeExtensionKindInterface
	TypeExtensionKindUnion
	TypeExtensionKindEnum
	TypeExtensionKindInputObject
)

type TypeExtensionKind int8

type TypeExtension struct {
	Directives            *Directives
	ImplementsInterface   *Types // Only allow "TypeKindNamed" kind NamedType.
	FieldsDefinition      *FieldDefinitions
	UnionMemberTypes      *Types // Only allow "TypeKindNamed" kind NamedType.
	EnumValuesDefinition  *EnumValueDefinitions
	InputFieldsDefinition *InputValueDefinitions
	Name                  string
	Kind                  TypeExtensionKind
}

// IsScalarTypeExtension ...
func IsScalarTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindScalar
}

// IsObjectTypeExtension ...
func IsObjectTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindObject
}

// IsInterfaceTypeExtension ...
func IsInterfaceTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindInterface
}

// IsUnionTypeExtension ...
func IsUnionTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindUnion
}

// IsEnumTypeExtension ...
func IsEnumTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindEnum
}

// IsInputObjectTypeExtension ...
func IsInputObjectTypeExtension(def *TypeExtension) bool {
	return def.Kind == TypeExtensionKindInputObject
}

// 3.6 Objects
// http://facebook.github.io/graphql/June2018/#sec-Objects

type InputValueDefinition struct {
	Description  string
	Name         string
	Type         Type
	Directives   *Directives
	DefaultValue *Value
}

// 3.13 Directives
// http://facebook.github.io/graphql/June2018/#sec-Type-System.Directives

var DirectiveLocationsByName = map[string]DirectiveLocations{
	"QUERY":                  DirectiveLocationKindQuery,
	"MUTATION":               DirectiveLocationKindMutation,
	"SUBSCRIPTION":           DirectiveLocationKindSubscription,
	"FIELD":                  DirectiveLocationKindField,
	"FRAGMENT_DEFINITION":    DirectiveLocationKindFragmentDefinition,
	"FRAGMENT_SPREAD":        DirectiveLocationKindFragmentSpread,
	"INLINE_FRAGMENT":        DirectiveLocationKindInlineFragment,
	"SCHEMA":                 DirectiveLocationKindSchema,
	"SCALAR":                 DirectiveLocationKindScalar,
	"OBJECT":                 DirectiveLocationKindObject,
	"FIELD_DEFINITION":       DirectiveLocationKindFieldDefinition,
	"ARGUMENT_DEFINITION":    DirectiveLocationKindArgumentDefinition,
	"INTERFACE":              DirectiveLocationKindInterface,
	"UNION":                  DirectiveLocationKindUnion,
	"ENUM":                   DirectiveLocationKindEnum,
	"ENUM_VALUE":             DirectiveLocationKindEnumValue,
	"INPUT_OBJECT":           DirectiveLocationKindInputObject,
	"INPUT_FIELD_DEFINITION": DirectiveLocationKindInputFieldDefinition,
}

var NamesByDirectiveLocations = map[DirectiveLocations]string{
	DirectiveLocationKindQuery:                "QUERY",
	DirectiveLocationKindMutation:             "MUTATION",
	DirectiveLocationKindSubscription:         "SUBSCRIPTION",
	DirectiveLocationKindField:                "FIELD",
	DirectiveLocationKindFragmentDefinition:   "FRAGMENT_DEFINITION",
	DirectiveLocationKindFragmentSpread:       "FRAGMENT_SPREAD",
	DirectiveLocationKindInlineFragment:       "INLINE_FRAGMENT",
	DirectiveLocationKindSchema:               "SCHEMA",
	DirectiveLocationKindScalar:               "SCALAR",
	DirectiveLocationKindObject:               "OBJECT",
	DirectiveLocationKindFieldDefinition:      "FIELD_DEFINITION",
	DirectiveLocationKindArgumentDefinition:   "ARGUMENT_DEFINITION",
	DirectiveLocationKindInterface:            "INTERFACE",
	DirectiveLocationKindUnion:                "UNION",
	DirectiveLocationKindEnum:                 "ENUM",
	DirectiveLocationKindEnumValue:            "ENUM_VALUE",
	DirectiveLocationKindInputObject:          "INPUT_OBJECT",
	DirectiveLocationKindInputFieldDefinition: "INPUT_FIELD_DEFINITION",
}

const (
	DirectiveLocationKindQuery DirectiveLocations = 1 << iota
	DirectiveLocationKindMutation
	DirectiveLocationKindSubscription
	DirectiveLocationKindField
	DirectiveLocationKindFragmentDefinition
	DirectiveLocationKindFragmentSpread
	DirectiveLocationKindInlineFragment
	DirectiveLocationKindSchema
	DirectiveLocationKindScalar
	DirectiveLocationKindObject
	DirectiveLocationKindFieldDefinition
	DirectiveLocationKindArgumentDefinition
	DirectiveLocationKindInterface
	DirectiveLocationKindUnion
	DirectiveLocationKindEnum
	DirectiveLocationKindEnumValue
	DirectiveLocationKindInputObject
	DirectiveLocationKindInputFieldDefinition
)

type DirectiveLocations int32

func (l DirectiveLocations) String() string {
	return NamesByDirectiveLocations[l]
}

type DirectiveDefinition struct {
	Description         string
	Name                string
	ArgumentsDefinition *InputValueDefinitions
	DirectiveLocations  DirectiveLocations
}
