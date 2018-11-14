package ast

import "github.com/bucketd/go-graphqlparser/graphql"

// 2.2 Document
// http://facebook.github.io/graphql/June2018/#sec-Language.Document

// @wg:ignore
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

type FragmentDefinition struct {
	Name          string
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  *Selections
}

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
	Name          string // but not "on"
	Alias         string
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
	StringValue  string
	ListValue    []Value
	ObjectValue  []ObjectField
	BooleanValue bool
	Kind         ValueKind
}

// @wg:ignore
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
	TypeSystemDefinitionKindSchema TypeSystemDefinitionKind = iota
	TypeSystemDefinitionKindType
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
	TypeSystemExtensionKindSchema TypeSystemExtensionKind = iota
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
	NamedType     Type // Only allow "TypeKindNamed" kind NamedType.
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

const (
	DirectiveLocationKindQuery DirectiveLocation = iota
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

type DirectiveLocation int8

func (l DirectiveLocation) String() string {
	switch l {
	case DirectiveLocationKindQuery:
		return "QUERY"
	case DirectiveLocationKindMutation:
		return "MUTATION"
	case DirectiveLocationKindSubscription:
		return "SUBSCRIPTION"
	case DirectiveLocationKindField:
		return "FIELD"
	case DirectiveLocationKindFragmentDefinition:
		return "FRAGMENT_DEFINITION"
	case DirectiveLocationKindFragmentSpread:
		return "FRAGMENT_SPREAD"
	case DirectiveLocationKindInlineFragment:
		return "INLINE_FRAGMENT"
	case DirectiveLocationKindSchema:
		return "SCHEMA"
	case DirectiveLocationKindScalar:
		return "SCALAR"
	case DirectiveLocationKindObject:
		return "OBJECT"
	case DirectiveLocationKindFieldDefinition:
		return "FIELD_DEFINITION"
	case DirectiveLocationKindArgumentDefinition:
		return "ARGUMENT_DEFINITION"
	case DirectiveLocationKindInterface:
		return "INTERFACE"
	case DirectiveLocationKindUnion:
		return "UNION"
	case DirectiveLocationKindEnum:
		return "ENUM"
	case DirectiveLocationKindEnumValue:
		return "ENUM_VALUE"
	case DirectiveLocationKindInputObject:
		return "INPUT_OBJECT"
	case DirectiveLocationKindInputFieldDefinition:
		return "INPUT_FIELD_DEFINITION"
	}

	return "invalid"
}

type DirectiveDefinition struct {
	Description         string
	Name                string
	ArgumentsDefinition *InputValueDefinitions
	DirectiveLocations  *DirectiveLocations
}
