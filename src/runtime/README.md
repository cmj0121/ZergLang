# Zerg Runtime

> Minimal interpreter implementation in Go

This is the Zerg bootstrap runtime that can read and compile the Zerg compiler itself.

It implements the minimal set of features and grammar rules to compile the Zerg compiler, and nothing more.
It also handles the subset of the standard library required to compile the compiler.

## Overview

The Zerg runtime is a minimal interpreter implementation in Go that can read and compile the Zerg compiler itself. It
is designed to be as simple as possible, while still being able to compile the Zerg compiler and the necessary subset
of the standard library.

The supported grammar rules and features are limited to only those required to compile the Zerg compiler. This means
that some syntax sugar and features not required for compilation are not supported in the bootstrap runtime.

## Components

The Zerg runtime consists of three main components:

| Component | Description                                               |
| --------- | --------------------------------------------------------- |
| lexer     | Responsible for tokenizing the input source code.         |
| parser    | Responsible for parsing the tokens into an AST (IR).      |
| evaluator | Responsible for evaluating the IR and executing the code. |

## Bootstrap Types

The following types are supported in the bootstrap runtime:

| Type     | Description                    |
| -------- | ------------------------------ |
| `bool`   | Boolean value (`true`/`false`) |
| `int`    | 64-bit signed integer          |
| `string` | UTF-8 encoded string           |
| `list`   | Ordered collection of elements |
| `map`    | Key-value pairs collection     |

## Supported Grammar Subset

The bootstrap runtime supports a minimal subset of Zerg grammar. The detailed support list will be updated as features
are implemented.

## Standard Library

The bootstrap runtime includes a minimal standard library with builtin functions and collection methods.

### Builtin Functions

| Function | Signature          | Description                             |
| -------- | ------------------ | --------------------------------------- |
| `print`  | `print(v...)`      | Output arguments to stdout with newline |
| `len`    | `len(c) -> int`    | Return length of string, list, or map   |
| `str`    | `str(v) -> string` | Convert any value to string             |
| `int`    | `int(v) -> int`    | Convert string or bool to integer       |

**Examples:**

```zerg
print("hello", "world")  # Output: hello world
print(len([1, 2, 3]))    # Output: 3
print(str(42))           # Output: 42
print(int("123"))        # Output: 123
print(int(true))         # Output: 1
```

### List Methods

| Method   | Signature                   | Description                                           |
| -------- | --------------------------- | ----------------------------------------------------- |
| `append` | `list.append(item) -> list` | Return new list with item added                       |
| `pop`    | `list.pop() -> item`        | Return last element                                   |
| `filter` | `list.filter(fn) -> list`   | Return new list with elements passing predicate       |
| `map`    | `list.map(fn) -> list`      | Return new list with function applied to each element |

**Examples:**

```zerg
nums := [1, 2, 3]
print(nums.append(4))                      # Output: [1, 2, 3, 4]
print(nums.pop())                          # Output: 3
print(nums.filter(fn(x) { return x > 1 })) # Output: [2, 3]
print(nums.map(fn(x) { return x * 2 }))    # Output: [2, 4, 6]
```

Note: `append`, `filter`, and `map` return new lists without modifying the original.

### Map Methods

| Method     | Signature                   | Description                              |
| ---------- | --------------------------- | ---------------------------------------- |
| `keys`     | `map.keys() -> list`        | Return list of all keys (sorted)         |
| `values`   | `map.values() -> list`      | Return list of all values (in key order) |
| `contains` | `map.contains(key) -> bool` | Check if key exists in map               |

**Examples:**

```zerg
m := {"a": 1, "b": 2}
print(m.keys())        # Output: [a, b]
print(m.values())      # Output: [1, 2]
print(m.contains("a")) # Output: true
print(m.contains("c")) # Output: false
```

## Project Structure

```txt
src/runtime/
├── go.mod
├── main.go           # Entry point
├── lexer/
│   ├── token.go      # Token types and keyword lookup
│   └── lexer.go      # Tokenizer (source → tokens)
├── parser/
│   ├── ast.go        # AST node definitions
│   └── parser.go     # Recursive descent parser (tokens → AST)
└── evaluator/
    ├── object.go     # Runtime value representations
    ├── evaluator.go  # Tree-walking interpreter (AST → result)
    └── builtins.go   # Built-in functions and collection methods
```
