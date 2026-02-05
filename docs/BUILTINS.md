# Built-ins

Built-in functions and statements are available in every module without importing. They are provided by the
Zerg runtime and cannot be shadowed or redefined.

## Overview

| Built-in         | Kind      | Description                               |
| ---------------- | --------- | ----------------------------------------- |
| `print(args...)` | function  | Write to stdout                           |
| `len(c)`         | function  | Length of string/list/map/set             |
| `int(v)`         | function  | Convert to int                            |
| `float(v)`       | function  | Convert to float                          |
| `str(v)`         | function  | Convert to string (calls `string()`)      |
| `input(prompt?)` | function  | Read line from stdin                      |
| `typeof(v)`      | function  | Return the type of v                      |
| `assert`         | statement | Raise `AssertionError` if condition false |
| `range`          | type      | Built-in iterable range type              |

## print

```txt
print(args...)
```

Writes its arguments to standard output followed by a newline. Multiple arguments are separated by a single
space. Each argument is converted to its string representation by calling its `string()` method (from the
`Stringable` spec).

```txt
print("hello")                  # hello
print("x =", 42)               # x = 42
print(1, 2, 3)                  # 1 2 3
print(true, nil)                # true nil
```

**Signature**: `fn print(args...: any)`

- **Parameters**: zero or more values of any type
- **Returns**: nothing

## len

```txt
len(c)
```

Returns the number of elements in a collection or the number of bytes in a string.

```txt
len([1, 2, 3])            # 3
len({"a": 1, "b": 2})     # 2
len({1, 2, 3})             # 3
len("hello")               # 5
len("")                    # 0
```

**Signature**: `fn len(c: any) ->int`

- **Parameters**: a `string`, `list`, `map`, or `set`
- **Returns**: `int` -- the number of elements (or bytes for strings)
- **Raises**: `TypeError` if the argument does not support `len`

## int

```txt
int(v)
```

Converts a value to `int`. Accepts `float` (truncates toward zero), `string` (parses decimal), and `bool`
(`true` = 1, `false` = 0).

```txt
int(3.9)       # 3
int("42")      # 42
int(true)      # 1
```

**Signature**: `fn int(v: any) ->int`

- **Raises**: `ValueError` if the value cannot be converted

## float

```txt
float(v)
```

Converts a value to `float`. Accepts `int`, `string` (parses floating-point), and `bool` (`true` = 1.0,
`false` = 0.0).

```txt
float(42)       # 42.0
float("3.14")   # 3.14
```

**Signature**: `fn float(v: any) ->float`

- **Raises**: `ValueError` if the value cannot be converted

## str

```txt
str(v)
```

Converts a value to its string representation by calling its `string()` method (from the `Stringable` spec).
All types implement `Stringable` via the `object` root class, so `str()` never fails.

```txt
str(42)         # "42"
str(3.14)       # "3.14"
str(true)       # "true"
str(nil)        # "nil"
```

**Signature**: `fn str(v: any) ->string`

## input

```txt
input(prompt?)
```

Reads a line of text from standard input. If a prompt string is provided, it is written to stdout before
reading (without a trailing newline).

```txt
name := input("What is your name? ")
line := input()
```

**Signature**: `fn input(prompt: string?) ->string`

- **Parameters**: optional prompt string
- **Returns**: `string` -- the line read (without trailing newline)

## typeof

```txt
typeof(v)
```

Returns the runtime type of a value. The returned value can be compared using `==` or used with `is`.

```txt
typeof(42) == int            # true
typeof("hello") == string    # true
```

**Signature**: `fn typeof(v: any) ->type`

## assert

`assert` is a **statement**, not a function. It evaluates an expression and raises `AssertionError` if the
result is `false`. An optional second expression provides a custom error message.

```txt
assert x > 0
assert x > 0, "x must be positive"
assert len(items) == 3, "expected 3 items, got {len(items)}"
```

- If the condition is `true`, execution continues normally.
- If the condition is `false`, raises `AssertionError` with the provided message (or a default one).

See [ERRORS.md](ERRORS.md) for the full exception hierarchy.

## range

`range` is a built-in **type** representing a sequence of integers. Range values are created using the `..`
(exclusive end) and `..=` (inclusive end) operators.

```txt
1..5         # range [1, 5) -- produces 1, 2, 3, 4
1..=5        # range [1, 5] -- produces 1, 2, 3, 4, 5
0..0         # empty range
```

Ranges implement `Iterable[int]` and can be used directly with `for` loops:

```txt
for i in 1..10 {
    print(i)
}
for i in 0..=100 {
    print(i)
}
```

Ranges are immutable values. The start must be less than or equal to the end; otherwise the range is empty.
