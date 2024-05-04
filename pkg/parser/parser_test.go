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
	checkParserErrors(t, ps)

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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	
	return 5;
	return 10;
	
	return 993322;
	
	`
	lx := lexer.New(input)
	ps := New(lx)
	prog := ps.ParseProgramme()
	checkParserErrors(t, ps)

	if len(prog.Statements) != 3 {
		t.Fatalf("Expected 3 RETURN statements, got %d", len(prog.Statements))
	}
	for _, stm := range prog.Statements {
		rstm, ok := stm.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got %T", stm)
			continue
		}
		if rstm.TokenLiteral() != "return" {
			t.Errorf(
				"return statement token literal not 'return', got %q",
				rstm.TokenLiteral(),
			)
		}
	}
}

func TestIdentExpression(t *testing.T) {
	input := "foobar;"
	lx := lexer.New(input)
	ps := New(lx)
	prog := ps.ParseProgramme()
	checkParserErrors(t, ps)
	if len(prog.Statements) != 1 {
		t.Fatalf(
			"Programme has not enough statements, got %d.",
			len(prog.Statements),
		)
	}
	stm, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statement is not ast.ExpressionStatement, got %T",
			prog.Statements[0],
		)
	}
	ident, ok := stm.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected *ast.Identifier, got %T", stm.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar, got %s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("token literal not 'foobar', got %s", ident.TokenLiteral())
	}
}
