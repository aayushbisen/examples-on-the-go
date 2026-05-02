# My Go Learning Journey

A personal repository documenting my journey learning Go (Golang), from the basics to concurrency patterns.

## Structure

```
.
├── main.go                          # Hello World & basic Go syntax
├── examples/
│   ├── channels/                    # Worker pool with channels
│   ├── channels-all/                # 9 channel patterns (unbuffered, buffered, fan-in, fan-out, pipeline, etc.)
│   ├── context/                     # 10 context patterns (cancellation, timeout, values, HTTP, graceful shutdown, etc.)
│   └── goroutines/                  # 5 goroutine patterns (basic, WaitGroup, worker pool, select, leak prevention)
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

# Goroutines - All Patterns (5 examples)
go run examples/goroutines/goroutines.go
```

## What I've Learned So Far

- [x] Go basics (packages, imports, main function)
- [x] Goroutines (lightweight concurrency)
- [x] WaitGroup (waiting for goroutines)
- [x] Channels (unbuffered & buffered)
- [x] Directional channels (send-only & receive-only)
- [x] Select statement (multiple channels, timeouts, non-blocking)
- [x] Fan-in / Fan-out patterns
- [x] Pipeline pattern
- [x] Avoiding goroutine leaks
- [x] Context cancellation (stop long-running work)
- [x] Context timeout & deadline (auto-cancel after duration or at a specific time)
- [x] Context values (pass request-scoped data like auth tokens)
- [x] Context propagation (flow through function chains)
- [x] Context in HTTP servers (client disconnect detection)
- [x] Graceful shutdown (stop all goroutines with one cancel)
- [x] Context with DB/external calls (query timeouts)
- [x] Context chain (parent cancels all children)
- [x] TODO vs Background (semantic difference)

## Resources

- [A Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
