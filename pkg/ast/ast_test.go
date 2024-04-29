package ast

import (
	"monkey/pkg/token"
	"testing"
)

func TestString(t *testing.T) {
	prog := &Programme{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "yourVar"},
					Value: "yourVar",
				},
			},
		},
	}
	if prog.String() != "let myVar = yourVar;" {
		t.Errorf("Programme.String() wrong, got %q", prog.String())
	}
}
