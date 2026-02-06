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
	case *parser.IfStatement:
		return evalIfStatement(node, env)
	case *parser.ForInStatement:
		return evalForInStatement(node, env)
	case *parser.ForConditionStatement:
		return evalForConditionStatement(node, env)
	case *parser.BreakStatement:
		return BREAK
	case *parser.ContinueStatement:
		return CONTINUE
	case *parser.NopStatement:
		return NULL
	case *parser.ListLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return &List{Elements: elements}
	case *parser.MapLiteral:
		return evalMapLiteral(node, env)
	case *parser.IndexExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *parser.MemberExpression:
		obj := Eval(node.Object, env)
		if IsError(obj) {
			return obj
		}
		return evalMemberExpression(obj, node.Member.Value)
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
			if rt == RETURN_VALUE_OBJ || rt == "ERROR" || rt == BREAK_OBJ || rt == CONTINUE_OBJ {
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
	if ok {
		if !b.mutable {
			return fmt.Errorf("cannot assign to immutable variable: %s", name)
		}
		b.value = val
		return nil
	}
	// Search in outer environments
	if e.outer != nil {
		return e.outer.Assign(name, val)
	}
	return fmt.Errorf("identifier not found: %s", name)
}

// Set stores a variable in the environment (for backward compatibility).
func (e *Environment) Set(name string, val Object) Object {
	return e.Declare(name, val, false)
}

// evalIfStatement evaluates an if statement.
func evalIfStatement(is *parser.IfStatement, env *Environment) Object {
	condition := Eval(is.Condition, env)
	if IsError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(is.Consequence, env)
	} else if is.Alternative != nil {
		return Eval(is.Alternative, env)
	}
	return NULL
}

// evalForInStatement evaluates a for-in loop.
func evalForInStatement(fis *parser.ForInStatement, env *Environment) Object {
	iterable := Eval(fis.Iterable, env)
	if IsError(iterable) {
		return iterable
	}

	switch obj := iterable.(type) {
	case *String:
		return evalForInString(fis, obj.Value, env)
	case *List:
		return evalForInList(fis, obj, env)
	default:
		return newError("cannot iterate over %s", iterable.Type())
	}
}

func evalForInString(fis *parser.ForInStatement, str string, env *Environment) Object {
	var result Object = NULL

	loopEnv := NewEnclosedEnvironment(env)

	for _, ch := range str {
		loopEnv.Declare(fis.Variable.Value, &String{Value: string(ch)}, false)

		result = Eval(fis.Body, loopEnv)

		if result != nil {
			switch result.Type() {
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				continue
			case RETURN_VALUE_OBJ, "ERROR":
				return result
			}
		}
	}

	return result
}

// evalForConditionStatement evaluates a for loop with condition.
func evalForConditionStatement(fcs *parser.ForConditionStatement, env *Environment) Object {
	var result Object = NULL

	for {
		// Check condition if present (nil means infinite loop)
		// Condition is evaluated in the outer environment
		if fcs.Condition != nil {
			condition := Eval(fcs.Condition, env)
			if IsError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Body gets its own scope for each iteration
		loopEnv := NewEnclosedEnvironment(env)
		result = Eval(fcs.Body, loopEnv)

		if result != nil {
			switch result.Type() {
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				continue
			case RETURN_VALUE_OBJ, "ERROR":
				return result
			}
		}
	}

	return result
}

func evalForInList(fis *parser.ForInStatement, list *List, env *Environment) Object {
	var result Object = NULL

	loopEnv := NewEnclosedEnvironment(env)

	for _, el := range list.Elements {
		loopEnv.Declare(fis.Variable.Value, el, false)

		result = Eval(fis.Body, loopEnv)

		if result != nil {
			switch result.Type() {
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				continue
			case RETURN_VALUE_OBJ, "ERROR":
				return result
			}
		}
	}

	return result
}

func evalMapLiteral(node *parser.MapLiteral, env *Environment) Object {
	pairs := make(map[HashKey]MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if IsError(key) {
			return key
		}

		hashKey, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if IsError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = MapPair{Key: key, Value: value}
	}

	return &Map{Pairs: pairs}
}

func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == LIST_OBJ && index.Type() == INTEGER_OBJ:
		return evalListIndexExpression(left, index)
	case left.Type() == MAP_OBJ:
		return evalMapIndexExpression(left, index)
	case left.Type() == STRING_OBJ && index.Type() == INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalListIndexExpression(list, index Object) Object {
	listObj := list.(*List)
	idx := index.(*Integer).Value
	max := int64(len(listObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return listObj.Elements[idx]
}

func evalMapIndexExpression(m, index Object) Object {
	mapObj := m.(*Map)

	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := mapObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalStringIndexExpression(str, index Object) Object {
	strObj := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(strObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &String{Value: string(strObj.Value[idx])}
}

func evalMemberExpression(obj Object, member string) Object {
	switch o := obj.(type) {
	case *Map:
		// Try to access map with string key
		key := &String{Value: member}
		pair, ok := o.Pairs[key.HashKey()]
		if !ok {
			return NULL
		}
		return pair.Value
	case *List:
		// List built-in methods
		switch member {
		case "length":
			return &Integer{Value: int64(len(o.Elements))}
		}
	case *String:
		// String built-in methods
		switch member {
		case "length":
			return &Integer{Value: int64(len(o.Value))}
		}
	}
	return newError("no member '%s' on type %s", member, obj.Type())
}
