# Zerg Memory Management Specification

Zerg uses **scope-based ownership** with copy-by-value semantics.
No garbage collector. No manual memory management. Memory is freed
deterministically when the owning scope exits.

## 1. Core Principles

| #   | Principle            | Rule                                           |
| --- | -------------------- | ---------------------------------------------- |
| P1  | Copy-by-value        | assignment produces independent copies         |
| P2  | Immutable by default | shared access is safe without synchronization  |
| P3  | Caller owns memory   | caller allocates, callees borrow, caller frees |
| P4  | No GC                | no tracing GC; refcount only for concurrency   |
| P5  | Deterministic        | `defer` for cleanup; refcount frees at zero    |

## 2. Ownership

Every value has exactly **one owner**: the scope that declared it.
Ownership transfers only through `return` (copy to caller). There is
no shared ownership and no tracing GC. The only exception: immutable
data crossing task boundaries uses reference counting (see Section 6).

```zerg
fn make_list() -> list[int] {
    result := [1, 2, 3]      # make_list owns result
    return result            # Return slot pointer - copy to caller's scope
}

items := make_list()         # caller owns this copy
process(items)               # process borrows items (immutable ref)
# items freed here — caller's scope exits
```

### Ownership rules

| Rule         | Description                                          |
| ------------ | ---------------------------------------------------- |
| Single owner | every value has exactly one owning scope             |
| Scope-bound  | values are freed when their owning scope exits       |
| No transfer  | ownership does not transfer to callees (they borrow) |
| Return slot  | `return` writes directly into caller's memory (sret) |

### Module-level variables

Module-level variables are **always immutable**. `mut` is forbidden
at module scope — mutable state must live inside function scopes.

```zerg
# module.zg
name := "MyModule"          # OK — immutable, lives for program lifetime
config := load_defaults()   # OK — immutable, initialized once

mut counter := 0            # ERROR: mut not allowed at module scope
```

Module variables are owned by the runtime and freed at program exit.
Since they are immutable, they are safe to read from any task without
synchronization. This prevents data races on global state entirely.

| Rule             | Description                             |
| ---------------- | --------------------------------------- |
| Always immutable | `mut` forbidden at module scope         |
| Runtime-owned    | freed at program exit, never early      |
| Concurrent-safe  | immutable — no synchronization needed   |
| Initialized once | `_init()` then `init()` on first import |

## 3. Calling Conventions

### 3.1 Immutable parameters — pass by reference

When a function declares `fn f(x: T)`, the caller passes an immutable
reference. The callee cannot modify the data. This is safe because
immutable references cannot cause mutation.

```zerg
fn sum(items: list[int]) -> int {
    mut total := 0
    for x in items {
        total += x
    }
    return total
}

nums := [1, 2, 3]
result := sum(nums)    # no copy — passes immutable ref to nums
```

### 3.2 Mutable parameters — pass by explicit reference

When a function declares `fn f(x: &mut T)`, the caller passes a mutable
reference with `&mut` at the call site. Changes are visible to the caller.

```zerg
fn append(items: &mut list[int], val: int) {
    items.push(val)    # modifies caller's list
}

mut nums := [1, 2, 3]
append(&mut nums, 4)   # explicit — caller knows mutation happens
# nums is now [1, 2, 3, 4]
```

### 3.3 Return values — caller-allocated slot

When the caller assigns a return value (`x := f()`), the compiler
passes a hidden pointer to the callee. The callee constructs the
result directly in the caller's memory — zero copy.

```zerg
# What the programmer writes:
x := make_list(1000)

# What the compiler generates:
#   x: list[int] = <uninitialized>
#   make_list(&x, 1000)    # hidden pointer — build in caller's space
```

### 3.4 Summary

| Declaration       | Call site    | Mechanism             | Copy cost |
| ----------------- | ------------ | --------------------- | --------- |
| `fn f(x: T)`      | `f(y)`       | immutable reference   | zero      |
| `fn f(x: &mut T)` | `f(&mut y)`  | mutable reference     | zero      |
| `-> T` (return)   | `x := f()`   | return slot pointer   | zero      |
| `-> T` (return)   | `f()` (temp) | stack temporary + ref | zero      |

For types containing `ptr[T]`, the same rules apply — the borrow
follows the containing struct. The `ptr` chain is not traversed at
the call boundary; it is just part of the borrowed data.

## 4. Assignment Semantics

### 4.1 Immutable to immutable — share reference

```zerg
a := [1, 2, 3]
b := a              # shares a's data (just a pointer copy)
# both immutable — no mutation possible, sharing is always safe
```

Since both variables are immutable, the compiler simply copies the
reference (8 bytes). No deep copy, no COW, no refcount.

**Lifetime constraint**: immutable reference sharing only occurs between
variables in the **same scope**. The original owner is freed at scope
exit; all shared references have the same lifetime. This is guaranteed
because all other contexts produce copies, not shares:

| Context             | Mechanism                       | Share or copy? |
| ------------------- | ------------------------------- | -------------- |
| Same-scope `b := a` | share reference                 | share          |
| `return a`          | return slot pointer             | copy           |
| `fn(a)` param       | borrow (caller outlives callee) | borrow         |
| `ch <- a` (channel) | refcount                        | refcount       |
| `rush` capture      | refcount                        | refcount       |
| Escaping closure    | deep copy                       | copy           |

The shared reference `b` can never outlive `a` because no mechanism
produces a shared reference across scope boundaries. For types with
`ptr[T]` fields, the entire heap chain is covered by this rule — `b`
shares the root, which owns the chain, and both die at the same time.

### 4.2 Immutable to mutable — deep copy

```zerg
a := [1, 2, 3]
mut b := a          # deep copy — b needs independent storage
b.push(4)           # safe — b has its own data
# a is still [1, 2, 3]
```

Mutable variables always get independent storage via deep copy.
The compiler may warn if a variable is declared `mut` but never mutated.

When the source variable is never used after the assignment, the
compiler may optimize `mut b := a` into a **move** (transfer the
backing storage instead of copying). The programmer sees copy-by-value
semantics — the optimization is transparent.

```zerg
a := [1, 2, 3]
mut b := a          # compiler sees a is never used after this
                    # optimized to move — zero copy
b.push(4)
# a is no longer accessible (consumed by move)
```

### 4.3 Reassignment — deep copy

```zerg
mut x := [1, 2]
x = [3, 4, 5]      # old value freed, new value deep-copied
```

### 4.4 Summary

| Assignment                | Behavior        | Reason                         |
| ------------------------- | --------------- | ------------------------------ |
| `b := a` (both immutable) | share reference | no mutation — sharing is safe  |
| `mut b := a`              | deep copy       | mutable needs independent data |
| `b = expr` (reassignment) | deep copy       | replaces existing value        |

For types with `ptr[T]` fields, deep copy recursively copies the
entire ownership chain. Immutable sharing follows the pointer chain.

## 5. Closure Capture

Closures can **only capture immutable variables**. Capturing a mutable
variable is a **compile error**.

```zerg
x := 42
f := |y| x + y          # OK — x is immutable

mut count := 0
g := || count + 1        # ERROR: cannot capture mutable variable
```

### 5.1 Why no mutable capture

| Problem             | Description                                       |
| ------------------- | ------------------------------------------------- |
| Ownership confusion | who owns the data — the closure or the caller?    |
| Dangling references | escaping closure may outlive the mutable variable |

Mutable data flows through explicit `&mut` parameters instead:

```zerg
mut count := 0
inc := |c: &mut int| { c += 1 }
inc(&mut count)          # explicit — caller controls mutation
```

### 5.2 Capture lifetime (escape analysis)

The compiler determines whether a closure escapes its declaring scope:

| Closure lifetime                  | Capture mechanism     | Cost |
| --------------------------------- | --------------------- | ---- |
| Non-escaping (used in same scope) | reference to original | zero |
| Escaping (returned from function) | deep copy at creation | O(n) |
| Escaping (sent to channel)        | share ref + refcount  | O(1) |
| Escaping (captured by `rush`)     | share ref + refcount  | O(1) |

A closure **escapes** if it is: returned from a function, stored in a
data structure that outlives the scope, sent through a channel, or
captured by a `rush` task.

```zerg
# Non-escaping — holds reference, zero cost
x := [1, 2, 3]
doubled := x.map(|v| v * 2)

# Escaping — deep copies x into the closure
fn make_adder(x: int) -> fn(int) -> int {
    return |y| x + y
}
```

## 6. Concurrency

Data crossing task boundaries uses **reference counting for immutable
data** and **deep copy for mutable data**. This is the only place in
Zerg where reference counting is used.

### 6.1 Immutable data — reference counted sharing

Immutable data can safely be shared across tasks because no task can
mutate it. A reference count tracks how many tasks hold a reference.
The data is freed when the last reference drops.

```zerg
big := make_huge_list(1_000_000)    # immutable, refcount = 1
ch <- big                           # share ref, refcount = 2
# caller scope exits → refcount = 1 (data survives)

rush || {
    data := <- ch                   # refcount stays 1 (moved from ch)
    print data.length()             # read-only access
    # task exits → refcount = 0 → data freed
}()
```

Immutable data may outlive its declaring scope when shared via
channels or `rush` captures. The runtime frees it when the last
reference drops — still deterministic, just not scope-bound.

### 6.2 Mutable data — always deep copy

Mutable data must be deep-copied because each task needs independent
storage for mutation.

```zerg
mut items := [1, 2, 3]
rush || {
    # items is deep-copied — task gets independent data
    items.push(4)                   # ERROR: can't capture mut
}()
```

Since closures cannot capture mutable variables, mutable data only
crosses task boundaries through channels:

```zerg
mut items := [1, 2, 3]
ch <- items                         # deep copy into channel
# items still accessible and mutable
```

### 6.3 Summary

| Operation         | Immutable data       | Mutable data  |
| ----------------- | -------------------- | ------------- |
| `ch <- x` (send)  | share ref (refcount) | deep copy     |
| `<- ch` (receive) | transfer ref         | transfer copy |
| `rush` capture    | share ref (refcount) | compile error |

For types with `ptr[T]` chains: immutable sharing follows the entire
chain (refcount on the root). Mutable deep copy recurses the chain.

### 6.4 Reference counting rules

Reference counting applies **only** at task boundaries:

| Rule      | Description                                     |
| --------- | ----------------------------------------------- |
| Scope     | refcount only for immutable data crossing tasks |
| Increment | on channel send or `rush` capture               |
| Decrement | on scope exit or channel receive                |
| Free      | when refcount reaches zero                      |
| Not used  | for same-scope assignment (share ref, no count) |
| Not used  | for mutable data (always deep copy)             |

### Channel ownership

`chan[T]` is a **runtime resource**, not a value. It cannot be assigned
or copied. It is created once and passed by reference to tasks.

```zerg
ch := chan[int](10)            # caller owns the channel
rush |c| { c <- 42 }(ch)      # task borrows ch (immutable ref)
value := <- ch                 # caller receives a copy
```

### Rush task exceptions

`rush` is fire-and-forget. If a task raises an unhandled exception:

1. The exception is logged to stderr
2. All `defer` blocks in the task run (LIFO)
3. The task exits silently
4. The caller is not affected

```zerg
rush || {
    defer cleanup()                # runs even on unhandled exception
    raise IOError { msg: "fail", path: "/tmp" }
    # exception logged to stderr, task dies, defer runs
}()
# caller continues — unaware of the failure
```

Tasks that need error reporting should send errors through channels:

```zerg
err_ch := chan[Exception]()
rush |ch| {
    try {
        do_risky_work()
    } except Exception as e {
        ch <- e                    # report error to caller
    }
}(err_ch)
```

### Main task exit

When the main task exits, all running `rush` tasks are terminated.
`defer` blocks in terminated tasks do **not** run. To ensure cleanup,
the main task should wait on a done channel before exiting:

```zerg
done := chan[bool]()
rush || {
    defer close_resources()
    do_work()
    done <- true
}()
<- done                            # wait for task to finish
```

## 7. Type Categories

### 7.1 Primitives

Stack-allocated. Bitwise copy. No heap involvement.

| Type    | Size    |
| ------- | ------- |
| `int`   | 8 bytes |
| `float` | 8 bytes |
| `bool`  | 1 byte  |
| `byte`  | 1 byte  |
| `rune`  | 4 bytes |

### 7.2 Strings

`str` is **always immutable**. The compiler may share backing buffers
(interning). Strings are safe to share in all contexts, including
across task boundaries — no deep copy needed.

### 7.3 Collections

| Type         | Immutable assign | Mutable assign | Concurrency |
| ------------ | ---------------- | -------------- | ----------- |
| `list[T]`    | share ref        | deep copy      | deep copy   |
| `map[K, V]`  | share ref        | deep copy      | deep copy   |
| `set[T]`     | share ref        | deep copy      | deep copy   |
| `tuple[...]` | share ref        | deep copy      | deep copy   |

### 7.4 Structs

Structs follow their field rules recursively. A struct with only
primitive fields is a trivial bitwise copy. A struct containing
collections triggers deep copy of those fields when required.

```zerg
struct Point { x: int, y: int }        # 16 bytes — trivial copy
struct Config { items: list[str] }      # follows list copy rules
```

### 7.5 Owned pointers

`ptr[T]` is an owned heap pointer. Enables recursive types (trees,
linked lists). The pointer owns the pointed-to value — freed when
the owner's scope exits.

```zerg
struct Node {
    value: int
    next: ptr[Node]?       # fixed size (pointer = 8 bytes)
}

n := Node { value: 1, next: ptr(Node { value: 2, next: nil }) }
print n.next?.value         # 2 (auto-deref)
```

| Assignment              | Behavior                            |
| ----------------------- | ----------------------------------- |
| `b := a` (immutable)    | share reference (safe)              |
| `mut b := a`            | deep copy (entire chain)            |
| Concurrency (immutable) | share ref + refcount (entire chain) |
| Concurrency (mutable)   | deep copy (entire chain)            |

No cycles possible — each `ptr` owns an independent copy.
Trees and linked lists work. Shared-node graphs do not.

For graphs and other shared-node structures, use index-based patterns:

```zerg
# ECS-style: entities are indices, components are parallel arrays
struct World {
    positions: list[Vec2]
    velocities: list[Vec2]
    healths: list[int]
}

# Graph: adjacency list with integer node IDs
struct Graph {
    nodes: list[NodeData]
    edges: map[int, list[int]]
}

fn neighbors(g: Graph, id: int) -> list[int] {
    return g.edges[id] ?? []
}
```

This pattern avoids shared ownership entirely. Nodes are values in a
flat collection, referenced by index. Used by ECS game engines,
relational databases, and compilers. No `ptr` needed — just `list`
and `map`.

### 7.6 Channels

`chan[T]` is a runtime resource. Cannot be copied or assigned.
Created once, passed by reference, freed when the owning scope exits.

## 8. Compiler Optimizations

All optimizations are transparent — the programmer sees copy-by-value.

| Optimization            | Description                                      |
| ----------------------- | ------------------------------------------------ |
| **Return slot pointer** | callee writes directly into caller's memory      |
| **Pass-by-reference**   | immutable params passed as references            |
| **Reference sharing**   | immutable-to-immutable shares via pointer copy   |
| **Move on last use**    | `mut b := a` becomes move when `a` unused after  |
| **Stack allocation**    | non-escaping values stay on the stack            |
| **String interning**    | immutable string literals shared at compile time |
| **Escape analysis**     | determines closure capture and allocation site   |
| **Recursive free**      | `ptr[T]` chains freed bottom-up at scope exit    |

### Return slot pointer detail

```zerg
# Chained calls — temporaries on caller's stack
process(make_list(100))
# compiler: tmp: list[int]; make_list(&tmp, 100); process(&tmp)
```

| Code                   | Generated                        | Copy |
| ---------------------- | -------------------------------- | ---- |
| `x := run()`           | `run(&x)`                        | zero |
| `x := Point{x:1, y:2}` | construct in-place               | zero |
| `return result`        | write to caller's slot           | zero |
| `print run()`          | `tmp: T; run(&tmp); print(&tmp)` | zero |

## 9. Borrowing Rules

The compiler enforces at compile time (zero runtime overhead):

| Rule                            | Enforcement                       |
| ------------------------------- | --------------------------------- |
| Immutable params cannot mutate  | type system                       |
| `&mut` is exclusive             | no other refs while `&mut` active |
| `&mut` requires `mut` variable  | cannot reference immutable        |
| Borrows cannot outlive owner    | scope-checked                     |
| Closures capture immutable only | mutable capture is compile error  |
| Concurrency always deep copies  | cannot borrow across tasks        |
| `ptr[T]` owns its target        | no shared pointers, no cycles     |

## 10. Resource Cleanup

Since there is no GC, deterministic cleanup uses `defer`:

```zerg
fn read_file(path: str) -> str {
    fd := open(path)
    defer close(fd)          # runs when scope exits
    return fd.read_all()
}
```

`defer` runs in LIFO order at scope exit, regardless of how the scope
is left (normal return, early return via guard, exception). This replaces
GC finalizers, destructors, and context managers.

## 11. Architecture

```text
┌─────────────────────────────────────────────────────┐
│                    Zerg Runtime                     │
│                                                     │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐         │
│  │  Task 1  │   │  Task 2  │   │  Task N  │         │
│  │  stack   │   │  stack   │   │  stack   │         │
│  │  locals  │   │  locals  │   │  locals  │         │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘         │
│       │              │              │               │
│       └───────── channels ──────────┘               │
│          (refcount immut / deep copy mut)           │
│                                                     │
│  ┌───────────────────────────────────────────────┐  │
│  │  Heap (scope-owned allocations)               │  │
│  │  - collection backing storage                 │  │
│  │  - ptr[T] targets (owned heap pointers)       │  │
│  │  - large structs                              │  │
│  │  - closure captures (deep copy for escaping)  │  │
│  │  - channel buffers                            │  │
│  │                                               │  │
│  │  Freed by: scope exit (deterministic)         │  │
│  │  No GC. No finalizers.                        │  │
│  │  Refcount only for cross-task immutable data  │  │
│  └───────────────────────────────────────────────┘  │
│                                                     │
│  ┌───────────────────────────────────────────────┐  │
│  │  Compile-time analysis (transparent)          │  │
│  │  - Escape analysis for closures + allocation  │  │
│  │  - Borrow checking for &mut exclusivity       │  │
│  │  - Return slot pointer insertion              │  │
│  │  - Reference sharing for immutable assignment │  │
│  │  - Move optimization on last use              │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### Data flow for function calls

```text
  Caller                            Callee
  ┌──────────┐                     ┌──────────┐
  │ x: T     │─── immutable ref ──▶│ x: T     │  (read-only)
  │          │                     │          │
  │ mut y: T │─── &mut ref ───────▶│ y: &mut T│  (read-write)
  │          │                     │          │
  │ z: T     │◀── return slot ─────│ result   │  (callee writes into z)
  └──────────┘                     └──────────┘
```

### Data flow for concurrency

```text
  Caller Task                     Rush Task
  ┌──────────┐                    ┌──────────┐
  │ x (immut)│── share ref+rc ───▶│ x (immut)│  (refcounted)
  │ y (mut)  │── deep copy ──────▶│ y (copy) │  (independent)
  │          │                    │          │
  │ ch: chan │◀── share/copy ─────│ result   │  (via channel)
  └──────────┘                    └──────────┘
```

## 12. Comparison

|             | Zerg              | Go           | Rust             | Swift        |
| ----------- | ----------------- | ------------ | ---------------- | ------------ |
| Memory      | scope ownership   | tracing GC   | ownership+borrow | ARC          |
| Params      | borrow            | copy/pointer | borrow/move      | copy/inout   |
| Mutation    | `&mut` explicit   | pointers     | `&mut` borrow ck | `inout`      |
| Assignment  | share / deep copy | shallow+GC   | move             | COW          |
| Concurrency | refcount+copy     | GC+detector  | Send/Sync        | actors       |
| Cleanup     | `defer`           | `defer`+GC   | `Drop`           | `deinit`+ARC |
| GC pauses   | none              | yes (sub-ms) | none             | none (ARC)   |
| Cycles      | impossible        | GC handles   | borrow ck        | weak refs    |
