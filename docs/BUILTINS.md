# Built-in Functions

Built-in functions are available in every module without importing. They are provided by the Zerg runtime and
cannot be shadowed or redefined.

## print

```txt
print(args...)
```

Writes its arguments to standard output followed by a newline. Multiple arguments are separated by a single
space. Each argument is converted to its string representation by calling its `string()` method (from the
`Stringable` spec).

### Examples

```txt
print("hello")                  # hello
print("x =", 42)               # x = 42
print(1, 2, 3)                  # 1 2 3
print(true, nil)                # true nil
```

### Signature

```txt
fn print(args...: any)
```

- **Parameters**: zero or more values of any type
- **Returns**: nothing
- **Output**: writes to stdout, terminates with a newline
