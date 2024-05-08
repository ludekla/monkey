package parser

import (
	"fmt"
	"monkey/pkg/ast"
	"monkey/pkg/lexer"
	"monkey/pkg/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunc(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lx        *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
	prefixers map[token.TokenType]prefixParserFn
	infixers  map[token.TokenType]infixParserFn
}

func New(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx, errors: []string{}}
	p.prefixers = make(map[token.TokenType]prefixParserFn)
	// register parser functions
	p.prefixers[token.IDENT] = p.parseIdentifier
	p.prefixers[token.INT] = p.parseIntegerLiteral
	p.prefixers[token.BANG] = p.parsePrefixExpression
	p.prefixers[token.MINUS] = p.parsePrefixExpression
	p.infixers = make(map[token.TokenType]infixParserFn)
	p.infixers[token.PLUS] = p.parseInfixExpression
	p.infixers[token.MINUS] = p.parseInfixExpression
	p.infixers[token.SLASH] = p.parseInfixExpression
	p.infixers[token.ASTERISK] = p.parseInfixExpression
	p.infixers[token.EQ] = p.parseInfixExpression
	p.infixers[token.NEQ] = p.parseInfixExpression
	p.infixers[token.LT] = p.parseInfixExpression
	p.infixers[token.GT] = p.parseInfixExpression
	// Read two tokens so that current and peek
	// tokens are initialised.
	p.Next()
	p.Next()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// Getter for error messages.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Next() {
	p.currToken = p.peekToken
	p.peekToken = p.lx.Next()
}

func (p *Parser) ParseProgramme() *ast.Programme {
	prog := &ast.Programme{}
	prog.Statements = make([]ast.Statement, 0, 10)
	for !p.currTokenIs(token.EOF) {
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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
	for !p.currTokenIs(token.SEMICOLON) {
		p.Next()
	}
	return stm
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.Next()
		return true
	}
	p.peekError(t)
	return false
}

// Helpers.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stm := &ast.ReturnStatement{Token: p.currToken}
	p.Next()
	// TODO
	for !p.currTokenIs(token.SEMICOLON) {
		p.Next()
	}
	return stm
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stm := &ast.ExpressionStatement{Token: p.currToken}
	stm.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.Next()
	}
	return stm
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as int", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixers[p.currToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf(
			"no prefix parse function for %s found",
			p.currToken.Type,
		)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpr := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixers[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}
		p.Next()
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.Next()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) peekPrecedence() int {
	if pd, ok := precedences[p.peekToken.Type]; ok {
		return pd
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if pd, ok := precedences[p.currToken.Type]; ok {
		return pd
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	prec := p.currPrecedence()
	p.Next()
	expr.Right = p.parseExpression(prec)
	return expr
}
