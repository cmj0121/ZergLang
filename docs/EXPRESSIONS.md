# Expressions

Expressions are constructs that produce a value. Every expression in Zerg has a type determined at compile
time. This document describes the available operators, their behavior, and their precedence.

## Operator Precedence

Operators are listed from **lowest** to **highest** precedence. Operators at the same precedence level are
evaluated **left to right** (left-associative), except for unary operators which are right-associative.

| Precedence   | Operator                         | Description            | Associativity |
| ------------ | -------------------------------- | ---------------------- | ------------- |
| 1 (lowest)   | `??`                             | Nil coalescing         | Left          |
| 2            | `or`                             | Logical OR             | Left          |
| 3            | `xor`                            | Logical XOR            | Left          |
| 4            | `and`                            | Logical AND            | Left          |
| 5            | `==` `!=` `<` `>` `<=` `>=` `is` | Comparison (chainable) | Left          |
| 6            | `\|`                             | Bitwise OR             | Left          |
| 7            | `^`                              | Bitwise XOR            | Left          |
| 8            | `&`                              | Bitwise AND            | Left          |
| 9            | `<<` `>>`                        | Bit shift              | Left          |
| 10           | `..` `..=`                       | Range                  | None          |
| 11           | `+` `-`                          | Addition, Subtraction  | Left          |
| 12           | `*` `/` `//` `%`                 | Multiply, Divide, Mod  | Left          |
| 13           | `-` `not` `~` `<-`               | Negation, NOT, Receive | Right         |
| 14           | `**`                             | Power                  | Right         |
| 15 (highest) | `()` `[]` `.` `?.` `?[]`         | Call, Index, Member    | Left          |

Parentheses `( )` can be used to override the default precedence.

## Arithmetic Operators

Arithmetic operators work on numeric types (`int` and `float`).

| Operator | Operation      | Operand Types     | Result Type |
| -------- | -------------- | ----------------- | ----------- |
| `+`      | Addition       | `int + int`       | `int`       |
|          |                | `float + float`   | `float`     |
|          | Concatenation  | `string + string` | `string`    |
| `-`      | Subtraction    | `int - int`       | `int`       |
|          |                | `float - float`   | `float`     |
| `*`      | Multiplication | `int * int`       | `int`       |
|          |                | `float * float`   | `float`     |
| `/`      | Division       | `int / int`       | `float`     |
|          |                | `float / float`   | `float`     |
| `//`     | Floor division | `int // int`      | `int`       |
|          |                | `float // float`  | `float`     |
| `%`      | Modulo         | `int % int`       | `int`       |
| `**`     | Power          | `int ** int`      | `int`       |
|          |                | `float ** float`  | `float`     |

Division (`/`) always returns a `float`, even when both operands are `int`. For example, `7 / 2` produces
`3.5`. Floor division (`//`) rounds the result down toward negative infinity and preserves the operand type.
For example, `7 // 2` produces `3` and `-7 // 2` produces `-4`.

The power operator (`**`) is **right-associative**: `2 ** 3 ** 2` is evaluated as `2 ** (3 ** 2)` = `512`.
It binds tighter than unary negation, so `-2 ** 3` is evaluated as `-(2 ** 3)` = `-8`. The exponent for
`int ** int` must be non-negative; a negative exponent raises an exception (use `float` for fractional
results).

Division, floor division, or modulo by zero raises an exception.

There is **no implicit type coercion** between `int` and `float`. Mixing `int` and `float` operands in an
arithmetic expression is a compile-time error. Use explicit conversion functions to convert between types.

## Comparison Operators

Comparison operators return a `bool` value. All comparison operators share the same precedence level.

| Operator | Description              |
| -------- | ------------------------ |
| `==`     | Equal to                 |
| `!=`     | Not equal to             |
| `<`      | Less than                |
| `>`      | Greater than             |
| `<=`     | Less than or equal to    |
| `>=`     | Greater than or equal to |

Equality (`==`, `!=`) is determined by the `Equatable` spec. All types embed a default structural equality
from `object`, but classes can override `equals(other)` in their `impl` block to customize behavior.

Ordering operators (`<`, `>`, `<=`, `>=`) require both operands to be of the same type and the type must
support ordering. The built-in numeric types (`int`, `float`) and `string` (lexicographic) support ordering
natively.

### Comparison Chaining

Comparisons can be **chained** to express range checks and multi-way equality concisely. A chained
comparison `a op1 b op2 c` is equivalent to `a op1 b and b op2 c`, except that each intermediate operand
is evaluated only once.

```txt
1 <= x < 10          # equivalent to: 1 <= x and x < 10
a < b <= c < d       # equivalent to: a < b and b <= c and c < d
a == b == c          # equivalent to: a == b and b == c
```

The chain short-circuits: if any comparison in the chain is `false`, the remaining comparisons are not
evaluated. For example, in `a < b < c`, if `a < b` is `false`, `c` is never evaluated.

### Type Checking with `is`

The `is` operator checks whether a value is of a given type or implements a given spec. It returns a `bool`
and shares the same precedence level as comparison operators.

```txt
x is int              # true if x is type int
x is Comparable       # true if x implements Comparable
```

The right-hand side of `is` must be a type name or spec name (not an expression). The `is` operator can be
combined with `not` for negative checks:

```txt
if not x is string {
    print("x is not a string")
}
```

## Logical Operators

Logical operators work on `bool` values and return a `bool`.

| Operator | Description | Behavior                                                      |
| -------- | ----------- | ------------------------------------------------------------- |
| `or`     | Logical OR  | Returns `true` if either operand is `true`                    |
| `xor`    | Logical XOR | Returns `true` if exactly one operand is `true`, but not both |
| `and`    | Logical AND | Returns `true` if both operands are `true`                    |
| `not`    | Logical NOT | Returns `true` if the operand is `false`                      |

Both `or` and `and` use **short-circuit evaluation**: the right operand is not evaluated if the result can
be determined from the left operand alone. Specifically:

- `a or b` -- if `a` is `true`, `b` is not evaluated.
- `a and b` -- if `a` is `false`, `b` is not evaluated.

The `xor` operator always evaluates both operands, since the result cannot be determined from one side alone.

## Bitwise Operators

Bitwise operators work on `int` values and perform bit-level manipulation.

| Operator | Operation   | Description                                     |
| -------- | ----------- | ----------------------------------------------- |
| `\|`     | Bitwise OR  | Sets each bit to 1 if either bit is 1           |
| `^`      | Bitwise XOR | Sets each bit to 1 if exactly one bit is 1      |
| `&`      | Bitwise AND | Sets each bit to 1 only if both bits are 1      |
| `~`      | Bitwise NOT | Flips every bit (unary)                         |
| `<<`     | Left shift  | Shifts bits left, filling with zeros            |
| `>>`     | Right shift | Arithmetic right shift, preserving the sign bit |

All bitwise operators require `int` operands and produce an `int` result. They are not defined for `float`,
`bool`, or any other type.

The shift amount must be a non-negative `int`. Shifting by a negative amount or by more than 63 bits raises
an exception.

## Range Operators

The range operators create `range` values representing sequences of integers.

| Operator | Description   | Example | Result        |
| -------- | ------------- | ------- | ------------- |
| `..`     | Exclusive end | `1..5`  | 1, 2, 3, 4    |
| `..=`    | Inclusive end | `1..=5` | 1, 2, 3, 4, 5 |

Both operands must be `int`. Range operators are non-associative -- `a..b..c` is a syntax error.
See [BUILTINS.md](BUILTINS.md#range) for details on the `range` type.

## Unary Operators

| Operator | Description      | Operand Type | Result Type |
| -------- | ---------------- | ------------ | ----------- |
| `-`      | Numeric negation | `int`        | `int`       |
|          |                  | `float`      | `float`     |
| `not`    | Logical negation | `bool`       | `bool`      |
| `~`      | Bitwise NOT      | `int`        | `int`       |
| `<-`     | Channel receive  | `chan[T]`    | `T`         |

## Lambda Expressions

A lambda is a lightweight anonymous function that evaluates a single expression. The syntax uses `|` to
delimit parameters and `=>` to introduce the body:

```txt
|x| => x ** 2
|a, b| => a + b
|| => 42
```

Parameter types can be omitted when the compiler can infer them from context:

```txt
numbers.map(|x| => x ** 2)              # type of x inferred from list element
numbers.filter(|x| => x > 0)
```

Explicit type annotations are also allowed:

```txt
|x: int, y: int| => x + y
```

The body must be a single expression or `nop` (a no-operation that returns nothing):

```txt
callbacks.each(|_| => nop)
```

Lambdas are expressions and can be assigned to variables, passed as arguments, or returned from functions.
For multi-statement bodies, use `fn` instead.

## Function Calls

A function call applies arguments to a callable expression. The syntax is:

```txt
expression(arg1, arg2, ...)
```

A trailing comma after the last argument is permitted. The number and types of arguments must match the
function's parameter list at compile time.

Since functions are first-class values, any expression that evaluates to a function type can be called.

### Named Arguments

Arguments can be passed by name using `name=value` syntax. Named arguments can appear in any order, but all
positional arguments must come before any named arguments. This works for all function calls, including class
constructors.

```txt
print("hello", end="\n")
Animal(name="Rex", age=3)
Animal("Rex", age=3)           # mixed: positional first, then named
```

A named argument binds the value to the parameter with the matching name. It is a compile-time error to pass
a named argument that does not correspond to any parameter, or to provide the same parameter both positionally
and by name.

## Member Access

The dot operator `.` accesses a property or method on an object:

```txt
expression.name
```

This is used to access class properties and call methods defined in `impl` blocks.

## Subscript Access

The bracket operator `[]` accesses an element by index or key:

```txt
items[0]           # list element by index
table["key"]       # map value by key
matrix[i][j]       # chained subscript
```

The index expression can be any expression. Out-of-bounds access on a `list` or a missing key on a `map`
raises an exception.

## Collection Literals

Zerg supports inline literals for `list`, `map`, and `set`. The type of a `{}` literal is disambiguated by
its contents:

```txt
# List literals
[1, 2, 3]                    # list[int]
["a", "b", "c"]              # list[string]
[]                           # empty list (requires type annotation)

# Map literals (key: value pairs)
{"name": "zerg", "ver": "1"} # map[string, string]
{:}                          # empty map (requires type annotation)

# Set literals (values only, no colons)
{1, 2, 3}                   # set[int]
{"a", "b", "c"}              # set[string]
{}                           # empty set (requires type annotation)
```

The distinction between `set` and `map` is syntactic: if elements contain `:` separators, it is a `map`;
otherwise it is a `set`. The special `{:}` syntax creates an empty map, while `{}` creates an empty set.
Both require an explicit type annotation on the variable.

Type constructors can also be used to create collections:

```txt
list[int]()
map[string, int]()
set[string]()
chan[int](10)                 # buffered channel with capacity 10
```

A trailing comma after the last element is permitted. The element type for lists, the key/value types for
maps, and the element type for sets are inferred from the contents.

## Null-Safe Operators

For nullable types (declared with `?`), Zerg provides three special operators:

| Operator | Name            | Description                                                                      |
| -------- | --------------- | -------------------------------------------------------------------------------- |
| `?.`     | Safe navigation | Accesses a member only if the receiver is not `nil`; returns `nil` otherwise     |
| `?[]`    | Safe subscript  | Accesses an element only if the receiver is not `nil`; returns `nil` otherwise   |
| `??`     | Nil coalescing  | Returns the left operand if it is not `nil`; otherwise returns the right operand |

The safe navigation operator `?.` and safe subscript operator `?[]` short-circuit the entire chain. If any
part of a chain like `a?.b?[0]?.c` evaluates to `nil`, the rest of the chain is skipped and the result is
`nil`.

The nil coalescing operator `??` provides a default value for nullable expressions:

```txt
name ?? "anonymous"
items?[0] ?? default_item
```

The right operand of `??` is only evaluated if the left operand is `nil`.

## Grouping

Parentheses `( )` override the default operator precedence:

```txt
(a + b) * c
```

A grouped expression evaluates to the same type and value as the inner expression.
