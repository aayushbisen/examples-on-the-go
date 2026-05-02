# My Go Learning Journey

A personal repository documenting my journey learning Go (Golang), from the basics to advanced concurrency patterns. Each example is a standalone, runnable program with detailed beginner-friendly comments.

## Structure

```
.
├── main.go                          # Hello World & basic Go syntax
├── examples/
│   ├── channels/                    # Worker pool with channels
│   ├── channels-all/                # 9 channel patterns (unbuffered, buffered, fan-in, fan-out, pipeline, etc.)
│   ├── context/                     # 10 context patterns (cancellation, timeout, values, HTTP, graceful shutdown, etc.)
│   ├── go-error-handling/           # Error handling — custom errors, wrapping, defer, panic/recover
│   ├── goroutines/                  # 5 goroutine patterns (basic, WaitGroup, worker pool, select, leak prevention)
│   ├── interfaces/                  # Interfaces — shapes, payment system, type assertions, nil gotcha, composition
│   ├── learn-pointers/              # Pointers — basics, nil, functions, struct fields, common patterns
│   └── structs-and-composition/     # Structs — embedding, receivers, JSON tags, anonymous structs, best practices
└── go.mod                           # Go module definition
```

## How to Run

Each example is a standalone program. Run any of them with:

```bash
# Hello World
go run main.go

# Channels - Worker Pool
go run examples/channels/channels.go

# Channels - All Patterns (9 examples)
go run examples/channels-all/channels.go

# Context - All Patterns (10 examples)
go run examples/context/context.go

# Error Handling
go run examples/go-error-handling/errors.go

# Goroutines - All Patterns (5 examples)
go run examples/goroutines/goroutines.go

# Interfaces
go run examples/interfaces/interfaces.go

# Pointers
go run examples/learn-pointers/pointers_complete.go

# Structs & Composition
go run examples/structs-and-composition/structs_and_composition.go
```

## What I've Learned So Far

### Basics
- [x] Go basics (packages, imports, main function)
- [x] Pointers (address-of, dereference, nil pointers, passing to functions)
- [x] Structs (zero values, initialization, factory functions, JSON tags)
- [x] Anonymous structs (inline definitions, table-driven tests)

### Methods & Receivers
- [x] Value receivers (read-only, small structs)
- [x] Pointer receivers (mutating, large structs)
- [x] Consistency rule (don't mix receivers on the same type)

### Composition
- [x] Struct embedding (anonymous fields, promoted fields/methods)
- [x] Method shadowing (not true overriding)
- [x] Named fields vs embedded fields (when to use which)
- [x] Interface embedding (composition of behaviors)

### Concurrency
- [x] Goroutines (lightweight concurrency, `go` keyword)
- [x] WaitGroup (waiting for goroutines to finish)
- [x] Channels (unbuffered & buffered)
- [x] Directional channels (send-only & receive-only)
- [x] Select statement (multiple channels, timeouts, non-blocking)
- [x] Fan-in / Fan-out patterns
- [x] Pipeline pattern
- [x] Avoiding goroutine leaks (time.After, context)

### Context
- [x] Context cancellation (stop long-running work)
- [x] Context timeout & deadline (auto-cancel after duration or at a specific time)
- [x] Context values (pass request-scoped data like auth tokens)
- [x] Context propagation (flow through function chains)
- [x] Context in HTTP servers (client disconnect detection)
- [x] Graceful shutdown (stop all goroutines with one cancel)
- [x] Context with DB/external calls (query timeouts)
- [x] Context chain (parent cancels all children)
- [x] TODO vs Background (semantic difference)

### Interfaces
- [x] Implicit implementation (no `implements` keyword)
- [x] Empty interface / `any` (accept any type)
- [x] Type assertions (safe vs unsafe)
- [x] Type switches (handling multiple types)
- [x] Interface composition (embedding interfaces)
- [x] Nil interface gotcha (type + value must both be nil)
- [x] Best practices (accept interfaces, return structs; keep interfaces small)

### Error Handling
- [x] Error as a value (not exceptions)
- [x] Custom error types
- [x] Error wrapping (`fmt.Errorf` with `%w`)
- [x] `defer` for cleanup
- [x] Panic and recover

## Resources

- [A Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Blog - Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- [Go Blog - Errors are Values](https://go.dev/blog/errors-are-values)
