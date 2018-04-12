package ast

// 2.2 Query Document
// http://facebook.github.io/graphql/October2016/#sec-Language.Query-Document

type Document struct {
	OperationDefinitions []OperationDefinition
	FragmentDefinitions  []FragmentDefinition
}

// 2.3 Operations
// http://facebook.github.io/graphql/October2016/#sec-Language.Operations

const (
	OperationTypeQuery OperationType = iota
	OperationTypeMutation
)

type OperationType int

type OperationDefinition struct {
	Type                OperationType
	Name                string
	VariableDefinitions []VariableDefinition
	Directives          []Directive
	SelectionSet        []Selection
}

// 2.4 Selection Sets
// http://facebook.github.io/graphql/October2016/#sec-Selection-Sets

type Selection struct {
	Field          Field
	FragmentSpread FragmentSpread
	InlineFragment InlineFragment
}

// 2.5 Fields
// http://facebook.github.io/graphql/October2016/#sec-Language.Fields

type Field struct {
	Alias        Alias
	Name         string
	Arguments    []Argument
	Directives   []Directive
	SelectionSet []Selection
}

// 2.6 Arguments
// http://facebook.github.io/graphql/October2016/#sec-Language.Arguments

// Argument :
type Argument struct {
	Name  string
	Value Value
}

// 2.7 Field Alias
// http://facebook.github.io/graphql/October2016/#sec-Field-Alias

type Alias struct {
	Name string
}

// 2.8 Fragments
// http://facebook.github.io/graphql/October2016/#sec-Language.Fragments

type FragmentSpread struct {
	Name       string // but not "on"
	Directives []Directive
}

type FragmentDefinition struct {
	Name          string // but not "on"
	TypeCondition TypeCondition
	Directives    []Directive
	SelectionSet  []Selection
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
// http://facebook.github.io/graphql/October2016/#sec-Input-Values

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

type Value struct {
	Kind         ValueKind
	IntValue     int
	FloatValue   float64
	StringValue  string
	BooleanValue bool
	EnumValue    string // Name, but not `true`, `false`, or `null`.
	ListValue    []Value
	ObjectValue  []ObjectField
}

type ObjectField struct {
	Name  string
	Value Value
}

// 2.10 Variables
// http://facebook.github.io/graphql/October2016/#sec-Language.Variables

type VariableDefinition struct {
	Name         string
	Type         Type
	DefaultValue Value
}

// 2.11 Input Types
// http://facebook.github.io/graphql/October2016/#sec-Input-Types

// Type :
const (
	TypeKindNamedType TypeKind = iota
	TypeKindListType
	TypeKindNonNullType
)

type TypeKind int

type Type struct {
	Kind        TypeKind
	NamedType   string
	ListType    []Type
	NonNullable bool
}

// 2.12 Directives
// http://facebook.github.io/graphql/October2016/#sec-Language.Directives

// Directive :
type Directive struct {
	Name      string
	Arguments []Argument
}
