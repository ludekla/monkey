package parser

import (
	"monkey/pkg/ast"
	"monkey/pkg/lexer"
	"monkey/pkg/token"
)

type Parser struct {
	lx        *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

func New(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx}
	// Read two tokens so that current and peek
	// tokens are initialised.
	p.Next()
	p.Next()
	return p
}

func (p *Parser) Next() {
	p.currToken = p.peekToken
	p.peekToken = p.lx.Next()
}

func (p *Parser) ParseProgramme() *ast.Programme {
	prog := &ast.Programme{}
	prog.Statements = make([]ast.Statement, 0, 10)
	for p.currToken.Type != token.EOF {
		stm := p.ParseStatement()
		if stm != nil {
			prog.Statements = append(prog.Statements, stm)
		}
		p.Next()
	}
	return prog
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stm := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stm.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO
	for p.currToken.Type != token.SEMICOLON {
		p.Next()
	}
	return stm
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.Next()
		return true
	}
	return false
}
