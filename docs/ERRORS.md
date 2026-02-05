# Exception Hierarchy

Zerg uses a flat exception hierarchy rooted at `Exception`. All exceptions carry a `message: string` property
describing the error. User-defined exceptions should extend `Exception` using class embedding.

## Base Exception

```txt
class Exception {
    pub message: string
}
```

All built-in exceptions embed `Exception` and can be caught by catching `Exception` as a catch-all.

## Built-in Exceptions

| Exception             | Raised When                                         |
| --------------------- | --------------------------------------------------- |
| `AssertionError`      | `assert` condition evaluates to `false`             |
| `ValueError`          | Invalid value for a conversion (e.g., `int("abc")`) |
| `IndexError`          | List index out of bounds                            |
| `KeyError`            | Map key not found                                   |
| `DivisionByZeroError` | Division, floor division, or modulo by zero         |
| `TypeError`           | Operation applied to incompatible type              |
| `StopIteration`       | Iterator exhausted (used internally by `for` loops) |
| `ChannelClosedError`  | Send on a closed channel                            |

## Handling Exceptions

Exceptions are handled using `try-expect-finally`. The `expect` clause catches exceptions by type using
the `as` keyword to bind the exception to a variable:

```txt
try {
    items[100]
} expect IndexError as e {
    print(e.message)
} expect Exception as e {
    print("unexpected: {e.message}")
} finally {
    cleanup()
}
```

Multiple `expect` clauses are checked in order. The first matching type handles the exception. Place more
specific exception types before more general ones -- if `Exception` is listed first, it will catch everything
and subsequent clauses will never execute.

The `finally` block is optional and runs regardless of whether an exception was raised. Inside an `expect`
block, a bare `raise` (with no expression) re-raises the current exception.

## User-Defined Exceptions

Create custom exceptions by embedding `Exception`:

```txt
class HttpError {
    Exception
    pub status_code: int
}

raise HttpError(message="Not Found", status_code=404)
```

User-defined exceptions can be caught by their own type or by `Exception`.
