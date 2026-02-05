# Concepts

The following are the core concepts of the Zerg programming language, how to use it, and why it is designed
this way.

## Variables

The variable in Zerg is immutable by default, which means once a variable is assigned, its value cannot be
changed. This helps you to avoid unintended side effects and make the code more predictable and easier to
understand.

You can declare a variable in any scope. You can also re-declare a variable with the same name in the same
scope or an inner scope, which creates a new variable and shadows the previous one. The old name binding is
released. This is not mutation -- it is a new variable that happens to reuse the same name.

Zerg supports two ways to declare a variable:

- `:=` declares a variable with automatic type inference, where the compiler detects the type from the
  assigned value.
- `: type =` explicitly declares the type of the variable.

The `=` symbol alone is used to update a mutable variable that has already been declared.

### Mutable

You can specify the variable with `mut` keyword to make it mutable, which means you can change the value of
the variable after it is assigned. However, the mutable variable should be used sparingly, and you should
prefer to use the immutable variable whenever possible.

### Null-safety

Zerg is designed to be a null-safe programming language. Variables are non-nullable by default, which means
a variable cannot hold `nil` unless its type is explicitly declared as nullable using the `?` suffix.

You can access a nullable variable after checking it is not `nil`, or using the safe navigation operator `?.`
to access the property or method of the nullable variable. You can also use `??` to provide the default value
for a nullable variable, which means if the variable is `nil`, it will return the default value instead of
throwing an exception.

### Assignment Semantics

All assignment in Zerg produces an independent copy of the value (value semantics). Whether the target is
nullable or non-nullable does not change this behavior -- assignment always copies.

This uniform copy-on-assign rule keeps the mental model simple: changing or deleting one variable never
affects another. The compiler and runtime are free to optimize copies away (e.g., copy-on-write or move
semantics for large data) as long as the observable behavior remains the same.

**Exception -- `chan` and `with` resources**: A `chan` is a shared handle -- assigning or passing a `chan`
shares the underlying conduit rather than copying it, so that multiple tasks can communicate through the
same channel. Resources bound in a `with` statement transfer ownership into the scope rather than copy.
See [Concurrency](#concurrency) and [Resource Management](#resource-management) for details.

## Comments

Zerg supports line comments starting with `#`. Everything from `#` to the end of the line is ignored by the
compiler. The `#!` shebang on the first line of a script is treated as a regular comment.

```txt
# this is a line comment
x := 42  # inline comment
```

A bare string literal used as a statement is also discarded by the compiler and can serve as a documentation
comment:

```txt
"This function computes the factorial of n."
fn factorial(n: int) : int {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
```

There are no block comments. Use multiple `#` lines for longer explanations.

## Strings

Strings in Zerg are UTF-8 encoded and delimited by double quotes. They support escape sequences and string
interpolation.

### String Interpolation

Any expression enclosed in `{}` inside a string is evaluated and its string representation is inserted
in place. The expression's type must implement `Stringable` (all types do via `object`).

```txt
name := "Zerg"
print("Hello, {name}!")           # Hello, Zerg!
print("1 + 2 = {1 + 2}")         # 1 + 2 = 3
```

To include a literal `{` or `}` in a string, escape it with a backslash:

```txt
print("Use \{braces\} for interpolation")  # Use {braces} for interpolation
```

### Raw Strings

A raw string is prefixed with `r` and does not process escape sequences or interpolation. Everything
between the quotes is taken literally.

```txt
path := r"C:\Users\name\docs"
pattern := r"(\d+)\s+{not interpolated}"
```

Raw strings cannot contain a double-quote character, since there is no escape mechanism inside them.

## Functions

Zerg is a procedural-first language. Functions are defined at the module level as standalone routines, not
inside classes. Classes and OOP features are optional tools for structuring data, not a requirement.

Code written directly in a file (outside any function) is wrapped by the compiler into a built-in `_init()`
function. This `_init()` is executed once and only once when the module is first imported. The Zerg compiler
detects circular imports at compile time and rejects them.

When run as a script, Zerg executes the file sequentially: `_init()` runs first (top-level code), then
`main()` is called if it exists. The `main` function is **not required** -- a script with only top-level
code is valid and will execute via `_init()` alone.

Functions can be called with parameters and return values. They are considered as first-class citizens in Zerg,
which means functions can be passed as arguments, returned from other functions, and assigned to variables.

### Built-in Functions

Zerg provides built-in functions and statements that are available in every module without importing. See
[BUILTINS.md](BUILTINS.md) for the complete reference, including `print`, `len`, `str`, type conversion
functions, `assert`, and `range`.

### Return Type

A function that declares a return type (`: type` after the parameter list) must return a value of that type
on every code path. A function with no return type annotation returns nothing -- attempting to use its result
(e.g., `x := foo()`) is a compile-time error. A bare `return` (with no expression) is valid inside a
void-returning function.

### Function Types

Function types are expressed using the `fn` keyword directly, describing the parameter types and return type
of a callable. Any function or method whose signature matches the declared `fn` type can be assigned or passed
directly.

As a convenience, a standalone function implicitly satisfies any single-method `spec` if its signature
matches -- no explicit `implement` declaration is needed. This allows functions to be used interchangeably
with single-method interfaces.

## Typing System

Zerg is a strongly-typed programming language, which means each variable has a specific type and cannot be
changed to another type. The type of a variable is determined at compile-time, either inferred by the compiler
through `:=` or explicitly declared through `: type =` notation. The type system is designed to be simple and
easy to use, with a small set of built-in types. In ambiguous cases (such as `nil` or empty collections), an
explicit type declaration is required.

### Built-in Types

Zerg only provides the necessary built-in types to cover most of the common use cases, and you can create your
own types by `class` with your properties and `impl` methods. The built-in types are:

| Type     | Description                                   | Bootstrap Supported |
| -------- | --------------------------------------------- | ------------------- |
| `bool`   | Boolean value (true or false)                 | Yes                 |
| `int`    | 64-bit signed integer                         | Yes                 |
| `float`  | 64-bit floating-point (IEEE 754)              | No                  |
| `string` | UTF-8 encoded string                          | Yes                 |
| `list`   | Ordered collection of elements                | Yes                 |
| `map`    | Key-value pairs collection                    | Yes                 |
| `set`    | Unordered unique elements                     | No                  |
| `chan`   | Typed channel for concurrency (shared handle) | No                  |
| `iter`   | Iterator over a sequence of values            | No                  |
| `range`  | Integer range (created with `..` or `..=`)    | No                  |

Note: `nil` is a value representing the absence of data, not a type. It can only be assigned to variables
whose type is explicitly declared as nullable.

### Enum Types

Zerg supports `enum` as an algebraic data type (sum type). An enum defines a type that can hold one of
several distinct variants, where each variant can optionally carry associated data. This enables type-safe
modeling of values that have multiple possible forms.

Enum variants are handled using the `match` statement, which requires all variants to be covered
exhaustively. The compiler rejects any `match` that does not handle every variant, ensuring no case is
silently ignored. Each branch can destructure the variant's associated data and bind it to local variables.

The built-in `Result[T, E]` enum has two variants: `Ok(value: T)` representing success, and `Err(error: E)`
representing failure. The caller must `match` both variants before accessing the inner value.

### Type Parameters

Zerg supports type parameters using `[T]` syntax. Collection types must be declared with explicit element
types, such as `list[int]`, `map[string, int]`, or `set[string]`. User-defined classes can also accept type
parameters in the same way, allowing reusable and type-safe data structures while preserving compile-time type
checking.

## Memory Management

Zerg uses value semantics with garbage collection (see [Assignment Semantics](#assignment-semantics) for
details on copy behavior and exceptions). You cannot directly access memory addresses.

### Garbage Collection

Zerg handles memory management automatically through garbage collection. You do not need to manually allocate
or deallocate memory. A value is collected when no variable holds it.

### Variable Deletion

You can use the `del` keyword to explicitly remove a variable binding from the current scope. After `del x`,
the name `x` is no longer accessible in the remaining code of that scope. It is useful for signaling intent
that a large piece of data is no longer needed.

If the deleted variable implements `Disposable`, `del` will call `close()` on it before removing the binding.
This means `del` on a `chan` closes the channel, and `del` on a file handle releases the resource. For
non-`Disposable` types, `del` simply removes the name binding and the GC handles deallocation.

### Resource Management

You can use the `with` statement to manage a scoped resource. The resource must implement the `Disposable`
spec, which defines how the resource is acquired (`open()`) and released (`close()`). The `with` statement
calls `open()` when entering the scope and `close()` when exiting. This is useful for managing non-memory
resources such as file handles, network connections, etc. The `with` statement transfers ownership of the
resource into the scope -- the resource cannot be used after the scope exits.

## Visibility

By default, all functions, classes, and variables are `private` to the current module, and you need to specify
the `pub` keyword to make them `public` to other modules. This helps you to encapsulate your code and avoid
name conflicts and unintended interactions between different modules.

## OOP (Object-Oriented Programming)

Zerg is a semi-OOP programming language with the smallest core possible. Each class instance is an object. The
`class` body defines **properties only** (data shape), and all **methods** are implemented inside a separate
`impl` block. This separates what the object holds from what the object can do.

### Embedding (Composition-based)

A class can embed one or more existing types by listing the type name without a field name in the class body.
Embedding promotes the embedded type's `pub` properties and methods to the outer class, allowing direct access
without explicit delegation. Private members of the embedded type remain invisible.

```txt
class Animal {
    pub name: string
    pub age: int
}

class Dog {
    Animal                  # embed Animal -- promotes name and age
    pub breed: string
}
```

The embedded type is stored as an anonymous field. Its promoted members can be accessed directly on the outer
class (e.g. `dog.name`), and the embedded value itself can be accessed using the type name as a field
(e.g. `dog.Animal`).

If multiple embedded types provide public members with the same name, the compiler raises an error. The
embedding class must resolve the conflict by overriding the ambiguous member explicitly in its own `impl`
block.

### Class Instantiation

Class instances are created by calling the class name as a constructor. Arguments are matched to properties
by name using named-argument syntax (`name=value`). Positional arguments are matched in declaration order.
Positional and named arguments can be mixed, but all positional arguments must come before any named ones.

```txt
animal := Animal(name="Rex", age=3)
dog := Dog(Animal=Animal(name="Rex", age=3), breed="Labrador")

# Positional, named, or mixed (positional first)
Animal("Rex", 3)
Animal("Rex", age=3)
```

Embedded types are initialized by passing the embedded value using the type name as the argument name.

### Methods and `this`

Methods are defined inside `impl` blocks and have implicit access to the current instance through the `this`
keyword. Inside a method, `this` refers to the receiver object and can be used to access its properties and
call its other methods.

By default, methods cannot modify the receiver. A method that needs to mutate `this` must be declared with
`mut fn`. The compiler enforces this -- calling a `mut fn` method on an immutable variable is a compile-time
error.

Parameters can also be declared mutable using `fn foo(x: mut int)`, allowing the function body to reassign
the parameter. This only affects the local copy -- it does not modify the caller's value (consistent with
Zerg's value semantics).

```txt
impl Animal {
    fn greet() : string {
        return "Hi, I am {this.name}"
    }
    mut fn birthday() {
        this.age = this.age + 1
    }
}
```

### Spec (Interface)

Zerg uses `spec` to define the interface contract. A class must explicitly declare its spec implementation
using `impl ClassName for SpecName`, and provide all the methods defined in the `spec`. Regular class methods
that do not belong to any spec are defined using `impl ClassName` without the `for` clause. This provides
interface-based polymorphism without class hierarchy complexity.

You can check whether a value implements a spec at runtime using the `is` operator (e.g., `x is Comparable`).
See [EXPRESSIONS.md](EXPRESSIONS.md#type-checking-with-is) for details. For the complete list of built-in
specs, see [SPECS.md](SPECS.md).

### Object Root

All types implicitly embed the `object` root class. The `object` class holds no properties, and provides
default implementations for the basic specs:

| Spec         | Method          | Default Behavior                     |
| ------------ | --------------- | ------------------------------------ |
| `Stringable` | `string()`      | Returns the class name               |
| `Equatable`  | `equals(other)` | Structural equality (field-by-field) |
| `Hashable`   | `hash()`        | Structural hash (field-by-field)     |

The `Disposable` spec (`open()` and `close()`) is not universal -- it is an opt-in spec that only resource
types implement. See [Resource Management](#resource-management) for details.

Built-in types (`int`, `bool`, `string`, etc.) also embed `object` implicitly but are sealed -- users cannot
embed or modify them. The compiler may optimize their internal representation, but they behave as objects in
all other respects.

## Error Handling

Zerg provides two mechanisms for error handling: **Result type** and **exceptions**, each suited for
different situations.

### Result Type

For expected and recoverable errors, use the `Result[T, E]` enum. The caller must `match` both the `Ok`
and `Err` variants before accessing the value, ensuring errors are never silently ignored.

### Exceptions

For unexpected or unrecoverable errors, Zerg supports exceptions. You can `raise` an exception to interrupt
the current execution, and handle it using the `try-expect-finally` statement. The `try` block contains the
code that may raise an exception, the `expect` block handles the raised exception, and the `finally` block
runs regardless of whether an exception was raised, typically used for cleanup. Inside an `expect` block, a
bare `raise` (with no expression) re-raises the current exception, allowing partial handling or logging
before propagating the error upward.

**Guideline**: Prefer `Result` for operations where failure is a normal, expected outcome (file not found,
invalid input, network timeout). Use exceptions for programming errors or truly unexpected situations
(assertion failures, corrupted state). Note: `StopIteration` is a special built-in exception used
internally by the `for` loop to signal iterator exhaustion -- it is not intended for general use and should
not be raised or caught by user code.

## Concurrency

Zerg supports concurrency as a core principle. You can spawn a lightweight concurrent task by using the `go`
keyword, which runs the given function concurrently without blocking the current execution. These concurrent
tasks are managed by the Zerg runtime and are multiplexed onto system threads automatically, making them
cheap to create.

### Channels

Concurrent tasks communicate through `chan[T]`, which are typed conduits for sending and receiving values
between tasks. A channel can be unbuffered (sender blocks until receiver is ready) or buffered (sender blocks
only when the buffer is full).

The `chan` itself is a shared handle -- assigning or passing a `chan` shares the underlying conduit, not
copies it. However, the data sent through a `chan` is **copied** (consistent with Zerg's value semantics).
The sender retains its own copy, and the receiver gets an independent copy. The compiler may optimize this
to a move when it can prove the sender does not use the value afterward.

A `chan` implements `Disposable` and can be closed by calling `close()`, or by using `del` on the channel
variable. Closing a channel signals to receivers that no more values will be sent -- any `for` loop
iterating over the channel will terminate, and further sends will raise an exception.

By default, data should be passed between concurrent tasks through channels rather than shared mutable state.
This follows the principle: **share memory by communicating, do not communicate by sharing memory**.

### Coroutines

Zerg supports coroutines as cooperative, suspendable functions using the `yield` keyword. A function that
contains `yield` is a coroutine. Calling a coroutine does not execute its body immediately -- instead, it
returns an `iter[T]` that produces values lazily, one at a time, each time the caller requests the next
value.

When `yield` is reached, the coroutine suspends its execution and produces a value to the caller. The
coroutine resumes from where it left off when the next value is requested. When the coroutine function
returns (or its body ends), the sequence is exhausted and produces no more values.

Yielded values follow Zerg's value semantics -- each yielded value is an independent copy delivered to the
caller.

Coroutines and `go` tasks serve different purposes. Use `go` when you want a task to run independently and
communicate results through `chan`. Use coroutines with `yield` when you want a function to produce a
stream of values on demand, driven by the caller's pace. Coroutines are single-threaded and cooperative --
they do not run in parallel with the caller.

## Iteration

Zerg uses two specs to support iteration: `Iterable` and `Iterator`.

### Iterable and Iterator Specs

The `Iterable` spec defines a single method `iterator()` that returns an `iter[T]`. Any type that implements
`Iterable` can be used directly with a `for` loop.

The `Iterator` spec defines a single method `next()` that returns the next value in the sequence. When the
sequence is exhausted, `next()` raises a `StopIteration` exception. This avoids using `nil` as a sentinel,
since `nil` may be a valid value in the sequence (e.g., iterating over a `list[string?]`). The `for` loop
internally handles `StopIteration` to terminate cleanly -- user code does not need to catch it manually.

The built-in `iter[T]` type implements both `Iterator` and `Iterable` (returning itself). This means an
`iter[T]` can be used directly with `for`, and can also be consumed manually by calling `next()`.

### For Loop

The `for` loop is the standard way to iterate. When given an `Iterable`, it calls `iterator()` to obtain
an `iter[T]`, then repeatedly calls `next()` and binds each value to the loop variable until `StopIteration`
is raised.

The following built-in types implement `Iterable` and can be used directly with `for`:

| Type        | Iterates Over                                       |
| ----------- | --------------------------------------------------- |
| `list[T]`   | Each element in order                               |
| `map[K, V]` | Each key-value pair                                 |
| `set[T]`    | Each element (unordered)                            |
| `string`    | Each character                                      |
| `iter[T]`   | Itself (already an iterator)                        |
| `chan[T]`   | Values received from the channel until it is closed |
| `range`     | Each integer in the range                           |

Because coroutines return `iter[T]`, a coroutine acts as a generator -- the `for` loop drives the coroutine,
resuming it to produce the next value on each iteration. Similarly, a `for` loop over a `chan` receives
values one at a time, blocking until a value is available, and terminates when the channel is closed.

User-defined classes can implement `Iterable` to make their instances usable with `for` loops, or implement
`Iterator` directly to create custom iterators.

## Packages and Modules

Zerg is a compiled programming language that can also be executed like an interpreted language. Zerg can run
source code directly without a separate compilation step -- the compiler compiles and runs in one action. It
detects the `main` function and uses it as the entry point.

A directory of `.zg` files forms a **package**. Each file in a package is a **module**. You can expose your
code as a package either as a local directory or through a remote git repository. Packages are imported using
string paths:

```txt
import "io"                  # stdlib or local package
import "example.com/utils"   # remote package (fetched from git repository)
```

A bare name like `"io"` resolves to the standard library or a local package. A domain-prefixed path like
`"example.com/utils"` refers to a remote package. The last segment of the path becomes the local name used
to access the package's public members (e.g. `utils.read()`).

Circular imports between packages are detected at compile time and rejected by the compiler.
