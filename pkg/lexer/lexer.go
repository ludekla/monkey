package lexer

import "monkey/pkg/token"

type Lexer struct {
	input    string
	ch       byte // Current char under examination.
	position int  // Position of ch in the input string.
	readPos  int  // Current reading position, one after ch.
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

// readChar fills the char field with the byte the
// reading cursor currently points to and then advances
// the cursor to the next char within the input string.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos++
}

// Next inspects the current byte and returns the
// corresponding token if it is a one-character token
// or advances down to the end of the sequence of
// non-whitespace characters and then returns the token
// if it is a legal one.
func (l *Lexer) Next() token.Token {
	var tok token.Token
	// skip whitepspace
	l.skipWhiteSpace()
	// inspect
	switch l.ch {
	case '=':
		if l.peekAhead() == '=' {
			l.readChar()
			return token.Token{Type: token.EQ, Literal: "=="}
		}
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
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case '<':
		tok = token.New(token.LT, l.ch)
	case '!':
		if l.peekAhead() == '=' {
			l.readChar()
			return token.Token{Type: token.NEQ, Literal: "!="}
		}
		tok = token.New(token.BANG, l.ch)
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

// peekAhead returns the next character without
// advancing the cursor for inspection purposes.
func (l *Lexer) peekAhead() byte {
	if l.readPos == len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}
