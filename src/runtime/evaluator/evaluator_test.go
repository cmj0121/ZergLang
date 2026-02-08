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

func TestReferenceExpression(t *testing.T) {
	input := `x := 42
ref := &x
ref`
	evaluated := testEval(input)

	ref, ok := evaluated.(*Reference)
	if !ok {
		t.Fatalf("expected Reference, got %T", evaluated)
	}

	intVal, ok := (*ref.Value).(*Integer)
	if !ok {
		t.Fatalf("expected referenced Integer, got %T", *ref.Value)
	}

	if intVal.Value != 42 {
		t.Fatalf("expected referenced value 42, got %d", intVal.Value)
	}
}

func TestAssertSuccess(t *testing.T) {
	input := `x := 10
assert x > 5
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestAssertFailure(t *testing.T) {
	input := `x := 3
assert x > 5`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "assertion failed" {
		t.Fatalf("expected 'assertion failed', got %s", err.Message)
	}
}

func TestAssertWithMessage(t *testing.T) {
	input := `x := 3
assert x > 5, "x must be greater than 5"`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "x must be greater than 5" {
		t.Fatalf("expected custom message, got %s", err.Message)
	}
}

func testEval(input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := NewEnvironmentWithBuiltins()
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

// Stdlib tests

func TestBuiltinLen(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`len("")`, int64(0)},
		{`len([1, 2, 3])`, int64(3)},
		{`len([])`, int64(0)},
		{`len({"a": 1, "b": 2})`, int64(2)},
		{`len({})`, int64(0)},
		{`len(1)`, "len() argument must be string, list, or map, not INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			err, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("expected Error, got %T (%+v)", evaluated, evaluated)
				continue
			}
			if err.Message != expected {
				t.Errorf("expected error %q, got %q", expected, err.Message)
			}
		}
	}
}

func TestBuiltinString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`string(42)`, "42"},
		{`string(true)`, "true"},
		{`string(false)`, "false"},
		{`string(nil)`, "nil"},
		{`string("hello")`, "hello"},
		{`string([1, 2, 3])`, "[1, 2, 3]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("expected String, got %T (%+v)", evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, str.Value)
		}
	}
}

func TestBuiltinInt(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`int("123")`, int64(123)},
		{`int("-456")`, int64(-456)},
		{`int(42)`, int64(42)},
		{`int(true)`, int64(1)},
		{`int(false)`, int64(0)},
		{`int("abc")`, "int() argument is not a valid integer: abc"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			err, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("expected Error, got %T (%+v)", evaluated, evaluated)
				continue
			}
			if err.Message != expected {
				t.Errorf("expected error %q, got %q", expected, err.Message)
			}
		}
	}
}

func TestListAppend(t *testing.T) {
	input := `nums := [1, 2, 3]
nums.append(4)`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(list.Elements))
	}

	testIntegerObject(t, list.Elements[3], 4)
}

func TestListAppendImmutable(t *testing.T) {
	// Verify append returns new list without modifying original
	input := `nums := [1, 2]
new_nums := nums.append(3)
nums.length`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 2)
}

func TestListPop(t *testing.T) {
	input := `[1, 2, 3].pop()`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 3)
}

func TestListPopEmpty(t *testing.T) {
	input := `[].pop()`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "pop() from empty list" {
		t.Fatalf("expected 'pop() from empty list', got %s", err.Message)
	}
}

func TestListFilter(t *testing.T) {
	input := `[1, 2, 3, 4, 5].filter(fn(x) { return x > 2 })`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	testIntegerObject(t, list.Elements[0], 3)
	testIntegerObject(t, list.Elements[1], 4)
	testIntegerObject(t, list.Elements[2], 5)
}

func TestListMap(t *testing.T) {
	input := `[1, 2, 3].map(fn(x) { return x * 2 })`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	testIntegerObject(t, list.Elements[0], 2)
	testIntegerObject(t, list.Elements[1], 4)
	testIntegerObject(t, list.Elements[2], 6)
}

func TestMapKeys(t *testing.T) {
	input := `{"a": 1, "b": 2}.keys()`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	// Keys are sorted alphabetically
	str1, ok := list.Elements[0].(*String)
	if !ok || str1.Value != "a" {
		t.Errorf("expected first key 'a', got %v", list.Elements[0])
	}
	str2, ok := list.Elements[1].(*String)
	if !ok || str2.Value != "b" {
		t.Errorf("expected second key 'b', got %v", list.Elements[1])
	}
}

func TestMapValues(t *testing.T) {
	input := `{"a": 1, "b": 2}.values()`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	// Values are in key-sorted order (a=1, b=2)
	testIntegerObject(t, list.Elements[0], 1)
	testIntegerObject(t, list.Elements[1], 2)
}

func TestMapContains(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"a": 1, "b": 2}.contains("a")`, true},
		{`{"a": 1, "b": 2}.contains("c")`, false},
		{`{1: "one", 2: "two"}.contains(1)`, true},
		{`{1: "one", 2: "two"}.contains(3)`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// ============================================================================
// Module tests
// ============================================================================

func TestSysModule(t *testing.T) {
	// Test _sys.os() (internal module, user should use src/std/sys.zg)
	osResult := testEval(`_sys.os()`)
	osStr, ok := osResult.(*String)
	if !ok {
		t.Fatalf("expected String for _sys.os(), got %T", osResult)
	}
	if osStr.Value != "darwin" && osStr.Value != "linux" && osStr.Value != "windows" {
		t.Fatalf("unexpected os: %s", osStr.Value)
	}

	// Test _sys.arch()
	archResult := testEval(`_sys.arch()`)
	archStr, ok := archResult.(*String)
	if !ok {
		t.Fatalf("expected String for _sys.arch(), got %T", archResult)
	}
	if archStr.Value != "amd64" && archStr.Value != "arm64" && archStr.Value != "386" {
		t.Fatalf("unexpected arch: %s", archStr.Value)
	}

	// Test _sys.args()
	argsResult := testEval(`_sys.args()`)
	_, ok = argsResult.(*List)
	if !ok {
		t.Fatalf("expected List for _sys.args(), got %T", argsResult)
	}

	// Test _sys.env()
	envResult := testEval(`_sys.env("PATH")`)
	_, ok = envResult.(*String)
	if !ok {
		t.Fatalf("expected String for _sys.env(), got %T", envResult)
	}
}

func TestStrModule(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`str.split("a,b,c", ",")`, []string{"a", "b", "c"}},
		{`str.join(["a", "b", "c"], ",")`, "a,b,c"},
		{`str.trim("  hello  ")`, "hello"},
		{`str.find("hello", "l")`, int64(2)},
		{`str.find("hello", "x")`, int64(-1)},
		{`str.replace("hello", "l", "L")`, "heLLo"},
		{`str.substring("hello", 1, 4)`, "ell"},
		{`str.starts_with("hello", "he")`, true},
		{`str.starts_with("hello", "lo")`, false},
		{`str.ends_with("hello", "lo")`, true},
		{`str.ends_with("hello", "he")`, false},
		{`str.upper("hello")`, "HELLO"},
		{`str.lower("HELLO")`, "hello"},
		{`str.contains("hello", "ell")`, true},
		{`str.contains("hello", "xyz")`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("for %s: expected String, got %T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("for %s: expected %q, got %q", tt.input, expected, str.Value)
			}
		case int64:
			testIntegerObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case []string:
			list, ok := evaluated.(*List)
			if !ok {
				t.Errorf("for %s: expected List, got %T", tt.input, evaluated)
				continue
			}
			if len(list.Elements) != len(expected) {
				t.Errorf("for %s: expected %d elements, got %d", tt.input, len(expected), len(list.Elements))
				continue
			}
			for i, exp := range expected {
				str, ok := list.Elements[i].(*String)
				if !ok || str.Value != exp {
					t.Errorf("for %s: element %d expected %q, got %v", tt.input, i, exp, list.Elements[i])
				}
			}
		}
	}
}

func TestCharModule(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`char.ord("A")`, int64(65)},
		{`char.ord("a")`, int64(97)},
		{`char.chr(65)`, "A"},
		{`char.chr(97)`, "a"},
		{`char.is_digit("5")`, true},
		{`char.is_digit("a")`, false},
		{`char.is_alpha("a")`, true},
		{`char.is_alpha("5")`, false},
		{`char.is_space(" ")`, true},
		{`char.is_space("a")`, false},
		{`char.is_alnum("a")`, true},
		{`char.is_alnum("5")`, true},
		{`char.is_alnum("!")`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("for %s: expected String, got %T", tt.input, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("for %s: expected %q, got %q", tt.input, expected, str.Value)
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

// ============================================================================
// Unsafe/asm tests
// ============================================================================

func TestUnsafeBlock(t *testing.T) {
	input := `unsafe {
		x := 42
		x
	}`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 42)
}

func TestAsmOutsideUnsafeError(t *testing.T) {
	input := `asm("sys_os")`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "asm() can only be used inside an unsafe block" {
		t.Fatalf("unexpected error: %s", err.Message)
	}
}

func TestAsmSysOs(t *testing.T) {
	input := `unsafe {
		asm("sys_os")
	}`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "darwin" && str.Value != "linux" && str.Value != "windows" {
		t.Fatalf("unexpected os: %s", str.Value)
	}
}

func TestAsmSysArch(t *testing.T) {
	input := `unsafe {
		asm("sys_arch")
	}`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "amd64" && str.Value != "arm64" && str.Value != "386" {
		t.Fatalf("unexpected arch: %s", str.Value)
	}
}

func TestAsmStrFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`unsafe { asm("str_upper", "hello") }`, "HELLO"},
		{`unsafe { asm("str_lower", "HELLO") }`, "hello"},
		{`unsafe { asm("str_trim", "  hello  ") }`, "hello"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("for %s: expected String, got %T", tt.input, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected %q, got %q", tt.input, tt.expected, str.Value)
		}
	}
}

func TestAsmUnknownFunction(t *testing.T) {
	input := `unsafe {
		asm("unknown_func")
	}`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message != "unknown asm function: unknown_func" {
		t.Fatalf("unexpected error: %s", err.Message)
	}
}

// ============================================================================
// Additional list method tests
// ============================================================================

func TestListJoin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`["a", "b", "c"].join(",")`, "a,b,c"},
		{`["hello", "world"].join(" ")`, "hello world"},
		{`[1, 2, 3].join("-")`, "1-2-3"},
		{`[].join(",")`, ""},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("for %s: expected String, got %T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("for %s: expected %q, got %q", tt.input, tt.expected, str.Value)
		}
	}
}

func TestListSlice(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{`[1, 2, 3, 4, 5].slice(1, 4)`, []int64{2, 3, 4}},
		{`[1, 2, 3].slice(0, 2)`, []int64{1, 2}},
		{`[1, 2, 3].slice(0, 0)`, []int64{}},
		{`[1, 2, 3].slice(5, 10)`, []int64{}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		list, ok := evaluated.(*List)
		if !ok {
			t.Errorf("for %s: expected List, got %T", tt.input, evaluated)
			continue
		}
		if len(list.Elements) != len(tt.expected) {
			t.Errorf("for %s: expected %d elements, got %d", tt.input, len(tt.expected), len(list.Elements))
			continue
		}
		for i, exp := range tt.expected {
			testIntegerObject(t, list.Elements[i], exp)
		}
	}
}

func TestListIndexMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`[1, 2, 3].index(2)`, 1},
		{`[1, 2, 3].index(1)`, 0},
		{`[1, 2, 3].index(3)`, 2},
		{`[1, 2, 3].index(5)`, -1},
		{`["a", "b", "c"].index("b")`, 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestListReverse(t *testing.T) {
	input := `[1, 2, 3].reverse()`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	testIntegerObject(t, list.Elements[0], 3)
	testIntegerObject(t, list.Elements[1], 2)
	testIntegerObject(t, list.Elements[2], 1)
}

func TestListSort(t *testing.T) {
	input := `[3, 1, 2].sort()`
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

func TestListSortStrings(t *testing.T) {
	input := `["c", "a", "b"].sort()`
	evaluated := testEval(input)

	list, ok := evaluated.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", evaluated)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	expected := []string{"a", "b", "c"}
	for i, exp := range expected {
		str, ok := list.Elements[i].(*String)
		if !ok || str.Value != exp {
			t.Errorf("element %d: expected %q, got %v", i, exp, list.Elements[i])
		}
	}
}

// ============================================================================
// Integration tests
// ============================================================================

func TestPlatformDetection(t *testing.T) {
	input := `os := _sys.os()
if os == "darwin" {
	result := "macOS"
} else {
	result := "other"
}
os`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	// Just verify it returned a valid OS string
	if str.Value != "darwin" && str.Value != "linux" && str.Value != "windows" {
		t.Fatalf("unexpected os: %s", str.Value)
	}
}

func TestModuleMethodChaining(t *testing.T) {
	input := `str.upper(str.trim("  hello  "))`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "HELLO" {
		t.Fatalf("expected HELLO, got %s", str.Value)
	}
}

// ============================================================================
// Enum and Match tests
// ============================================================================

func TestEnumDeclaration(t *testing.T) {
	input := `enum TokenType {
		INT
		STRING
		PLUS
	}
	TokenType`
	evaluated := testEval(input)

	enumType, ok := evaluated.(*EnumType)
	if !ok {
		t.Fatalf("expected EnumType, got %T", evaluated)
	}

	if enumType.Name != "TokenType" {
		t.Fatalf("expected name 'TokenType', got %s", enumType.Name)
	}

	if len(enumType.Variants) != 3 {
		t.Fatalf("expected 3 variants, got %d", len(enumType.Variants))
	}

	expected := []string{"INT", "STRING", "PLUS"}
	for i, v := range expected {
		if enumType.Variants[i] != v {
			t.Errorf("variant %d: expected %s, got %s", i, v, enumType.Variants[i])
		}
	}
}

func TestEnumVariantAccess(t *testing.T) {
	input := `enum Color {
		RED
		GREEN
		BLUE
	}
	Color.RED`
	evaluated := testEval(input)

	enumVal, ok := evaluated.(*EnumValue)
	if !ok {
		t.Fatalf("expected EnumValue, got %T", evaluated)
	}

	if enumVal.EnumName != "Color" {
		t.Fatalf("expected enum name 'Color', got %s", enumVal.EnumName)
	}

	if enumVal.VariantName != "RED" {
		t.Fatalf("expected variant 'RED', got %s", enumVal.VariantName)
	}
}

func TestEnumComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`enum Color { RED GREEN }
		Color.RED == Color.RED`, true},
		{`enum Color { RED GREEN }
		Color.RED == Color.GREEN`, false},
		{`enum Color { RED GREEN }
		c := Color.RED
		c == Color.RED`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestMatchLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`match 1 {
			1 => { "one" }
			2 => { "two" }
			_ => { "other" }
		}`, "one"},
		{`match 2 {
			1 => { "one" }
			2 => { "two" }
			_ => { "other" }
		}`, "two"},
		{`match 3 {
			1 => { "one" }
			2 => { "two" }
			_ => { "other" }
		}`, "other"},
		{`match "hello" {
			"world" => { 1 }
			"hello" => { 2 }
			_ => { 3 }
		}`, int64(2)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*String)
			if !ok {
				t.Errorf("expected String, got %T (%+v)", evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("expected %q, got %q", expected, str.Value)
			}
		case int64:
			testIntegerObject(t, evaluated, expected)
		}
	}
}

func TestMatchEnum(t *testing.T) {
	input := `enum Status {
		OK
		ERROR
		PENDING
	}
	status := Status.OK
	match status {
		Status.OK => { "success" }
		Status.ERROR => { "failure" }
		_ => { "unknown" }
	}`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "success" {
		t.Fatalf("expected 'success', got %s", str.Value)
	}
}

func TestMatchWildcard(t *testing.T) {
	input := `match 42 {
		1 => { "one" }
		_ => { "anything else" }
	}`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	if str.Value != "anything else" {
		t.Fatalf("expected 'anything else', got %s", str.Value)
	}
}

func TestMatchAlternatives(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`match 1 {
			1 | 2 | 3 => { "small" }
			_ => { "big" }
		}`, "small"},
		{`match 2 {
			1 | 2 | 3 => { "small" }
			_ => { "big" }
		}`, "small"},
		{`match 10 {
			1 | 2 | 3 => { "small" }
			_ => { "big" }
		}`, "big"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("expected String, got %T (%+v)", evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, str.Value)
		}
	}
}

func TestMatchWithGuard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`x := 5
		match x {
			_ if x > 10 => { "big" }
			_ if x > 0 => { "positive" }
			_ => { "zero or negative" }
		}`, "positive"},
		{`x := 15
		match x {
			_ if x > 10 => { "big" }
			_ if x > 0 => { "positive" }
			_ => { "zero or negative" }
		}`, "big"},
		{`x := -5
		match x {
			_ if x > 10 => { "big" }
			_ if x > 0 => { "positive" }
			_ => { "zero or negative" }
		}`, "zero or negative"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("expected String, got %T (%+v)", evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, str.Value)
		}
	}
}

func TestResultOkErr(t *testing.T) {
	// Test Ok() creation
	okResult := testEval(`Ok(42)`)
	resultOk, ok := okResult.(*ResultOk)
	if !ok {
		t.Fatalf("expected ResultOk, got %T", okResult)
	}
	testIntegerObject(t, resultOk.Value, 42)

	// Test Err() creation
	errResult := testEval(`Err("not found")`)
	resultErr, ok := errResult.(*ResultErr)
	if !ok {
		t.Fatalf("expected ResultErr, got %T", errResult)
	}
	str, ok := resultErr.Error.(*String)
	if !ok || str.Value != "not found" {
		t.Fatalf("expected error 'not found', got %v", resultErr.Error)
	}
}

func TestResultValueAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Ok(42).value`, int64(42)},
		{`Err("oops").error`, "oops"},
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
				t.Errorf("expected %q, got %q", expected, str.Value)
			}
		}
	}
}

func TestIsExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`Ok(42) is Ok`, true},
		{`Ok(42) is Err`, false},
		{`Err("error") is Err`, true},
		{`Err("error") is Ok`, false},
		{`42 is int`, true},
		{`42 is string`, false},
		{`"hello" is string`, true},
		{`[1, 2, 3] is list`, true},
		{`{"a": 1} is map`, true},
		{`true is bool`, true},
		{`nil is nil`, true},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			t.Errorf("failed at test case %d: %s", i, tt.input)
		}
	}
}

func TestIsWithEnum(t *testing.T) {
	input := `enum Color { RED GREEN BLUE }
	c := Color.RED
	c is Color`
	evaluated := testEval(input)
	testBooleanObject(t, evaluated, true)
}

func TestMatchOnResult(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`result := Ok(42)
		match result {
			Ok => { "success" }
			Err => { "failure" }
		}`, "success"},
		{`result := Err("oops")
		match result {
			Ok => { "success" }
			Err => { "failure" }
		}`, "failure"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("expected String, got %T (%+v)", evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, str.Value)
		}
	}
}

func TestMatchNoMatch(t *testing.T) {
	input := `match 5 {
		1 => { "one" }
		2 => { "two" }
	}`
	evaluated := testEval(input)

	if evaluated != NULL {
		t.Fatalf("expected NULL when no match, got %T (%+v)", evaluated, evaluated)
	}
}

// ============================================================================
// Float tests
// ============================================================================

func TestFloatLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"123.456", 123.456},
		{"1_000.5", 1000.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestFloatArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.5 + 2.5", 4.0},
		{"5.0 - 2.0", 3.0},
		{"2.5 * 4.0", 10.0},
		{"10.0 / 4.0", 2.5},
		{"-3.5", -3.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestFloatIntMixedArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.5 + 2", 3.5},
		{"5 - 2.5", 2.5},
		{"2 * 3.5", 7.0},
		{"10 / 4.0", 2.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestFloatComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"3.14 < 3.15", true},
		{"3.14 > 3.15", false},
		{"3.14 == 3.14", true},
		{"3.14 != 3.15", true},
		{"3.14 <= 3.14", true},
		{"3.14 >= 3.14", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestFloatBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"float(42)", 42.0},
		{`float("3.14")`, 3.14},
		{"float(3.14)", 3.14},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestIntFromFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"int(3.7)", 3},
		{"int(3.14)", 3},
		{"int(-2.9)", -2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFloatIsType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"3.14 is float", true},
		{"3.14 is int", false},
		{"42 is float", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testFloatObject(t *testing.T, obj Object, expected float64) bool {
	result, ok := obj.(*Float)
	if !ok {
		t.Errorf("expected Float, got %T (%+v)", obj, obj)
		return false
	}
	// Use small epsilon for float comparison
	if diff := result.Value - expected; diff < -0.0001 || diff > 0.0001 {
		t.Errorf("expected %f, got %f", expected, result.Value)
		return false
	}
	return true
}

// ============================================================================
// Range tests
// ============================================================================

func TestRangeExpression(t *testing.T) {
	input := `1..5`
	evaluated := testEval(input)

	r, ok := evaluated.(*Range)
	if !ok {
		t.Fatalf("expected Range, got %T", evaluated)
	}

	if r.Start != 1 {
		t.Errorf("expected start 1, got %d", r.Start)
	}
	if r.End != 5 {
		t.Errorf("expected end 5, got %d", r.End)
	}
	if r.Inclusive {
		t.Errorf("expected exclusive range")
	}
}

func TestRangeInclusiveExpression(t *testing.T) {
	input := `1..=5`
	evaluated := testEval(input)

	r, ok := evaluated.(*Range)
	if !ok {
		t.Fatalf("expected Range, got %T", evaluated)
	}

	if r.Start != 1 {
		t.Errorf("expected start 1, got %d", r.Start)
	}
	if r.End != 5 {
		t.Errorf("expected end 5, got %d", r.End)
	}
	if !r.Inclusive {
		t.Errorf("expected inclusive range")
	}
}

func TestForInRange(t *testing.T) {
	input := `mut sum := 0
for i in 1..5 {
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10) // 1+2+3+4 = 10
}

func TestForInRangeInclusive(t *testing.T) {
	input := `mut sum := 0
for i in 1..=5 {
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15) // 1+2+3+4+5 = 15
}

func TestForInRangeWithVariables(t *testing.T) {
	input := `start := 2
end := 6
mut sum := 0
for i in start..end {
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 14) // 2+3+4+5 = 14
}

func TestForInRangeBreak(t *testing.T) {
	input := `mut sum := 0
for i in 1..10 {
    if i == 5 {
        break
    }
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10) // 1+2+3+4 = 10
}

func TestRangeInspect(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1..5", "1..5"},
		{"1..=5", "1..=5"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		r, ok := evaluated.(*Range)
		if !ok {
			t.Errorf("expected Range, got %T", evaluated)
			continue
		}
		if r.Inspect() != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, r.Inspect())
		}
	}
}

// ============================================================================
// Compound assignment tests
// ============================================================================

func TestCompoundPlusAssign(t *testing.T) {
	input := `mut x := 10
x += 5
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestCompoundMinusAssign(t *testing.T) {
	input := `mut x := 10
x -= 3
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 7)
}

func TestCompoundAsteriskAssign(t *testing.T) {
	input := `mut x := 10
x *= 2
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestCompoundSlashAssign(t *testing.T) {
	input := `mut x := 10
x /= 4
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 2)
}

func TestCompoundPercentAssign(t *testing.T) {
	input := `mut x := 10
x %= 3
x`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 1)
}

func TestCompoundAssignWithFloat(t *testing.T) {
	input := `mut x := 10.0
x += 2.5
x`
	evaluated := testEval(input)
	testFloatObject(t, evaluated, 12.5)
}

func TestCompoundAssignStringConcat(t *testing.T) {
	input := `mut s := "hello"
s += " world"
s`
	evaluated := testEval(input)
	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}
	if str.Value != "hello world" {
		t.Fatalf("expected 'hello world', got %s", str.Value)
	}
}

func TestCompoundAssignImmutableError(t *testing.T) {
	input := `x := 10
x += 5`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	expected := "cannot assign to immutable variable: x"
	if err.Message != expected {
		t.Fatalf("expected %q, got %q", expected, err.Message)
	}
}

func TestCompoundAssignInLoop(t *testing.T) {
	input := `mut sum := 0
for i in 1..=5 {
    sum += i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

// ============================================================================
// String interpolation tests
// ============================================================================

func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`name := "World"
"Hello, {name}!"`, "Hello, World!"},
		{`x := 42
"The answer is {x}"`, "The answer is 42"},
		{`a := 2
b := 3
"{a} + {b} = {a + b}"`, "2 + 3 = 5"},
		{`"no interpolation"`, "no interpolation"},
		{`"{1 + 1}"`, "2"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*String)
		if !ok {
			t.Errorf("for %q: expected String, got %T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("for %q: expected %q, got %q", tt.input, tt.expected, str.Value)
		}
	}
}

func TestStringInterpolationMultiple(t *testing.T) {
	input := `first := "John"
last := "Doe"
age := 30
"Name: {first} {last}, Age: {age}"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "Name: John Doe, Age: 30"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

func TestStringInterpolationWithExpressions(t *testing.T) {
	input := `x := 10
y := 5
"{x} * {y} = {x * y}"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "10 * 5 = 50"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

func TestStringInterpolationEscape(t *testing.T) {
	// Test escaping braces with \{
	input := `"literal \{brace\}"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "literal {brace}"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

// ========== Import/Module System Tests ==========

func TestImportModule(t *testing.T) {
	// Set up the module loader to find testdata
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "math"
result := math.add(2, 3)
result`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestImportModuleWithAlias(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "math" as m
result := m.multiply(4, 5)
result`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestImportModuleFunctionCall(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "greet"
msg := greet.hello("World")
msg`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "Hello, World"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

func TestImportModuleConstant(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "math"
math.PI`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 3)
}

func TestImportModuleConstantWithAlias(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "greet" as g
g.DEFAULT_GREETING`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "Hi there"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

func TestImportModuleNotFound(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "nonexistent"`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestImportModuleMemberNotFound(t *testing.T) {
	DefaultLoader.SetCurrentDir("testdata")
	defer DefaultLoader.SetCurrentDir(".")

	input := `import "math"
math.nonexistent`
	evaluated := testEval(input)

	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("expected Error, got %T", evaluated)
	}

	if err.Message == "" {
		t.Fatalf("expected error message about missing member")
	}
}

func TestImportNestedModule(t *testing.T) {
	// Clear cache to ensure fresh load
	DefaultLoader = NewModuleLoader()
	DefaultLoader.SetCurrentDir("testdata")
	defer func() {
		DefaultLoader = NewModuleLoader()
	}()

	input := `import "utils"
result := utils.square(5)
result`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 25)
}

func TestImportNestedModuleDouble(t *testing.T) {
	// Clear cache to ensure fresh load
	DefaultLoader = NewModuleLoader()
	DefaultLoader.SetCurrentDir("testdata")
	defer func() {
		DefaultLoader = NewModuleLoader()
	}()

	input := `import "utils"
result := utils.double(7)
result`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 14)
}

func TestImportNestedModuleConstant(t *testing.T) {
	// Clear cache to ensure fresh load
	DefaultLoader = NewModuleLoader()
	DefaultLoader.SetCurrentDir("testdata")
	defer func() {
		DefaultLoader = NewModuleLoader()
	}()

	input := `import "utils"
utils.VERSION`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "1.0"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}

func TestImportModuleCaching(t *testing.T) {
	// Clear cache
	DefaultLoader = NewModuleLoader()
	DefaultLoader.SetCurrentDir("testdata")
	defer func() {
		DefaultLoader = NewModuleLoader()
	}()

	// Import same module twice, should use cached version
	input := `import "math"
import "math" as m2
result := math.add(1, 2) + m2.add(3, 4)
result`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10) // 3 + 7
}

func TestImportMultipleModules(t *testing.T) {
	DefaultLoader = NewModuleLoader()
	DefaultLoader.SetCurrentDir("testdata")
	defer func() {
		DefaultLoader = NewModuleLoader()
	}()

	input := `import "math"
import "greet"
sum := math.add(10, 20)
msg := greet.hello("Zerg")
"{msg} - Sum is {sum}"`
	evaluated := testEval(input)

	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", evaluated)
	}

	expected := "Hello, Zerg - Sum is 30"
	if str.Value != expected {
		t.Fatalf("expected %q, got %q", expected, str.Value)
	}
}
