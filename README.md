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
│   ├── panic_defer_recover/         # Defer, panic, recover — comprehensive single-file example
│   ├── waitgroup/                   # Deep-dive WaitGroup tutorial (concepts, examples, common mistakes)
│   ├── interfaces/                  # Interfaces — shapes, payment system, type assertions, nil gotcha, composition
│   ├── json/                       # JSON — marshaling, struct tags, omitempty, custom marshalers, streaming, RawMessage, unknown structures
│   ├── learn-pointers/              # Pointers — basics, nil, functions, struct fields, common patterns
│   ├── mutex/                       # Mutex — race conditions, sync.Mutex, sync/RWMutex, sync/atomic, real-world examples
│   ├── slices-and-maps/             # Arrays, slices — declaration, initialization, iteration, value types, make
│   └── structs-and-composition/     # Structs — embedding, receivers, JSON tags, anonymous structs, best practices
├── exercises/                      # Practice exercises for learning concepts
│   ├── safe-divider/              # Error handling — custom sentinel errors
│   ├── converter/                # Temperature conversion functions
│   ├── shape-area/                # Interfaces — Shape polymorphism
│   ├── config-loader/             # JSON parsing — encoding/json
│   ├── tasks/                     # Slice operations
│   ├── website-checker/            # Concurrency — goroutines, WaitGroup, channels
│   ├── log-processor/             # HTTP handler with channel pipeline
│   ├── sentinel/                 # Package structure — internal packages
│   └── logstream/                 # Pipeline pattern — producer, transformer, consumer
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

# Defer, Panic, Recover - Comprehensive Example
go run examples/panic_defer_recover/main.go

# Goroutines - All Patterns (5 examples)
go run examples/goroutines/goroutines.go

# Interfaces
go run examples/interfaces/interfaces.go

# JSON
go run examples/json/json.go

# Pointers
go run examples/learn-pointers/pointers_complete.go

# Structs & Composition
go run examples/structs-and-composition/structs_and_composition.go

# Mutex & Synchronization
go run examples/mutex/main.go

# Slices & Maps
go run examples/slices-and-maps/main.go

# WaitGroup Tutorial (deep-dive with examples & common mistakes)
go run examples/waitgroup/main.go

# Run with race detector (to see race conditions)
go run -race examples/mutex/main.go

# Exercises
cd exercises/safe-divider && go run divider.go
cd exercises/converter && go run converter.go
cd exercises/shape-area && go run shape-area.go
cd exercises/config-loader && go run loader.go
cd exercises/tasks && go run tasks.go
cd exercises/website-checker && go run checker.go
cd exercises/log-processor && go run processor.go
cd exercises/sentinel/cmd/sentinel && go run main.go
cd exercises/logstream/cmd/logstream && go run main.go
```

## What I've Learned So Far

### Basics
- [x] Go basics (packages, imports, main function)
- [x] Pointers (address-of, dereference, nil pointers, passing to functions)
- [x] Structs (zero values, initialization, factory functions, JSON tags)
- [x] Anonymous structs (inline definitions, table-driven tests)
- [x] Arrays (declaration, initialization, compiler-inferred length, value types)
- [x] Slices (declaration, nil value, make function, underlying array/pointer/length/capacity)

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
- [x] WaitGroup (waiting for goroutines to finish) — see `examples/waitgroup/` for deep-dive tutorial
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
- [x] Defer, panic, recover comprehensive example (`examples/panic_defer_recover/`)

### JSON (encoding/json)
- [x] Marshal/Unmarshal (structs, maps, primitives)
- [x] Struct tags (`json:"name"`, `json:"-"`, `json:",omitempty"`)
- [x] Custom marshaling (MarshalJSON/UnmarshalJSON)
- [x] `json.RawMessage` (delayed parsing, embed raw JSON)
- [x] `json.Encoder` / `json.Decoder` (streaming)
- [x] Decoding unknown structures (`map[string]any`)
- [x] `json.Number` (preserve precision vs float64)
- [x] Embedded structs (promoted fields)

### Synchronization & Mutex
- [x] Race conditions (what happens without synchronization)
- [x] `sync.Mutex` (Lock/Unlock to protect shared state)
- [x] `sync.RWMutex` (multiple readers, single writer)
- [x] `sync/atomic` (preferred for simple counters)
- [x] Lock hygiene (don't sleep or do I/O inside locks)
- [x] Real-world examples: bank account, ticket booking, safe cache
- [x] Best practices (validate before lock, keep critical sections short, use `-race` detector)

### Exercises Completed
- [x] safe-divider — custom sentinel errors (ErrNegative, ErrZero)
- [x] converter — temperature conversion functions
- [x] shape-area — interfaces with Shape polymorphism
- [x] config-loader — JSON parsing with encoding/json
- [x] tasks — slice operations (AddTask, PrintOptions)
- [x] website-checker — goroutines and WaitGroup for concurrent HTTP checks
- [x] log-processor — HTTP handler with channel-based logging pipeline
- [x] sentinel — internal package structure (dashboard, monitor, storage, types)
- [x] logstream — pipeline pattern (Producer, Transformer, Consumer)



## Resources

- [A Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Blog - Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- [Go Blog - Errors are Values](https://go.dev/blog/errors-are-values)
