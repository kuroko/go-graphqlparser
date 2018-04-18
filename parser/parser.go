package parser

import (
	"errors"

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

	// Read first token in before any expectations are set.
	p.scan()

	for {
		definition, err := p.parseDefinition()
		if err != nil {
			return document, err
		}

		document.Definitions = append(document.Definitions, definition)

		done, err := p.skip(token.EOF)
		if err != nil {
			return document, err
		}

		if done {
			break
		}
	}

	return document, nil
}

func (p *Parser) parseDefinition() (ast.Definition, error) {
	// TODO(seeruk): parseTypeSystemDefinition things.
	// TODO(seeruk): Maybe expand this to have parseExecutableDefinition?

	var definition ast.Definition

	if p.assert(token.Name) {
		isQuery := p.token.Literal == "query"
		isMutation := p.token.Literal == "mutation"
		if isQuery || isMutation {
			return p.parseOperationDefinition()
		}

		if p.token.Literal == "fragment" {
			// TODO(seeruk): siufhisuhdf
		}
	}

	if p.assert(token.Punctuator) && p.token.Literal == "{" {
		return ast.Definition{
			Kind:          ast.DefinitionKindOperation,
			OperationType: ast.OperationTypeQuery,
			// SelectionSet: ...
			// TODO(seeruk): Location.
		}, nil
	}

	return definition, errors.New("unexpected: todo")
}

func (p *Parser) parseOperationDefinition() (ast.Definition, error) {
	var definition ast.Definition

	opType, err := p.parseOperationType()
	if err != nil {
		return definition, err
	}

	var name string
	if p.assert(token.Name) {
		name = p.token.Literal
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
	token, err := p.expect(token.Name)
	if err != nil {
		// TODO(seeruk): More error context?
		return -1, err
	}

	switch token.Literal {
	case "query":
		return ast.OperationTypeQuery, nil
	case "mutation":
		return ast.OperationTypeMutation, nil
	}

	return -1, errors.New("unexpected: todo")
}

func (p *Parser) assert(t token.Type) bool {
	return p.token.Type == t
}

func (p *Parser) expect(t token.Type) (lexer.Token, error) {
	token := p.token
	match, err := p.skip(t)
	if err != nil {
		return token, err
	}

	if match {
		return token, nil
	}

	return token, errors.New("syntax error: todo")
}

func (p *Parser) scan() (err error) {
	p.token, err = p.lexer.Scan()
	return err
}

func (p *Parser) skip(t token.Type) (bool, error) {
	var err error
	match := p.token.Type == t
	if match {
		p.token, err = p.lexer.Scan()
	}

	return match, err
}
