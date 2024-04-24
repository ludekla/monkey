package lexer

import (
	"monkey/pkg/token"
	"testing"
)

func TestNext(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.Next()
		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - token type wrong, expected = %q, got = %q",
				i, tt.expectedType, tok.Type,
			)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - token literal wrong, expected = %q, got = %q",
				i, tt.expectedLiteral, tok.Literal,
			)
		}
	}
}
