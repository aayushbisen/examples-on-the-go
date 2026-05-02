// Package main demonstrates the context package — one of Go's most important
// tools for managing goroutine lifecycles, timeouts, and request-scoped data.
//
// ============================================================================
// WHAT IS CONTEXT?
// ============================================================================
// context.Context is a Go type that carries:
//   1. DEADLINES — "Stop working by this time"
//   2. CANCELLATION — "Stop working now"
//   3. REQUEST VALUES — "Here's some data for this request"
//
// Why does it exist?
// When you launch a goroutine, it runs independently. If the user cancels
// their request, the goroutine keeps running, wasting resources. Context
// gives you a way to tell goroutines: "Hey, stop what you're doing."
//
// It flows through your program like a river — every function in a request
// chain receives the same context and respects its signals.
//
// Run this with: go run examples/context/context.go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// Example 1: Context Cancellation (The Most Common Use Case)
// ============================================================================
// This is the #1 use of context: canceling long-running work.
//
// Scenario: You start a background task, but the user navigates away.
// Without context, the task runs to completion wasting CPU.
// With context, you call cancel() and the task stops immediately.
//
// Key functions:
//   context.Background()          — the root context, never canceled
//   context.WithCancel(parent)    — returns a child context + cancel function
//   ctx.Done()                    — returns a channel that closes on cancel
//   ctx.Err()                     — returns why the context was canceled
func contextCancellation() {
	// context.Background() is the root of all context trees.
	// It is NEVER canceled, has NO deadline, and carries NO values.
	// Use it as the starting point when no other context exists.
	ctx := context.Background()

	// context.WithCancel creates a child context that CAN be canceled.
	// It returns:
	//   - ctx: the new context (carries the cancel signal)
	//   - cancel: a function that, when called, cancels the context
	ctx, cancel := context.WithCancel(ctx)

	// IMPORTANT: Always call cancel() when done, usually with defer.
	// If you don't, you leak the goroutine and its memory.
	// Think of it like closing a file — defer ensures it happens.
	defer cancel()

	var wg sync.WaitGroup

	// Launch a goroutine that does long-running work.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// for-select is the standard pattern for listening to context.Done().
		// It checks: "Has anyone told me to stop?"
		for {
			select {
			case <-ctx.Done():
				// ctx.Done() returns a channel. When cancel() is called,
				// this channel CLOSES. Receiving from a closed channel
				// returns immediately (the zero value).
				// This is how the goroutine learns it should stop.
				fmt.Printf("  Goroutine stopped! Reason: %v\n", ctx.Err())
				return // Exit the goroutine cleanly
			default:
				// Not canceled yet — do some work.
				fmt.Println("  Working...")
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Let the goroutine work for 1 second, then cancel it.
	time.Sleep(1 * time.Second)
	fmt.Println("  [Main] Canceling the goroutine...")
	cancel() // This closes ctx.Done(), signaling the goroutine to stop

	wg.Wait() // Wait for goroutine to finish
}

// ============================================================================
// Example 2: Context Timeout (Auto-Cancel After a Duration)
// ============================================================================
// context.WithTimeout automatically cancels after a set duration.
//
// Scenario: An API call should fail after 2 seconds instead of hanging forever.
// This is used everywhere in production: HTTP clients, DB queries, etc.
func contextTimeout() {
	ctx := context.Background()

	// context.WithTimeout creates a context that auto-cancels after 2 seconds.
	// It's equivalent to WithCancel + a timer that calls cancel() automatically.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel() // Still call it! The timer calls it too, but this cleans up early if we finish before the timeout.

	// Simulate a slow operation (takes 5 seconds).
	resultChan := make(chan string)
	go func() {
		fmt.Println("  Starting slow operation (takes 5s)...")
		time.Sleep(5 * time.Second)
		resultChan <- "Done!"
	}()

	// Wait for either the result OR the timeout.
	select {
	case result := <-resultChan:
		fmt.Println("  Result:", result)
	case <-ctx.Done():
		// After 2 seconds, ctx.Done() fires automatically.
		// ctx.Err() tells us WHY: "context deadline exceeded"
		fmt.Printf("  Operation timed out! Reason: %v\n", ctx.Err())
	}
}

// ============================================================================
// Example 3: Context Deadline (Auto-Cancel at a Specific Time)
// ============================================================================
// context.WithDeadline is like WithTimeout but uses an absolute time
// instead of a duration.
//
// Scenario: "This task must complete before 3:00 PM" (not "in 2 hours").
func contextDeadline() {
	ctx := context.Background()

	// Set a deadline: 2 seconds from now.
	// WithTimeout(duration) is just a convenience wrapper around WithDeadline(time.Now().Add(duration)).
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	// Simulate work that takes longer than the deadline.
	resultChan := make(chan string)
	go func() {
		fmt.Println("  Working... (deadline is in 2s)")
		time.Sleep(3 * time.Second)
		resultChan <- "Finished"
	}()

	select {
	case result := <-resultChan:
		fmt.Println("  Result:", result)
	case <-ctx.Done():
		fmt.Printf("  Deadline exceeded! Reason: %v\n", ctx.Err())
	}
}

// ============================================================================
// Example 4: Context Values (Passing Request-Scoped Data)
// ============================================================================
// context.WithValue attaches key-value pairs to a context.
//
// Scenario: Passing a request ID, auth token, or user info through
// a chain of function calls without adding parameters to every function.
//
// IMPORTANT RULES:
//   - Only use context for request-scoped data (not for passing optional params)
//   - Use custom types as keys (not strings) to avoid key collisions
//   - Values are immutable — you can't change them, only create a new context

// ctxKey is a custom type for context keys.
// This prevents key collisions between different packages.
// If two packages both use "userID" as a string key, they'd overwrite each other.
// Using a custom type (instead of a raw string) makes each package's keys unique.
type ctxKey string

const userIDKey ctxKey = "userID"
const requestIDKey ctxKey = "requestID"

func contextValues() {
	ctx := context.Background()

	// context.WithValue returns a NEW context with the key-value pair added.
	// The original ctx is unchanged (contexts are immutable).
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	ctx = context.WithValue(ctx, requestIDKey, "req-abc")

	// Simulate a chain of function calls.
	processRequest(ctx)
}

// processRequest receives the context and passes it to deeper functions.
func processRequest(ctx context.Context) {
	// Extract values from context.
	// ctx.Value(key) returns interface{}, so we type-assert to get the real type.
	// If the key doesn't exist, ctx.Value returns nil.
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		fmt.Println("  User ID from context:", userID)
	}
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Println("  Request ID from context:", requestID)
	}

	// Pass context to the next function in the chain.
	validateToken(ctx)
}

func validateToken(ctx context.Context) {
	// The context flows through the call chain. Every function can access the same values.
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Println("  Token validated for request:", requestID)
	}
}

// ============================================================================
// Example 5: Context Propagation (Passing Through Function Chains)
// ============================================================================
// This shows how context flows through multiple layers of function calls
// and goroutines, allowing cancellation at any level to propagate everywhere.
//
// This is the MOST IMPORTANT pattern in real Go applications.
// Every HTTP handler, DB call, and API request should accept a context.
func contextPropagation() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("  Starting request chain...")

	// Pass context through the chain.
	// If any layer cancels, ALL downstream operations stop.
	layer1(ctx)
}

func layer1(ctx context.Context) {
	fmt.Println("  Layer 1: Received request, calling layer 2...")
	layer2(ctx)
}

func layer2(ctx context.Context) {
	fmt.Println("  Layer 2: Processing, calling layer 3...")
	layer3(ctx)
}

func layer3(ctx context.Context) {
	fmt.Println("  Layer 3: Starting slow operation...")

	resultChan := make(chan string)
	go func() {
		// Simulate a slow database query (3 seconds).
		time.Sleep(3 * time.Second)
		resultChan <- "Database result"
	}()

	select {
	case result := <-resultChan:
		fmt.Println("  Layer 3:", result)
	case <-ctx.Done():
		// The timeout (1s) fires before the query (3s) completes.
		// This cancellation propagates all the way up.
		fmt.Printf("  Layer 3 canceled! Reason: %v\n", ctx.Err())
	}
}

// ============================================================================
// Example 6: Context with HTTP Server (Request Cancellation)
// ============================================================================
// In an HTTP server, context is automatically created for each request.
// It is canceled when:
//   - The client disconnects
//   - The request times out
//   - The server shuts down
//
// This prevents wasted work when users abandon requests.
func httpServerWithContext() {
	// Define a handler function.
	// r.Context() gives us the request's context.
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Log when the request starts.
		fmt.Println("  [Server] Request received")

		// Simulate slow processing.
		select {
		case <-time.After(5 * time.Second):
			fmt.Fprintf(w, "Done!")
		case <-ctx.Done():
			// The client disconnected or the request timed out.
			// We stop processing immediately.
			fmt.Println("  [Server] Client disconnected, stopping work")
		}
	}

	// Start the server in a goroutine.
	server := &http.Server{Addr: ":8081", Handler: http.HandlerFunc(handler)}
	go func() {
		fmt.Println("  [Server] Listening on :8081")
		// ListenAndServe returns an error when the server stops.
		// http.ErrServerClosed is the expected error from Shutdown().
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("  [Server] Error: %v\n", err)
		}
	}()

	// Give the server time to start.
	time.Sleep(200 * time.Millisecond)

	// Simulate a client request that cancels after 1 second.
	clientCtx, clientCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer clientCancel()

	req, _ := http.NewRequestWithContext(clientCtx, "GET", "http://localhost:8081/", nil)
	client := &http.Client{}

	fmt.Println("  [Client] Sending request (will cancel after 1s)...")
	_, err := client.Do(req)
	if err != nil {
		fmt.Printf("  [Client] Request canceled: %v\n", err)
	}

	// Give server time to process the cancellation.
	time.Sleep(500 * time.Millisecond)

	// Shut down the server gracefully.
	serverCtx, serverCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer serverCancel()
	server.Shutdown(serverCtx)
}

// ============================================================================
// Example 7: Context with Multiple Goroutines (Graceful Shutdown)
// ============================================================================
// This is the production pattern for shutting down services cleanly.
// One cancel() call stops ALL goroutines at once.
func gracefulShutdown() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// Launch 5 worker goroutines.
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			workerWithCtx(ctx, id)
		}(i)
	}

	// Let workers run for 1 second.
	fmt.Println("  Workers running... (shutting down in 1s)")
	time.Sleep(1 * time.Second)

	// ONE cancel() call stops ALL 5 workers.
	fmt.Println("  [Shutdown] Canceling all workers...")
	cancel()

	wg.Wait()
	fmt.Println("  [Shutdown] All workers stopped cleanly")
}

// workerWithCtx does work until the context is canceled.
func workerWithCtx(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  Worker %d shutting down: %v\n", id, ctx.Err())
			return
		default:
			// Simulate work.
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// ============================================================================
// Example 8: Context in Database/External Service Calls
// ============================================================================
// Real database drivers and HTTP clients accept context.
// This lets you set timeouts per-query and cancel slow operations.
func contextWithExternalCalls() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Simulate a database query.
	fmt.Println("  Querying database (will timeout)...")
	result, err := fakeDBQuery(ctx, "SELECT * FROM users")
	if err != nil {
		fmt.Printf("  DB Error: %v\n", err)
	} else {
		fmt.Println("  DB Result:", result)
	}

	fmt.Println()

	// Same query with more time.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	fmt.Println("  Querying database (enough time)...")
	result2, err2 := fakeDBQuery(ctx2, "SELECT * FROM users")
	if err2 != nil {
		fmt.Printf("  DB Error: %v\n", err2)
	} else {
		fmt.Println("  DB Result:", result2)
	}
}

// fakeDBQuery simulates a database call that takes 1 second.
// In real code, this would be db.QueryContext(ctx, "SELECT ...").
func fakeDBQuery(ctx context.Context, query string) (string, error) {
	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		// Simulate a 1-second database query.
		time.Sleep(1 * time.Second)

		// 30% chance of a random error (to simulate real DB behavior).
		if rand.Intn(10) < 3 {
			errChan <- fmt.Errorf("connection lost")
			return
		}
		resultChan <- "user data from DB"
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		// The query was canceled (timeout or manual cancel).
		// In real code, the DB driver handles this and stops the query.
		return "", fmt.Errorf("query canceled: %w", ctx.Err())
	}
}

// ============================================================================
// Example 9: Context TODO vs Background
// ============================================================================
// context.TODO() and context.Background() are identical in behavior.
// The difference is SEMANTIC — they signal intent to other developers.
func contextTODO() {
	// context.Background() — Use this when you KNOW this is the root context.
	// Example: In main(), in tests, or when starting a new request.
	rootCtx := context.Background()

	// context.TODO() — Use this when you're NOT SURE which context to use,
	// or you're refactoring and need a placeholder. It says "I know I should
	// use a proper context here, but I haven't figured it out yet."
	//
	// It is NOT a different type. It's the same empty context as Background().
	todoCtx := context.TODO()

	fmt.Println("  Background() and TODO() are both empty contexts:")
	fmt.Printf("  Background: canceled=%v, deadline=%v, err=%v\n",
		rootCtx.Done() != nil, false, rootCtx.Err())
	fmt.Printf("  TODO:       canceled=%v, deadline=%v, err=%v\n",
		todoCtx.Done() != nil, false, todoCtx.Err())
}

// ============================================================================
// Example 10: Context Chain (Cancellation Propagation)
// ============================================================================
// When you create a child context from a parent, canceling the parent
// automatically cancels ALL children. This is how cancellation cascades.
func contextChain() {
	// Root context with a 3-second timeout.
	rootCtx, rootCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer rootCancel()

	// Child context with a 1-second timeout (stricter than parent).
	childCtx, childCancel := context.WithTimeout(rootCtx, 1*time.Second)
	defer childCancel()

	// Grandchild context with cancellation.
	grandchildCtx, grandchildCancel := context.WithCancel(childCtx)
	defer grandchildCancel()

	// Launch a goroutine using the grandchild context.
	go func() {
		<-grandchildCtx.Done()
		fmt.Printf("  Grandchild stopped! Reason: %v\n", grandchildCtx.Err())
	}()

	// Wait 500ms, then cancel the CHILD.
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  [Main] Canceling child context...")
	childCancel()

	// The child's cancellation automatically cancels the grandchild too!
	// You don't need to cancel each one individually.
	time.Sleep(200 * time.Millisecond)
}

// main runs all context examples.
func main() {
	fmt.Println("=== 1. Context Cancellation ===")
	contextCancellation()

	fmt.Println("\n=== 2. Context Timeout ===")
	contextTimeout()

	fmt.Println("\n=== 3. Context Deadline ===")
	contextDeadline()

	fmt.Println("\n=== 4. Context Values ===")
	contextValues()

	fmt.Println("\n=== 5. Context Propagation ===")
	contextPropagation()

	fmt.Println("\n=== 6. HTTP Server with Context ===")
	httpServerWithContext()

	fmt.Println("\n=== 7. Graceful Shutdown ===")
	gracefulShutdown()

	fmt.Println("\n=== 8. Context with External Calls (DB/API) ===")
	contextWithExternalCalls()

	fmt.Println("\n=== 9. TODO vs Background ===")
	contextTODO()

	fmt.Println("\n=== 10. Context Chain (Cancellation Propagation) ===")
	contextChain()
}
