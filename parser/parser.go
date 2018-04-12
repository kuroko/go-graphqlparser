package parser

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
}

func New(input []byte) *Parser {
	return &Parser{
		lexer: lexer.New(input),
	}
}

func (p *Parser) Parse() ast.Document {
	return ast.Document{}
}
