// Package main demonstrates goroutines — Go's lightweight concurrency mechanism.
// This file shows 5 common goroutine patterns with detailed explanations.
//
// What is a goroutine?
// A goroutine is a function running concurrently (at the same time) with other functions.
// They are managed by Go's runtime, NOT the operating system.
// Key differences from OS threads:
//   - Goroutines start with only ~2KB of stack (threads use ~1MB)
//   - Go's runtime schedules them efficiently across CPU cores
//   - You can run hundreds of thousands of goroutines on a single machine
//   - Use `go` before a function call to launch one
//
// Run this with: go run examples/goroutines/goroutines.go
package main

import (
	"fmt"  // Formatted I/O — printing and formatting text
	"sync" // sync provides synchronization primitives like WaitGroup
	"time" // time provides timers and duration utilities
)

// ============================================================================
// Example 1: Basic Goroutine
// ============================================================================
// The simplest use of goroutines. The `go` keyword launches an anonymous
// function (also called a "closure" or "lambda") to run in the background.
//
// Key concept: Go does NOT wait for goroutines to finish.
// If main() exits, the program ends and ALL goroutines are killed instantly.
// That's why we use time.Sleep here — to give the goroutine time to print.
func basicGoroutine() {
	// go func() { ... }() creates and launches an anonymous goroutine.
	// The () at the end calls (invokes) the function immediately with `go`.
	go func() {
		fmt.Println("  I am a goroutine running in the background!")
	}()

	// Without this sleep, main() would exit immediately and kill the goroutine
	// before it has a chance to print. We'll learn proper ways to wait for
	// goroutines in the next examples (WaitGroup and channels).
	time.Sleep(100 * time.Millisecond)
}

// ============================================================================
// Example 2: Multiple Goroutines with WaitGroup
// ============================================================================
// sync.WaitGroup is Go's way of saying "wait for these goroutines to finish."
// This is the PROPER way to wait for goroutines instead of using time.Sleep.
//
// WaitGroup has 3 methods:
//   - Add(n): say how many goroutines to wait for
//   - Done(): signal that one goroutine finished (called inside the goroutine)
//   - Wait(): block until all goroutines have called Done()
func multipleGoroutines() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		// wg.Add(1) says "I'm about to launch 1 more goroutine to wait for."
		// Must be called BEFORE launching the goroutine.
		wg.Add(1)

		// go func(id int) launches a goroutine with id as a parameter.
		// NOTE: We pass `i` as (i) at the end. This is CRITICAL in Go loops!
		// Without it, all goroutines might see the final value of i (a common Go gotcha).
		// The variable `i` is captured by reference in the closure, so we pass it
		// as a parameter to capture its current value.
		go func(id int) {
			// defer schedules a function call to run when the surrounding function returns.
			// wg.Done() signals "this goroutine is finished."
			// defer is perfect here because it runs even if the function panics.
			defer wg.Done()

			fmt.Printf("  Goroutine %d is running\n", id)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("  Goroutine %d finished\n", id)
		}(i) // ← Important: pass `i` as an argument to capture its value
	}

	// wg.Wait() blocks (pauses) this function until all goroutines call wg.Done().
	// The count must reach zero (we called Add(3) and 3 goroutines each call Done()).
	wg.Wait()
}

// ============================================================================
// Example 3: Worker Pool (Goroutines + Channels)
// ============================================================================
// This combines goroutines with channels to create a worker pool pattern.
// Instead of sharing data with locks, goroutines communicate by sending
// messages through channels. This is Go's philosophy:
// "Don't communicate by sharing memory; share memory by communicating."
func workerPool() {
	// Buffered channels that can hold 5 integers each.
	// jobs: where we send work items
	// results: where workers send back their results
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// WaitGroup to track when all workers are done processing.
	var wg sync.WaitGroup

	// Launch 2 worker goroutines.
	// They'll compete to pick up jobs from the jobs channel.
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// range jobs blocks until a job arrives or the channel closes.
			// This is the core of the worker pattern — keep processing jobs forever.
			for job := range jobs {
				fmt.Printf("  Worker %d processing job %d\n", id, job)
				time.Sleep(300 * time.Millisecond) // Simulate work
				results <- job * 2                 // Send result back
			}
		}(w)
	}

	// Send 4 jobs to be processed.
	for j := 1; j <= 4; j++ {
		jobs <- j
	}
	close(jobs) // Tell workers: "No more jobs incoming."

	// We need to close the results channel so `range results` below can exit.
	// But we can't close it until ALL workers are done sending.
	// So we launch a goroutine that waits for workers, then closes results.
	go func() {
		wg.Wait()    // Wait for all workers to finish
		close(results) // Now it's safe to close
	}()

	// Read all results. range results blocks until a result arrives.
	// When results is closed, this loop exits automatically.
	for r := range results {
		fmt.Printf("  Result: %d\n", r)
	}
}

// ============================================================================
// Example 4: Select with Multiple Channels
// ============================================================================
// `select` is like a `switch` statement but for channels.
// It waits until ONE of its cases is ready, then executes that case.
// If multiple cases are ready, it picks one RANDOMLY (fair scheduling).
//
// Use cases:
//   - Waiting on multiple channels at once
//   - Implementing timeouts
//   - Non-blocking channel operations (with `default`)
func selectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// This goroutine takes 100ms before sending.
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Message from channel 1"
	}()

	// This goroutine takes only 50ms — it will be ready first!
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "Message from channel 2"
	}()

	// select waits for EITHER ch1 or ch2 to have a message.
	// It blocks until at least one case can proceed.
	// ch2's goroutine finishes first (50ms < 100ms), so that case runs first.
	// Then the second iteration picks up ch1's message.
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("  Received:", msg)
		case msg := <-ch2:
			fmt.Println("  Received:", msg)
		}
	}
}

// ============================================================================
// Example 5: Avoiding Goroutine Leaks
// ============================================================================
// A goroutine leak happens when a goroutine is blocked forever and can never exit.
// This wastes memory and can crash your program if it happens repeatedly.
//
// Common causes:
//   - Waiting to receive from a channel that nobody sends to
//   - Waiting to send to a channel that nobody receives from
//   - Deadlock (two goroutines waiting on each other)
//
// Solution: Always use timeouts or context cancellation!
func goroutineLeak() {
	ch := make(chan int)

	go func() {
		// select with time.After is the standard Go pattern to prevent leaks.
		// It says: "Wait for a value from ch, BUT if 500ms passes, give up."
		select {
		case val := <-ch:
			fmt.Println("  Got:", val)
		case <-time.After(500 * time.Millisecond):
			// time.After returns a channel that sends the current time after the duration.
			// We don't care about the time value, so we use <- without assigning.
			fmt.Println("  Timeout: goroutine would have leaked, but time.After saved us!")
		}
		// After the select completes (either case), this goroutine exits. No leak!
	}()

	// We intentionally NEVER send to `ch`. Without the timeout, the goroutine
	// would be stuck on `case val := <-ch` forever — a classic goroutine leak.
	// time.After acts as an escape hatch.
	time.Sleep(600 * time.Millisecond) // Give the goroutine time to hit the timeout
}

// main runs all 5 examples in sequence.
func main() {
	fmt.Println("=== 1. Basic Goroutine ===")
	basicGoroutine()

	fmt.Println("\n=== 2. Multiple Goroutines ===")
	multipleGoroutines()

	fmt.Println("\n=== 3. Worker Pool (Goroutines + Channels) ===")
	workerPool()

	fmt.Println("\n=== 4. Select with Multiple Channels ===")
	selectExample()

	fmt.Println("\n=== 5. Avoiding Goroutine Leaks ===")
	goroutineLeak()
}
