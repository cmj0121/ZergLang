package parser

import (
	"testing"

	"github.com/xrspace/zerglang/runtime/lexer"
)

func TestDeclarationStatement(t *testing.T) {
	input := `x := 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*DeclarationStatement)
	if !ok {
		t.Fatalf("expected DeclarationStatement, got %T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Fatalf("expected x, got %s", stmt.Name.Value)
	}

	lit, ok := stmt.Value.(*IntegerLiteral)
	if !ok {
		t.Fatalf("expected IntegerLiteral, got %T", stmt.Value)
	}

	if lit.Value != 42 {
		t.Fatalf("expected 42, got %d", lit.Value)
	}
}

func TestMultipleDeclarations(t *testing.T) {
	input := `x := 1
y := 2
z := 3`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}
}

func TestLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"42", int64(42)},
		{`"hello"`, "hello"},
		{"true", true},
		{"false", false},
		{"nil", nil},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
		}

		switch expected := tt.expected.(type) {
		case int64:
			lit, ok := stmt.Expression.(*IntegerLiteral)
			if !ok {
				t.Fatalf("expected IntegerLiteral, got %T", stmt.Expression)
			}
			if lit.Value != expected {
				t.Fatalf("expected %d, got %d", expected, lit.Value)
			}
		case string:
			lit, ok := stmt.Expression.(*StringLiteral)
			if !ok {
				t.Fatalf("expected StringLiteral, got %T", stmt.Expression)
			}
			if lit.Value != expected {
				t.Fatalf("expected %s, got %s", expected, lit.Value)
			}
		case bool:
			lit, ok := stmt.Expression.(*BooleanLiteral)
			if !ok {
				t.Fatalf("expected BooleanLiteral, got %T", stmt.Expression)
			}
			if lit.Value != expected {
				t.Fatalf("expected %t, got %t", expected, lit.Value)
			}
		case nil:
			_, ok := stmt.Expression.(*NilLiteral)
			if !ok {
				t.Fatalf("expected NilLiteral, got %T", stmt.Expression)
			}
		}
	}
}

func TestIdentifier(t *testing.T) {
	input := `myVariable`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ExpressionStatement)
	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Fatalf("expected Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "myVariable" {
		t.Fatalf("expected myVariable, got %s", ident.Value)
	}
}

func TestMutableDeclaration(t *testing.T) {
	input := `mut x := 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*DeclarationStatement)
	if !ok {
		t.Fatalf("expected DeclarationStatement, got %T", program.Statements[0])
	}

	if !stmt.Mutable {
		t.Fatal("expected mutable declaration")
	}

	if stmt.Name.Value != "x" {
		t.Fatalf("expected x, got %s", stmt.Name.Value)
	}
}

func TestAssignmentStatement(t *testing.T) {
	input := `x = 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("expected AssignmentStatement, got %T", program.Statements[0])
	}

	if len(stmt.Names) != 1 {
		t.Fatalf("expected 1 name, got %d", len(stmt.Names))
	}

	if stmt.Names[0].Value != "x" {
		t.Fatalf("expected x, got %s", stmt.Names[0].Value)
	}
}

func TestMultiValueAssignment(t *testing.T) {
	input := `a, b = 1, 2`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("expected AssignmentStatement, got %T", program.Statements[0])
	}

	if len(stmt.Names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(stmt.Names))
	}

	if len(stmt.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(stmt.Values))
	}

	if stmt.Names[0].Value != "a" || stmt.Names[1].Value != "b" {
		t.Fatalf("expected a, b got %s, %s", stmt.Names[0].Value, stmt.Names[1].Value)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
