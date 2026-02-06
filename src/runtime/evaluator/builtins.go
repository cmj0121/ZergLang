package evaluator

import (
	"fmt"
	"sort"
	"strconv"
)

// Builtins is the map of all builtin functions.
var Builtins = map[string]*Builtin{
	"print": {Name: "print", Fn: builtinPrint},
	"len":   {Name: "len", Fn: builtinLen},
	"str":   {Name: "str", Fn: builtinStr},
	"int":   {Name: "int", Fn: builtinInt},
}

// builtinPrint outputs arguments to stdout with a newline.
func builtinPrint(args ...Object) Object {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		values[i] = arg.Inspect()
	}
	fmt.Println(values...)
	return NULL
}

// builtinLen returns the length of a string, list, or map.
func builtinLen(args ...Object) Object {
	if len(args) != 1 {
		return newError("len() takes exactly 1 argument (%d given)", len(args))
	}

	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *List:
		return &Integer{Value: int64(len(arg.Elements))}
	case *Map:
		return &Integer{Value: int64(len(arg.Pairs))}
	default:
		return newError("len() argument must be string, list, or map, not %s", args[0].Type())
	}
}

// builtinStr converts any value to a string via Inspect().
func builtinStr(args ...Object) Object {
	if len(args) != 1 {
		return newError("str() takes exactly 1 argument (%d given)", len(args))
	}

	return &String{Value: args[0].Inspect()}
}

// builtinInt converts a string or bool to an integer.
func builtinInt(args ...Object) Object {
	if len(args) != 1 {
		return newError("int() takes exactly 1 argument (%d given)", len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return arg
	case *String:
		val, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return newError("int() argument is not a valid integer: %s", arg.Value)
		}
		return &Integer{Value: val}
	case *Boolean:
		if arg.Value {
			return &Integer{Value: 1}
		}
		return &Integer{Value: 0}
	default:
		return newError("int() argument must be string, integer, or bool, not %s", args[0].Type())
	}
}

// List method implementations

// listAppend returns a new list with item added.
func listAppend(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 1 {
		return newError("append() takes exactly 1 argument (%d given)", len(args))
	}

	newElements := make([]Object, len(list.Elements)+1)
	copy(newElements, list.Elements)
	newElements[len(list.Elements)] = args[0]
	return &List{Elements: newElements}
}

// listPop returns the last element of the list.
func listPop(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 0 {
		return newError("pop() takes no arguments (%d given)", len(args))
	}

	if len(list.Elements) == 0 {
		return newError("pop() from empty list")
	}

	return list.Elements[len(list.Elements)-1]
}

// listFilter returns a new list with elements that pass the predicate function.
func listFilter(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 1 {
		return newError("filter() takes exactly 1 argument (%d given)", len(args))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		return newError("filter() argument must be a function, not %s", args[0].Type())
	}

	var result []Object
	for _, el := range list.Elements {
		// Apply function to element
		extendedEnv := NewEnclosedEnvironment(fn.Env)
		if len(fn.Parameters) > 0 {
			extendedEnv.Declare(fn.Parameters[0].Name.Value, el, false)
		}
		evaluated := Eval(fn.Body, extendedEnv)
		evaluated = unwrapReturnValue(evaluated)

		if IsError(evaluated) {
			return evaluated
		}

		if isTruthy(evaluated) {
			result = append(result, el)
		}
	}

	return &List{Elements: result}
}

// listMap returns a new list with the function applied to each element.
func listMap(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 1 {
		return newError("map() takes exactly 1 argument (%d given)", len(args))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		return newError("map() argument must be a function, not %s", args[0].Type())
	}

	result := make([]Object, len(list.Elements))
	for i, el := range list.Elements {
		// Apply function to element
		extendedEnv := NewEnclosedEnvironment(fn.Env)
		if len(fn.Parameters) > 0 {
			extendedEnv.Declare(fn.Parameters[0].Name.Value, el, false)
		}
		evaluated := Eval(fn.Body, extendedEnv)
		evaluated = unwrapReturnValue(evaluated)

		if IsError(evaluated) {
			return evaluated
		}

		result[i] = evaluated
	}

	return &List{Elements: result}
}

// Map method implementations

// mapKeys returns a list of all keys in the map.
func mapKeys(receiver Object, args ...Object) Object {
	m := receiver.(*Map)
	if len(args) != 0 {
		return newError("keys() takes no arguments (%d given)", len(args))
	}

	keys := make([]Object, 0, len(m.Pairs))
	for _, pair := range m.Pairs {
		keys = append(keys, pair.Key)
	}

	// Sort for deterministic output
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Inspect() < keys[j].Inspect()
	})

	return &List{Elements: keys}
}

// mapValues returns a list of all values in the map.
func mapValues(receiver Object, args ...Object) Object {
	m := receiver.(*Map)
	if len(args) != 0 {
		return newError("values() takes no arguments (%d given)", len(args))
	}

	// Get keys and sort them for deterministic order
	keys := make([]Object, 0, len(m.Pairs))
	for _, pair := range m.Pairs {
		keys = append(keys, pair.Key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Inspect() < keys[j].Inspect()
	})

	values := make([]Object, 0, len(m.Pairs))
	for _, key := range keys {
		hashable := key.(Hashable)
		values = append(values, m.Pairs[hashable.HashKey()].Value)
	}

	return &List{Elements: values}
}

// mapContains checks if a key exists in the map.
func mapContains(receiver Object, args ...Object) Object {
	m := receiver.(*Map)
	if len(args) != 1 {
		return newError("contains() takes exactly 1 argument (%d given)", len(args))
	}

	key, ok := args[0].(Hashable)
	if !ok {
		return newError("contains() argument must be hashable, not %s", args[0].Type())
	}

	_, exists := m.Pairs[key.HashKey()]
	return nativeBoolToBooleanObject(exists)
}

// GetListMethod returns the builtin method for lists by name, or nil if not found.
func GetListMethod(name string) BoundBuiltinFn {
	switch name {
	case "append":
		return listAppend
	case "pop":
		return listPop
	case "filter":
		return listFilter
	case "map":
		return listMap
	default:
		return nil
	}
}

// GetMapMethod returns the builtin method for maps by name, or nil if not found.
func GetMapMethod(name string) BoundBuiltinFn {
	switch name {
	case "keys":
		return mapKeys
	case "values":
		return mapValues
	case "contains":
		return mapContains
	default:
		return nil
	}
}
