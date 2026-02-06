package evaluator

import (
	"testing"

	"github.com/xrspace/zerglang/runtime/lexer"
	"github.com/xrspace/zerglang/runtime/parser"
)

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"42", 42},
		{"1_000", 1000},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	input := `"hello world"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "hello world" {
		t.Fatalf("expected hello world, got %s", str.Value)
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalNilLiteral(t *testing.T) {
	evaluated := testEval("nil")

	if evaluated != NULL {
		t.Fatalf("expected NULL, got %T", evaluated)
	}
}

func TestDeclarationStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"x := 5", 5},
		{"x := 5\nx", 5},
		{"x := 5\ny := x\ny", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIdentifierError(t *testing.T) {
	evaluated := testEval("x")

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "identifier not found: x" {
		t.Fatalf("wrong error message: %s", err.Message)
	}
}

func TestMutableDeclaration(t *testing.T) {
	input := `mut x := 10
x = 20
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestImmutableAssignmentError(t *testing.T) {
	input := `x := 10
x = 20`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	expected := "cannot assign to immutable variable: x"
	if err.Message != expected {
		t.Fatalf("wrong error message: got %q, want %q", err.Message, expected)
	}
}

func TestMultiValueAssignment(t *testing.T) {
	input := `mut a := 10
mut b := 20
a, b = b, a
a`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestMultiValueAssignmentSwap(t *testing.T) {
	input := `mut a := 10
mut b := 20
a, b = b, a
b`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func testEval(input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("expected Integer, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("expected %d, got %d", expected, result.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("expected Boolean, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("expected %t, got %t", expected, result.Value)
		return false
	}
	return true
}
