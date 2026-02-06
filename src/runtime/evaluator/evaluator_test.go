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

func TestArithmeticOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"10 / 2", 5},
		{"10 % 3", 1},
		{"2 ** 3", 8},
		{"10 + 5 * 2", 20},
		{"(10 + 5) * 2", 30},
		{"-5", -5},
		{"-10 + 20", 10},
		{"2 ** 3 ** 2", 512},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 < 10", true},
		{"5 > 10", false},
		{"5 <= 5", true},
		{"5 >= 5", true},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 10", false},
		{"5 != 10", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true and true", true},
		{"true and false", false},
		{"false or true", true},
		{"false or false", false},
		{"not true", false},
		{"not false", true},
		{"10 > 5 and 3 < 7", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" + " " + "world"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "hello world" {
		t.Fatalf("expected 'hello world', got %s", str.Value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"-5", int64(-5)},
		{"-10", int64(-10)},
		{"not true", false},
		{"not false", true},
		{"not nil", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := `fn add(a, b) {
    return a + b
}
add(10, 20)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 30)
}

func TestAnonymousFunction(t *testing.T) {
	input := `double := fn(x) { return x * 2 }
double(5)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestClosures(t *testing.T) {
	input := `fn makeAdder(x) {
    return fn(y) { return x + y }
}
add5 := makeAdder(5)
add5(10)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestDefaultParameters(t *testing.T) {
	input := `fn greet(name = "world") {
    return "hello " + name
}
greet()`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}
	if str.Value != "hello world" {
		t.Fatalf("expected 'hello world', got %s", str.Value)
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", int64(10)},
		{"if false { 10 }", nil},
		{"if 1 < 2 { 10 }", int64(10)},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", int64(20)},
		{"if 1 < 2 { 10 } else { 20 }", int64(10)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case nil:
			if evaluated != NULL {
				t.Errorf("expected NULL, got %T", evaluated)
			}
		}
	}
}

func TestForConditionLoop(t *testing.T) {
	input := `mut count := 0
for count < 5 {
    count = count + 1
}
count`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestForInLoop(t *testing.T) {
	input := `mut result := ""
for ch in "abc" {
    result = result + ch
}
result`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}
	if str.Value != "abc" {
		t.Fatalf("expected 'abc', got %s", str.Value)
	}
}

func TestBreakStatement(t *testing.T) {
	input := `mut sum := 0
mut i := 0
for i < 10 {
    if i == 5 {
        break
    }
    sum = sum + i
    i = i + 1
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10) // 0+1+2+3+4 = 10
}

func TestContinueStatement(t *testing.T) {
	input := `mut sum := 0
mut i := 0
for i < 5 {
    i = i + 1
    if i == 3 {
        continue
    }
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 12) // 1+2+4+5 = 12
}

func TestNopStatement(t *testing.T) {
	input := `x := 5
nop
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
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
