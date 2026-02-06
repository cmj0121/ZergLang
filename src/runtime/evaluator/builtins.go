package evaluator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Builtins is the map of all builtin functions.
var Builtins = map[string]*Builtin{
	"print":  {Name: "print", Fn: builtinPrint},
	"len":    {Name: "len", Fn: builtinLen},
	"string": {Name: "string", Fn: builtinStr},
	"int":    {Name: "int", Fn: builtinInt},
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

// listJoin joins list elements with a separator.
func listJoin(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 1 {
		return newError("join() takes exactly 1 argument (%d given)", len(args))
	}

	sep, ok := args[0].(*String)
	if !ok {
		return newError("join() argument must be a string, not %s", args[0].Type())
	}

	strs := make([]string, len(list.Elements))
	for i, el := range list.Elements {
		if str, ok := el.(*String); ok {
			strs[i] = str.Value
		} else {
			strs[i] = el.Inspect()
		}
	}

	return &String{Value: strings.Join(strs, sep.Value)}
}

// listSlice returns a sublist from start to end index.
func listSlice(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 2 {
		return newError("slice() takes exactly 2 arguments (%d given)", len(args))
	}

	start, ok := args[0].(*Integer)
	if !ok {
		return newError("slice() start must be an integer, not %s", args[0].Type())
	}

	end, ok := args[1].(*Integer)
	if !ok {
		return newError("slice() end must be an integer, not %s", args[1].Type())
	}

	length := int64(len(list.Elements))
	startIdx := start.Value
	endIdx := end.Value

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > length {
		endIdx = length
	}
	if startIdx > endIdx {
		return &List{Elements: []Object{}}
	}

	newElements := make([]Object, endIdx-startIdx)
	copy(newElements, list.Elements[startIdx:endIdx])
	return &List{Elements: newElements}
}

// listIndex returns the index of an item, or -1 if not found.
func listIndex(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 1 {
		return newError("index() takes exactly 1 argument (%d given)", len(args))
	}

	target := args[0]
	for i, el := range list.Elements {
		if objectsEqual(el, target) {
			return &Integer{Value: int64(i)}
		}
	}

	return &Integer{Value: -1}
}

// listReverse returns a new list with elements in reverse order.
func listReverse(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 0 {
		return newError("reverse() takes no arguments (%d given)", len(args))
	}

	length := len(list.Elements)
	newElements := make([]Object, length)
	for i, el := range list.Elements {
		newElements[length-1-i] = el
	}

	return &List{Elements: newElements}
}

// listSort returns a new list with elements sorted.
func listSort(receiver Object, args ...Object) Object {
	list := receiver.(*List)
	if len(args) != 0 {
		return newError("sort() takes no arguments (%d given)", len(args))
	}

	// Make a copy
	newElements := make([]Object, len(list.Elements))
	copy(newElements, list.Elements)

	// Sort using Inspect() for comparison
	sort.Slice(newElements, func(i, j int) bool {
		// Try numeric comparison first
		if intI, okI := newElements[i].(*Integer); okI {
			if intJ, okJ := newElements[j].(*Integer); okJ {
				return intI.Value < intJ.Value
			}
		}
		// Fall back to string comparison
		return newElements[i].Inspect() < newElements[j].Inspect()
	})

	return &List{Elements: newElements}
}

// objectsEqual checks if two objects are equal.
func objectsEqual(a, b Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch aVal := a.(type) {
	case *Integer:
		return aVal.Value == b.(*Integer).Value
	case *String:
		return aVal.Value == b.(*String).Value
	case *Boolean:
		return aVal.Value == b.(*Boolean).Value
	case *Null:
		return true
	default:
		return a == b
	}
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
	case "join":
		return listJoin
	case "slice":
		return listSlice
	case "index":
		return listIndex
	case "reverse":
		return listReverse
	case "sort":
		return listSort
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
