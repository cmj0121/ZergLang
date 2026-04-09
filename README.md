# Zerg

> Write the code as you think — one way, and only one way, to do it.

Like a zerg rush, Zerg programs are fast to write, easy to read, and overwhelmingly straightforward
— swarm your problems with simplicity.

More examples in [`examples/`](examples/).

## Design Principles

| Principle         | Description                            |
| ----------------- | -------------------------------------- |
| small and crisp   | minimal syntax (105 grammar rules)     |
| procedural-first  | straightforward, top-down control flow |
| concurrent        | built-in support for concurrency       |
| garbage-collected | no manual memory management            |
| strongly typed    | catch errors at compile time           |
| null-safe         | no billion-dollar mistakes             |

## Language Highlights

| Feature              | Syntax                                                      |
| -------------------- | ----------------------------------------------------------- |
| variables            | `x := 1` / `mut x := 1`                                     |
| explicit type        | `x: int = 1`                                                |
| print statement      | `print "hello"`                                             |
| string interpolation | `"hello {name}"`                                            |
| null safety          | `T?`, `?.`, `??`                                            |
| pattern matching     | `match x { ... }`                                           |
| specs (interfaces)   | `spec Printable { ... }`                                    |
| generics             | `fn sort[T: Comparable](...)`                               |
| range                | `1..5` (exclusive) `1..=5` (inclusive)                      |
| collections          | `list[int]`, `map[str, int]`, `set[int]`, `tuple[int, str]` |
| guard                | `return x if condition`                                     |
| enum with data       | `Token.Ident(str)`                                          |
| lambda               | `fn(x: int) -> int { return x * 2 }`                        |
| pass by reference    | `fn inc(x: &mut int)`, `inc(&mut n)`                        |
| defer                | `defer close(fd)`                                           |
| raw strings          | `r"no \escapes"`                                            |

## DDD (Dream-Driven Development)

This project is based on the DDD (dream-driven development) methodology which means
the project is based on what I dream of.

All the features are based on my needs and my dreams.
