package ast

// 2.2 Document
// http://facebook.github.io/graphql/June2018/#sec-Language.Document

type Document struct {
	Definitions *Definitions
}

const (
	DefinitionKindExecutable DefinitionKind = iota
	DefinitionKindTypeSystem
)

type DefinitionKind int8

func (k DefinitionKind) String() string {
	switch k {
	case DefinitionKindExecutable:
		return "executable"
	case DefinitionKindTypeSystem:
		return "type system"
	}

	return "invalid"
}

type Definition struct {
	ExecutableDefinition *ExecutableDefinition
	TypeSystemDefinition *TypeSystemDefinition
	Kind                 DefinitionKind
}

// 2.3 Operations
// http://facebook.github.io/graphql/June2018/#sec-Language.Operations
// 2.8 Fragments
// http://facebook.github.io/graphql/June2018/#sec-Language.Fragments

// +Kind?
const (
	OperationTypeQuery OperationType = iota
	OperationTypeMutation
	OperationTypeSubscription
)

type OperationType int8

func (t OperationType) String() string {
	switch t {
	case OperationTypeQuery:
		return "query"
	case OperationTypeMutation:
		return "mutation"
	case OperationTypeSubscription:
		return "subscription"
	}

	return "invalid"
}

const (
	ExecutableDefinitionKindOperation ExecutableDefinitionKind = iota
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
	Name                string         // but not "on" if is FragmentDefinition kind.
	TypeCondition       *TypeCondition // not on operation definition.
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        *Selections
	OperationType       OperationType
	Kind                ExecutableDefinitionKind
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
	Kind          SelectionKind
}

// 2.5 Fields
// http://facebook.github.io/graphql/June2018/#sec-Language.Fields

//type Field struct {
//	Alias        string
//	Name         string
//	Arguments    []*Argument
//	Directives   *Directives
//	SelectionSet *Selections
//}

// 2.6 Arguments
// http://facebook.github.io/graphql/June2018/#sec-Language.Arguments

// Argument :
type Argument struct {
	Name  string
	Value Value
}

// 2.8 Fragments
// http://facebook.github.io/graphql/June2018/#sec-Language.Fragments

//type FragmentSpread struct {
//	Name       string // but not "on"
//	Directives *Directives
//}

type TypeCondition struct {
	NamedType Type // Only allow "TypeKindNamedType" kind NamedType.
}

//type InlineFragment struct {
//	TypeCondition TypeCondition
//	Directives    *Directives
//	SelectionSet  *Selections
//}

// 2.9 Input Values
// http://facebook.github.io/graphql/June2018/#sec-Input-Values

// Value :
const (
	ValueKindVariable ValueKind = iota
	ValueKindIntValue
	ValueKindFloatValue
	ValueKindStringValue
	ValueKindBooleanValue
	ValueKindNullValue
	ValueKindEnumValue
	ValueKindListValue
	ValueKindObjectValue
)

type ValueKind int8

func (k ValueKind) String() string {
	switch k {
	case ValueKindVariable:
		return "variable value"
	case ValueKindIntValue:
		return "int8 value"
	case ValueKindFloatValue:
		return "float value"
	case ValueKindStringValue:
		return "string value"
	case ValueKindBooleanValue:
		return "boolean value"
	case ValueKindNullValue:
		return "null value"
	case ValueKindEnumValue:
		return "enum value"
	case ValueKindListValue:
		return "list value"
	case ValueKindObjectValue:
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
	TypeKindNamedType TypeKind = iota
	TypeKindListType
)

type TypeKind int8

func (k TypeKind) String() string {
	switch k {
	case TypeKindNamedType:
		return "NamedType"
	case TypeKindListType:
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
	TypeSystemExtensionKindSchema = iota
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
	NamedType     Type // Only allow "TypeKindNamedType" kind NamedType.
	OperationType OperationType
}

// 3.2.2 Schema Extension
type SchemaExtension struct {
	Directives               *Directives
	OperationTypeDefinitions *OperationTypeDefinitions
}

type OperationTypeDefinition struct {
	NamedType     string
	OperationType OperationType
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
	ImplementsInterface   *Types // Only allow "TypeKindNamedType" kind NamedType.
	FieldsDefinition      *FieldDefinitions
	UnionMemberTypes      *Types // Only allow "TypeKindNamedType" kind NamedType.
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
	TypeExtensionKindScalar = iota
	TypeExtensionKindObject
	TypeExtensionKindInterface
	TypeExtensionKindUnion
	TypeExtensionKindEnum
	TypeExtensionKindInputObject
)

type TypeExtensionKind int8

type TypeExtension struct {
	Directives            *Directives
	ImplementsInterface   *Types // Only allow "TypeKindNamedType" kind NamedType.
	FieldsDefinition      *FieldDefinitions
	UnionMemberTypes      *Types // Only allow "TypeKindNamedType" kind NamedType.
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

// +Kind?
const (
	DirectiveLocationQuery DirectiveLocation = iota
	DirectiveLocationMutation
	DirectiveLocationSubscription
	DirectiveLocationField
	DirectiveLocationFragmentDefinition
	DirectiveLocationFragmentSpread
	DirectiveLocationInlineFragment
	DirectiveLocationSchema
	DirectiveLocationScalar
	DirectiveLocationObject
	DirectiveLocationFieldDefinition
	DirectiveLocationArgumentDefinition
	DirectiveLocationInterface
	DirectiveLocationUnion
	DirectiveLocationEnum
	DirectiveLocationEnumValue
	DirectiveLocationInputObject
	DirectiveLocationInputFieldDefinition
)

type DirectiveLocation int8

func (l DirectiveLocation) String() string {
	switch l {
	case DirectiveLocationQuery:
		return "QUERY"
	case DirectiveLocationMutation:
		return "MUTATION"
	case DirectiveLocationSubscription:
		return "SUBSCRIPTION"
	case DirectiveLocationField:
		return "FIELD"
	case DirectiveLocationFragmentDefinition:
		return "FRAGMENT_DEFINITION"
	case DirectiveLocationFragmentSpread:
		return "FRAGMENT_SPREAD"
	case DirectiveLocationInlineFragment:
		return "INLINE_FRAGMENT"
	case DirectiveLocationSchema:
		return "SCHEMA"
	case DirectiveLocationScalar:
		return "SCALAR"
	case DirectiveLocationObject:
		return "OBJECT"
	case DirectiveLocationFieldDefinition:
		return "FIELD_DEFINITION"
	case DirectiveLocationArgumentDefinition:
		return "ARGUMENT_DEFINITION"
	case DirectiveLocationInterface:
		return "INTERFACE"
	case DirectiveLocationUnion:
		return "UNION"
	case DirectiveLocationEnum:
		return "ENUM"
	case DirectiveLocationEnumValue:
		return "ENUM_VALUE"
	case DirectiveLocationInputObject:
		return "INPUT_OBJECT"
	case DirectiveLocationInputFieldDefinition:
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
