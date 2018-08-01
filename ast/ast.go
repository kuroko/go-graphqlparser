package ast

// 2.2 Document
// http://facebook.github.io/graphql/June2018/#sec-Language.Document

type Document struct {
	Definitions []*Definition
}

const (
	DefinitionKindExecutable DefinitionKind = iota
	DefinitionKindTypeSystem
)

type DefinitionKind int

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
	Kind                 DefinitionKind
	ExecutableDefinition *ExecutableDefinition
	TypeSystemDefinition *TypeSystemDefinition
}

// 2.3 Operations
// http://facebook.github.io/graphql/June2018/#sec-Language.Operations
// 2.8 Fragments
// http://facebook.github.io/graphql/June2018/#sec-Language.Fragments

const (
	OperationTypeQuery OperationType = iota
	OperationTypeMutation
	OperationTypeSubscription
)

type OperationType int

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

type ExecutableDefinitionKind int

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
	Kind                ExecutableDefinitionKind
	OperationType       OperationType
	Name                string // but not "on" if is FragmentDefinition kind.
	TypeCondition       *TypeCondition
	VariableDefinitions []*VariableDefinition
	Directives          []*Directive
	SelectionSet        []*Selection
}

// 2.4 Selection Sets
// http://facebook.github.io/graphql/June2018/#sec-Selection-Sets

const (
	SelectionKindField SelectionKind = iota
	SelectionKindFragmentSpread
	SelectionKindInlineFragment
)

type SelectionKind int

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
	Kind          SelectionKind
	Name          string // but not "on"
	Alias         string
	TypeCondition *TypeCondition
	Arguments     []*Argument
	Directives    []*Directive
	SelectionSet  []*Selection
}

// 2.5 Fields
// http://facebook.github.io/graphql/June2018/#sec-Language.Fields

type Field struct {
	Alias        string
	Name         string
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet []*Selection
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

//type FragmentSpread struct {
//	Name       string // but not "on"
//	Directives []*Directive
//}

type TypeCondition struct {
	NamedType Type // Only allow "TypeKindNamedType" kind NamedType.
}

//type InlineFragment struct {
//	TypeCondition TypeCondition
//	Directives    []*Directive
//	SelectionSet  []*Selection
//}

// 2.9 Input Values
// http://facebook.github.io/graphql/June2018/#sec-Input-Values

// Value [Const] :
const (
	ValueKindVariable ValueKind = iota // [~Const]
	ValueKindIntValue
	ValueKindFloatValue
	ValueKindStringValue
	ValueKindBooleanValue
	ValueKindNullValue
	ValueKindEnumValue
	ValueKindListValue   // [?Const]
	ValueKindObjectValue // [?Const]
)

type ValueKind int

func (k ValueKind) String() string {
	switch k {
	case ValueKindVariable:
		return "variable value"
	case ValueKindIntValue:
		return "int value"
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
	Kind          ValueKind
	VariableValue string
	IntValue      int
	FloatValue    float64
	StringValue   string
	EnumValue     string // Name, but not `true`, `false`, or `null`.
	ListValue     []Value
	ObjectValue   []ObjectField
	BooleanValue  bool
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

type TypeKind int

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
	Kind        TypeKind
	NamedType   string
	ListType    *Type
	NonNullable bool
}

// 2.12 Directives
// http://facebook.github.io/graphql/June2018/#sec-Language.Directives

// Directive :
type Directive struct {
	Name      string
	Arguments []*Argument
}

// 3.0 NamedType System
// http://facebook.github.io/graphql/June2018/#TypeSystemDefinition

const (
	TypeSystemDefinitionKindSchema TypeSystemDefinitionKind = iota
	TypeSystemDefinitionKindType
	TypeSystemDefinitionKindDirective
)

type TypeSystemDefinitionKind int

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
	Kind                TypeSystemDefinitionKind
	SchemaDefinition    SchemaDefinition
	TypeDefinition      TypeDefinition
	DirectiveDefinition DirectiveDefinition
}

// 3.2 Schema
// http://facebook.github.io/graphql/June2018/#sec-Schema

type SchemaDefinition struct {
	Directives                   []*Directive
	RootOperationTypeDefinitions []RootOperationTypeDefinition
}

type RootOperationTypeDefinition struct {
	OperationType OperationType
	NamedType     Type // Only allow "TypeKindNamedType" kind NamedType.
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

type TypeDefinitionKind int

type TypeDefinition struct {
	Kind                  TypeDefinitionKind
	Description           string
	Name                  string
	Directives            []*Directive
	ImplementsInterface   []Type // Only allow "TypeKindNamedType" kind NamedType.
	FieldsDefinition      []FieldDefinition
	UnionMemberTypes      []Type // Only allow "TypeKindNamedType" kind NamedType.
	EnumValuesDefinition  []EnumValueDefinition
	InputFieldsDefinition []FieldDefinition
}

type FieldDefinition struct {
	Description         string
	Name                string
	Directives          []*Directive
	Type                Type
	ArgumentsDefinition []InputValueDefinition
}

type EnumValueDefinition struct {
	Description string
	Directives  []*Directive
	EnumValue   string // Name but not true or false or null.
}

// 3.6 Objects
// http://facebook.github.io/graphql/June2018/#sec-Objects

type InputValueDefinition struct {
	Description  string
	Name         string
	Type         Type
	Directives   []*Directive
	DefaultValue *Value
}

// 3.13 Directives
// http://facebook.github.io/graphql/June2018/#sec-Type-System.Directives

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

type DirectiveLocation int

type DirectiveDefinition struct {
	Description         string
	Name                string
	ArgumentsDefinition []InputValueDefinition
	DirectiveLocations  []*DirectiveLocation
}
