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

// dumpDefinitions ...
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

// dumpDefinition ...
func (d *dumper) dumpDefinition(definition Definition) {
	switch definition.Kind {
	case DefinitionKindExecutable:
		d.dumpExecutableDefinition(definition.ExecutableDefinition)
	case DefinitionKindTypeSystem:
		// TODO
	}
}

// dumpExecutableDefinition ...
func (d *dumper) dumpExecutableDefinition(def *ExecutableDefinition) {
	switch def.Kind {
	case ExecutableDefinitionKindOperation:
		d.dumpOperationDefinition(def)
	case ExecutableDefinitionKindFragment:
		d.dumpFragmentDefinition(def)
	}
}

// dumpOperationDefinition ...
func (d *dumper) dumpOperationDefinition(def *ExecutableDefinition) {
	var shorthand bool
	if d.defs == 1 {
		shorthand = true
	}

	switch def.OperationType {
	case OperationTypeQuery:
		if !shorthand || def.Name != "" {
			io.WriteString(d.w, "query")
			io.WriteString(d.w, " ")
		}
	case OperationTypeMutation:
		io.WriteString(d.w, "mutation")
		io.WriteString(d.w, " ")
	case OperationTypeSubscription:
		io.WriteString(d.w, "subscription")
		io.WriteString(d.w, " ")
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

// dumpVariableDefinitions ...
func (d *dumper) dumpVariableDefinitions(definitions *VariableDefinitions) {
	definitionsLen := definitions.Len()
	if definitionsLen == 0 {
		return
	}

	io.WriteString(d.w, "(")

	definitions.ForEach(func(definition VariableDefinition, i int) {
		d.dumpVariableDefinition(definition)

		if i < definitionsLen-1 {
			io.WriteString(d.w, ", ")
		}
	})

	io.WriteString(d.w, ")")
}

// dumpVariableDefinition ...
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

// dumpSelections ...
func (d *dumper) dumpSelections(selections *Selections) {
	selectionsLen := selections.Len()
	if selectionsLen == 0 {
		return
	}

	io.WriteString(d.w, "{\n")

	selections.ForEach(func(selection Selection, i int) {
		d.dumpSelection(selection)
		io.WriteString(d.w, "\n")
	})

	io.WriteString(d.w, strings.Repeat(indentation, d.depth))
	io.WriteString(d.w, "}")
}

// dumpSelection ...
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

// dumpFieldSelection ...
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

// dumpArguments ...
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

// dumpArgument ...
func (d *dumper) dumpArgument(arg Argument) {
	io.WriteString(d.w, arg.Name)
	io.WriteString(d.w, ": ")

	d.dumpValue(arg.Value)
}

// dumpFragmentSpread ...
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

// dumpInlineFragment ...
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

// dumpFragmentDefinition ...
func (d *dumper) dumpFragmentDefinition(def *ExecutableDefinition) {
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

// dumpValue ...
func (d *dumper) dumpValue(value Value) {
	switch value.Kind {
	case ValueKindVariable:
		io.WriteString(d.w, "$")
		io.WriteString(d.w, value.StringValue)
	case ValueKindIntValue:
		io.WriteString(d.w, strconv.Itoa(value.IntValue))
	case ValueKindFloatValue:
		io.WriteString(d.w, fmt.Sprintf("%g", value.FloatValue))
	case ValueKindStringValue:
		hasLF := strings.Contains(value.StringValue, "\n")
		hasCR := strings.Contains(value.StringValue, "\r")

		// If the string contains a new line, we'll print it out as a multi-line string.
		if hasLF || hasCR {
			indent := strings.Repeat(indentation, d.depth)

			escaped := escapeGraphQLBlockString(value.StringValue)
			lines := strings.Split(escaped, "\n")

			buf := bytes.Buffer{}
			for i, line := range lines {
				buf.WriteString(indent)
				buf.WriteString(indentation)
				buf.WriteString(line)

				if i != len(lines) - 1 {
					buf.WriteRune('\n')
				}
			}

			io.WriteString(d.w, `"`)
			io.WriteString(d.w, `"`)
			io.WriteString(d.w, "\"\n")
			io.WriteString(d.w, buf.String())
			io.WriteString(d.w, indent)
			io.WriteString(d.w, "\n\"")
			io.WriteString(d.w, `"`)
			io.WriteString(d.w, `"`)
		} else {
			io.WriteString(d.w, `"`)
			io.WriteString(d.w, escapeGraphQLString(value.StringValue))
			io.WriteString(d.w, `"`)
		}

	case ValueKindBooleanValue:
		if value.BooleanValue {
			io.WriteString(d.w, "true")
		} else {
			io.WriteString(d.w, "false")
		}
	case ValueKindNullValue:
		io.WriteString(d.w, "null")
	case ValueKindEnumValue:
		io.WriteString(d.w, value.StringValue)
	case ValueKindListValue:
		io.WriteString(d.w, "[")

		vals := len(value.ListValue)
		for i, v := range value.ListValue {
			d.dumpValue(v)
			if i < vals-1 {
				io.WriteString(d.w, ", ")
			}
		}

		io.WriteString(d.w, "]")
	case ValueKindObjectValue:
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

// dumpType ...
func (d *dumper) dumpType(astType Type) {
	switch astType.Kind {
	case TypeKindNamedType:
		io.WriteString(d.w, astType.NamedType)
	case TypeKindListType:
		io.WriteString(d.w, "[")
		d.dumpType(*astType.ListType)
		io.WriteString(d.w, "]")
	}

	if astType.NonNullable {
		io.WriteString(d.w, "!")
	}
}

// dumpDirectives ...
func (d *dumper) dumpDirectives(directives *Directives) {
	directivesLen := directives.Len()
	if directivesLen == 0 {
		return
	}

	directives.ForEach(func(directive Directive, i int) {
		if i != 0 {
			io.WriteString(d.w, " ")
		}

		d.dumpDirective(directive)
	})
}

// dumpDirective ..
func (d *dumper) dumpDirective(directive Directive) {
	io.WriteString(d.w, "@")
	io.WriteString(d.w, directive.Name)

	d.dumpArguments(directive.Arguments)
}

// escapeGraphQLString takes a single-line GraphQL string and escapes all special characters that
// need to be escapes in it, returning the result.
func escapeGraphQLString(in string) string {
	buf := bytes.Buffer{}

	for _, r := range in {
		switch {
		case r >= utf8.RuneSelf && r <= '\uFFFF':
			escUni := fmt.Sprintf(`%x`, r)
			padding := strings.Repeat("0", utf8.UTFMax - len(escUni))

			buf.WriteString(fmt.Sprintf(`\u%s`, padding + escUni))
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
	return in
}
