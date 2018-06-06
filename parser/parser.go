package parser

import (
	"fmt"
	"strings"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/davecgh/go-spew/spew"
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

	// Read first token in before any expectations are set.
	p.scan()

	for {
		definition, err := p.parseDefinition()
		if err != nil {
			return document, err
		}

		document.Definitions = append(document.Definitions, definition)

		if p.next(token.Illegal) {
			return document, p.unexpected(p.token, token.EOF)
		}

		if p.next(token.EOF) {
			return document, nil
		}
	}
}

func (p *Parser) parseDefinition() (ast.Definition, error) {
	// TODO(seeruk): parseTypeSystemDefinition things.
	// TODO(seeruk): Maybe expand this to have parseExecutableDefinition?

	var definition ast.Definition

	if p.next(token.Name, "query", "mutation") || p.next(token.Punctuator, "{") {
		spew.Dump(p.token)
		return p.parseOperationDefinition(p.token.Literal == "{")
	}

	if p.next(token.Name, "fragment") {
		// TODO(seeruk): Implement.
	}

	return definition, p.unexpected(p.token, token.Name, "query", "mutation", "fragment")
}

func (p *Parser) parseOperationDefinition(isQuery bool) (ast.Definition, error) {
	var definition ast.Definition

	var opType ast.OperationType
	var name string
	var err error

	if !isQuery {
		opType, err = p.parseOperationType()
		if err != nil {
			return definition, err
		}

		if tok, err := p.consume(token.Name); err == nil {
			name = tok.Literal
		}
	}

	if _, err = p.consume(token.Punctuator, "{"); err != nil {
		return definition, err
	}

	return ast.Definition{
		Kind:          ast.DefinitionKindOperation,
		OperationType: opType,
		Name:          name,
		// VariableDefinitions: ...
		// Directives: ...
		// SelectionSet: ...
		// TODO(seeruk): Location.
	}, nil
}

func (p *Parser) parseOperationType() (ast.OperationType, error) {
	tok, err := p.consume(token.Name, "query", "mutation")
	if err != nil {
		return -1, err
	}

	if tok.Literal == "query" {
		return ast.OperationTypeQuery, nil
	}

	// Only other thing it can be at this point...
	return ast.OperationTypeMutation, nil
}

// Parser utilities:

func (p *Parser) expectAll(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) consume(t token.Type, ls ...string) (lexer.Token, error) {
	tok := p.token
	if tok.Type != t {
		return tok, p.unexpected(tok, t, ls...)
	}

	if len(ls) == 0 {
		p.scan()
		return tok, nil
	}

	for _, l := range ls {
		if tok.Literal != l {
			continue
		}

		p.scan()
		return tok, nil
	}

	return tok, p.unexpected(tok, t, ls...)
}

func (p *Parser) expect(t token.Type, ls ...string) error {
	if !p.next(t, ls...) {
		return p.unexpected(p.token, t, ls...)
	}

	return nil
}

func (p *Parser) expectFn(t token.Type, l string) func() error {
	return func() error {
		return p.expect(t, l)
	}
}

func (p *Parser) next(t token.Type, ls ...string) bool {
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
	_, err := p.consume(t, ls...)
	if err != nil {
		return false
	}

	return true
}

func (p *Parser) scan() {
	p.token = p.lexer.Scan()
}

func (p *Parser) unexpected(token lexer.Token, t token.Type, ls ...string) error {
	if len(ls) == 0 {
		ls = []string{"N/A"}
	}

	return fmt.Errorf(
		"parser error: unexpected token found: %s (%q). Wanted: %s (%q). Line: %d. Column: %d",
		token.Type.String(),
		token.Literal,
		t.String(),
		strings.Join(ls, "|"),
		token.Line,
		token.Position,
	)
}
