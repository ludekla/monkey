package lexer

import "monkey/pkg/token"

type Lexer struct {
	input    string
	position int  // position in input string
	readPos  int  // current reading position, one after ch
	ch       byte // current char under examination
}

// isLetter is a helper function to check for digits.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit is a helper function to check for digits.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// New is the Lexer factory.
func New(input string) *Lexer {
	lx := &Lexer{input: input}
	lx.readChar()
	return lx
}

// readChar advanced the reader pointer to the next char
// within the input string.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos++
}

// Next inspects the current byte and returns its
// correspnding token.
func (l *Lexer) Next() token.Token {
	var tok token.Token
	// skip whitepspace
	l.skipWhiteSpace()
	// inspect
	switch l.ch {
	case '=':
		tok = token.New(token.ASSIGN, l.ch)
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// readIdent reads a identifier.
func (l *Lexer) readIdent() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

// skipWhiteSpace moves the cursor to the end
// of the whitespace part.
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads as many digits as possible.
func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}
