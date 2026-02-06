# Unsafe and Low-Level Operations

Zerg is designed to be safe by default, but sometimes low-level operations are necessary for system programming,
FFI (Foreign Function Interface), or performance-critical code. The `unsafe` block and `asm` expression provide
controlled escape hatches for these scenarios.

## Overview

| Construct            | Purpose                                                | Example          |
| -------------------- | ------------------------------------------------------ | ---------------- |
| `unsafe { }`         | Marks a block where low-level operations are permitted | `unsafe { ... }` |
| `asm("fn", args...)` | Calls a registered runtime function directly           | `asm("sys_os")`  |

## Unsafe Block

The `unsafe` keyword introduces a block where certain low-level operations are permitted. Code inside an `unsafe`
block can use `asm()` expressions to call runtime functions directly.

```zerg
unsafe {
    os := asm("sys_os")
    arch := asm("sys_arch")
    print("Running on " + os + "/" + arch)
}
```

### Why Unsafe?

The `unsafe` block serves several purposes:

1. **Visibility**: Makes low-level code easy to locate and audit
2. **Intentionality**: Requires explicit opt-in for potentially dangerous operations
3. **Containment**: Limits the scope where special operations can occur
4. **Documentation**: Signals to readers that careful attention is needed

### Rules

- `asm()` expressions can **only** be used inside `unsafe` blocks
- Calling `asm()` outside an `unsafe` block is a compile-time error
- `unsafe` blocks can contain any valid Zerg code, not just `asm` calls
- Variables declared inside `unsafe` are scoped to the block (standard scoping rules apply)

## Asm Expression

The `asm()` expression calls a registered runtime function by name. The first argument is a string literal
specifying the function name, followed by any arguments the function requires.

```zerg
unsafe {
    result := asm("function_name", arg1, arg2, ...)
}
```

### Available Functions

The following functions are available through `asm()` in the bootstrap runtime:

#### System Functions

| Function   | Arguments      | Returns  | Description                                          |
| ---------- | -------------- | -------- | ---------------------------------------------------- |
| `sys_os`   | none           | `string` | Operating system: `"linux"`, `"darwin"`, `"windows"` |
| `sys_arch` | none           | `string` | Architecture: `"amd64"`, `"arm64"`, `"386"`          |
| `sys_args` | none           | `list`   | Command-line arguments                               |
| `sys_exit` | `code: int`    | never    | Exit program with code                               |
| `sys_env`  | `name: string` | `string` | Get environment variable                             |

#### File Functions

| Function      | Arguments                      | Returns  | Description                             |
| ------------- | ------------------------------ | -------- | --------------------------------------- |
| `file_open`   | `path: string`, `mode: string` | `handle` | Open file (`"r"`, `"w"`, `"a"`, `"rw"`) |
| `file_read`   | `handle`                       | `string` | Read entire file contents               |
| `file_write`  | `handle`, `data: string`       | `int`    | Write string, return bytes written      |
| `file_close`  | `handle`                       | `nil`    | Close file handle                       |
| `file_exists` | `path: string`                 | `bool`   | Check if file exists                    |
| `read_file`   | `path: string`                 | `string` | Read file contents directly             |
| `write_file`  | `path: string`, `data: string` | `nil`    | Write file contents directly            |

#### String Functions

| Function          | Arguments                                 | Returns  | Description                            |
| ----------------- | ----------------------------------------- | -------- | -------------------------------------- |
| `str_split`       | `s: string`, `sep: string`                | `list`   | Split string by separator              |
| `str_join`        | `list`, `sep: string`                     | `string` | Join list elements with separator      |
| `str_trim`        | `s: string`                               | `string` | Remove leading/trailing whitespace     |
| `str_find`        | `s: string`, `sub: string`                | `int`    | Find substring index (-1 if not found) |
| `str_replace`     | `s: string`, `old: string`, `new: string` | `string` | Replace all occurrences                |
| `str_substring`   | `s: string`, `start: int`, `end: int`     | `string` | Extract substring                      |
| `str_starts_with` | `s: string`, `prefix: string`             | `bool`   | Check if string starts with prefix     |
| `str_ends_with`   | `s: string`, `suffix: string`             | `bool`   | Check if string ends with suffix       |
| `str_upper`       | `s: string`                               | `string` | Convert to uppercase                   |
| `str_lower`       | `s: string`                               | `string` | Convert to lowercase                   |
| `str_contains`    | `s: string`, `sub: string`                | `bool`   | Check if string contains substring     |

#### Character Functions

| Function        | Arguments   | Returns  | Description                               |
| --------------- | ----------- | -------- | ----------------------------------------- |
| `char_ord`      | `c: string` | `int`    | Get ASCII/Unicode code of first character |
| `char_chr`      | `code: int` | `string` | Convert code to character                 |
| `char_is_digit` | `c: string` | `bool`   | Check if character is a digit             |
| `char_is_alpha` | `c: string` | `bool`   | Check if character is a letter            |
| `char_is_space` | `c: string` | `bool`   | Check if character is whitespace          |
| `char_is_alnum` | `c: string` | `bool`   | Check if character is alphanumeric        |

### Error Handling

- Calling an unknown function name returns an error
- Type mismatches in arguments return an error
- File operation failures return descriptive error messages

## Modules vs Unsafe

For most operations, prefer the built-in modules (`sys`, `io`, `str`, `char`) over `unsafe` + `asm`. The modules
provide a safe, high-level API for common operations.

| Use Case               | Preferred Approach            | When to Use `unsafe`        |
| ---------------------- | ----------------------------- | --------------------------- |
| Read a file            | `io.read_file(path)`          | Never for basic I/O         |
| Get OS name            | `sys.os()`                    | Never for system info       |
| String manipulation    | `str.split()`, `str.upper()`  | Never for string ops        |
| Platform-specific code | Module + `if sys.os() == ...` | Raw syscalls, FFI           |
| Performance-critical   | Use modules first             | Profile-driven optimization |

### Example: Safe vs Unsafe

```zerg
# Preferred: Using modules (safe, readable)
os := sys.os()
content := io.read_file("data.txt")
parts := str.split(content, "\n")

# Alternative: Using unsafe (explicit, low-level)
unsafe {
    os := asm("sys_os")
    content := asm("read_file", "data.txt")
    parts := asm("str_split", content, "\n")
}
```

Both approaches produce identical results, but the module version is clearer and doesn't require `unsafe`.

## Use Cases for Unsafe

The `unsafe` block is intended for:

1. **FFI/Interop**: Calling external C libraries or system APIs
2. **Low-level optimization**: When profiling shows a hot path needs direct access
3. **Platform-specific hacks**: Workarounds for OS-specific behaviors
4. **Bootstrap/Compiler development**: Building the compiler in Zerg itself

## Future Extensions

In the self-hosted Zerg compiler, `unsafe` blocks may enable additional capabilities:

- Direct memory access
- Raw pointer manipulation
- Inline assembly for specific architectures
- Custom syscall wrappers

For the bootstrap interpreter, `asm` provides access to Go runtime functions. The self-hosted compiler will
extend this to native code generation.

## Best Practices

1. **Minimize unsafe code**: Keep `unsafe` blocks as small as possible
2. **Document intent**: Explain why `unsafe` is necessary with comments
3. **Wrap in safe functions**: Expose safe APIs that internally use `unsafe`
4. **Test thoroughly**: Unsafe code requires extra testing attention
5. **Prefer modules**: Use `sys`, `io`, `str`, `char` modules when possible

```zerg
# Good: Small, documented unsafe block wrapped in a safe function
fn get_platform() -> string {
    # Using unsafe for direct runtime access
    unsafe {
        os := asm("sys_os")
        arch := asm("sys_arch")
        return os + "-" + arch
    }
}

# Usage is completely safe
platform := get_platform()
print("Platform: " + platform)
```
