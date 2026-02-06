# Built-ins

Built-in functions and statements are available in every module without importing. They are provided by the
Zerg runtime and cannot be shadowed or redefined.

## Overview

| Built-in         | Kind      | Description                               |
| ---------------- | --------- | ----------------------------------------- |
| `print(v)`       | function  | Write to stdout                           |
| `len(c)`         | function  | Length of string/list/map/set             |
| `int(v)`         | function  | Convert to int                            |
| `float(v)`       | function  | Convert to float                          |
| `str(v)`         | function  | Convert to string (calls `string()`)      |
| `input(prompt?)` | function  | Read line from stdin                      |
| `assert`         | statement | Raise `AssertionError` if condition false |
| `range`          | type      | Built-in iterable range type              |

## print

```txt
print(v)
```

Writes a value to standard output followed by a newline. The value is converted to its string representation
by calling its `string()` method (from the `Stringable` spec). Use string interpolation for multiple values.

```txt
print("hello")                  # hello
print(42)                       # 42
print("x = {x}, y = {y}")      # x = 1, y = 2
```

**Signature**: `fn print(v: Stringable)`

- **Parameters**: a value that implements `Stringable` (all types do)
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

- **Parameters**: a `string`, `list[T]`, `map[K, V]`, or `set[T]`
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

- **Parameters**: `float`, `string`, or `bool`
- **Returns**: `int`
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

- **Parameters**: `int`, `string`, or `bool`
- **Returns**: `float`
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

**Signature**: `fn str(v: Stringable) -> string`

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

## Built-in Modules

Zerg provides built-in modules that are globally available without import. These modules provide safe, high-level
APIs for common operations.

### sys Module

The `sys` module provides system-related functions.

| Method           | Signature          | Description                                         |
| ---------------- | ------------------ | --------------------------------------------------- |
| `sys.os()`       | `-> string`        | Returns OS name: `"linux"`, `"darwin"`, `"windows"` |
| `sys.arch()`     | `-> string`        | Returns architecture: `"amd64"`, `"arm64"`, `"386"` |
| `sys.args()`     | `-> list`          | Returns command-line arguments                      |
| `sys.exit(code)` | `int -> nil`       | Exits program with given code                       |
| `sys.env(name)`  | `string -> string` | Returns environment variable value                  |

```zerg
if sys.os() == "darwin" {
    print("Running on macOS")
}

args := sys.args()
if len(args) < 2 {
    print("Usage: program <file>")
    sys.exit(1)
}
```

### io Module

The `io` module provides file I/O operations.

| Method                      | Signature                  | Description                                    |
| --------------------------- | -------------------------- | ---------------------------------------------- |
| `io.open(path, mode)`       | `string, string -> handle` | Open file (modes: `"r"`, `"w"`, `"a"`, `"rw"`) |
| `io.read(handle)`           | `handle -> string`         | Read entire file contents                      |
| `io.read_lines(handle)`     | `handle -> list`           | Read file as list of lines                     |
| `io.write(handle, data)`    | `handle, string -> int`    | Write string, return bytes written             |
| `io.close(handle)`          | `handle -> nil`            | Close file handle                              |
| `io.exists(path)`           | `string -> bool`           | Check if file exists                           |
| `io.read_file(path)`        | `string -> string`         | Read file contents directly                    |
| `io.write_file(path, data)` | `string, string -> nil`    | Write file contents directly                   |

```zerg
# Simple file reading
content := io.read_file("data.txt")
lines := str.split(content, "\n")

# Using file handles
handle := io.open("output.txt", "w")
io.write(handle, "Hello, World!\n")
io.close(handle)
```

### str Module

The `str` module provides string manipulation functions.

| Method                         | Signature                          | Description                            |
| ------------------------------ | ---------------------------------- | -------------------------------------- |
| `str.split(s, sep)`            | `string, string -> list`           | Split string by separator              |
| `str.join(list, sep)`          | `list, string -> string`           | Join list elements with separator      |
| `str.trim(s)`                  | `string -> string`                 | Remove leading/trailing whitespace     |
| `str.find(s, sub)`             | `string, string -> int`            | Find substring index (-1 if not found) |
| `str.replace(s, old, new)`     | `string, string, string -> string` | Replace all occurrences                |
| `str.substring(s, start, end)` | `string, int, int -> string`       | Extract substring                      |
| `str.starts_with(s, prefix)`   | `string, string -> bool`           | Check if starts with prefix            |
| `str.ends_with(s, suffix)`     | `string, string -> bool`           | Check if ends with suffix              |
| `str.upper(s)`                 | `string -> string`                 | Convert to uppercase                   |
| `str.lower(s)`                 | `string -> string`                 | Convert to lowercase                   |
| `str.contains(s, sub)`         | `string, string -> bool`           | Check if contains substring            |

```zerg
text := "  Hello, World!  "
trimmed := str.trim(text)           # "Hello, World!"
upper := str.upper(trimmed)         # "HELLO, WORLD!"

csv := "a,b,c"
parts := str.split(csv, ",")        # ["a", "b", "c"]
joined := str.join(parts, " | ")    # "a | b | c"
```

### char Module

The `char` module provides character-level operations.

| Method             | Signature        | Description                               |
| ------------------ | ---------------- | ----------------------------------------- |
| `char.ord(c)`      | `string -> int`  | Get ASCII/Unicode code of first character |
| `char.chr(code)`   | `int -> string`  | Convert code to character                 |
| `char.is_digit(c)` | `string -> bool` | Check if character is a digit (0-9)       |
| `char.is_alpha(c)` | `string -> bool` | Check if character is a letter            |
| `char.is_space(c)` | `string -> bool` | Check if character is whitespace          |
| `char.is_alnum(c)` | `string -> bool` | Check if character is alphanumeric        |

```zerg
# Character checking
if char.is_digit("5") {
    print("It's a digit")
}

# Character conversion
code := char.ord("A")   # 65
ch := char.chr(97)      # "a"
```

These modules are the recommended way to perform system operations. For low-level access, see
[UNSAFE.md](UNSAFE.md).
