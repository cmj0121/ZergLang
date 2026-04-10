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

## Pointer Types

| Type     | Description                     | Constructor |
| -------- | ------------------------------- | ----------- |
| `ptr[T]` | Owned heap pointer to a value T | `ptr(expr)` |

`ptr[T]` enables recursive data structures (linked lists, trees).
The pointer owns the heap-allocated value — freed when the owner's
scope exits (recursive: freeing a struct frees all its `ptr` fields).

```zerg
struct Node {
    value: int
    next: ptr[Node]?
}

root := Node {
    value: 1,
    next: ptr(Node { value: 2, next: nil })
}
print root.next?.value    # 2 (auto-deref)
```

Copy semantics:

| Assignment               | Behavior                              |
| ------------------------ | ------------------------------------- |
| `b := a` (immutable)     | share heap allocation (safe, no copy) |
| `mut b := a`             | deep copy (entire chain)              |
| Concurrency (`rush`, ch) | deep copy (entire chain)              |

Limitations: `ptr[T]` owns its target — no shared references, no cycles.
Trees and linked lists work. Graphs with shared nodes do not — use
index-based patterns with `list[T]` + `map[K, V]` instead.

## Concurrency Types

| Type      | Description                       | Constructor                                        |
| --------- | --------------------------------- | -------------------------------------------------- |
| `chan[T]` | Typed channel for message passing | `chan[int]()` unbuffered, `chan[int](10)` buffered |
| `sync[T]` | Mutex-protected shared value      | `sync[int](0)`                                     |

Unbuffered channels block send until a receiver is ready.
Buffered channels block send only when the buffer is full.
Channels implement `Iterable[T]` — iterating receives until closed.

### sync

`sync[T]` wraps a value with a read-write lock. The data is only
accessible through the lock API — impossible to bypass. Passed by
immutable reference to tasks (like `chan[T]`).

```zerg
counter := sync[int](0)
rush |c| {
    c.lock(|v: &mut int| { v += 1 })
}(counter)
print counter.read()    # immutable copy, read-lock
```

| Method                     | Description               | Lock       |
| -------------------------- | ------------------------- | ---------- |
| `sync[T](val)`             | create with initial value | —          |
| `.lock(\|v: &mut T\| { })` | exclusive write access    | write lock |
| `.read() -> T`             | return immutable copy     | read lock  |

The compiler may optimize `sync[T]` for primitive types (int, bool,
byte, rune) to use CPU atomic instructions instead of a mutex.

Both `chan[T]` and `sync[T]` are **runtime resources**: cannot be
copied or assigned, passed by reference, freed at owning scope exit.

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

The `for` loop borrows the iterable (immutable reference). The iterator
state is managed internally. The original value is not modified.

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

### Operator Specs

Operator overloading is implemented through specs. The compiler rewrites
operator expressions into method calls. Built-in types implement these
specs automatically.

#### Arithmetic

| Operator    | Spec     | Method                  | Rewrite              |
| ----------- | -------- | ----------------------- | -------------------- |
| `+`         | `Add[T]` | `fn add(other: T) -> T` | `a + b` = `a.add(b)` |
| `-`         | `Sub[T]` | `fn sub(other: T) -> T` | `a - b` = `a.sub(b)` |
| `*`         | `Mul[T]` | `fn mul(other: T) -> T` | `a * b` = `a.mul(b)` |
| `/`         | `Div[T]` | `fn div(other: T) -> T` | `a / b` = `a.div(b)` |
| `%`         | `Mod[T]` | `fn mod(other: T) -> T` | `a % b` = `a.mod(b)` |
| `-` (unary) | `Neg`    | `fn neg() -> this`      | `-a` = `a.neg()`     |

```zerg
spec Add[T] { fn add(other: T) -> T }
spec Sub[T] { fn sub(other: T) -> T }
spec Mul[T] { fn mul(other: T) -> T }
spec Div[T] { fn div(other: T) -> T }
spec Mod[T] { fn mod(other: T) -> T }
spec Neg    { fn neg() -> this }
```

#### Comparison

| Operator             | Spec         | Method                           |
| -------------------- | ------------ | -------------------------------- |
| `==`, `!=`           | `Eq`         | `fn eq(other: this) -> bool`     |
| `<`, `>`, `<=`, `>=` | `Comparable` | `fn compare(other: this) -> int` |

`!=` is the negation of `eq()`. `>`, `<=`, `>=` are derived from `compare()`:
returns negative (less), zero (equal), or positive (greater).

```zerg
spec Eq {
    fn eq(other: this) -> bool
}

spec Comparable: Eq {
    fn compare(other: this) -> int
}
```

#### Bitwise

| Operator | Spec        | Method                      |
| -------- | ----------- | --------------------------- |
| `&`      | `BitAnd[T]` | `fn bitand(other: T) -> T`  |
| `\|`     | `BitOr[T]`  | `fn bitor(other: T) -> T`   |
| `^`      | `BitXor[T]` | `fn bitxor(other: T) -> T`  |
| `~`      | `BitNot`    | `fn bitnot() -> this`       |
| `<<`     | `Shl`       | `fn shl(bits: int) -> this` |
| `>>`     | `Shr`       | `fn shr(bits: int) -> this` |

```zerg
spec BitAnd[T] { fn bitand(other: T) -> T }
spec BitOr[T]  { fn bitor(other: T) -> T }
spec BitXor[T] { fn bitxor(other: T) -> T }
spec BitNot    { fn bitnot() -> this }
spec Shl       { fn shl(bits: int) -> this }
spec Shr       { fn shr(bits: int) -> this }
```

#### Indexing

| Operator | Spec          | Method                  |
| -------- | ------------- | ----------------------- |
| `a[i]`   | `Index[K, V]` | `fn index(key: K) -> V` |

```zerg
spec Index[K, V] {
    fn index(key: K) -> V
}
```

#### Hashable

Required for use as `map` keys or `set` elements.

```zerg
spec Hashable {
    fn hash() -> int
}
```

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
