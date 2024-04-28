package parser

import (
	"monkey/pkg/ast"
	"monkey/pkg/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;

	let y = 10;
	let foobar = 838383;`
	lx := lexer.New(input)
	ps := New(lx)

	prog := ps.ParseProgramme()
	if prog == nil {
		t.Fatalf("ParseProgramm() returned nil")
	}
	if len(prog.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(prog.Statements))
	}
	tests := []string{
		"x",
		"y",
		"foobar",
	}
	for i, tt := range tests {
		stm := prog.Statements[i]
		testLetStatement(t, stm, tt)
	}
}

// Helper function.
func testLetStatement(t *testing.T, s ast.Statement, name string) {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral not 'let'. Got %q", s.TokenLiteral())
	}
	letSmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. Got %T", s)
	}
	if letSmt.Name.Value != name {
		t.Fatalf("LetStatement.Name.Value not %q, got %q", name, letSmt.Name.Value)
	}
}
