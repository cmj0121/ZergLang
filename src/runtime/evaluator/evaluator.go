package evaluator

import (
	"github.com/xrspace/zerglang/runtime/parser"
)

// Eval evaluates an AST node and returns the result.
func Eval(node parser.Node, env *Environment) Object {
	// TODO: Implement evaluation
	return nil
}

// Environment stores variable bindings.
type Environment struct {
	store map[string]Object
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set stores a variable in the environment.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
