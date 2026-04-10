# Zerg Patterns and Limitations

This document lists programming patterns that Zerg's memory model
does not support directly, and the idiomatic workarounds.

## 1. Shared-Node Structures

Zerg's `ptr[T]` is single-owner. Multiple parents cannot share a
child node. This affects graphs, DAGs, and doubly-linked lists.

### Graph with shared nodes

```zerg
# IMPOSSIBLE: shared node between two parents
# nodeA -> nodeC
# nodeB -> nodeC    <-- can't share nodeC

# WORKAROUND: index-based adjacency list
struct Graph {
    nodes: list[NodeData]
    edges: map[int, list[int]]
}

fn add_edge(g: &mut Graph, from: int, to: int) {
    mut neighbors := g.edges[from] ?? []
    neighbors.push(to)
    g.edges[from] = neighbors
}
```

### Doubly-linked list

```zerg
# IMPOSSIBLE: prev and next can't both point to same node
# struct DNode { next: ptr[DNode]?, prev: ptr[DNode]? }

# WORKAROUND: use list[T] — better cache locality anyway
items: list[int] = [1, 2, 3, 4, 5]
```

### Tree with parent pointer

```zerg
# IMPOSSIBLE: child can't reference parent

# WORKAROUND: store parent index
struct TreeNode {
    value: int
    parent_idx: int         # -1 for root
    children: list[int]     # indices into node list
}

struct Tree {
    nodes: list[TreeNode]
}
```

## 2. Mutable Shared State

Closures cannot capture `mut` variables. Shared mutable state
requires explicit concurrency primitives.

### Observer / event listener

```zerg
# IMPOSSIBLE: observers holding back-references to subject
# struct Subject { observers: list[&Observer] }

# WORKAROUND: channel-based events
event_ch := chan[Event](100)

# Observer task
rush |ch| {
    for event in ch {
        handle(event)
    }
}(event_ch)

# Subject sends events
event_ch <- Event { kind: "click", data: "button1" }
```

### Global mutable singleton

```zerg
# IMPOSSIBLE: pub mut at module scope
# pub mut instance := Config { ... }

# WORKAROUND: sync[T] at module scope
instance := sync[Config](Config { debug: false })

pub fn get_config() -> Config {
    return instance.read()
}

pub fn set_debug(on: bool) {
    instance.lock(|c: &mut Config| { c.debug = on })
}
```

### Connection pool

```zerg
# WORKAROUND: channel as pool
fn make_pool(size: int) -> chan[Connection] {
    pool := chan[Connection](size)
    for _ in 0..size {
        pool <- create_connection()
    }
    return pool
}

# Acquire and release
conn := <- pool         # blocks if empty
defer { pool <- conn }  # return to pool on scope exit
use(conn)
```

## 3. Closure Patterns

### Mutable accumulator

```zerg
# IMPOSSIBLE: closure can't capture mut
# mut total := 0
# nums.for_each(|| { total += 1 })

# WORKAROUND: explicit loop
mut total := 0
for n in nums {
    total += n
}

# WORKAROUND: concurrent accumulation with sync
total := sync[int](0)
for _ in 0..4 {
    rush |t, data| {
        t.lock(|v: &mut int| { v += sum(data) })
    }(total, chunk)
}
```

### Stateful callback

```zerg
# IMPOSSIBLE: callback can't hold mutable state
# mut count := 0
# on_click := || { count += 1 }

# WORKAROUND: callback receives &mut
fn on_click(count: &mut int) {
    count += 1
}

mut click_count := 0
on_click(&mut click_count)
on_click(&mut click_count)
print click_count              # 2
```

### Builder pattern

```zerg
# IMPOSSIBLE: fluent builder needs mutation + return self
# builder.set_x(1).set_y(2).build()

# WORKAROUND: &mut builder
struct ConfigBuilder {
    host: str
    port: int
}

fn set_host(b: &mut ConfigBuilder, host: str) {
    b.host = host
}

fn set_port(b: &mut ConfigBuilder, port: int) {
    b.port = port
}

mut b := ConfigBuilder { host: "", port: 0 }
set_host(&mut b, "localhost")
set_port(&mut b, 8080)
config := b     # immutable copy of final state
```

## 4. Performance Patterns

### Zero-copy slice

```zerg
# IMPOSSIBLE: no borrowed sub-range type
# fn first_three(items: list[int]) -> list[int] {
#     return items[0..3]    # creates a NEW list
# }

# WORKAROUND: pass range indices
fn sum_range(items: list[int], start: int, end: int) -> int {
    mut total := 0
    for i in start..end {
        total += items[i]
    }
    return total
}
```

### Object pool (reuse allocations)

```zerg
# IMPOSSIBLE: objects freed at scope exit, can't return to pool

# WORKAROUND: channel as pool (see connection pool above)
buf_pool := chan[list[byte]](16)
for _ in 0..16 {
    buf_pool <- make_buffer(4096)
}

# Use and return
buf := <- buf_pool
defer {
    buf_pool <- buf
}
write_data(&mut buf)
```

## 5. Design Pattern Alternatives

| Classic OOP pattern | Zerg idiomatic alternative              |
| ------------------- | --------------------------------------- |
| Inheritance         | `spec` composition + `impl`             |
| Visitor             | function taking `&mut Node`             |
| Observer            | channel-based event dispatch            |
| Singleton           | `sync[T]` at module scope               |
| Factory             | function returning struct               |
| Iterator            | `Iterable[T]` spec                      |
| Strategy            | function type `fn(T) -> U` or lambda    |
| Command (undo)      | store data, replay from command list    |
| Mediator            | `rush` task with channel hub            |
| Decorator           | wrapper struct implementing same `spec` |

## 6. Summary

| Category         | Limitation          | Root cause            | Workaround             |
| ---------------- | ------------------- | --------------------- | ---------------------- |
| Shared nodes     | no shared ownership | `ptr[T]` single-owner | index-based patterns   |
| Mutable closures | can't capture `mut` | ownership safety      | `&mut` params or loops |
| Shared mut state | no global mutable   | concurrency safety    | `sync[T]` or `chan[T]` |
| Back-references  | no cycles possible  | no shared ownership   | parent index           |
| Zero-copy views  | no slice/view type  | value semantics       | pass indices           |
| Object reuse     | freed at scope exit | deterministic cleanup | channel-based pool     |

These are conscious trade-offs for:
no GC pauses, no data races, no dangling pointers, no borrow checker
complexity, and fully deterministic memory management.
