# Concepts

The following are the core concepts of the Zerg programming language, how to use it, and why it is designed
this way.

## Reserved Keywords

The following identifiers are reserved and cannot be used as variable, function, or type names.

### Declarations

| Keyword | Description                                    |
| ------- | ---------------------------------------------- |
| `fn`    | Declares a function                            |
| `class` | Declares a class (data structure)              |
| `impl`  | Implements methods for a class or spec         |
| `spec`  | Declares a specification (interface)           |
| `enum`  | Declares an enumeration (sum type)             |
| `type`  | Declares a type alias                          |
| `const` | Declares a compile-time constant               |
| `pub`   | Makes a declaration visible outside its module |
| `mut`   | Marks a variable or method as mutable          |

### Control Flow Keywords

| Keyword    | Description                                     |
| ---------- | ----------------------------------------------- |
| `if`       | Conditional branch                              |
| `else`     | Alternative branch for `if`                     |
| `for`      | Loop construct (iteration, condition, infinite) |
| `match`    | Pattern matching statement                      |
| `return`   | Returns a value from a function                 |
| `break`    | Exits the innermost loop                        |
| `continue` | Skips to the next loop iteration                |

### Concurrency Keywords

| Keyword | Description                               |
| ------- | ----------------------------------------- |
| `go`    | Spawns a concurrent task                  |
| `yield` | Suspends a coroutine and produces a value |

### Error Handling Keywords

| Keyword   | Description                                   |
| --------- | --------------------------------------------- |
| `try`     | Begins an exception-handling block            |
| `expect`  | Catches exceptions by type                    |
| `finally` | Runs cleanup code regardless of exceptions    |
| `raise`   | Raises an exception                           |
| `assert`  | Raises `AssertionError` if condition is false |

### Resource Management Keywords

| Keyword | Description                             |
| ------- | --------------------------------------- |
| `with`  | Manages scoped resources (`Disposable`) |
| `del`   | Deletes a variable binding              |

### Operators

| Keyword | Description                                    |
| ------- | ---------------------------------------------- |
| `and`   | Logical AND (short-circuit)                    |
| `or`    | Logical OR (short-circuit)                     |
| `xor`   | Logical XOR                                    |
| `not`   | Logical NOT                                    |
| `is`    | Type or spec check                             |
| `in`    | Membership test or `for` loop iterator         |
| `as`    | Binds exception to variable in `expect` clause |
| `&`     | Reference (prefix) or bitwise AND (infix)      |

### Literals and Values

| Keyword | Description                            |
| ------- | -------------------------------------- |
| `true`  | Boolean true value                     |
| `false` | Boolean false value                    |
| `nil`   | Absence of value (nullable types only) |
| `nop`   | No-operation statement                 |

### Special

| Keyword  | Description                               |
| -------- | ----------------------------------------- |
| `import` | Imports a module or package               |
| `Self`   | Refers to the implementing type in specs  |
| `this`   | Refers to the current instance in methods |

### Built-in Result Variants

| Keyword | Description                       |
| ------- | --------------------------------- |
| `Ok`    | Success variant of `Result[T, E]` |
| `Err`   | Error variant of `Result[T, E]`   |

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

### Multi-Value Assignment

Multiple variables can be assigned in a single statement using comma-separated lists. Both sides must have
the same number of values. All right-hand-side expressions are evaluated before any assignment occurs,
enabling idiomatic swap without a temporary variable.

```txt
mut a := 10
mut b := 20

# Multi-value assignment (swap)
a, b = b, a
```

Multi-value declaration is not supported -- declare each variable on its own line.

Zerg has no tuple type. For functions that need to return multiple values, return a `list` and unpack it:

```txt
fn divmod(a: int, b: int) -> list[int] {
    return [a // b, a % b]
}

q, r := divmod(17, 5)    # unpack: q = 3, r = 2
```

### Mutable

You can specify the variable with `mut` keyword to make it mutable, which means you can change the value of
the variable after it is assigned. However, the mutable variable should be used sparingly, and you should
prefer to use the immutable variable whenever possible.

### Constants

The `const` keyword declares a compile-time constant. The value must be known at compile time -- literals,
arithmetic on literals, or references to other constants. Constants cannot be reassigned or shadowed.

```txt
const PI = 3.14159
const MAX_SIZE = 1024
pub const VERSION = "1.0.0"
```

### Null-safety

Zerg is designed to be a null-safe programming language. Variables are non-nullable by default, which means
a variable cannot hold `nil` unless its type is explicitly declared as nullable using the `?` suffix.

You can access a nullable variable after checking it is not `nil`, or using the safe navigation operator `?.`
to access the property or method of the nullable variable. You can also use `??` to provide the default value
for a nullable variable, which means if the variable is `nil`, it will return the default value instead of
throwing an exception.

### Nullable Collections

The `?` suffix applies to the complete type. Its position relative to type parameters determines what is
nullable:

```txt
list[string]?          # nullable list -- the list itself can be nil
list[string?]          # list of nullable strings -- elements can be nil
list[string?]?         # nullable list of nullable strings
map[string, int?]      # map with nullable values
```

Placing `?` between the type name and its parameters is a syntax error:

```txt
list?[string]          # syntax error -- ? must come after the complete type
```

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
fn factorial(n: int) -> int {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
```

There are no block comments. Use multiple `#` lines for longer explanations.

## Strings

Strings in Zerg are UTF-8 encoded and delimited by double quotes. They support escape sequences and string
interpolation.

### Escape Sequences

The following escape sequences are recognized inside regular strings:

| Escape     | Description            |
| ---------- | ---------------------- |
| `\\`       | Backslash              |
| `\"`       | Double quote           |
| `\n`       | Newline                |
| `\r`       | Carriage return        |
| `\t`       | Tab                    |
| `\0`       | Null character         |
| `\{`       | Literal `{`            |
| `\}`       | Literal `}`            |
| `\xHH`     | Hex byte (e.g. `\x41`) |
| `\u{XXXX}` | Unicode codepoint      |

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

### Default Parameters

Parameters can have default values. When a default is provided, the caller may omit that argument. Parameters
with defaults must come after all parameters without defaults.

```txt
fn greet(name: string = "world") {
    print("Hello, {name}!")
}

greet()          # Hello, world!
greet("Alice")   # Hello, Alice!
```

### Return Type

A function that declares a return type (`-> type` after the parameter list) must return a value of that type
on every code path. A function with no return type annotation returns nothing -- attempting to use its result
(e.g., `x := foo()`) is a compile-time error. A bare `return` (with no expression) is valid inside a
void-returning function.

### Anonymous Functions

Anonymous functions are function literals without a name. They use the same `fn` syntax as named functions
but omit the function name. Anonymous functions can be assigned to variables, passed as arguments, or
invoked immediately.

```txt
square := fn(x: int) -> int { return x ** 2 }
items.filter(fn(x: int) -> bool { return x > 0 })
go fn() { done <- true }()
```

Parameter types must be declared explicitly. Anonymous functions are first-class values and can capture
variables from the enclosing scope.

### Generic Functions

Functions can have type parameters using `[T]` syntax after the function name. This enables writing reusable
functions that work with multiple types while preserving type safety.

```txt
fn identity[T](x: T) -> T {
    return x
}

fn first[T](items: list[T]) -> T {
    return items[0]
}
```

Type parameters can have constraints by specifying a spec:

```txt
fn max[T: Comparable[T]](a: T, b: T) -> T {
    if a.compare(b) > 0 { return a }
    return b
}
```

### Function Types

Function types are expressed using the `fn` keyword directly, describing the parameter types and return type
of a callable. Any function or method whose signature matches the declared `fn` type can be assigned or passed
directly.

As a convenience, a standalone function implicitly satisfies any single-method `spec` if its signature
matches -- no explicit `implement` declaration is needed. This allows functions to be used interchangeably
with single-method interfaces.

## Control Flow

Zerg provides minimal control flow constructs that favor explicitness over convenience.

### If Statement

The `if` statement executes a block when a condition is true. It is a **statement**, not an expression --
it cannot produce a value. Zerg does not support `else if` or `elif`; use `match` for multi-branch logic.

```txt
if condition {
    # executed when condition is true
}

if condition {
    # true branch
} else {
    # false branch
}
```

For more than two branches, use `match`:

```txt
match value {
    1 => { print("one") }
    2 => { print("two") }
    _ => { print("other") }
}
```

### Break and Continue

The `break` statement exits the innermost loop immediately. The `continue` statement skips the rest of the
current iteration and proceeds to the next one.

```txt
for i in 0..10 {
    if i == 5 { break }       # exit loop when i is 5
    if i % 2 == 0 { continue } # skip even numbers
    print(i)                   # prints 1, 3
}
```

### Nop

The `nop` statement is a no-operation placeholder that does nothing. Use it where a statement is required
but no action is needed.

```txt
if debug { nop }                              # placeholder for future code
callbacks.each(fn(_: int) { nop })            # ignore each element
```

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

```txt
enum Status {
    Pending
    Active(since: int)
    Completed(result: string, code: int)
}
```

Variants without associated data are written as plain identifiers. Variants with data list their fields in
parentheses, similar to function parameters. Enums can also have type parameters:

```txt
enum Option[T] {
    Some(value: T)
    None
}
```

The `enum` keyword is syntax sugar -- the underlying data representation is an implementation detail managed
by the compiler. User code should treat enums as opaque and access variants only through pattern matching.

The `match` statement handles values by pattern matching. It is a **statement**, not an expression -- it
cannot be assigned to a variable. Internally, `match` is syntax sugar for `if-elif-else` chains.

When matching enum variants, the compiler requires all variants to be covered exhaustively. The compiler
rejects any `match` that does not handle every variant, ensuring no case is silently ignored. Each branch
can destructure the variant's associated data and bind it to local variables.

```txt
match result {
    Ok(v) if v > 0 => { print("positive: {v}") }
    Ok(v) => { print("non-positive: {v}") }
    Err(e) => { print(e.message) }
}
```

A match branch can include a **guard** (`if` condition) to further constrain when the branch matches. The
guard can reference variables bound by the pattern. The branch matches only if the pattern matches AND the
guard is true.

The `match` statement also works with non-enum types. Use `_` as the wildcard pattern to match any value:

```txt
match code {
    200 => { print("ok") }
    404 => { print("not found") }
    _ => { print("other: {code}") }
}
```

For non-enum matches, the wildcard `_` is required to ensure exhaustiveness.

The built-in `Result[T, E]` enum has two variants: `Ok(value: T)` representing success, and `Err(error: E)`
representing failure. `Ok` and `Err` are built-in keywords -- they do not need to be imported or defined.
The caller must `match` both variants before accessing the inner value.

### Type Parameters

Zerg supports type parameters using `[T]` syntax. Collection types must be declared with explicit element
types, such as `list[int]`, `map[string, int]`, or `set[string]`. User-defined classes can also accept type
parameters in the same way, allowing reusable and type-safe data structures while preserving compile-time type
checking.

### Type Aliases

The `type` keyword creates a named alias for an existing type. The alias and the original type are fully
interchangeable -- no new type is created.

```txt
type Handler = fn(int) -> string
type StringMap = map[string, string]
```

Type aliases can be generic using type parameters:

```txt
type Pair[T] = list[T]
type Result2[T] = Result[T, string]
```

Aliases can be marked `pub` to export them from the module.

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

```txt
# With assignment -- bind the resource to a variable
with file := open("data.txt") {
    content := file.read()
}

# With expression -- resource is managed but not bound to a name
with open("log.txt") {
    # useful when the resource is used implicitly
}
```

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

Methods are private by default -- they can only be called from within the instance. Use `pub` to make a
method accessible from outside:

```txt
impl Animal {
    fn internal_helper() {         # private -- only callable within Animal methods
        # ...
    }
    pub fn greet() -> string {     # public -- callable from anywhere
        return "Hi, I am {this.name}"
    }
}
```

By default, methods cannot modify the receiver. A method that needs to mutate `this` must be declared with
`mut fn`. The compiler enforces this -- calling a `mut fn` method on an immutable variable is a compile-time
error.

Parameters can also be declared mutable using `mut` before the parameter name:

```txt
fn process(mut x: int) {
    x = x + 1              # OK -- x is mutable inside the function
}                          # caller's value is unchanged (copy semantics)
```

To modify the caller's value, use a reference parameter with `&`:

```txt
fn increment(mut x: &int) {
    x = x + 1              # modifies the caller's value
}

mut n := 10
increment(&n)              # pass reference with &
print(n)                   # 11
```

The caller must explicitly pass `&x` to allow modification, making mutation visible at the call site.

```txt
impl Animal {
    pub fn greet() -> string {
        return "Hi, I am {this.name}"
    }
    pub mut fn birthday() {
        this.age = this.age + 1
    }
}
```

### Spec (Interface)

Zerg uses `spec` to define the interface contract. A class must explicitly declare its spec implementation
using `impl ClassName for SpecName`, and provide all the methods defined in the `spec`. Regular class methods
that do not belong to any spec are defined using `impl ClassName` without the `for` clause. This provides
interface-based polymorphism without class hierarchy complexity.

All methods in a spec are implicitly `pub` -- no visibility modifier is needed or allowed. Methods can be
marked `mut fn` to indicate they mutate the receiver. The special type `Self` refers to the type of the
current instance -- a method using `Self` is always an instance method with access to `this`. See
[SPECS.md](SPECS.md#self-type) for details.

For generic specs, the `impl` declaration specifies the concrete type argument:

```txt
impl Dog for Comparable[Dog] {
    fn compare(other: Dog) -> int {
        return this.age - other.age
    }
}
```

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
the current execution, and handle it using the `try-expect-finally` statement. The `expect` clause catches
exceptions by type using the `as` keyword:

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

The `finally` block runs regardless of whether an exception was raised. Inside an `expect` block, a bare
`raise` (with no expression) re-raises the current exception. See [ERRORS.md](ERRORS.md) for the full
exception hierarchy and user-defined exceptions.

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

### Task Lifetime

The program exits when `main()` returns. All `go` tasks still running at that point are terminated
immediately -- they do not run to completion. If a task must finish its work before the program exits, the
caller should wait for it explicitly using a channel or other synchronization mechanism.

### Channels

Concurrent tasks communicate through `chan[T]`, which are typed conduits for sending and receiving values
between tasks. A channel can be unbuffered (sender blocks until receiver is ready) or buffered (sender blocks
only when the buffer is full). Buffered channels are created with a capacity: `chan[int](10)`.

The `<-` operator is used for channel send and receive:

```txt
done := chan[bool]()
go fn() {
    done <- true           # send a value into the channel
}()
result := <-done           # receive a value from the channel
```

- **Send**: `ch <- value` sends `value` into channel `ch`. Blocks if the channel is full (or unbuffered and
  no receiver is ready). Raises `ChannelClosedError` if the channel is closed.
- **Receive**: `<-ch` receives a value from channel `ch`. Blocks until a value is available. Returns the
  received value.

Sending and receiving do not modify the channel itself -- they pass data _through_ it. Therefore a `chan`
does not require `mut` to be used. An immutable channel variable can still send and receive freely.

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

The `for` loop supports multiple loop variables separated by commas. The number of variables determines what
is bound on each iteration:

```txt
for item in items {}           # element only
for i, item in items {}        # index + element
for key in scores {}           # key only (maps)
for i, ch in "hello" {}        # index + character
```

The following built-in types implement `Iterable` and can be used directly with `for`:

| Type        | 1 Variable | 2 Variables         |
| ----------- | ---------- | ------------------- |
| `list[T]`   | element    | index, element      |
| `map[K, V]` | key        | _(compile error)_   |
| `set[T]`    | element    | _(compile error)_   |
| `string`    | character  | index, character    |
| `range`     | value      | index, value        |
| `chan[T]`   | value      | _(compile error)_   |
| `iter[T]`   | value      | _(depends on iter)_ |

Zerg has no `while` loop. For conditional looping, use `for` with a condition, and for infinite loops use
`for` with an empty clause:

```txt
for condition {}           # loop while condition is true
for {}                     # infinite loop (use break to exit)
```

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
