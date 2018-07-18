package ast

// 2.2 Document
// http://facebook.github.io/graphql/June2018/#sec-Language.Document

type Document struct {
	Definitions []Definition
}

type Definition struct {
	ExecutableDefinition ExecutableDefinition
	TypeSystemDefinition TypeSystemDefinition
	TypeSystemExtension  TypeSystemExtension
}

type TypeSystemDefinition struct{}
type TypeSystemExtension struct{}

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
	DefinitionKindOperation DefinitionKind = iota
	DefinitionKindFragment
)

type DefinitionKind int

func (k DefinitionKind) String() string {
	switch k {
	case DefinitionKindOperation:
		return "operation"
	case DefinitionKindFragment:
		return "fragment"
	}

	return "invalid"
}

type ExecutableDefinition struct {
	Kind                DefinitionKind
	OperationType       OperationType
	Name                string // but not "on" if is FragmentDefinition kind.
	TypeCondition       TypeCondition
	VariableDefinitions []VariableDefinition
	Directives          []Directive
	SelectionSet        []Selection
}

// 2.4 Selection Sets
// http://facebook.github.io/graphql/June2018/#sec-Selection-Sets

type Selection struct {
	Field          Field
	FragmentSpread FragmentSpread
	InlineFragment InlineFragment
}

// 2.5 Fields
// http://facebook.github.io/graphql/June2018/#sec-Language.Fields

type Field struct {
	Alias        string
	Name         string
	Arguments    []Argument
	Directives   []Directive
	SelectionSet []Selection
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

type FragmentSpread struct {
	Name       string // but not "on"
	Directives []Directive
}

type TypeCondition struct {
	Type Type // Only allow "TypeKindNamedType" kind Type.
}

type InlineFragment struct {
	TypeCondition TypeCondition
	Directives    []Directive
	SelectionSet  []Selection
}

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
	BooleanValue  bool
	EnumValue     string // Name, but not `true`, `false`, or `null`.
	ListValue     []Value
	ObjectValue   []ObjectField
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

// Type :
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
	Arguments []Argument
}
