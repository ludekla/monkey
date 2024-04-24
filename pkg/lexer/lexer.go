package lexer

import "monkey/pkg/token"

type Lexer struct {
	input    string
	position int  // position in input string
	readPos  int  // current reading position, one after ch
	ch       byte // current char under examination
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
		l.ch = 0 // go to the beginning
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
	}
	l.readChar()
	return tok
}
