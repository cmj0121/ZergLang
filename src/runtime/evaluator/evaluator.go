package evaluator

import (
	"fmt"
	"os"

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
	case *parser.CompoundAssignmentStatement:
		return evalCompoundAssignmentStatement(node, env)
	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)
	case *parser.Identifier:
		return evalIdentifier(node, env)
	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *parser.FloatLiteral:
		return &Float{Value: node.Value}
	case *parser.StringLiteral:
		return &String{Value: node.Value}
	case *parser.InterpolatedString:
		return evalInterpolatedString(node, env)
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
		// Check postfix condition first
		if node.Condition != nil {
			cond := Eval(node.Condition, env)
			if IsError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				return NULL // condition not met, don't return
			}
		}
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
		// Evaluate named arguments
		namedArgs := make(map[string]Object)
		for _, na := range node.NamedArgs {
			val := Eval(na.Value, env)
			if IsError(val) {
				return val
			}
			namedArgs[na.Name] = val
		}
		return applyFunctionWithNamedArgs(function, args, namedArgs)
	case *parser.ChainedAssignment:
		// Evaluate the left side (the object)
		obj := Eval(node.Left, env)
		if IsError(obj) {
			return obj
		}
		// Must be an instance
		instance, ok := obj.(*Instance)
		if !ok {
			return &Error{Message: "builder syntax (..) can only be used on class instances"}
		}
		// Evaluate the value
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		// Check if field exists and is mutable
		field, exists := instance.Class.Fields[node.Name]
		if !exists {
			return &Error{Message: fmt.Sprintf("unknown field: %s", node.Name)}
		}
		if !field.Mutable {
			return &Error{Message: fmt.Sprintf("cannot assign to immutable field: %s", node.Name)}
		}
		// Assign the value
		instance.Fields[node.Name] = val
		// Return the instance for chaining
		return instance
	case *parser.IfStatement:
		return evalIfStatement(node, env)
	case *parser.ForInStatement:
		return evalForInStatement(node, env)
	case *parser.ForConditionStatement:
		return evalForConditionStatement(node, env)
	case *parser.BreakStatement:
		// Check postfix condition
		if node.Condition != nil {
			cond := Eval(node.Condition, env)
			if IsError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				return NULL // condition not met, don't break
			}
		}
		return BREAK
	case *parser.ContinueStatement:
		// Check postfix condition
		if node.Condition != nil {
			cond := Eval(node.Condition, env)
			if IsError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				return NULL // condition not met, don't continue
			}
		}
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
	case *parser.ClassDeclaration:
		return evalClassDeclaration(node, env)
	case *parser.ImplDeclaration:
		return evalImplDeclaration(node, env)
	case *parser.ThisExpression:
		return evalThis(env)
	case *parser.MemberAssignmentStatement:
		return evalMemberAssignment(node, env)
	case *parser.IndexAssignmentStatement:
		return evalIndexAssignment(node, env)
	case *parser.SpecDeclaration:
		return evalSpecDeclaration(node, env)
	case *parser.ImplForDeclaration:
		return evalImplForDeclaration(node, env)
	case *parser.SelfExpression:
		return evalSelf(env)
	case *parser.ReferenceExpression:
		return evalReferenceExpression(node, env)
	case *parser.AssertStatement:
		return evalAssertStatement(node, env)
	case *parser.UnsafeBlock:
		return evalUnsafeBlock(node, env)
	case *parser.AsmExpression:
		return evalAsmExpression(node, env)
	case *parser.EnumDeclaration:
		return evalEnumDeclaration(node, env)
	case *parser.MatchStatement:
		return evalMatchStatement(node, env)
	case *parser.WildcardPattern:
		return &String{Value: "_"} // Wildcard marker
	case *parser.IsExpression:
		return evalIsExpression(node, env)
	case *parser.RangeExpression:
		return evalRangeExpression(node, env)
	case *parser.ImportStatement:
		return evalImportStatement(node, env)
	case *parser.WithStatement:
		return evalWithStatement(node, env)
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
	return applyFunctionWithNamedArgs(fn, args, nil)
}

func applyFunctionWithNamedArgs(fn Object, args []Object, namedArgs map[string]Object) Object {
	switch f := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnvWithNamedArgs(f, args, namedArgs)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Class:
		return instantiateClassWithNamedArgs(f, args, namedArgs)
	case *BoundMethod:
		return applyMethodWithNamedArgs(f, args, namedArgs)
	case *Builtin:
		return f.Fn(args...)
	case *BoundBuiltin:
		return f.Fn(f.Receiver, args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	return extendFunctionEnvWithNamedArgs(fn, args, nil)
}

func extendFunctionEnvWithNamedArgs(fn *Function, args []Object, namedArgs map[string]Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		paramName := param.Name.Value
		// Check if this parameter was provided as named argument
		if namedArgs != nil {
			if val, ok := namedArgs[paramName]; ok {
				env.Declare(paramName, val, false)
				continue
			}
		}
		// Otherwise use positional argument
		if i < len(args) {
			env.Declare(paramName, args[i], false)
		} else if param.Default != nil {
			// Evaluate default value in function's closure environment
			defaultVal := Eval(param.Default, fn.Env)
			env.Declare(paramName, defaultVal, false)
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

func evalCompoundAssignmentStatement(cas *parser.CompoundAssignmentStatement, env *Environment) Object {
	// Get current value
	currentVal, ok := env.Get(cas.Name.Value)
	if !ok {
		return newError("identifier not found: %s", cas.Name.Value)
	}

	// Evaluate the right-hand side
	rightVal := Eval(cas.Value, env)
	if IsError(rightVal) {
		return rightVal
	}

	// Apply the operation
	result := evalInfixExpression(cas.Operator, currentVal, rightVal)
	if IsError(result) {
		return result
	}

	// Assign the result
	if err := env.Assign(cas.Name.Value, result); err != nil {
		return newError("%s", err.Error())
	}

	return result
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
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: -right.Value}
	case *Float:
		return &Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
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
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		// Promote integer to float
		rightFloat := float64(right.(*Integer).Value)
		return evalFloatInfixExpression(operator, left, &Float{Value: rightFloat})
	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		// Promote integer to float
		leftFloat := float64(left.(*Integer).Value)
		return evalFloatInfixExpression(operator, &Float{Value: leftFloat}, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == ENUM_VALUE_OBJ && right.Type() == ENUM_VALUE_OBJ:
		return evalEnumInfixExpression(operator, left, right)
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

func evalEnumInfixExpression(operator string, left, right Object) Object {
	leftEnum := left.(*EnumValue)
	rightEnum := right.(*EnumValue)

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(
			leftEnum.EnumName == rightEnum.EnumName &&
				leftEnum.VariantName == rightEnum.VariantName,
		)
	case "!=":
		return nativeBoolToBooleanObject(
			leftEnum.EnumName != rightEnum.EnumName ||
				leftEnum.VariantName != rightEnum.VariantName,
		)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Float).Value
	rightVal := right.(*Float).Value

	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "**":
		return &Float{Value: floatPow(leftVal, rightVal)}
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

func floatPow(base, exp float64) float64 {
	// Use math.Pow for float exponentiation
	result := 1.0
	if exp == 0 {
		return 1.0
	}
	if exp < 0 {
		base = 1 / base
		exp = -exp
	}
	for exp >= 1 {
		if int(exp)%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
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

// NewEnvironmentWithBuiltins creates a new Environment with all builtins pre-declared.
func NewEnvironmentWithBuiltins() *Environment {
	env := NewEnvironment()
	// Register builtin functions
	for name, builtin := range Builtins {
		env.Declare(name, builtin, false)
	}
	// Register builtin modules
	for name, module := range BuiltinModules {
		env.Declare(name, module, false)
	}
	return env
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
	case *Range:
		return evalForInRange(fis, obj, env)
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

func evalForInRange(fis *parser.ForInStatement, r *Range, env *Environment) Object {
	var result Object = NULL

	loopEnv := NewEnclosedEnvironment(env)

	end := r.End
	if r.Inclusive {
		end++
	}

	for i := r.Start; i < end; i++ {
		loopEnv.Declare(fis.Variable.Value, &Integer{Value: i}, false)

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

	return &String{Value: string([]byte{strObj.Value[idx]})}
}

func evalMemberExpression(obj Object, member string) Object {
	switch o := obj.(type) {
	case *EnumType:
		// Enum variant access: EnumType.VARIANT
		for _, variant := range o.Variants {
			if variant == member {
				return &EnumValue{EnumName: o.Name, VariantName: member}
			}
		}
		return newError("enum '%s' has no variant '%s'", o.Name, member)
	case *ResultOk:
		if member == "value" {
			return o.Value
		}
		return newError("Ok has no member '%s'", member)
	case *ResultErr:
		if member == "error" {
			return o.Error
		}
		return newError("Err has no member '%s'", member)
	case *Module:
		if method, ok := o.Methods[member]; ok {
			return method
		}
		return newError("module '%s' has no method '%s'", o.Name, member)
	case *UserModule:
		// Access member from the module's environment
		val, ok := o.Env.Get(member)
		if !ok {
			return newError("module '%s' has no member '%s'", o.Name, member)
		}
		return val
	case *Map:
		// Check for builtin methods first
		if methodFn := GetMapMethod(member); methodFn != nil {
			return &BoundBuiltin{Name: member, Receiver: o, Fn: methodFn}
		}
		// Try to access map with string key
		key := &String{Value: member}
		pair, ok := o.Pairs[key.HashKey()]
		if !ok {
			return NULL
		}
		return pair.Value
	case *List:
		// Check for builtin methods first
		if methodFn := GetListMethod(member); methodFn != nil {
			return &BoundBuiltin{Name: member, Receiver: o, Fn: methodFn}
		}
		// List built-in properties
		switch member {
		case "length":
			return &Integer{Value: int64(len(o.Elements))}
		}
	case *File:
		// Check for File methods (read, write, seek, tell, close)
		if methodFn := GetFileMethod(member); methodFn != nil {
			return &BoundBuiltin{Name: member, Receiver: o, Fn: methodFn}
		}
		return newError("file has no method '%s'", member)
	case *String:
		// String built-in methods
		switch member {
		case "length":
			return &Integer{Value: int64(len(o.Value))}
		}
	case *Instance:
		// First check for field
		if val, ok := o.Fields[member]; ok {
			return val
		}
		// Then check for method
		if method, ok := o.Class.Methods[member]; ok {
			return &BoundMethod{Instance: o, Method: method}
		}
		return newError("no member '%s' on instance of %s", member, o.Class.Name)
	case *Class:
		// Static method access
		if method, ok := o.StaticMethods[member]; ok {
			return &BoundMethod{Instance: nil, Method: method}
		}
		return newError("no static member '%s' on class %s", member, o.Name)
	}
	return newError("no member '%s' on type %s", member, obj.Type())
}

func evalClassDeclaration(cd *parser.ClassDeclaration, env *Environment) Object {
	class := &Class{
		Name:          cd.Name.Value,
		Fields:        make(map[string]*ClassField),
		Methods:       make(map[string]*ClassMethod),
		StaticMethods: make(map[string]*ClassMethod),
		Implements:    make(map[string]*Spec),
	}

	for _, field := range cd.Fields {
		var defaultVal Object
		if field.Default != nil {
			defaultVal = Eval(field.Default, env)
			if IsError(defaultVal) {
				return defaultVal
			}
		}

		class.Fields[field.Name.Value] = &ClassField{
			Name:    field.Name.Value,
			Default: defaultVal,
			Public:  field.Public,
			Mutable: field.Mutable,
		}
	}

	env.Declare(cd.Name.Value, class, false)
	return class
}

func evalImplDeclaration(id *parser.ImplDeclaration, env *Environment) Object {
	classObj, ok := env.Get(id.Class.Value)
	if !ok {
		return newError("class not found: %s", id.Class.Value)
	}

	class, ok := classObj.(*Class)
	if !ok {
		return newError("%s is not a class", id.Class.Value)
	}

	for _, method := range id.Methods {
		cm := &ClassMethod{
			Name:       method.Name.Value,
			Parameters: []string{},
			Body:       method.Body,
			Public:     method.Public,
			Static:     method.Static,
			Mutable:    method.Mutable,
			Env:        env,
		}

		for _, param := range method.Parameters {
			cm.Parameters = append(cm.Parameters, param.Name.Value)
		}

		if method.Static {
			class.StaticMethods[method.Name.Value] = cm
		} else {
			class.Methods[method.Name.Value] = cm
		}
	}

	return NULL
}

func evalThis(env *Environment) Object {
	val, ok := env.Get("this")
	if !ok {
		return newError("'this' used outside of method")
	}
	return val
}

func instantiateClass(class *Class, args []Object) Object {
	instance := &Instance{
		Class:  class,
		Fields: make(map[string]Object),
	}

	// Initialize fields with defaults
	for name, field := range class.Fields {
		if field.Default != nil {
			instance.Fields[name] = field.Default
		} else {
			instance.Fields[name] = NULL
		}
	}

	// Call init method if it exists
	if initMethod, ok := class.Methods["init"]; ok {
		applyMethodOnInstance(initMethod, instance, args)
	}

	return instance
}

func applyMethod(bm *BoundMethod, args []Object) Object {
	if bm.Instance == nil {
		// Static method
		return applyStaticMethod(bm.Method, args)
	}
	return applyMethodOnInstance(bm.Method, bm.Instance, args)
}

func applyMethodOnInstance(method *ClassMethod, instance *Instance, args []Object) Object {
	env := NewEnclosedEnvironment(method.Env)

	// Bind 'this' to the instance
	env.Declare("this", instance, false)

	// Bind parameters
	for i, param := range method.Parameters {
		if i < len(args) {
			env.Declare(param, args[i], false)
		}
	}

	body := method.Body.(*parser.BlockStatement)
	evaluated := Eval(body, env)
	return unwrapReturnValue(evaluated)
}

func applyStaticMethod(method *ClassMethod, args []Object) Object {
	env := NewEnclosedEnvironment(method.Env)

	// Bind parameters
	for i, param := range method.Parameters {
		if i < len(args) {
			env.Declare(param, args[i], false)
		}
	}

	body := method.Body.(*parser.BlockStatement)
	evaluated := Eval(body, env)
	return unwrapReturnValue(evaluated)
}

// instantiateClassWithNamedArgs creates a new instance and calls init with named args
func instantiateClassWithNamedArgs(class *Class, args []Object, namedArgs map[string]Object) Object {
	instance := &Instance{
		Class:  class,
		Fields: make(map[string]Object),
	}

	// Initialize fields with defaults
	for name, field := range class.Fields {
		if field.Default != nil {
			instance.Fields[name] = field.Default
		} else {
			instance.Fields[name] = NULL
		}
	}

	// Call init method if it exists with named args support
	if initMethod, ok := class.Methods["init"]; ok {
		applyMethodOnInstanceWithNamedArgs(initMethod, instance, args, namedArgs)
	}

	return instance
}

// applyMethodWithNamedArgs applies a bound method with named argument support
func applyMethodWithNamedArgs(bm *BoundMethod, args []Object, namedArgs map[string]Object) Object {
	if bm.Instance == nil {
		// Static method - named args not yet supported for static methods
		return applyStaticMethod(bm.Method, args)
	}
	return applyMethodOnInstanceWithNamedArgs(bm.Method, bm.Instance, args, namedArgs)
}

// applyMethodOnInstanceWithNamedArgs applies a method with named args support
func applyMethodOnInstanceWithNamedArgs(method *ClassMethod, instance *Instance, args []Object, namedArgs map[string]Object) Object {
	env := NewEnclosedEnvironment(method.Env)

	// Bind 'this' to the instance
	env.Declare("this", instance, false)

	// Bind parameters - check named args first, then positional
	for i, param := range method.Parameters {
		// Check if this parameter was provided as named argument
		if namedArgs != nil {
			if val, ok := namedArgs[param]; ok {
				env.Declare(param, val, false)
				continue
			}
		}
		// Otherwise use positional argument
		if i < len(args) {
			env.Declare(param, args[i], false)
		}
	}

	body := method.Body.(*parser.BlockStatement)
	evaluated := Eval(body, env)
	return unwrapReturnValue(evaluated)
}

func evalMemberAssignment(mas *parser.MemberAssignmentStatement, env *Environment) Object {
	obj := Eval(mas.Object, env)
	if IsError(obj) {
		return obj
	}

	value := Eval(mas.Value, env)
	if IsError(value) {
		return value
	}

	switch o := obj.(type) {
	case *Instance:
		// Check if field exists and is mutable
		field, ok := o.Class.Fields[mas.Member.Value]
		if !ok {
			return newError("no field '%s' on class %s", mas.Member.Value, o.Class.Name)
		}
		if !field.Mutable {
			return newError("cannot assign to immutable field '%s'", mas.Member.Value)
		}
		o.Fields[mas.Member.Value] = value
		return value
	case *Map:
		// Allow map member assignment via dot notation
		key := &String{Value: mas.Member.Value}
		o.Pairs[key.HashKey()] = MapPair{Key: key, Value: value}
		return value
	default:
		return newError("cannot assign member on type %s", obj.Type())
	}
}

func evalIndexAssignment(ias *parser.IndexAssignmentStatement, env *Environment) Object {
	left := Eval(ias.Left, env)
	if IsError(left) {
		return left
	}

	index := Eval(ias.Index, env)
	if IsError(index) {
		return index
	}

	value := Eval(ias.Value, env)
	if IsError(value) {
		return value
	}

	switch o := left.(type) {
	case *List:
		idx, ok := index.(*Integer)
		if !ok {
			return newError("list index must be integer, got %s", index.Type())
		}
		if idx.Value < 0 || idx.Value >= int64(len(o.Elements)) {
			return newError("list index out of bounds: %d", idx.Value)
		}
		o.Elements[idx.Value] = value
		return value
	case *Map:
		hashKey, ok := index.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", index.Type())
		}
		o.Pairs[hashKey.HashKey()] = MapPair{Key: index, Value: value}
		return value
	default:
		return newError("cannot assign index on type %s", left.Type())
	}
}

func evalSpecDeclaration(sd *parser.SpecDeclaration, env *Environment) Object {
	spec := &Spec{
		Name:    sd.Name.Value,
		Methods: make(map[string]*SpecMethod),
	}

	for _, method := range sd.Methods {
		params := []string{}
		for _, param := range method.Parameters {
			params = append(params, param.Value)
		}

		spec.Methods[method.Name.Value] = &SpecMethod{
			Name:       method.Name.Value,
			Parameters: params,
			Public:     method.Public,
			Mutable:    method.Mutable,
		}
	}

	env.Declare(sd.Name.Value, spec, false)
	return spec
}

func evalImplForDeclaration(ifd *parser.ImplForDeclaration, env *Environment) Object {
	// Get the class
	classObj, ok := env.Get(ifd.Class.Value)
	if !ok {
		return newError("class not found: %s", ifd.Class.Value)
	}

	class, ok := classObj.(*Class)
	if !ok {
		return newError("%s is not a class", ifd.Class.Value)
	}

	// Get the spec
	specObj, ok := env.Get(ifd.Spec.Value)
	if !ok {
		return newError("spec not found: %s", ifd.Spec.Value)
	}

	spec, ok := specObj.(*Spec)
	if !ok {
		return newError("%s is not a spec", ifd.Spec.Value)
	}

	// Add methods to the class
	for _, method := range ifd.Methods {
		cm := &ClassMethod{
			Name:       method.Name.Value,
			Parameters: []string{},
			Body:       method.Body,
			Public:     method.Public,
			Static:     method.Static,
			Mutable:    method.Mutable,
			Env:        env,
		}

		for _, param := range method.Parameters {
			cm.Parameters = append(cm.Parameters, param.Name.Value)
		}

		if method.Static {
			class.StaticMethods[method.Name.Value] = cm
		} else {
			class.Methods[method.Name.Value] = cm
		}
	}

	// Verify all spec methods are implemented
	for name, specMethod := range spec.Methods {
		classMethod, ok := class.Methods[name]
		if !ok {
			return newError("class %s does not implement method '%s' from spec %s",
				class.Name, name, spec.Name)
		}

		// Check parameter count matches
		if len(classMethod.Parameters) != len(specMethod.Parameters) {
			return newError("method '%s' has %d parameters, spec requires %d",
				name, len(classMethod.Parameters), len(specMethod.Parameters))
		}
	}

	// Record that class implements this spec
	class.Implements[spec.Name] = spec

	return NULL
}

func evalSelf(env *Environment) Object {
	// Self is used in spec contexts to refer to the implementing type
	// In this bootstrap, we just return a placeholder
	return &String{Value: "Self"}
}

func evalReferenceExpression(node *parser.ReferenceExpression, env *Environment) Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	return &Reference{Value: &val}
}

func evalAssertStatement(node *parser.AssertStatement, env *Environment) Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if !isTruthy(condition) {
		msg := "assertion failed"
		if node.Message != nil {
			msgObj := Eval(node.Message, env)
			if str, ok := msgObj.(*String); ok {
				msg = str.Value
			} else {
				msg = msgObj.Inspect()
			}
		}
		return &Error{Message: msg}
	}

	return NULL
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// evalUnsafeBlock evaluates an unsafe block.
// Currently, this just evaluates the body - in a future compiler,
// this would enable access to low-level operations.
func evalUnsafeBlock(ub *parser.UnsafeBlock, env *Environment) Object {
	// Create a new environment to mark that we're in an unsafe context
	unsafeEnv := NewEnclosedEnvironment(env)
	unsafeEnv.Declare("__unsafe__", TRUE, false)

	return Eval(ub.Body, unsafeEnv)
}

// evalWithStatement evaluates a with statement for automatic resource management.
// The resource is automatically closed after the body executes.
func evalWithStatement(ws *parser.WithStatement, env *Environment) Object {
	// Evaluate the resource expression
	resource := Eval(ws.Resource, env)
	if IsError(resource) {
		return resource
	}

	// Create a new scope and bind the resource
	innerEnv := NewEnclosedEnvironment(env)
	innerEnv.Declare(ws.Name.Value, resource, false)

	// Evaluate the body
	result := Eval(ws.Body, innerEnv)

	// Auto-close: call close() if resource is a File
	if f, ok := resource.(*File); ok {
		if file, ok := f.Handle.(*os.File); ok {
			file.Close()
		}
	}

	return result
}

// AsmFunction is the signature for asm-callable functions.
type AsmFunction func(args ...Object) Object

// AsmRegistry maps function names to their Go implementations.
var AsmRegistry = map[string]AsmFunction{
	// sys functions
	"sys_os":   func(args ...Object) Object { return sysOs(args...) },
	"sys_arch": func(args ...Object) Object { return sysArch(args...) },
	"sys_args": func(args ...Object) Object { return sysArgs(args...) },
	"sys_exit": func(args ...Object) Object { return sysExit(args...) },
	"sys_env":  func(args ...Object) Object { return sysEnv(args...) },

	// io functions
	"file_open":   func(args ...Object) Object { return ioOpen(args...) },
	"file_read":   func(args ...Object) Object { return ioRead(args...) },
	"file_read_n": func(args ...Object) Object { return asmFileReadN(args...) },
	"file_write":  func(args ...Object) Object { return ioWrite(args...) },
	"file_close":  func(args ...Object) Object { return ioClose(args...) },
	"file_seek":   func(args ...Object) Object { return asmFileSeek(args...) },
	"file_tell":   func(args ...Object) Object { return asmFileTell(args...) },
	"file_exists": func(args ...Object) Object { return ioExists(args...) },
	"read_file":   func(args ...Object) Object { return ioReadFile(args...) },
	"write_file":  func(args ...Object) Object { return ioWriteFile(args...) },
	"byte_chr":    func(args ...Object) Object { return asmByteChr(args...) },

	// str functions
	"str_split":       func(args ...Object) Object { return strSplit(args...) },
	"str_join":        func(args ...Object) Object { return strJoin(args...) },
	"str_trim":        func(args ...Object) Object { return strTrim(args...) },
	"str_find":        func(args ...Object) Object { return strFind(args...) },
	"str_replace":     func(args ...Object) Object { return strReplace(args...) },
	"str_substring":   func(args ...Object) Object { return strSubstring(args...) },
	"str_starts_with": func(args ...Object) Object { return strStartsWith(args...) },
	"str_ends_with":   func(args ...Object) Object { return strEndsWith(args...) },
	"str_upper":       func(args ...Object) Object { return strUpper(args...) },
	"str_lower":       func(args ...Object) Object { return strLower(args...) },
	"str_contains":    func(args ...Object) Object { return strContains(args...) },

	// char functions
	"char_ord":      func(args ...Object) Object { return charOrd(args...) },
	"char_chr":      func(args ...Object) Object { return charChr(args...) },
	"char_is_digit": func(args ...Object) Object { return charIsDigit(args...) },
	"char_is_alpha": func(args ...Object) Object { return charIsAlpha(args...) },
	"char_is_space": func(args ...Object) Object { return charIsSpace(args...) },
	"char_is_alnum": func(args ...Object) Object { return charIsAlnum(args...) },
}

// evalAsmExpression evaluates an asm expression.
func evalAsmExpression(ae *parser.AsmExpression, env *Environment) Object {
	// Check if we're in an unsafe context
	if _, ok := env.Get("__unsafe__"); !ok {
		return newError("asm() can only be used inside an unsafe block")
	}

	// Look up the function in the registry
	fn, ok := AsmRegistry[ae.Function]
	if !ok {
		return newError("unknown asm function: %s", ae.Function)
	}

	// Evaluate arguments
	args := evalExpressions(ae.Args, env)
	if len(args) == 1 && IsError(args[0]) {
		return args[0]
	}

	return fn(args...)
}

func evalEnumDeclaration(ed *parser.EnumDeclaration, env *Environment) Object {
	enumType := &EnumType{
		Name:     ed.Name.Value,
		Variants: ed.Variants,
	}

	env.Declare(ed.Name.Value, enumType, false)
	return enumType
}

func evalMatchStatement(ms *parser.MatchStatement, env *Environment) Object {
	value := Eval(ms.Value, env)
	if IsError(value) {
		return value
	}

	for _, arm := range ms.Arms {
		matched := false

		// Check all patterns in this arm (| alternatives)
		for _, pattern := range arm.Patterns {
			if matchPattern(pattern, value, env) {
				matched = true
				break
			}
		}

		if matched {
			// Check guard condition if present
			if arm.Guard != nil {
				guardResult := Eval(arm.Guard, env)
				if IsError(guardResult) {
					return guardResult
				}
				if !isTruthy(guardResult) {
					continue // Guard failed, try next arm
				}
			}

			// Execute the body
			result := Eval(arm.Body, env)
			return unwrapReturnValue(result)
		}
	}

	// No arm matched
	return NULL
}

func matchPattern(pattern parser.Expression, value Object, env *Environment) bool {
	switch p := pattern.(type) {
	case *parser.WildcardPattern:
		// Wildcard matches everything
		return true

	case *parser.IntegerLiteral:
		if intVal, ok := value.(*Integer); ok {
			return intVal.Value == p.Value
		}
		return false

	case *parser.StringLiteral:
		if strVal, ok := value.(*String); ok {
			return strVal.Value == p.Value
		}
		return false

	case *parser.BooleanLiteral:
		if boolVal, ok := value.(*Boolean); ok {
			return boolVal.Value == p.Value
		}
		return false

	case *parser.NilLiteral:
		return value == NULL

	case *parser.Identifier:
		// Check if this is Ok or Err
		if p.Value == "Ok" {
			_, ok := value.(*ResultOk)
			return ok
		}
		if p.Value == "Err" {
			_, ok := value.(*ResultErr)
			return ok
		}
		// Otherwise evaluate as expression and compare
		patternVal := Eval(pattern, env)
		if IsError(patternVal) {
			return false
		}
		return objectsEqual(patternVal, value)

	case *parser.MemberExpression:
		// Enum variant access: EnumType.VARIANT
		patternVal := Eval(pattern, env)
		if IsError(patternVal) {
			return false
		}

		// Compare enum values
		if patternEnum, ok := patternVal.(*EnumValue); ok {
			if valueEnum, ok := value.(*EnumValue); ok {
				return patternEnum.EnumName == valueEnum.EnumName &&
					patternEnum.VariantName == valueEnum.VariantName
			}
		}

		return objectsEqual(patternVal, value)

	default:
		// Evaluate the pattern and compare
		patternVal := Eval(pattern, env)
		if IsError(patternVal) {
			return false
		}
		return objectsEqual(patternVal, value)
	}
}

func evalInterpolatedString(is *parser.InterpolatedString, env *Environment) Object {
	var result string

	for _, part := range is.Parts {
		val := Eval(part, env)
		if IsError(val) {
			return val
		}
		result += val.Inspect()
	}

	return &String{Value: result}
}

func evalRangeExpression(re *parser.RangeExpression, env *Environment) Object {
	start := Eval(re.Start, env)
	if IsError(start) {
		return start
	}

	end := Eval(re.End, env)
	if IsError(end) {
		return end
	}

	startInt, ok := start.(*Integer)
	if !ok {
		return newError("range start must be an integer, got %s", start.Type())
	}

	endInt, ok := end.(*Integer)
	if !ok {
		return newError("range end must be an integer, got %s", end.Type())
	}

	return &Range{
		Start:     startInt.Value,
		End:       endInt.Value,
		Inclusive: re.Inclusive,
	}
}

func evalIsExpression(ie *parser.IsExpression, env *Environment) Object {
	left := Eval(ie.Left, env)
	if IsError(left) {
		return left
	}

	// Get the type name from the right side
	var typeName string
	switch t := ie.Right.(type) {
	case *parser.Identifier:
		typeName = t.Value
	case *parser.NilLiteral:
		// Handle "x is nil"
		return nativeBoolToBooleanObject(left == NULL)
	case *parser.MemberExpression:
		// For enum variant checking: value is EnumType.VARIANT
		rightVal := Eval(ie.Right, env)
		if IsError(rightVal) {
			return rightVal
		}
		if enumVal, ok := rightVal.(*EnumValue); ok {
			if leftEnum, ok := left.(*EnumValue); ok {
				return nativeBoolToBooleanObject(
					leftEnum.EnumName == enumVal.EnumName &&
						leftEnum.VariantName == enumVal.VariantName,
				)
			}
			return FALSE
		}
		return newError("invalid type in is expression")
	default:
		// Try to check if the right side evaluates to a type
		return newError("invalid type in is expression: expected identifier or member expression")
	}

	// Check built-in type names
	switch typeName {
	case "Ok":
		_, ok := left.(*ResultOk)
		return nativeBoolToBooleanObject(ok)
	case "Err":
		_, ok := left.(*ResultErr)
		return nativeBoolToBooleanObject(ok)
	case "int", "Integer":
		_, ok := left.(*Integer)
		return nativeBoolToBooleanObject(ok)
	case "float", "Float":
		_, ok := left.(*Float)
		return nativeBoolToBooleanObject(ok)
	case "string", "String":
		_, ok := left.(*String)
		return nativeBoolToBooleanObject(ok)
	case "bool", "Boolean":
		_, ok := left.(*Boolean)
		return nativeBoolToBooleanObject(ok)
	case "list", "List":
		_, ok := left.(*List)
		return nativeBoolToBooleanObject(ok)
	case "map", "Map":
		_, ok := left.(*Map)
		return nativeBoolToBooleanObject(ok)
	case "nil", "Nil":
		return nativeBoolToBooleanObject(left == NULL)
	default:
		// Check if it's an instance of a class
		if instance, ok := left.(*Instance); ok {
			return nativeBoolToBooleanObject(instance.Class.Name == typeName)
		}
		// Check if it's an enum type
		if enumVal, ok := left.(*EnumValue); ok {
			return nativeBoolToBooleanObject(enumVal.EnumName == typeName)
		}
		return FALSE
	}
}
