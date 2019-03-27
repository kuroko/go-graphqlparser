package ast

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Fdump ...
func Fdump(w io.Writer, doc Document) {
	dmpr := dumper{defs: doc.Definitions.Len(), w: w}
	dmpr.dumpDefinitions(doc.Definitions)
}

// Dump ...
func Dump(doc Document) {
	Fdump(os.Stdout, doc)
}

// Sdump prints the given AST document as a GraphQL document string, allowing the parser to be
// easily validated against some given input. In fact, if formatted the same, the output of this
// function should match the input query given to the parser to produce the AST.
func Sdump(doc Document) string {
	buf := bytes.Buffer{}

	Fdump(&buf, doc)

	return buf.String()
}

// indentation is the string used as indentation, repeated as nesting becomes deeper.
const indentation = "  "

// dumper ...
type dumper struct {
	defs  int
	depth int
	w     io.Writer
}

// 2.2 Document
func (d *dumper) dumpDefinitions(definitions *Definitions) {
	definitions.ForEach(func(definition Definition, i int) {
		if i != 0 {
			io.WriteString(d.w, "\n")
		}

		d.dumpDefinition(definition)

		if i < d.defs-1 {
			io.WriteString(d.w, "\n")
		}
	})
}

func (d *dumper) dumpDefinition(definition Definition) {
	switch definition.Kind {
	case DefinitionKindExecutable:
		d.dumpExecutableDefinition(definition.ExecutableDefinition)
	case DefinitionKindTypeSystem:
		d.dumpTypeSystemDefinition(definition.TypeSystemDefinition)
	case DefinitionKindTypeSystemExtension:
		d.dumpTypeSystemExtension(definition.TypeSystemExtension)
	}
}

func (d *dumper) dumpExecutableDefinition(def *ExecutableDefinition) {
	switch def.Kind {
	case ExecutableDefinitionKindOperation:
		d.dumpOperationDefinition(def.OperationDefinition)
	case ExecutableDefinitionKindFragment:
		d.dumpFragmentDefinition(def.FragmentDefinition)
	}
}

// 2.3 Operations
func (d *dumper) dumpOperationDefinition(def *OperationDefinition) {
	var shorthand bool
	if d.defs == 1 {
		shorthand = true
	}

	switch def.Kind {
	case OperationDefinitionKindQuery:
		if !shorthand || def.Name != "" {
			io.WriteString(d.w, "query ")
		}
	case OperationDefinitionKindMutation:
		io.WriteString(d.w, "mutation ")
	case OperationDefinitionKindSubscription:
		io.WriteString(d.w, "subscription ")
	}

	if def.Name != "" {
		io.WriteString(d.w, def.Name)
	}

	if def.VariableDefinitions != nil {
		d.dumpVariableDefinitions(def.VariableDefinitions)
	}

	if def.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(def.Directives)
	}

	if !shorthand {
		io.WriteString(d.w, " ")
	}

	d.dumpSelections(def.SelectionSet)
}

// 2.4 Selection Sets
func (d *dumper) dumpSelections(selections *Selections) {
	io.WriteString(d.w, "{\n")

	selections.ForEach(func(selection Selection, _ int) {
		d.dumpSelection(selection)
		io.WriteString(d.w, "\n")
	})

	io.WriteString(d.w, strings.Repeat(indentation, d.depth))
	io.WriteString(d.w, "}")
}

func (d *dumper) dumpSelection(selection Selection) {
	d.depth++

	switch selection.Kind {
	case SelectionKindField:
		d.dumpFieldSelection(selection)
	case SelectionKindFragmentSpread:
		d.dumpFragmentSpread(selection)
	case SelectionKindInlineFragment:
		d.dumpInlineFragment(selection)
	}

	d.depth--
}

// 2.5 Fields
func (d *dumper) dumpFieldSelection(selection Selection) {
	indent := strings.Repeat(indentation, d.depth)

	io.WriteString(d.w, indent)

	if selection.Alias != "" {
		io.WriteString(d.w, selection.Alias)
		io.WriteString(d.w, ": ")
	}

	io.WriteString(d.w, selection.Name)

	d.dumpArguments(selection.Arguments)

	if selection.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(selection.Directives)
	}

	if selection.SelectionSet != nil {
		io.WriteString(d.w, " ")
		d.dumpSelections(selection.SelectionSet)
	}
}

// 2.6 Arguments
func (d *dumper) dumpArguments(args *Arguments) {
	argsLen := args.Len()
	if argsLen == 0 {
		return
	}

	io.WriteString(d.w, "(")

	args.ForEach(func(arg Argument, i int) {
		d.dumpArgument(arg)
		if i < argsLen-1 {
			io.WriteString(d.w, ", ")
		}
	})

	io.WriteString(d.w, ")")
}

func (d *dumper) dumpArgument(arg Argument) {
	io.WriteString(d.w, arg.Name)
	io.WriteString(d.w, ": ")

	d.dumpValue(arg.Value)
}

// 2.8 Fragments
func (d *dumper) dumpFragmentSpread(selection Selection) {
	indent := strings.Repeat(indentation, d.depth)

	io.WriteString(d.w, indent)

	io.WriteString(d.w, "...")
	io.WriteString(d.w, selection.Name)

	if selection.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(selection.Directives)
	}
}

// 2.8.2 Inline Fragments
func (d *dumper) dumpInlineFragment(selection Selection) {
	indent := strings.Repeat(indentation, d.depth)

	io.WriteString(d.w, indent)

	io.WriteString(d.w, "...")

	if selection.TypeCondition != nil {
		io.WriteString(d.w, " ")
		io.WriteString(d.w, "on")
		io.WriteString(d.w, " ")
		io.WriteString(d.w, selection.TypeCondition.NamedType.NamedType)
	}

	if selection.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(selection.Directives)
	}

	io.WriteString(d.w, " ")
	d.dumpSelections(selection.SelectionSet)
}

func (d *dumper) dumpFragmentDefinition(def *FragmentDefinition) {
	io.WriteString(d.w, "fragment")
	io.WriteString(d.w, " ")

	io.WriteString(d.w, def.Name)
	io.WriteString(d.w, " ")

	io.WriteString(d.w, "on")
	io.WriteString(d.w, " ")
	io.WriteString(d.w, def.TypeCondition.NamedType.NamedType)

	if def.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(def.Directives)
	}

	io.WriteString(d.w, " ")
	d.dumpSelections(def.SelectionSet)
}

// 2.9 Input Values
func (d *dumper) dumpValue(value Value) {
	switch value.Kind {
	case ValueKindVariable:
		io.WriteString(d.w, "$")
		io.WriteString(d.w, value.StringValue)
	case ValueKindInt:
		io.WriteString(d.w, strconv.Itoa(value.IntValue))
	case ValueKindFloat:
		io.WriteString(d.w, fmt.Sprintf("%g", value.FloatValue))
	case ValueKindString:
		hasLF := strings.Contains(value.StringValue, "\n")

		// If the string contains a new line, we'll print it out as a multi-line string.
		if hasLF {
			indent := strings.Repeat(indentation, d.depth)

			escaped := escapeGraphQLBlockString(value.StringValue)
			lines := strings.Split(escaped, "\n")

			buf := bytes.Buffer{}
			for i, line := range lines {
				buf.WriteString(indent)
				buf.WriteString(indentation) // Add one more level of indentation.
				buf.WriteString(line)

				if i != len(lines)-1 {
					buf.WriteRune('\n')
				}
			}

			io.WriteString(d.w, `"""`)
			io.WriteString(d.w, "\n")
			io.WriteString(d.w, buf.String())
			io.WriteString(d.w, "\n")
			io.WriteString(d.w, indent)
			io.WriteString(d.w, `"""`)
		} else {
			io.WriteString(d.w, `"`)
			io.WriteString(d.w, escapeGraphQLString(value.StringValue))
			io.WriteString(d.w, `"`)
		}

	case ValueKindBoolean:
		if value.BooleanValue {
			io.WriteString(d.w, "true")
		} else {
			io.WriteString(d.w, "false")
		}
	case ValueKindNull:
		io.WriteString(d.w, "null")
	case ValueKindEnum:
		io.WriteString(d.w, value.StringValue)
	case ValueKindList:
		io.WriteString(d.w, "[")

		vals := len(value.ListValue)
		for i, v := range value.ListValue {
			d.dumpValue(v)
			if i < vals-1 {
				io.WriteString(d.w, ", ")
			}
		}

		io.WriteString(d.w, "]")
	case ValueKindObject:
		io.WriteString(d.w, "{ ")

		vals := len(value.ObjectValue)
		for i, v := range value.ObjectValue {
			io.WriteString(d.w, v.Name)
			io.WriteString(d.w, ": ")

			d.dumpValue(v.Value)
			if i < vals-1 {
				io.WriteString(d.w, ", ")
			}
		}

		io.WriteString(d.w, " }")
	}
}

// 2.10 Variables
func (d *dumper) dumpVariableDefinitions(definitions *VariableDefinitions) {
	definitionsLen := definitions.Len()

	io.WriteString(d.w, "(")

	definitions.ForEach(func(definition VariableDefinition, i int) {
		d.dumpVariableDefinition(definition)

		if i < definitionsLen-1 {
			io.WriteString(d.w, ", ")
		}
	})

	io.WriteString(d.w, ")")
}

func (d *dumper) dumpVariableDefinition(definition VariableDefinition) {
	io.WriteString(d.w, "$")
	io.WriteString(d.w, definition.Name)
	io.WriteString(d.w, ": ")

	d.dumpType(definition.Type)

	if definition.DefaultValue != nil {
		io.WriteString(d.w, " = ")

		d.dumpValue(*definition.DefaultValue)
	}
}

// 2.11 Type References
func (d *dumper) dumpType(astType Type) {
	switch astType.Kind {
	case TypeKindNamed:
		io.WriteString(d.w, astType.NamedType)
	case TypeKindList:
		io.WriteString(d.w, "[")
		d.dumpType(*astType.ListType)
		io.WriteString(d.w, "]")
	}

	if astType.NonNullable {
		io.WriteString(d.w, "!")
	}
}

// 2.12 Directives
func (d *dumper) dumpDirectives(directives *Directives) {
	directives.ForEach(func(directive Directive, i int) {
		if i != 0 {
			io.WriteString(d.w, " ")
		}

		d.dumpDirective(directive)
	})
}

func (d *dumper) dumpDirective(directive Directive) {
	io.WriteString(d.w, "@")
	io.WriteString(d.w, directive.Name)

	d.dumpArguments(directive.Arguments)
}

// 3 Type System
func (d *dumper) dumpTypeSystemDefinition(def *TypeSystemDefinition) {
	switch def.Kind {
	case TypeSystemDefinitionKindSchema:
		d.dumpSchemaDefinition(def.SchemaDefinition)
	case TypeSystemDefinitionKindType:
		d.dumpTypeDefinition(def.TypeDefinition)
	case TypeSystemDefinitionKindDirective:
		d.dumpDirectiveDefinition(def.DirectiveDefinition)
	}
}

// 3.1 Type System Extensions
func (d *dumper) dumpTypeSystemExtension(ext *TypeSystemExtension) {
	switch ext.Kind {
	case TypeSystemExtensionKindSchema:
		d.dumpSchemaExtension(ext.SchemaExtension)
	case TypeSystemExtensionKindType:
		d.dumpTypeExtension(ext.TypeExtension)
	}
}

// 3.2 Schema
func (d *dumper) dumpSchemaDefinition(def *SchemaDefinition) {
	io.WriteString(d.w, "schema")

	if def.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(def.Directives)
	}

	io.WriteString(d.w, " {\n")
	d.depth++

	def.OperationTypeDefinitions.ForEach(func(opTypeDef OperationTypeDefinition, _ int) {
		d.dumpOperationTypeDefinition(opTypeDef)
		io.WriteString(d.w, "\n")
	})

	d.depth--
	io.WriteString(d.w, "}")
}

// 3.2.1 Root Operation Types
func (d *dumper) dumpOperationTypeDefinition(opTypeDef OperationTypeDefinition) {
	indent := strings.Repeat(indentation, d.depth)

	io.WriteString(d.w, indent)
	io.WriteString(d.w, opTypeDef.OperationType.String())
	io.WriteString(d.w, ": ")

	d.dumpType(opTypeDef.NamedType)
}

// 3.2.2 Schema Extension
func (d *dumper) dumpSchemaExtension(sext *SchemaExtension) {
	io.WriteString(d.w, "extend schema")

	if sext.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(sext.Directives)
	}

	if sext.OperationTypeDefinitions != nil {
		io.WriteString(d.w, " {\n")
		d.depth++

		sext.OperationTypeDefinitions.ForEach(func(otd OperationTypeDefinition, _ int) {
			d.dumpOperationTypeDefinition(otd)
			io.WriteString(d.w, "\n")
		})

		d.depth--
		io.WriteString(d.w, "}")
	}
}

// 3.3 Descriptions
func (d *dumper) dumpDescription(description string) {
	if description != "" {
		io.WriteString(d.w, "\"\"\"\n")
		io.WriteString(d.w, description)
		io.WriteString(d.w, "\n\"\"\"\n")
	}
}

// 3.4 Types
func (d *dumper) dumpTypeDefinition(td *TypeDefinition) {
	switch td.Kind {
	case TypeDefinitionKindScalar:
		d.dumpTypeDefinitionScalar(td)
	case TypeDefinitionKindObject:
		d.dumpTypeDefinitionObject(td)
	case TypeDefinitionKindInterface:
		d.dumpTypeDefinitionInterface(td)
	case TypeDefinitionKindUnion:
		d.dumpTypeDefinitionUnion(td)
	case TypeDefinitionKindEnum:
		d.dumpTypeDefinitionEnum(td)
	case TypeDefinitionKindInputObject:
		d.dumpTypeDefinitionInputObject(td)
	}
}

// 3.4.3 Type Extensions
func (d *dumper) dumpTypeExtension(te *TypeExtension) {
	switch te.Kind {
	case TypeExtensionKindScalar:
		d.dumpScalarTypeExtension(te)
	case TypeExtensionKindObject:
		d.dumpObjectTypeExtension(te)
	case TypeExtensionKindInterface:
		d.dumpInterfaceTypeExtension(te)
	case TypeExtensionKindUnion:
		d.dumpUnionTypeExtension(te)
	case TypeExtensionKindEnum:
		d.dumpEnumTypeExtension(te)
	case TypeExtensionKindInputObject:
		d.dumpInputObjectTypeExtension(te)
	}
}

// 3.5 Scalars
func (d *dumper) dumpTypeDefinitionScalar(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "scalar ")
	io.WriteString(d.w, td.Name)

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}
}

// 3.5.6 Scalar Extensions
func (d *dumper) dumpScalarTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend scalar ")
	io.WriteString(d.w, te.Name)

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}
}

// 3.6 Objects
func (d *dumper) dumpTypeDefinitionObject(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "type ")
	io.WriteString(d.w, td.Name)

	if td.ImplementsInterface != nil {
		io.WriteString(d.w, " ")
		d.dumpImplementsInterfaces(td.ImplementsInterface, td.FieldsDefinition != nil)
	}

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}

	if td.FieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpFieldsDefinition(td.FieldsDefinition)
	}
}

// dumpImplementsInterfaces dumps the ImplementsInterfaces. The formatting of
// ImplementsInterfaces is flexible, so the dumper uses a sensible pattern
// which you can see an example of below.
//
// type Climbing implements Ropes
// type Climbing implements Ropes & Rocks
// type Climbing implements
//   & Ropes
//   & Rocks
//   & Chalk
func (d *dumper) dumpImplementsInterfaces(ii *Types, hasFields bool) {
	io.WriteString(d.w, "implements")

	l := ii.Len()

	separator := " & "
	if l < 3 {
		io.WriteString(d.w, " ")
	} else {
		if hasFields {
			separator = indentation + separator
		}

		separator = "\n " + separator
		io.WriteString(d.w, separator)
	}

	ii.ForEach(func(t Type, i int) {
		io.WriteString(d.w, t.NamedType)
		if i < l-1 {
			io.WriteString(d.w, separator)
		}
	})
}

func (d *dumper) dumpFieldsDefinition(fields *FieldDefinitions) {
	io.WriteString(d.w, "{\n")

	fields.ForEach(func(field FieldDefinition, _ int) {
		d.dumpFieldDefinition(field)
		io.WriteString(d.w, "\n")
	})

	io.WriteString(d.w, "}")
}

func (d *dumper) dumpFieldDefinition(field FieldDefinition) {
	d.dumpDescription(field.Description)

	io.WriteString(d.w, indentation)
	io.WriteString(d.w, field.Name)

	if field.ArgumentsDefinition != nil {
		d.dumpArgumentsDefinition(field.ArgumentsDefinition)
	}

	io.WriteString(d.w, ": ")
	io.WriteString(d.w, field.Type.NamedType)
	if field.Type.NonNullable {
		io.WriteString(d.w, "!")
	}

	if field.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(field.Directives)
	}
}

// 3.6.1 Field Arguments
func (d *dumper) dumpArgumentsDefinition(arguments *InputValueDefinitions) {
	io.WriteString(d.w, "(")

	arguments.ForEach(func(ivd InputValueDefinition, i int) {
		if i > 0 {
			io.WriteString(d.w, ", ")
		}
		d.dumpInputValueDefinition(ivd)
	})

	io.WriteString(d.w, ")")
}

// TODO: with > 3 arguments, need line breaks in-between.
func (d *dumper) dumpInputValueDefinition(ivd InputValueDefinition) {
	d.dumpDescription(ivd.Description)

	io.WriteString(d.w, ivd.Name)
	io.WriteString(d.w, ": ")
	io.WriteString(d.w, ivd.Type.NamedType)
	if ivd.Type.NonNullable {
		io.WriteString(d.w, "!")
	}

	if ivd.DefaultValue != nil {
		io.WriteString(d.w, " = ")
		d.dumpValue(*ivd.DefaultValue)
	}

	if ivd.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(ivd.Directives)
	}
}

// 3.6.3 Interfaces
func (d *dumper) dumpObjectTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend type ")
	io.WriteString(d.w, te.Name)

	if te.ImplementsInterface != nil {
		io.WriteString(d.w, " ")
		d.dumpImplementsInterfaces(te.ImplementsInterface, te.FieldsDefinition != nil)
	}

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}

	if te.FieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpFieldsDefinition(te.FieldsDefinition)
	}
}

// 3.7
func (d *dumper) dumpTypeDefinitionInterface(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "interface ")
	io.WriteString(d.w, td.Name)

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}

	if td.FieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpFieldsDefinition(td.FieldsDefinition)
	}
}

// 3.7.1 Interface Extensions
func (d *dumper) dumpInterfaceTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend interface ")
	io.WriteString(d.w, te.Name)

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}

	if te.FieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpFieldsDefinition(te.FieldsDefinition)
	}
}

// 3.8 Unions
func (d *dumper) dumpTypeDefinitionUnion(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "union ")
	io.WriteString(d.w, td.Name)

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}

	if td.UnionMemberTypes != nil {
		io.WriteString(d.w, " ")
		d.dumpUnionMemberTypes(td.UnionMemberTypes)
	}
}

// dumpUnionMemberTypes dumps the UnionMemberTypes. The formatting of
// UnionMemberTypes is flexible, so the dumper uses a sensible pattern
// which you can see an example of below.
//
// union SearchResult = Photo
// union SearchResult = Photo | Person
// union SearchResult =
//   | Photo
//   | Person
//   | Plate
func (d *dumper) dumpUnionMemberTypes(umt *Types) {
	io.WriteString(d.w, "=")

	l := umt.Len()

	separator := " | "
	if l < 3 {
		io.WriteString(d.w, " ")
	} else {
		separator = "\n " + separator
		io.WriteString(d.w, separator)
	}

	umt.ForEach(func(t Type, i int) {
		io.WriteString(d.w, t.NamedType)
		if i < l-1 {
			io.WriteString(d.w, separator)
		}
	})
}

// 3.8.1 Union Extensions
func (d *dumper) dumpUnionTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend union ")
	io.WriteString(d.w, te.Name)

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}

	if te.UnionMemberTypes != nil {
		io.WriteString(d.w, " ")
		d.dumpUnionMemberTypes(te.UnionMemberTypes)
	}
}

// 3.9 Enums
func (d *dumper) dumpTypeDefinitionEnum(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "enum ")
	io.WriteString(d.w, td.Name)

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}

	if td.EnumValuesDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpEnumValuesDefinition(td.EnumValuesDefinition)
	}
}

func (d *dumper) dumpEnumValuesDefinition(evds *EnumValueDefinitions) {
	io.WriteString(d.w, "{\n")

	evds.ForEach(func(evd EnumValueDefinition, _ int) {
		d.dumpEnumValueDefinition(evd)
		io.WriteString(d.w, "\n")
	})

	io.WriteString(d.w, "}")
}

func (d *dumper) dumpEnumValueDefinition(evd EnumValueDefinition) {
	d.dumpDescription(evd.Description)

	io.WriteString(d.w, indentation)
	io.WriteString(d.w, evd.EnumValue)

	if evd.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(evd.Directives)
	}
}

// 3.9.1 Enum Extensions
func (d *dumper) dumpEnumTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend enum ")
	io.WriteString(d.w, te.Name)

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}

	if te.EnumValuesDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpEnumValuesDefinition(te.EnumValuesDefinition)
	}
}

// 3.10 Input Objects
func (d *dumper) dumpTypeDefinitionInputObject(td *TypeDefinition) {
	d.dumpDescription(td.Description)

	io.WriteString(d.w, "input ")
	io.WriteString(d.w, td.Name)

	if td.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(td.Directives)
	}

	if td.InputFieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpInputFieldsDefinition(td.InputFieldsDefinition)
	}
}

func (d *dumper) dumpInputFieldsDefinition(ivds *InputValueDefinitions) {
	io.WriteString(d.w, "{\n")

	ivds.ForEach(func(ivd InputValueDefinition, i int) {
		io.WriteString(d.w, indentation)
		d.dumpInputValueDefinition(ivd)
		io.WriteString(d.w, "\n")
	})

	io.WriteString(d.w, "}")
}

// 3.10.1 Input Object Extensions
func (d *dumper) dumpInputObjectTypeExtension(te *TypeExtension) {
	io.WriteString(d.w, "extend input ")
	io.WriteString(d.w, te.Name)

	if te.Directives != nil {
		io.WriteString(d.w, " ")
		d.dumpDirectives(te.Directives)
	}

	if te.InputFieldsDefinition != nil {
		io.WriteString(d.w, " ")
		d.dumpInputFieldsDefinition(te.InputFieldsDefinition)
	}
}

// 3.13 Directives
// dumpDirectiveDefinition ...
func (d *dumper) dumpDirectiveDefinition(def *DirectiveDefinition) {
	d.dumpDescription(def.Description)

	io.WriteString(d.w, "directive @")
	io.WriteString(d.w, def.Name)

	if def.ArgumentsDefinition != nil {
		d.dumpArgumentsDefinition(def.ArgumentsDefinition)
	}

	io.WriteString(d.w, " on")
	d.dumpDirectiveLocations(def.DirectiveLocations)
}

// dumpDirectiveLocations dumps the DirectiveLocation. The formatting of
// DirectiveLocation is flexible, so the dumper uses a sensible pattern
// which you can see an example of below.
//
// directive @example on FIELD
// directive @example on FIELD_DEFINITION | ARGUMENT_DEFINITION
// directive @example on
//   | FIELD
//   | FRAGMENT_SPREAD
//   | INLINE_FRAGMENT
func (d *dumper) dumpDirectiveLocations(dls DirectiveLocation) {
	dll := len(NamesByDirectiveLocations)

	var locs []string
	for i := 0; i <= dll; i++ {
		bit := dls & (1 << uint(i))
		if bit == 0 {
			continue
		}

		locs = append(locs, NamesByDirectiveLocations[bit])
	}

	l := len(locs)

	separator := " | "
	if l < 3 {
		io.WriteString(d.w, " ")
	} else {
		separator = "\n " + separator
		io.WriteString(d.w, separator)
	}

	for i, loc := range locs {
		io.WriteString(d.w, loc)
		if i < l-1 {
			io.WriteString(d.w, separator)
		}
	}
}

/*****************************************************************************
 * Utility functions                                                         *
 *****************************************************************************/

// escapeGraphQLString takes a single-line GraphQL string and escapes all special characters that
// need to be escapes in it, returning the result.
func escapeGraphQLString(in string) string {
	buf := bytes.Buffer{}

	for _, r := range in {
		switch {
		case r >= utf8.RuneSelf && r <= '\uFFFF':
			escUni := fmt.Sprintf(`%x`, r)
			padding := strings.Repeat("0", utf8.UTFMax-len(escUni))

			buf.WriteString(fmt.Sprintf(`\u%s`, padding+escUni))
		case r == '"':
			buf.WriteString(`\"`)
		case r == '\\':
			buf.WriteString(`\\`)
		case r == '/':
			buf.WriteString(`\/`)
		case r == '\b':
			buf.WriteString(`\b`)
		case r == '\f':
			buf.WriteString(`\f`)
		case r == '\t':
			buf.WriteString(`\t`)
		default:
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// escapeGraphQLBlockString takes a GraphQL block string and escapes all special characters that
// need to be escapes in it, returning the result.
func escapeGraphQLBlockString(in string) string {
	return strings.Replace(in, `"""`, `\"""`, -1)
}
