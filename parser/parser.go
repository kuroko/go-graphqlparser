package parser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unsafe"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
)

// Parser is a parser for GraphQL documents.
type Parser struct {
	lexer *lexer.Lexer
	token lexer.Token

	// Stateful checks.
	hasShorthandQuery bool
}

// New returns a new Parser instance.
func New(input []byte) *Parser {
	return &Parser{
		lexer: lexer.New(input),
	}
}

// Parse loops over the lexically analysed tokens produced by the lexer from the raw bytes of input
// and parses them into an AST of the GraphQL Document which it returns.
func (p *Parser) Parse() (ast.Document, error) {
	var document ast.Document

	p.token = p.lexer.Scan()

	var definitions *ast.Definitions

	for {
		// This should be set during the first iteration.
		if p.hasShorthandQuery {
			return ast.Document{}, p.unexpected(p.token, p.expected(token.EOF))
		}

		definition, err := p.parseDefinition(document)
		if err != nil {
			return ast.Document{}, err
		}

		definitions = definitions.Add(definition)

		if p.peek0(token.Illegal) {
			return ast.Document{}, p.unexpected(p.token)
		}

		if p.peek0(token.EOF) {
			break
		}
	}

	document.Definitions = definitions.Reverse()

	return document, nil
}

// parseDefinition ...
func (p *Parser) parseDefinition(document ast.Document) (ast.Definition, error) {
	var err error

	// We can only allow a shorthand query if it's the only definition.
	p.hasShorthandQuery = document.Definitions.Len() == 0 && p.token.Literal == "{"

	// ExecutableDefinition...
	if p.peekn(token.Name, "query", "mutation", "subscription") || p.peek1(token.Punctuator, "{") {
		definition := ast.Definition{}
		definition.Kind = ast.DefinitionKindExecutable
		definition.ExecutableDefinition, err = p.parseOperationDefinition(p.hasShorthandQuery)

		return definition, err
	}

	// ExecutableDefinition...
	if p.peek1(token.Name, "fragment") {
		definition := ast.Definition{}
		definition.Kind = ast.DefinitionKindExecutable
		definition.ExecutableDefinition, err = p.parseFragmentDefinition()

		return definition, err
	}

	// TypeSystemExtension...
	if p.peek1(token.Name, "extend") {
		definition := ast.Definition{}
		definition.Kind = ast.DefinitionKindTypeSystemExtension
		definition.TypeSystemExtension, err = p.parseTypeSystemExtension()

		return definition, err
	}

	var description string

	// If we have a description, then we should encounter a `scalar`, a `type`, a `interface`, a
	// `union`, an `enum`, an `input`, or a `directive`.
	if tok, ok := p.consume0(token.StringValue); ok {
		description = tok.Literal
	}

	typeDefLits := make([]string, 0, 8)
	typeDefLits = append(typeDefLits, "scalar", "type", "interface", "union", "enum", "input", "directive")

	if description == "" {
		typeDefLits = append(typeDefLits, "schema")
	}

	if p.peekn(token.Name, typeDefLits...) {
		definition := ast.Definition{}
		definition.Kind = ast.DefinitionKindTypeSystem
		definition.TypeSystemDefinition, err = p.parseTypeSystemDefinition(description)

		return definition, err
	}

	return ast.Definition{}, p.unexpected(p.token,
		p.expected(token.Name, "query", "mutation", "fragment"),
		p.expected(token.Punctuator, "{"),
	)

}

// parseOperationDefinition ...
func (p *Parser) parseOperationDefinition(isQuery bool) (*ast.ExecutableDefinition, error) {
	var variableDefinitions *ast.VariableDefinitions
	var directives *ast.Directives

	var name string
	var err error

	opType := ast.OperationTypeQuery

	if !isQuery {
		opType, err = p.parseOperationType()
		if err != nil {
			return nil, err
		}

		if tok, ok := p.consume0(token.Name); ok {
			name = tok.Literal
		}

		variableDefinitions, err = p.parseVariableDefinitions()
		if err != nil {
			return nil, err
		}

		directives, err = p.parseDirectives()
		if err != nil {
			return nil, err
		}
	}

	selectionSet, err := p.parseSelectionSet(false)
	if err != nil {
		return nil, err
	}

	return &ast.ExecutableDefinition{
		Kind:                ast.ExecutableDefinitionKindOperation,
		OperationType:       opType,
		Name:                name,
		VariableDefinitions: variableDefinitions,
		Directives:          directives,
		SelectionSet:        selectionSet,
	}, nil
}

// parseOperationType ...
func (p *Parser) parseOperationType() (ast.OperationType, error) {
	tok, err := p.mustConsumen(token.Name, "query", "mutation", "subscription")
	if err != nil {
		return -1, err
	}

	switch tok.Literal {
	case "query":
		return ast.OperationTypeQuery, nil
	case "mutation":
		return ast.OperationTypeMutation, nil
	default:
		return ast.OperationTypeSubscription, nil
	}
}

// parseFragmentDefinition ...
func (p *Parser) parseFragmentDefinition() (*ast.ExecutableDefinition, error) {
	if !p.skip1(token.Name, "fragment") {
		return nil, nil
	}

	tok, ok := p.consume0(token.Name)
	if !ok {
		return nil, nil
	}

	if tok.Literal == "on" {
		return nil, p.unexpected(p.token, p.expected(token.Name, "!on"))
	}

	condition, err := p.parseTypeCondition()
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	selections, err := p.parseSelectionSet(false)
	if err != nil {
		return nil, err
	}

	definition := &ast.ExecutableDefinition{}
	definition.Kind = ast.ExecutableDefinitionKindFragment
	definition.Name = tok.Literal
	definition.TypeCondition = condition
	definition.Directives = directives
	definition.SelectionSet = selections

	return definition, nil
}

// parseTypeCondition ...
func (p *Parser) parseTypeCondition() (*ast.TypeCondition, error) {
	_, err := p.mustConsume1(token.Name, "on")
	if err != nil {
		return nil, err
	}

	conType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	if conType.Kind != ast.TypeKindNamedType {
		return nil, p.unexpected(p.token, "NamedType")
	}

	condition := &ast.TypeCondition{}
	condition.NamedType = conType

	return condition, nil
}

// parseVariableDefinitions ...
func (p *Parser) parseVariableDefinitions() (*ast.VariableDefinitions, error) {
	if !p.skip1(token.Punctuator, "(") {
		return nil, nil
	}

	var definitions *ast.VariableDefinitions

	for {
		if _, err := p.mustConsume1(token.Punctuator, "$"); err != nil {
			return nil, err
		}

		tok, err := p.mustConsume0(token.Name)
		if err != nil {
			return nil, err
		}

		definition := ast.VariableDefinition{}
		definition.Name = tok.Literal

		if _, err := p.mustConsume1(token.Punctuator, ":"); err != nil {
			return nil, err
		}

		definition.Type, err = p.parseType()
		if err != nil {
			return nil, err
		}

		definition.DefaultValue, err = p.parseDefaultValue()
		if err != nil {
			return nil, err
		}

		definitions = definitions.Add(definition)

		if p.peek1(token.Punctuator, ")") {
			break
		}
	}

	if _, err := p.mustConsume1(token.Punctuator, ")"); err != nil {
		return nil, err
	}

	return definitions.Reverse(), nil
}

// parseDirectives ...
func (p *Parser) parseDirectives() (*ast.Directives, error) {
	var directives *ast.Directives

	for p.peek1(token.Punctuator, "@") {
		_, err := p.mustConsume1(token.Punctuator, "@")
		if err != nil {
			return nil, err
		}

		name, err := p.mustConsume0(token.Name)
		if err != nil {
			return nil, err
		}

		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		directive := ast.Directive{}
		directive.Name = name.Literal
		directive.Arguments = args

		directives = directives.Add(directive)
	}

	if directives != nil {
		return directives.Reverse(), nil
	}

	return nil, nil
}

// parseSelectionSet ...
func (p *Parser) parseSelectionSet(optional bool) (*ast.Selections, error) {
	if optional && !p.skip1(token.Punctuator, "{") {
		return nil, nil
	}

	if !optional && !p.skip1(token.Punctuator, "{") {
		return nil, p.unexpected(p.token, p.expected(token.Punctuator, "{"))
	}

	var selections *ast.Selections

	for {
		var selection ast.Selection
		var err error

		if p.skip1(token.Punctuator, "...") {
			if p.peek0(token.Name) && p.token.Literal != "on" {
				selection, err = p.parseFragmentSpread()
				if err != nil {
					return nil, err
				}
			} else {
				selection, err = p.parseInlineFragment()
				if err != nil {
					return nil, err
				}
			}
		} else {
			selection, err = p.parseField()
			if err != nil {
				return nil, err
			}
		}

		selections = selections.Add(selection)

		if p.peek1(token.Punctuator, "}") || p.peek0(token.EOF) {
			break
		}
	}

	_, err := p.mustConsume1(token.Punctuator, "}")
	if err != nil {
		return nil, err
	}

	return selections.Reverse(), nil
}

// parseFragmentSpread ...
func (p *Parser) parseFragmentSpread() (ast.Selection, error) {
	var selection ast.Selection

	tok, err := p.mustConsume0(token.Name)
	if err != nil {
		return selection, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return selection, err
	}

	selection.Kind = ast.SelectionKindFragmentSpread
	selection.Name = tok.Literal
	selection.Directives = directives

	return selection, nil
}

// parseInlineFragment ...
func (p *Parser) parseInlineFragment() (ast.Selection, error) {
	var selection ast.Selection
	var condition *ast.TypeCondition
	var err error

	if p.peek1(token.Name, "on") {
		condition, err = p.parseTypeCondition()
		if err != nil {
			return selection, err
		}
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return selection, err
	}

	selections, err := p.parseSelectionSet(false)
	if err != nil {
		return selection, err
	}

	selection.Kind = ast.SelectionKindInlineFragment
	selection.TypeCondition = condition
	selection.Directives = directives
	selection.SelectionSet = selections

	return selection, nil
}

// parseField ...
func (p *Parser) parseField() (ast.Selection, error) {
	var selection ast.Selection
	var name string
	var alias string

	nameTok, err := p.mustConsume0(token.Name)
	if err != nil {
		return selection, err
	}

	name = nameTok.Literal

	if p.skip1(token.Punctuator, ":") {
		alias = nameTok.Literal

		nameTok, err = p.mustConsume0(token.Name)
		if err != nil {
			return selection, err
		}

		name = nameTok.Literal
	}

	arguments, err := p.parseArguments()
	if err != nil {
		return selection, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return selection, err
	}

	selections, err := p.parseSelectionSet(true)
	if err != nil {
		return selection, err
	}

	selection.Kind = ast.SelectionKindField
	selection.Name = name
	selection.Alias = alias
	selection.Arguments = arguments
	selection.Directives = directives
	selection.SelectionSet = selections

	return selection, nil
}

// parseArguments ...
func (p *Parser) parseArguments() (*ast.Arguments, error) {
	if !p.skip1(token.Punctuator, "(") {
		return nil, nil
	}

	var arguments *ast.Arguments

	for !p.skip1(token.Punctuator, ")") {
		name, err := p.mustConsume0(token.Name)
		if err != nil {
			return nil, err
		}

		_, err = p.mustConsume1(token.Punctuator, ":")
		if err != nil {
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		argument := ast.Argument{}
		argument.Name = name.Literal
		argument.Value = value

		arguments = arguments.Add(argument)
	}

	return arguments.Reverse(), nil
}

// parseDefaultValue ...
func (p *Parser) parseDefaultValue() (*ast.Value, error) {
	if !p.skip1(token.Punctuator, "=") {
		return nil, nil
	}

	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	return &val, nil
}

// parseValue ...
func (p *Parser) parseValue() (ast.Value, error) {
	if p.skip1(token.Punctuator, "$") {
		tok, err := p.mustConsume0(token.Name)
		if err != nil {
			return ast.Value{}, err
		}

		return ast.Value{
			Kind:        ast.ValueKindVariable,
			StringValue: tok.Literal,
		}, nil
	}

	if tok, ok := p.consume0(token.IntValue); ok {
		iv, err := strconv.Atoi(tok.Literal)
		if err != nil {
			return ast.Value{}, err
		}

		return ast.Value{
			Kind:     ast.ValueKindIntValue,
			IntValue: iv,
		}, nil
	}

	if tok, ok := p.consume0(token.FloatValue); ok {
		fv, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return ast.Value{}, err
		}

		return ast.Value{
			Kind:       ast.ValueKindFloatValue,
			FloatValue: fv,
		}, nil
	}

	if tok, ok := p.consume0(token.StringValue); ok {
		return ast.Value{
			Kind:        ast.ValueKindStringValue,
			StringValue: tok.Literal,
		}, nil
	}

	if tok, ok := p.consumen(token.Name, "true", "false"); ok {
		return ast.Value{
			Kind:         ast.ValueKindBooleanValue,
			BooleanValue: tok.Literal == "true",
		}, nil
	}

	if p.skip1(token.Name, "null") {
		return ast.Value{
			Kind: ast.ValueKindNullValue,
		}, nil
	}

	if tok, ok := p.consume0(token.Name); ok {
		return ast.Value{
			Kind:        ast.ValueKindEnumValue,
			StringValue: tok.Literal,
		}, nil
	}

	if p.skip1(token.Punctuator, "[") {
		list := ast.Value{}
		list.Kind = ast.ValueKindListValue

		for !p.skip1(token.Punctuator, "]") {
			val, err := p.parseValue()
			if err != nil {
				return list, err
			}

			list.ListValue = append(list.ListValue, val)
		}

		return list, nil
	}

	if p.skip1(token.Punctuator, "{") {
		object := ast.Value{}
		object.Kind = ast.ValueKindObjectValue

		for !p.skip1(token.Punctuator, "}") {
			tok, err := p.mustConsume0(token.Name)
			if err != nil {
				return object, err
			}

			_, err = p.mustConsume1(token.Punctuator, ":")
			if err != nil {
				return object, err
			}

			value, err := p.parseValue()
			if err != nil {
				return object, err
			}

			field := ast.ObjectField{}
			field.Name = tok.Literal
			field.Value = value

			object.ObjectValue = append(object.ObjectValue, field)
		}

		return object, nil
	}

	return ast.Value{}, errors.New("TODO: see `parseDefinition`")
}

// parseType ...
func (p *Parser) parseType() (ast.Type, error) {
	var astType ast.Type

	// If we hit an opening square brace, we've got a list type, time to dive in.
	if p.skip1(token.Punctuator, "[") {
		astType.Kind = ast.TypeKindListType

		itemType, err := p.parseType()
		if err != nil {
			return astType, nil
		}

		astType.ListType = &itemType

		if _, err := p.mustConsume1(token.Punctuator, "]"); err != nil {
			return astType, err
		}
	} else {
		astType.Kind = ast.TypeKindNamedType

		tok, err := p.mustConsume0(token.Name)
		if err != nil {
			return astType, err
		}

		astType.NamedType = tok.Literal
	}

	if p.skip1(token.Punctuator, "!") {
		astType.NonNullable = true
	}

	return astType, nil
}

// 3.4.3 Type Extensions
// parseTypeExtension ...
func (p *Parser) parseTypeExtension() (*ast.TypeExtension, error) {
	kind, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	te := &ast.TypeExtension{}

	if kind.Literal == "type" && p.peek1(token.Name, "implements") {
		ii, err := p.parseImplementsInterfaces()
		if err != nil {
			return nil, err
		}
		te.ImplementsInterface = ii
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	switch kind.Literal {
	case "scalar":
		te.Kind = ast.TypeExtensionKindScalar

	case "type":
		fd, err := p.parseFieldsDefinition()
		if err != nil {
			return nil, err
		}
		te.FieldsDefinition = fd
		te.Kind = ast.TypeExtensionKindObject

	case "interface":
		fd, err := p.parseFieldsDefinition()
		if err != nil {
			return nil, err
		}
		te.FieldsDefinition = fd
		te.Kind = ast.TypeExtensionKindInterface

	case "union":
		umt, err := p.parseUnionMemberTypes()
		if err != nil {
			return nil, err
		}
		te.UnionMemberTypes = umt
		te.Kind = ast.TypeExtensionKindUnion

	case "enum":
		evd, err := p.parseEnumValuesDefinition()
		if err != nil {
			return nil, err
		}
		te.EnumValuesDefinition = evd
		te.Kind = ast.TypeExtensionKindEnum

	case "input":
		ifd, err := p.parseInputFieldsDefinition()
		if err != nil {
			return nil, err
		}
		te.InputFieldsDefinition = ifd
		te.Kind = ast.TypeExtensionKindInputObject
	}

	te.Name = name.Literal
	te.Directives = directives

	return te, nil
}

// parseTypeSystemDefinition ...
func (p *Parser) parseTypeSystemDefinition(description string) (*ast.TypeSystemDefinition, error) {
	// definition.SchemaDefinition
	if p.peek1(token.Name, "schema") {
		schemaDef, err := p.parseSchemaDefinition()
		if err != nil {
			return nil, err
		}

		tsDefinition := &ast.TypeSystemDefinition{}
		tsDefinition.Kind = ast.TypeSystemDefinitionKindSchema
		tsDefinition.SchemaDefinition = schemaDef

		return tsDefinition, nil
	}

	// definition.DirectiveDefinition
	if p.peek1(token.Name, "directive") {
		directiveDef, err := p.parseDirectiveDefinition(description)
		if err != nil {
			return nil, err
		}

		tsDefinition := &ast.TypeSystemDefinition{}
		tsDefinition.Kind = ast.TypeSystemDefinitionKindDirective
		tsDefinition.DirectiveDefinition = directiveDef

		return tsDefinition, nil
	}

	// definition.TypeDefinition
	if p.peekn(token.Name, "scalar", "type", "interface", "union", "enum", "input") {
		typeDef, err := p.parseTypeDefinition()
		if err != nil {
			return nil, err
		}

		if typeDef != nil {
			// Avoid passing description down several layers, copying along the way.
			typeDef.Description = description
		}

		tsDefinition := &ast.TypeSystemDefinition{}
		tsDefinition.Kind = ast.TypeSystemDefinitionKindType
		tsDefinition.TypeDefinition = typeDef

		return tsDefinition, nil
	}

	return &ast.TypeSystemDefinition{}, nil
}

// 3.4.3
// parseTypeSystemExtension ...
func (p *Parser) parseTypeSystemExtension() (*ast.TypeSystemExtension, error) {
	if !p.skip1(token.Name, "extend") {
		return nil, p.unexpected(p.token, p.expected(token.Name, "extend"))
	}

	// definition.SchemaDefinition
	if p.peek1(token.Name, "schema") {
		schemaExt, err := p.parseSchemaExtension()
		if err != nil {
			return nil, err
		}

		tsExtension := &ast.TypeSystemExtension{}
		tsExtension.Kind = ast.TypeSystemExtensionKindSchema
		tsExtension.SchemaExtension = schemaExt

		return tsExtension, nil
	}

	// definition.TypeExtension
	if p.peekn(token.Name, "scalar", "type", "interface", "union", "enum", "input") {
		typeExt, err := p.parseTypeExtension()
		if err != nil {
			return nil, err
		}

		tsExtension := &ast.TypeSystemExtension{}
		tsExtension.Kind = ast.TypeSystemExtensionKindType
		tsExtension.TypeExtension = typeExt

		return tsExtension, nil
	}

	return nil, nil
}

// parseSchemaDefinition ...
func (p *Parser) parseSchemaDefinition() (*ast.SchemaDefinition, error) {
	if !p.skip1(token.Name, "schema") {
		return nil, p.unexpected(p.token, p.expected(token.Name, "schema"))
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	if !p.skip1(token.Punctuator, "{") {
		return nil, p.unexpected(p.token, p.expected(token.Punctuator, "{"))
	}

	var rootOperationTypeDefinitions *ast.RootOperationTypeDefinitions

	for {
		opType, err := p.parseOperationType()
		if err != nil {
			return nil, err
		}

		if !p.skip1(token.Punctuator, ":") {
			return nil, p.unexpected(p.token, p.expected(token.Punctuator, ":"))
		}

		namedType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		if namedType.Kind != ast.TypeKindNamedType {
			return nil, p.unexpected(p.token, "NamedType")
		}

		rootOperationTypeDefinitions = rootOperationTypeDefinitions.Add(ast.RootOperationTypeDefinition{
			OperationType: opType,
			NamedType:     namedType,
		})

		if p.peek1(token.Punctuator, "}") || p.peek0(token.EOF) {
			break
		}
	}

	_, err = p.mustConsume1(token.Punctuator, "}")
	if err != nil {
		return nil, err
	}

	return &ast.SchemaDefinition{
		Directives:                   directives,
		RootOperationTypeDefinitions: rootOperationTypeDefinitions.Reverse(),
	}, nil
}

// 3.2.2
// parseSchemaExtension ...
func (p *Parser) parseSchemaExtension() (*ast.SchemaExtension, error) {
	if !p.skip1(token.Name, "schema") {
		return nil, p.unexpected(p.token, p.expected(token.Name, "schema"))
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	var operationTypeDefinitions *ast.OperationTypeDefinitions

	if p.peek1(token.Punctuator, "{") {
		operationTypeDefinitions, err = p.parseOperationTypeDefinitions()
		if err != nil {
			return nil, err
		}
	}

	return &ast.SchemaExtension{
		Directives:               directives,
		OperationTypeDefinitions: operationTypeDefinitions.Reverse(),
	}, nil
}

func (p *Parser) parseOperationTypeDefinitions() (*ast.OperationTypeDefinitions, error) {
	if !p.skip1(token.Punctuator, "{") {
		return nil, p.unexpected(p.token, p.expected(token.Punctuator, "{"))
	}

	var operationTypeDefinitions *ast.OperationTypeDefinitions

	for {
		opType, err := p.parseOperationType()
		if err != nil {
			return nil, err
		}

		if !p.skip1(token.Punctuator, ":") {
			return nil, p.unexpected(p.token, p.expected(token.Punctuator, ":"))
		}

		namedType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		if namedType.Kind != ast.TypeKindNamedType {
			return nil, p.unexpected(p.token, "NamedType")
		}

		operationTypeDefinitions = operationTypeDefinitions.Add(ast.OperationTypeDefinition{
			OperationType: opType,
			NamedType:     namedType,
		})

		if p.peek1(token.Punctuator, "}") || p.peek0(token.EOF) {
			break
		}
	}

	_, err := p.mustConsume1(token.Punctuator, "}")
	if err != nil {
		return nil, err
	}

	return operationTypeDefinitions, nil
}

// parseArgumentsDefinition ...
func (p *Parser) parseArgumentsDefinition() (*ast.InputValueDefinitions, error) {
	var defs *ast.InputValueDefinitions

	if !p.skip1(token.Punctuator, "(") {
		return nil, nil
	}

	for {
		def, err := p.parseInputValueDefinition()
		if err != nil {
			return nil, err
		}

		defs = defs.Add(def)

		if p.peek1(token.Punctuator, ")") || p.peek0(token.EOF) {
			break
		}
	}

	_, err := p.mustConsume1(token.Punctuator, ")")
	if err != nil {
		return nil, err
	}

	return defs.Reverse(), nil
}

// parseInputValueDefinition ...
func (p *Parser) parseInputValueDefinition() (ast.InputValueDefinition, error) {
	var description string

	descriptionTok, ok := p.consume0(token.StringValue)
	if ok {
		description = descriptionTok.Literal
	}

	nameTok, err := p.mustConsume0(token.Name)
	if err != nil {
		return ast.InputValueDefinition{}, err
	}

	if !p.skip1(token.Punctuator, ":") {
		return ast.InputValueDefinition{}, p.unexpected(p.token, p.expected(token.Punctuator, ":"))
	}

	inputValType, err := p.parseType()
	if err != nil {
		return ast.InputValueDefinition{}, err
	}

	defVal, err := p.parseDefaultValue()
	if err != nil {
		return ast.InputValueDefinition{}, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return ast.InputValueDefinition{}, err
	}

	def := ast.InputValueDefinition{
		Description:  description,
		Name:         nameTok.Literal,
		Type:         inputValType,
		DefaultValue: defVal,
		Directives:   directives,
	}

	return def, nil
}

// parseDirectiveDefinition ...
func (p *Parser) parseDirectiveDefinition(description string) (*ast.DirectiveDefinition, error) {
	if !p.skip1(token.Name, "directive") {
		return nil, nil
	}

	if !p.skip1(token.Punctuator, "@") {
		return nil, p.unexpected(p.token, p.expected(token.Punctuator, "@"))
	}

	nameTok, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	arguments, err := p.parseArgumentsDefinition()
	if err != nil {
		return nil, err
	}

	if !p.skip1(token.Name, "on") {
		return nil, p.unexpected(p.token, p.expected(token.Name, "on"))
	}

	locations, err := p.parseDirectiveLocations()
	if err != nil {
		return nil, err
	}

	return &ast.DirectiveDefinition{
		Description:         description,
		Name:                nameTok.Literal,
		DirectiveLocations:  locations,
		ArgumentsDefinition: arguments,
	}, nil
}

// parseTypeDefinition ...
func (p *Parser) parseTypeDefinition() (*ast.TypeDefinition, error) {
	defType, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	switch defType.Literal {
	case "scalar":
		return p.parseScalarTypeDefinition()
	case "type":
		return p.parseObjectDefinition()
	case "interface":
		return p.parseInterfaceDefinition()
	case "union":
		return p.parseUnionDefinition()
	case "enum":
		return p.parseEnumDefinition()
	case "input":
		return p.parseInputObjectDefinition()
	}

	return nil, nil
}

// parseScalarTypeDefinition ...
func (p *Parser) parseScalarTypeDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:       ast.TypeDefinitionKindScalar,
		Name:       name.Literal,
		Directives: directives,
	}, nil
}

// parseObjectDefinition ...
func (p *Parser) parseObjectDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	implements, err := p.parseImplementsInterfaces()
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	fieldDefs, err := p.parseFieldsDefinition()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:                ast.TypeDefinitionKindObject,
		Name:                name.Literal,
		ImplementsInterface: implements,
		Directives:          directives,
		FieldsDefinition:    fieldDefs,
	}, nil
}

// parseInterfaceDefinition ...
func (p *Parser) parseInterfaceDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	fieldDefs, err := p.parseFieldsDefinition()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:             ast.TypeDefinitionKindInterface,
		Name:             name.Literal,
		Directives:       directives,
		FieldsDefinition: fieldDefs,
	}, nil
}

// parseUnionDefinition ...
func (p *Parser) parseUnionDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	memberTypes, err := p.parseUnionMemberTypes()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:             ast.TypeDefinitionKindUnion,
		Name:             name.Literal,
		Directives:       directives,
		UnionMemberTypes: memberTypes,
	}, nil
}

// parseEnumDefinition ...
func (p *Parser) parseEnumDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	enumValues, err := p.parseEnumValuesDefinition()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:                 ast.TypeDefinitionKindEnum,
		Name:                 name.Literal,
		Directives:           directives,
		EnumValuesDefinition: enumValues,
	}, nil
}

// parseInputObjectDefinition ...
func (p *Parser) parseInputObjectDefinition() (*ast.TypeDefinition, error) {
	name, err := p.mustConsume0(token.Name)
	if err != nil {
		return nil, err
	}

	directives, err := p.parseDirectives()
	if err != nil {
		return nil, err
	}

	inputFields, err := p.parseInputFieldsDefinition()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDefinition{
		Kind:                  ast.TypeDefinitionKindInputObject,
		Name:                  name.Literal,
		Directives:            directives,
		InputFieldsDefinition: inputFields,
	}, nil
}

// parseImplementsInterfaces ...
func (p *Parser) parseImplementsInterfaces() (*ast.Types, error) {
	if !p.skip1(token.Name, "implements") {
		return nil, nil
	}

	p.skip1(token.Punctuator, "&")

	var interfaceTypes *ast.Types

	for {
		interfaceType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		interfaceTypes = interfaceTypes.Add(interfaceType)

		if !p.skip1(token.Punctuator, "&") {
			break
		}
	}

	return interfaceTypes.Reverse(), nil
}

// parseFieldsDefinition ...
func (p *Parser) parseFieldsDefinition() (*ast.FieldDefinitions, error) {
	if !p.skip1(token.Punctuator, "{") {
		return nil, nil
	}

	var fieldDefs *ast.FieldDefinitions

	for {
		var description string
		if tok, ok := p.consume0(token.StringValue); ok {
			description = tok.Literal
		}

		name, err := p.mustConsume0(token.Name)
		if err != nil {
			return nil, err
		}

		arguments, err := p.parseArgumentsDefinition()
		if err != nil {
			return nil, err
		}

		if !p.skip1(token.Punctuator, ":") {
			return nil, p.unexpected(p.token, p.expected(token.Punctuator, ":"))
		}

		fieldDefType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		directives, err := p.parseDirectives()
		if err != nil {
			return nil, err
		}

		fieldDef := ast.FieldDefinition{
			Description:         description,
			Name:                name.Literal,
			ArgumentsDefinition: arguments,
			Type:                fieldDefType,
			Directives:          directives,
		}

		fieldDefs = fieldDefs.Add(fieldDef)

		if p.skip1(token.Punctuator, "}") {
			break
		}
	}

	if fieldDefs != nil {
		return fieldDefs.Reverse(), nil
	}

	return nil, nil
}

// parseUnionMemberTypes ...
func (p *Parser) parseUnionMemberTypes() (*ast.Types, error) {
	if !p.skip1(token.Punctuator, "=") {
		return nil, nil
	}

	p.skip1(token.Punctuator, "|")

	var memberTypes *ast.Types

	for {
		memberType, err := p.parseType()
		if err != nil {
			return nil, err
		}

		memberTypes = memberTypes.Add(memberType)

		if !p.skip1(token.Punctuator, "|") {
			break
		}
	}

	return memberTypes.Reverse(), nil
}

// parseEnumValuesDefinition ...
func (p *Parser) parseEnumValuesDefinition() (*ast.EnumValueDefinitions, error) {
	if !p.skip1(token.Punctuator, "{") {
		return nil, nil
	}

	var valDefs *ast.EnumValueDefinitions

	for {
		var description string
		if tok, ok := p.consume0(token.StringValue); ok {
			description = tok.Literal
		}

		enumValue, err := p.mustConsume0(token.Name)
		if err != nil {
			return nil, err
		}

		directives, err := p.parseDirectives()
		if err != nil {
			return nil, err
		}

		valDef := ast.EnumValueDefinition{
			Description: description,
			EnumValue:   enumValue.Literal,
			Directives:  directives,
		}

		valDefs = valDefs.Add(valDef)

		if p.skip1(token.Punctuator, "}") {
			break
		}
	}

	return valDefs.Reverse(), nil
}

// parseInputFieldsDefinition ...
func (p *Parser) parseInputFieldsDefinition() (*ast.InputValueDefinitions, error) {
	if !p.skip1(token.Punctuator, "{") {
		return nil, nil
	}

	var valDefs *ast.InputValueDefinitions

	for {
		valDef, err := p.parseInputValueDefinition()
		if err != nil {
			return nil, err
		}

		valDefs = valDefs.Add(valDef)

		if p.skip1(token.Punctuator, "}") {
			break
		}
	}

	return valDefs.Reverse(), nil
}

var directiveLocations = []string{
	"QUERY",
	"MUTATION",
	"SUBSCRIPTION",
	"FIELD",
	"FRAGMENT_DEFINITION",
	"FRAGMENT_SPREAD",
	"INLINE_FRAGMENT",
	"SCHEMA",
	"SCALAR",
	"OBJECT",
	"FIELD_DEFINITION",
	"ARGUMENT_DEFINITION",
	"INTERFACE",
	"UNION",
	"ENUM",
	"ENUM_VALUE",
	"INPUT_OBJECT",
	"INPUT_FIELD_DEFINITION",
}

var directiveLocationsMap = map[string]ast.DirectiveLocation{
	"QUERY":                  ast.DirectiveLocationQuery,
	"MUTATION":               ast.DirectiveLocationMutation,
	"SUBSCRIPTION":           ast.DirectiveLocationSubscription,
	"FIELD":                  ast.DirectiveLocationField,
	"FRAGMENT_DEFINITION":    ast.DirectiveLocationFragmentDefinition,
	"FRAGMENT_SPREAD":        ast.DirectiveLocationFragmentSpread,
	"INLINE_FRAGMENT":        ast.DirectiveLocationInlineFragment,
	"SCHEMA":                 ast.DirectiveLocationSchema,
	"SCALAR":                 ast.DirectiveLocationScalar,
	"OBJECT":                 ast.DirectiveLocationObject,
	"FIELD_DEFINITION":       ast.DirectiveLocationFieldDefinition,
	"ARGUMENT_DEFINITION":    ast.DirectiveLocationArgumentDefinition,
	"INTERFACE":              ast.DirectiveLocationInterface,
	"UNION":                  ast.DirectiveLocationUnion,
	"ENUM":                   ast.DirectiveLocationEnum,
	"ENUM_VALUE":             ast.DirectiveLocationEnumValue,
	"INPUT_OBJECT":           ast.DirectiveLocationInputObject,
	"INPUT_FIELD_DEFINITION": ast.DirectiveLocationInputFieldDefinition,
}

// parseDirectiveLocations ...
func (p *Parser) parseDirectiveLocations() (*ast.DirectiveLocations, error) {
	p.skip1(token.Punctuator, "|") // this one is optional

	var locs *ast.DirectiveLocations

	for {
		tok, err := p.mustConsumen(token.Name, directiveLocations...)
		if err != nil {
			return nil, p.unexpected(p.token, p.expected(token.Name, directiveLocations...))
		}

		locs = locs.Add(directiveLocationsMap[tok.Literal])

		if !p.skip1(token.Punctuator, "|") {
			break
		}
	}

	return locs.Reverse(), nil
}

/*****************************************************************************
* Utility functions                                                          *
*****************************************************************************/

// consume0 ...
func (p *Parser) consume0(t token.Type) (lexer.Token, bool) {
	tok := p.token
	ok := p.token.Type == t

	if ok {
		p.token = p.lexer.Scan()
	}

	return tok, ok
}

// consume1 ...
func (p *Parser) consume1(t token.Type, l string) (lexer.Token, bool) {
	tok := p.token
	ok := p.token.Type == t && p.token.Literal == l

	if ok {
		p.token = p.lexer.Scan()
	}

	return tok, ok
}

// consumen ...
func (p *Parser) consumen(t token.Type, ls ...string) (lexer.Token, bool) {
	tok := p.token
	if tok.Type != t {
		return tok, false
	}

	if len(ls) == 0 {
		p.token = p.lexer.Scan()
		return tok, true
	}

	for _, l := range ls {
		if tok.Literal != l {
			continue
		}

		p.token = p.lexer.Scan()
		return tok, true
	}

	return tok, false
}

// mustConsume0 ...
func (p *Parser) mustConsume0(t token.Type) (lexer.Token, error) {
	tok := p.token

	if p.token.Type != t {
		return tok, p.unexpected(tok, p.expected(t))
	}

	p.token = p.lexer.Scan()

	return tok, nil
}

// mustConsume1 ...
func (p *Parser) mustConsume1(t token.Type, l string) (lexer.Token, error) {
	tok := p.token

	if p.token.Type != t || p.token.Literal != l {
		return tok, p.unexpected(tok, p.expected(t, l))
	}

	p.token = p.lexer.Scan()

	return tok, nil
}

// mustConsumen ...
func (p *Parser) mustConsumen(t token.Type, ls ...string) (lexer.Token, error) {
	tok, ok := p.consumen(t, ls...)
	if !ok {
		return tok, p.unexpected(tok, p.expected(t, ls...))
	}

	return tok, nil
}

// peek0 ...
func (p *Parser) peek0(t token.Type) bool {
	return p.token.Type == t
}

// peek1 ...
func (p *Parser) peek1(t token.Type, l string) bool {
	return p.token.Type == t && p.token.Literal == l
}

// peekn ...
func (p *Parser) peekn(t token.Type, ls ...string) bool {
	if p.token.Type != t {
		return false
	}

	if len(ls) == 0 {
		return true
	}

	for _, l := range ls {
		if p.token.Literal == l {
			return true
		}
	}

	return false
}

// skip0 ...
func (p *Parser) skip0(t token.Type) bool {
	if p.token.Type != t {
		return false
	}

	p.token = p.lexer.Scan()

	return true
}

// skip1 ...
func (p *Parser) skip1(t token.Type, l string) bool {
	if p.token.Type != t || p.token.Literal != l {
		return false
	}

	p.token = p.lexer.Scan()

	return true
}

// skip ...
func (p *Parser) skip(t token.Type, ls ...string) bool {
	match := p.peekn(t, ls...)
	if !match {
		return false
	}

	p.token = p.lexer.Scan()

	return true
}

// scan ...
func (p *Parser) scan() {
	p.token = p.lexer.Scan()
}

// expected ...
func (p *Parser) expected(t token.Type, ls ...string) string {
	buf := &bytes.Buffer{}
	buf.WriteString(t.String())
	buf.WriteString(" '")
	buf.WriteString(strings.Join(ls, "|"))
	buf.WriteString("'")
	return btos(buf.Bytes())
}

// unexpected ...
func (p *Parser) unexpected(token lexer.Token, wants ...string) error {
	//_, file, line, _ := runtime.Caller(2)
	//fmt.Println(file, line)

	buf := &bytes.Buffer{}
	buf.WriteString("parser error: unexpected token found at ")
	buf.WriteString("line: ")
	buf.WriteString(strconv.Itoa(token.Line))
	buf.WriteString(", column: ")
	buf.WriteString(strconv.Itoa(token.Position))
	buf.WriteString(". Found: ")
	buf.WriteString(token.Type.String())
	buf.WriteString(" '")
	buf.WriteString(token.Literal)
	buf.WriteString("'.")

	if len(wants) > 0 {
		buf.WriteString(" Wanted: ")
		for i, want := range wants {
			buf.WriteString(want)
			if i < len(wants)-1 {
				buf.WriteString(" or ")
			}
		}
		buf.WriteString(".")
	}

	return errors.New(btos(buf.Bytes()))
}

// btos takes the given bytes, and turns them into a string.
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
