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
	}

	return nil
}

func evalProgram(program *parser.Program, env *Environment) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
	}

	return result
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
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]*binding)}
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (Object, bool) {
	b, ok := e.store[name]
	if !ok {
		return nil, false
	}
	return b.value, true
}

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
