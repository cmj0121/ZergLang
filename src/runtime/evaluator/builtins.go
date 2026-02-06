package evaluator

import (
	"fmt"
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
