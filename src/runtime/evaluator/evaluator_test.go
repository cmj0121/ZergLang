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

func TestListLiteral(t *testing.T) {
	input := `[1, 2, 3]`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	testIntegerObject(t, list.Elements[0], 1)
	testIntegerObject(t, list.Elements[1], 2)
	testIntegerObject(t, list.Elements[2], 3)
}

func TestListIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", int64(1)},
		{"[1, 2, 3][1]", int64(2)},
		{"[1, 2, 3][2]", int64(3)},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
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

func TestMapLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2}`
	evaluated := testEval(input)

	m, ok := evaluated.(*Map)
	if !ok {
		t.Fatalf("expected Map, got %T", evaluated)
	}

	if len(m.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(m.Pairs))
	}
}

func TestMapIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, int64(5)},
		{`{"foo": 5}["bar"]`, nil},
		{`{1: "one"}[1]`, "one"},
		{`{true: "yes"}[true]`, "yes"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("expected String, got %T", evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("expected %s, got %s", expected, str.Value)
			}
		case nil:
			if evaluated != NULL {
				t.Errorf("expected NULL, got %T", evaluated)
			}
		}
	}
}

func TestMemberAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"name": "Alice"}.name`, "Alice"},
		{`[1, 2, 3].length`, int64(3)},
		{`"hello".length`, int64(5)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("expected String, got %T", evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("expected %s, got %s", expected, str.Value)
			}
		}
	}
}

func TestStringIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello"[0]`, "h"},
		{`"hello"[4]`, "o"},
		{`"hello"[5]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("expected String, got %T", evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("expected %s, got %s", expected, str.Value)
			}
		case nil:
			if evaluated != NULL {
				t.Errorf("expected NULL, got %T", evaluated)
			}
		}
	}
}

func TestForInList(t *testing.T) {
	input := `mut sum := 0
for n in [1, 2, 3, 4] {
    sum = sum + n
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestClassDeclaration(t *testing.T) {
	input := `class Point {
    pub mut x = 0
    pub mut y = 0
}
Point`
	evaluated := testEval(input)

	class, ok := evaluated.(*Class)
	if !ok {
		t.Fatalf("expected Class, got %T", evaluated)
	}

	if class.Name != "Point" {
		t.Fatalf("expected class name 'Point', got %s", class.Name)
	}

	if len(class.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(class.Fields))
	}
}

func TestClassInstantiation(t *testing.T) {
	input := `class Counter {
    pub mut value = 0
}
c := Counter()
c.value`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 0)
}

func TestClassMethods(t *testing.T) {
	input := `class Point {
    pub mut x = 0
    pub mut y = 0
}
impl Point {
    fn init(x, y) {
        this.x = x
        this.y = y
    }
    fn sum() {
        return this.x + this.y
    }
}
p := Point(3, 4)
p.sum()`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 7)
}

func TestStaticMethods(t *testing.T) {
	input := `class Math {
    pub mut dummy = 0
}
impl Math {
    pub static fn add(a, b) {
        return a + b
    }
}
Math.add(10, 20)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 30)
}

func TestMemberAssignment(t *testing.T) {
	input := `class Box {
    pub mut value = 0
}
b := Box()
b.value = 42
b.value`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 42)
}

func TestIndexAssignment(t *testing.T) {
	input := `mut arr := [1, 2, 3]
arr[1] = 99
arr[1]`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 99)
}

func TestSpecDeclaration(t *testing.T) {
	input := `spec Drawable {
    fn draw()
    fn area()
}
Drawable`
	evaluated := testEval(input)

	spec, ok := evaluated.(*Spec)
	if !ok {
		t.Fatalf("expected Spec, got %T", evaluated)
	}

	if spec.Name != "Drawable" {
		t.Fatalf("expected spec name 'Drawable', got %s", spec.Name)
	}

	if len(spec.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(spec.Methods))
	}
}

func TestSpecWithTypedParameters(t *testing.T) {
	input := `spec Container {
    fn get(index: int) -> int
    fn set(index: int, value: int)
}
Container`
	evaluated := testEval(input)

	spec, ok := evaluated.(*Spec)
	if !ok {
		t.Fatalf("expected Spec, got %T", evaluated)
	}

	if spec.Name != "Container" {
		t.Fatalf("expected spec name 'Container', got %s", spec.Name)
	}

	if len(spec.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(spec.Methods))
	}

	getMethod := spec.Methods["get"]
	if getMethod == nil {
		t.Fatalf("expected 'get' method")
	}
	if len(getMethod.Parameters) != 1 {
		t.Fatalf("expected 1 parameter for get, got %d", len(getMethod.Parameters))
	}

	setMethod := spec.Methods["set"]
	if setMethod == nil {
		t.Fatalf("expected 'set' method")
	}
	if len(setMethod.Parameters) != 2 {
		t.Fatalf("expected 2 parameters for set, got %d", len(setMethod.Parameters))
	}
}

func TestImplForSpec(t *testing.T) {
	input := `spec Measurable {
    fn size()
}
class Box {
    pub mut value = 0
}
impl Box for Measurable {
    fn size() {
        return this.value
    }
}
impl Box {
    fn init(v) {
        this.value = v
    }
}
b := Box(42)
b.size()`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 42)
}

func TestSpecNotImplementedError(t *testing.T) {
	input := `spec Required {
    fn mustHave()
}
class Empty {
    pub mut x = 0
}
impl Empty for Required {
    # Missing mustHave method
}`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message == "" {
		t.Fatalf("expected error message")
	}
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
