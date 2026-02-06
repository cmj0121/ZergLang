package evaluator

import "fmt"

// ObjectType represents the type of a runtime object.
type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	STRING_OBJ       ObjectType = "STRING"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	BREAK_OBJ        ObjectType = "BREAK"
	CONTINUE_OBJ     ObjectType = "CONTINUE"
	LIST_OBJ         ObjectType = "LIST"
	MAP_OBJ          ObjectType = "MAP"
	CLASS_OBJ        ObjectType = "CLASS"
	INSTANCE_OBJ     ObjectType = "INSTANCE"
	SPEC_OBJ         ObjectType = "SPEC"
	REFERENCE_OBJ    ObjectType = "REFERENCE"
	ERROR_OBJ        ObjectType = "ERROR"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
)

// Object is the interface for all runtime values.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer value.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// String represents a string value.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean represents a boolean value.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null represents the absence of a value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "nil" }

// Singleton null and boolean objects for efficiency.
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// ReturnValue wraps a return value to propagate it up the call stack.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// BreakSignal signals a break statement in a loop.
type BreakSignal struct{}

func (bs *BreakSignal) Type() ObjectType { return BREAK_OBJ }
func (bs *BreakSignal) Inspect() string  { return "break" }

// ContinueSignal signals a continue statement in a loop.
type ContinueSignal struct{}

func (cs *ContinueSignal) Type() ObjectType { return CONTINUE_OBJ }
func (cs *ContinueSignal) Inspect() string  { return "continue" }

// Singleton break and continue signals.
var (
	BREAK    = &BreakSignal{}
	CONTINUE = &ContinueSignal{}
)

// List represents a list value.
type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out string
	out += "["
	for i, el := range l.Elements {
		if i > 0 {
			out += ", "
		}
		out += el.Inspect()
	}
	out += "]"
	return out
}

// HashKey is used for map keys.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Hashable is the interface for objects that can be used as map keys.
type Hashable interface {
	HashKey() HashKey
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	var h uint64 = 5381
	for i := 0; i < len(s.Value); i++ {
		h = (h << 5) + h + uint64(s.Value[i])
	}
	return HashKey{Type: s.Type(), Value: h}
}

func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	}
	return HashKey{Type: b.Type(), Value: 0}
}

// MapPair represents a key-value pair in a map.
type MapPair struct {
	Key   Object
	Value Object
}

// Map represents a map value.
type Map struct {
	Pairs map[HashKey]MapPair
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out string
	out += "{"
	first := true
	for _, pair := range m.Pairs {
		if !first {
			out += ", "
		}
		out += pair.Key.Inspect() + ": " + pair.Value.Inspect()
		first = false
	}
	out += "}"
	return out
}

// ClassField represents a field definition in a class.
type ClassField struct {
	Name    string
	Default Object
	Public  bool
	Mutable bool
}

// ClassMethod represents a method definition in a class.
type ClassMethod struct {
	Name       string
	Parameters []string
	Body       interface{} // *parser.BlockStatement
	Public     bool
	Static     bool
	Mutable    bool // mutable receiver
	Env        *Environment
}

// Class represents a class definition.
type Class struct {
	Name          string
	Fields        map[string]*ClassField
	Methods       map[string]*ClassMethod
	StaticMethods map[string]*ClassMethod
	Implements    map[string]*Spec // specs this class implements
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return fmt.Sprintf("<class %s>", c.Name) }

// Instance represents an instance of a class.
type Instance struct {
	Class  *Class
	Fields map[string]Object
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string  { return fmt.Sprintf("<%s instance>", i.Class.Name) }

// BoundMethod represents a method bound to an instance.
type BoundMethod struct {
	Instance *Instance // nil for static methods
	Method   *ClassMethod
}

func (bm *BoundMethod) Type() ObjectType { return FUNCTION_OBJ }
func (bm *BoundMethod) Inspect() string  { return fmt.Sprintf("<method %s>", bm.Method.Name) }

// SpecMethod represents a method signature in a spec.
type SpecMethod struct {
	Name       string
	Parameters []string
	Public     bool
	Mutable    bool
}

// Spec represents a spec (interface) definition.
type Spec struct {
	Name    string
	Methods map[string]*SpecMethod
}

func (s *Spec) Type() ObjectType { return SPEC_OBJ }
func (s *Spec) Inspect() string  { return fmt.Sprintf("<spec %s>", s.Name) }

// Reference represents a reference to an object.
type Reference struct {
	Value *Object
}

func (r *Reference) Type() ObjectType { return REFERENCE_OBJ }
func (r *Reference) Inspect() string  { return fmt.Sprintf("&%s", (*r.Value).Inspect()) }

// BuiltinFn is the signature for builtin functions.
type BuiltinFn func(args ...Object) Object

// Builtin represents a builtin function.
type Builtin struct {
	Name string
	Fn   BuiltinFn
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return fmt.Sprintf("<builtin %s>", b.Name) }
