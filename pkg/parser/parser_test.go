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

func TestIntegerExpression(t *testing.T) {
	input := "5;"
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
	literal, ok := stm.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected *ast.IntegerLiteral, got %T", stm.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("ident.Value not 5, got %v", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("token literal not '5', got %s", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}
	for _, tt := range prefixTests {
		lx := lexer.New(tt.input)
		ps := New(lx)
		prog := ps.ParseProgramme()
		checkParserErrors(t, ps)
		if len(prog.Statements) != 1 {
			t.Fatalf(
				"Programme does not have 1 statement, got %d.",
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
		expr, ok := stm.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected *ast.PrefixExpression, got %T", stm.Expression)
		}
		if expr.Operator != tt.operator {
			t.Fatalf("Operator is not %q, got %s", tt.operator, expr.Operator)
		}
		if !testIntegerLiteral(t, expr.Right, tt.intValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got %T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d, got %d", value, integ.Value)
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tt := range infixTests {
		lx := lexer.New(tt.input)
		ps := New(lx)
		prog := ps.ParseProgramme()
		checkParserErrors(t, ps)
		if len(prog.Statements) != 1 {
			t.Fatalf(
				"Progamme.Statements does not contain 1 statement, got %d",
				len(prog.Statements),
			)
		}
		stm, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"Statement is not an ExpressionStatement, got %T",
				prog.Statements[0],
			)
		}
		expr, ok := stm.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"Statement is not an InfixExpression, got %T",
				stm.Expression,
			)
		}
		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {
			return
		}
		if expr.Operator != tt.operator {
			t.Fatalf("operator is not %q, got %s", tt.operator, expr.Operator)
		}
		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		lx := lexer.New(tt.input)
		ps := New(lx)
		prog := ps.ParseProgramme()
		checkParserErrors(t, ps)
		actual := prog.String()
		if actual != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, actual)
		}
	}
}
