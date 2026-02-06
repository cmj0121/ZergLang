package evaluator

import (
	"fmt"

	"github.com/xrspace/zerglang/runtime/parser"
)

// Eval evaluates an AST node and returns the result.
func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {
	case *parser.Program:
		return evalProgram(node, env)
	case *parser.DeclarationStatement:
		return evalDeclarationStatement(node, env)
	case *parser.AssignmentStatement:
		return evalAssignmentStatement(node, env)
	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)
	case *parser.Identifier:
		return evalIdentifier(node, env)
	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *parser.StringLiteral:
		return &String{Value: node.Value}
	case *parser.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *parser.NilLiteral:
		return NULL
	case *parser.PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *parser.InfixExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *parser.BlockStatement:
		return evalBlockStatement(node, env)
	case *parser.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *parser.FunctionLiteral:
		return &Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *parser.CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *parser.Program, env *Environment) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
		if IsError(result) {
			return result
		}
	}

	return result
}

func evalBlockStatement(block *parser.BlockStatement, env *Environment) Object {
	var result Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == "ERROR" {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exps []parser.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn Object, args []Object) Object {
	function, ok := fn.(*Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Declare(param.Name.Value, args[i], false)
		} else if param.Default != nil {
			// Evaluate default value in function's closure environment
			defaultVal := Eval(param.Default, fn.Env)
			env.Declare(param.Name.Value, defaultVal, false)
		}
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalDeclarationStatement(ds *parser.DeclarationStatement, env *Environment) Object {
	val := Eval(ds.Value, env)
	if val == nil {
		return nil
	}
	env.Declare(ds.Name.Value, val, ds.Mutable)
	return val
}

func evalAssignmentStatement(as *parser.AssignmentStatement, env *Environment) Object {
	// Evaluate all values first (for swap: a, b = b, a)
	values := make([]Object, len(as.Values))
	for i, expr := range as.Values {
		val := Eval(expr, env)
		if IsError(val) {
			return val
		}
		values[i] = val
	}

	// Assign all values
	for i, name := range as.Names {
		if err := env.Assign(name.Value, values[i]); err != nil {
			return newError("%s", err.Error())
		}
	}

	if len(values) == 1 {
		return values[0]
	}
	return values[len(values)-1]
}

func evalIdentifier(node *parser.Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "not":
		return evalNotOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalNotOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case operator == "and":
		return evalAndExpression(left, right)
	case operator == "or":
		return evalOrExpression(left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal % rightVal}
	case "**":
		return &Integer{Value: intPow(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalAndExpression(left, right Object) Object {
	if isTruthy(left) {
		return right
	}
	return left
}

func evalOrExpression(left, right Object) Object {
	if isTruthy(left) {
		return left
	}
	return right
}

func isTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func intPow(base, exp int64) int64 {
	if exp < 0 {
		return 0
	}
	result := int64(1)
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
}

// Error represents a runtime error.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return "ERROR" }
func (e *Error) Inspect() string  { return "error: " + e.Message }

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// IsError checks if an object is an error.
func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == "ERROR"
	}
	return false
}

// binding represents a variable binding with its value and mutability.
type binding struct {
	value   Object
	mutable bool
}

// Environment stores variable bindings.
type Environment struct {
	store map[string]*binding
	outer *Environment
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]*binding), outer: nil}
}

// NewEnclosedEnvironment creates a new Environment with an outer scope.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (Object, bool) {
	b, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	if !ok {
		return nil, false
	}
	return b.value, true
}

// Function represents a function value with closure.
type Function struct {
	Parameters []*parser.Parameter
	Body       *parser.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string  { return "fn(...) {...}" }

// Declare creates a new variable in the environment.
func (e *Environment) Declare(name string, val Object, mutable bool) Object {
	e.store[name] = &binding{value: val, mutable: mutable}
	return val
}

// Assign updates a mutable variable in the environment.
func (e *Environment) Assign(name string, val Object) error {
	b, ok := e.store[name]
	if !ok {
		return fmt.Errorf("identifier not found: %s", name)
	}
	if !b.mutable {
		return fmt.Errorf("cannot assign to immutable variable: %s", name)
	}
	b.value = val
	return nil
}

// Set stores a variable in the environment (for backward compatibility).
func (e *Environment) Set(name string, val Object) Object {
	return e.Declare(name, val, false)
}
