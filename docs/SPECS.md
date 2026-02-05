# Built-in Specs

Zerg provides a set of built-in specs (interfaces) that define standard behavior across types. These specs
are available without importing and form the foundation of Zerg's type system. See
[CONCEPTS.md](CONCEPTS.md) for how to implement specs using `impl ClassName for SpecName`.

## Stringable

Converts a value to its string representation. All types implement `Stringable` via the `object` root class.

```txt
spec Stringable {
    fn string() : string
}
```

Called implicitly by `print()`, string interpolation `{expr}`, and `str()`.

## Equatable

Defines equality comparison between values of the same type.

```txt
spec Equatable {
    fn equals(other: Self) : bool
}
```

The `==` and `!=` operators delegate to `equals()`. The default implementation from `object` performs
structural equality (field-by-field comparison).

## Hashable

Produces an integer hash for use in `map` keys and `set` elements.

```txt
spec Hashable {
    fn hash() : int
}
```

Two values that are equal (via `Equatable`) must produce the same hash. The default implementation from
`object` computes a structural hash from all fields.

## Comparable

Defines ordering between values. Uses a generic type parameter `T` to allow cross-type comparison when
needed, though most implementations compare against their own type.

```txt
spec Comparable[T] {
    fn compare(other: T) : int
}
```

Returns a negative integer if `this < other`, zero if equal, and a positive integer if `this > other`. The
ordering operators (`<`, `>`, `<=`, `>=`) delegate to `compare()`.

## Disposable

Manages resource acquisition and release. Types that hold external resources (files, network connections,
channels) implement this spec to ensure proper cleanup.

```txt
spec Disposable {
    fn open()
    fn close()
}
```

Used by the `with` statement (calls `open()` on entry, `close()` on exit) and `del` (calls `close()` before
removing the binding). See [Resource Management](CONCEPTS.md#resource-management) for details.

## Iterable

Produces an iterator for sequential access to elements. Any type implementing `Iterable` can be used with
`for` loops.

```txt
spec Iterable[T] {
    fn iterator() : iter[T]
}
```

The `for` loop calls `iterator()` once to obtain an `iter[T]`, then repeatedly calls `next()` on the
iterator. See [Iteration](CONCEPTS.md#iteration) for details.

## Iterator

Produces values one at a time from a sequence. When the sequence is exhausted, `next()` raises
`StopIteration`.

```txt
spec Iterator[T] {
    fn next() : T
}
```

The built-in `iter[T]` type implements both `Iterator` and `Iterable` (returning itself). Coroutines that
use `yield` return an `iter[T]` that implements this spec.
