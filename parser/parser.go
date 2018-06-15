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

type Parser struct {
	lexer *lexer.Lexer
	token lexer.Token
}

func New(input []byte) *Parser {
	return &Parser{
		lexer: lexer.New(input),
	}
}

func (p *Parser) Parse() (ast.Document, error) {
	var document ast.Document

	p.scan()

	for {
		definition, err := p.parseDefinition()
		if err != nil {
			return document, err
		}

		document.Definitions = append(document.Definitions, definition)

		if p.peek(token.Illegal) {
			return document, p.unexpected(p.token, token.EOF)
		}

		if p.peek(token.EOF) {
			return document, nil
		}
	}
}

func (p *Parser) parseDefinition() (ast.Definition, error) {
	var definition ast.Definition
	var err error

	// ExecutableDefinition...
	if p.peek(token.Name, "query", "mutation", "subscription") || p.peek(token.Punctuator, "{") {
		definition.ExecutableDefinition, err = p.parseOperationDefinition(p.token.Literal == "{")
		return definition, err
	}

	// ExecutableDefinition...
	if p.peek(token.Name, "fragment") {
		// TODO(seeruk): Implement.
	}

	// TODO(seeruk): We need unexpected to support multiple token types and literals. Maybe:
	// func unexpected(tok lexer.Token, wants ...func(t token.Type, ls ...string) []lexer.Token) {}
	//
	// So in the case below, we need to support what is there now, and also a punctuator with the
	// literal "{". The current error message doesn't highlight all of the things we might want.
	return definition, p.unexpected(p.token, token.Name, "query", "mutation", "fragment")
}

func (p *Parser) parseOperationDefinition(isQuery bool) (ast.ExecutableDefinition, error) {
	var definition ast.ExecutableDefinition

	var name string
	var err error

	opType := ast.OperationTypeQuery

	if !isQuery {
		opType, err = p.parseOperationType()
		if err != nil {
			return definition, err
		}

		if tok, ok := p.consume(token.Name); ok {
			name = tok.Literal
		}
	}

	variableDefinitions, err := p.parseVariableDefinitions()
	if err != nil {
		return definition, err
	}

	if _, err = p.mustConsume(token.Punctuator, "{"); err != nil {
		return definition, err
	}

	// TODO(seeruk): parseSelectionSet.

	if _, err = p.mustConsume(token.Punctuator, "}"); err != nil {
		return definition, err
	}

	return ast.ExecutableDefinition{
		Kind:                ast.DefinitionKindOperation,
		OperationType:       opType,
		Name:                name,
		VariableDefinitions: variableDefinitions,
		// Directives: ...
		// SelectionSet: ...
		// TODO(seeruk): Location.
	}, nil
}

func (p *Parser) parseOperationType() (ast.OperationType, error) {
	tok, err := p.mustConsume(token.Name, "query", "mutation")
	if err != nil {
		return -1, err
	}

	if tok.Literal == "query" {
		return ast.OperationTypeQuery, nil
	}

	// Only other thing it can be at this point...
	return ast.OperationTypeMutation, nil
}

func (p *Parser) parseVariableDefinitions() ([]ast.VariableDefinition, error) {
	var definitions []ast.VariableDefinition

	if _, err := p.mustConsume(token.Punctuator, "("); err != nil {
		return definitions, err
	}

	for {
		if _, err := p.mustConsume(token.Punctuator, "$"); err != nil {
			return definitions, err
		}

		tok, err := p.mustConsume(token.Name)
		if err != nil {
			return definitions, err
		}

		definition := ast.VariableDefinition{}
		definition.Name = tok.Literal

		if _, err := p.mustConsume(token.Punctuator, ":"); err != nil {
			return definitions, err
		}

		definition.Type, err = p.parseType()
		if err != nil {
			return definitions, err
		}

		definition.DefaultValue, err = p.parseDefaultValue()
		if err != nil {
			return definitions, err
		}

		definitions = append(definitions, definition)

		if p.peek(token.Punctuator, ")") {
			break
		}
	}

	if _, err := p.mustConsume(token.Punctuator, ")"); err != nil {
		return definitions, err
	}

	return definitions, nil
}

func (p *Parser) parseDefaultValue() (*ast.Value, error) {
	if !p.skip(token.Punctuator, "=") {
		return nil, nil
	}

	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (p *Parser) parseValue() (ast.Value, error) {
	if p.skip(token.Punctuator, "$") {
		return p.parseVariableValue()
	}

	if tok, ok := p.consume(token.IntValue); ok {
		return p.parseIntValue(tok)
	}

	if tok, ok := p.consume(token.FloatValue); ok {
		return p.parseFloatValue(tok)
	}

	if tok, ok := p.consume(token.StringValue); ok {
		return p.parseStringValue(tok)
	}

	if tok, ok := p.consume(token.Name, "true", "false"); ok {
		return p.parseBooleanValue(tok)
	}

	if p.skip(token.Name, "null") {
		return p.parseNullValue()
	}

	if tok, ok := p.consume(token.Name); ok {
		return p.parseEnumValue(tok)
	}

	if p.skip(token.Punctuator, "[") {
		return p.parseListValue()
	}

	if p.skip(token.Punctuator, "{") {
		return p.parseObjectValue()
	}

	return ast.Value{}, errors.New("TODO: see `parseDefinition`")
}

func (p *Parser) parseVariableValue() (ast.Value, error) {
	tok, err := p.mustConsume(token.Name)
	if err != nil {
		return ast.Value{}, err
	}

	return ast.Value{
		Kind:          ast.ValueKindVariable,
		VariableValue: tok.Literal,
	}, nil
}

func (p *Parser) parseIntValue(tok lexer.Token) (ast.Value, error) {
	iv, err := strconv.Atoi(tok.Literal)
	if err != nil {
		return ast.Value{}, err
	}

	return ast.Value{
		Kind:     ast.ValueKindIntValue,
		IntValue: iv,
	}, nil
}

func (p *Parser) parseFloatValue(tok lexer.Token) (ast.Value, error) {
	fv, err := strconv.ParseFloat(tok.Literal, 64)
	if err != nil {
		return ast.Value{}, err
	}

	return ast.Value{
		Kind:       ast.ValueKindFloatValue,
		FloatValue: fv,
	}, nil
}

func (p *Parser) parseStringValue(tok lexer.Token) (ast.Value, error) {
	return ast.Value{
		Kind:        ast.ValueKindStringValue,
		StringValue: tok.Literal,
	}, nil
}

func (p *Parser) parseBooleanValue(tok lexer.Token) (ast.Value, error) {
	return ast.Value{
		Kind:         ast.ValueKindBooleanValue,
		BooleanValue: tok.Literal == "true",
	}, nil
}

func (p *Parser) parseNullValue() (ast.Value, error) {
	return ast.Value{
		Kind: ast.ValueKindNullValue,
	}, nil
}

func (p *Parser) parseEnumValue(tok lexer.Token) (ast.Value, error) {
	return ast.Value{
		Kind:      ast.ValueKindEnumValue,
		EnumValue: tok.Literal,
	}, nil
}

func (p *Parser) parseListValue() (ast.Value, error) {
	list := ast.Value{}
	list.Kind = ast.ValueKindListValue

	for !p.skip(token.Punctuator, "]") {
		val, err := p.parseValue()
		if err != nil {
			return list, err
		}

		list.ListValue = append(list.ListValue, val)
	}

	return list, nil
}

func (p *Parser) parseObjectValue() (ast.Value, error) {
	object := ast.Value{}
	object.Kind = ast.ValueKindObjectValue

	for !p.skip(token.Punctuator, "}") {
		field, err := p.parseObjectField()
		if err != nil {
			return object, err
		}

		object.ObjectValue = append(object.ObjectValue, field)
	}

	return object, nil
}

func (p *Parser) parseObjectField() (ast.ObjectField, error) {
	var field ast.ObjectField

	tok, err := p.mustConsume(token.Name)
	if err != nil {
		return field, err
	}

	_, err = p.mustConsume(token.Punctuator, ":")
	if err != nil {
		return field, err
	}

	value, err := p.parseValue()
	if err != nil {
		return field, err
	}

	field.Name = tok.Literal
	field.Value = value

	return field, nil
}

func (p *Parser) parseType() (ast.Type, error) {
	var astType ast.Type

	// If we hit an opening square brace, we've got a list type, time to dive in.
	if p.skip(token.Punctuator, "[") {
		astType.Kind = ast.TypeKindListType

		itemType, err := p.parseType()
		if err != nil {
			return astType, nil
		}

		astType.ListType = &itemType

		if _, err := p.mustConsume(token.Punctuator, "]"); err != nil {
			return astType, err
		}
	} else {
		astType.Kind = ast.TypeKindNamedType

		tok, err := p.mustConsume(token.Name)
		if err != nil {
			return astType, err
		}

		astType.NamedType = tok.Literal
	}

	if p.skip(token.Punctuator, "!") {
		astType.NonNullable = true
	}

	return astType, nil
}

// Parser utilities:

func (p *Parser) consume(t token.Type, ls ...string) (lexer.Token, bool) {
	tok := p.token
	if tok.Type != t {
		return tok, false
	}

	if len(ls) == 0 {
		p.scan()
		return tok, true
	}

	for _, l := range ls {
		if tok.Literal != l {
			continue
		}

		p.scan()
		return tok, true
	}

	return tok, false
}

func (p *Parser) mustConsume(t token.Type, ls ...string) (lexer.Token, error) {
	tok, ok := p.consume(t, ls...)
	if !ok {
		return tok, p.unexpected(tok, t, ls...)
	}

	return tok, nil
}

func (p *Parser) peek(t token.Type, ls ...string) bool {
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

func (p *Parser) skip(t token.Type, ls ...string) bool {
	match := p.peek(t, ls...)
	if !match {
		return false
	}

	p.scan()

	return true
}

func (p *Parser) scan() {
	p.token = p.lexer.Scan()
}

func (p *Parser) unexpected(token lexer.Token, t token.Type, ls ...string) error {
	if len(ls) == 0 {
		ls = []string{"N/A"}
	}

	// This is as nasty as I'm willing to make this right now. But this is the slowest function in
	// the parser by far, because of the allocations it has to do, simply because it's generating
	// this message.
	// TODO(seeruk): Revisit this, it can almost definitely be improved.
	// TODO(seeruk): Don't call unexpected when it's not absolutely necessary. We can not pass
	// around errors if we don't need to (i.e. if we want to mustConsume without caring about the error,
	// like if we just care about whether or not we did mustConsume something).
	buf := bytes.Buffer{}
	buf.WriteString("parser error: unexpected token found: ")
	buf.WriteString(token.Type.String())
	buf.WriteString(" '")
	buf.WriteString(token.Literal)
	buf.WriteString("'. Wanted: ")
	buf.WriteString(t.String())
	buf.WriteString(" '")
	buf.WriteString(strings.Join(ls, "|"))
	buf.WriteString("'. Line: ")
	buf.WriteString(strconv.Itoa(token.Line))
	buf.WriteString(". Column: ")
	buf.WriteString(strconv.Itoa(token.Position))

	return errors.New(btos(buf.Bytes()))
}

// btos takes the given bytes, and turns them into a string.
// Q: naming btos or bbtos? :D
// TODO(seeruk): Is this actually portable then?
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
