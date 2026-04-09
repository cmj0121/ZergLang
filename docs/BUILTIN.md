# Zerg Built-in Types

Compiler-provided types, specs, and enums. These are not keywords — they
are resolved as IDENT by the lexer and pre-declared by the compiler.

## Primitive Types

| Type    | Description                     | Zero value |
| ------- | ------------------------------- | ---------- |
| `int`   | Platform-width signed integer   | `0`        |
| `float` | 64-bit IEEE 754                 | `0.0`      |
| `bool`  | `true` or `false`               | `false`    |
| `str`   | UTF-8 immutable string          | `""`       |
| `byte`  | 8-bit unsigned integer (0--255) | `0`        |
| `rune`  | Unicode code point              | `'\0'`     |

## Collection Types

| Type               | Description                       | Literal        | Empty |
| ------------------ | --------------------------------- | -------------- | ----- |
| `list[T]`          | Ordered, variable-length sequence | `[1, 2, 3]`    | `[]`  |
| `map[K, V]`        | Key-value associative container   | `{"a": 1}`     | `{:}` |
| `set[T]`           | Unordered unique values           | `{1, 2, 3}`    | `{}`  |
| `tuple[T, U, ...]` | Fixed-size heterogeneous sequence | `(1, "hello")` | N/A   |

All collections implement `Iterable[T]`. Maps iterate over keys.
Arity of generic arguments is checked semantically, not syntactically.

## Concurrency Types

| Type      | Description                       | Constructor                                        |
| --------- | --------------------------------- | -------------------------------------------------- |
| `chan[T]` | Typed channel for message passing | `chan[int]()` unbuffered, `chan[int](10)` buffered |

Unbuffered channels block send until a receiver is ready.
Buffered channels block send only when the buffer is full.
Channels implement `Iterable[T]` — iterating receives until closed.

## Built-in Specs

### Printable

Used internally by the `print` statement. Any type passed to `print`
must implement `Printable` (built-in types do so automatically).

```zerg
spec Printable {
    fn to_str() -> str
}
```

### Iterable

The iteration protocol. Any type implementing `Iterable[T]` can be
used with `for x in expr`. Returns `T?` (`Option[T]`) — `nil` signals
exhaustion.

The `for` loop creates a mutable copy of the iterable (copy-by-value),
then calls `next()` repeatedly. The original value is not modified.

```zerg
spec Iterable[T] {
    fn next() -> T?
}
```

Built-in iterables:

| Type         | Yields | Order           |
| ------------ | ------ | --------------- |
| `list[T]`    | `T`    | index order     |
| `map[K, V]`  | `K`    | insertion order |
| `set[T]`     | `T`    | unspecified     |
| `tuple[...]` | N/A    | not iterable    |
| `chan[T]`    | `T`    | receive order   |
| `range (..)` | `int`  | ascending       |

### Exception

The exception protocol. Any type implementing `Exception` can be
used with `raise` and caught by `try/except`.

```zerg
spec Exception {
    fn message() -> str
}
```

## Built-in Enums

### Result

The error-handling type. `Result.Ok(T)` for success, `Result.Err(E)`
for failure. The `?` postfix operator propagates errors.

```zerg
enum Result[T, E] {
    Ok(T)
    Err(E)
}
```

### Option (type alias)

Alias for `Result[T, nil]`. The `T?` syntax is sugar for `Option[T]`.

```zerg
type Option[T] = Result[T, nil]
```

Null safety operators work on `Option[T]`:

| Operator | Usage             | Meaning                              |
| -------- | ----------------- | ------------------------------------ |
| `T?`     | type position     | `Option[T]` = `Result[T, nil]`       |
| `?.`     | `expr?.field`     | safe navigation, nil if receiver nil |
| `??`     | `expr ?? default` | nil coalescing                       |
| `?`      | `expr?`           | propagate Err, unwrap Ok             |
