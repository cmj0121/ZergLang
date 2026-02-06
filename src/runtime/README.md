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
    └── builtins.go   # Built-in statement (assert)
```
