// Package main demonstrates channels — Go's mechanism for goroutine communication.
// This file shows 9 channel patterns with detailed explanations for beginners.
//
// What is a channel?
// A channel is a typed pipe that connects concurrent goroutines.
// You send values into one end and receive them from the other.
// Channels are Go's answer to the problem: "How do goroutines talk to each other safely?"
//
// Key concepts:
//   - make(chan Type)        — creates an unbuffered channel (blocks until receiver is ready)
//   - make(chan Type, size)  — creates a buffered channel (holds `size` values before blocking)
//   - ch <- value            — sends a value into the channel
//   - value := <-ch          — receives a value from the channel
//   - close(ch)              — closes the channel (no more sends allowed)
//   - range ch               — reads values until the channel is closed
//
// Run this with: go run examples/channels-all/channels.go
package main

import (
	"fmt"  // Formatted I/O — printing and formatting text
	"time" // time provides timers, sleep, and duration utilities
)

// ============================================================================
// Example 1: Unbuffered Channel
// ============================================================================
// An unbuffered channel has NO capacity. It blocks the sender until a receiver
// is ready, and blocks the receiver until a sender is ready.
//
// This guarantees synchronization — the sender knows the receiver got the value.
// Think of it like a phone call: both parties must be on the line at the same time.
//
// Visual:
//   Sender  ----|  |----  Receiver
//       (blocks until both are ready)
func unbufferedChannel() {
	// make(chan string) creates an unbuffered channel.
	// No buffer means: a send blocks until someone receives, and vice versa.
	ch := make(chan string)

	// Launch a goroutine that sends a message.
	// This send WILL BLOCK because the channel has no buffer and nobody is
	// receiving yet. It will unblock when main() receives below.
	go func() {
		ch <- "Hello from goroutine!" // Blocks until main() does <-ch
	}()

	// Receive from the channel. This blocks until the goroutine sends.
	// After this line, both the sender (goroutine) and receiver (main) proceed.
	msg := <-ch
	fmt.Println("  Received:", msg)
}

// ============================================================================
// Example 2: Buffered Channel
// ============================================================================
// A buffered channel has a capacity. Sends only block when the buffer is FULL.
// Receives only block when the buffer is EMPTY.
//
// Think of it like a mailbox with limited slots — you can drop letters in
// without the recipient being home, but once it's full, you have to wait.
//
// Visual:
//   Sender  ----[1][2][3]----  Receiver
//       (can send up to 3 without blocking)
func bufferedChannel() {
	// make(chan int, 3) creates a buffered channel that can hold 3 integers.
	// Sends to this channel won't block until 3 values are in the buffer.
	ch := make(chan int, 3)

	// These 3 sends succeed immediately — they go into the buffer.
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4 // ← Uncommenting this would cause a DEADLOCK!
	// The buffer is full (3/3), so this send blocks forever because nobody
	// is receiving yet. In a single-goroutine program, this is a deadlock
	// and Go's runtime will crash with "all goroutines are asleep".

	fmt.Println("  Buffer has 3 items")

	// Receive values. Each receive removes one item from the buffer.
	fmt.Println("  Received:", <-ch) // Gets 1, buffer now has [2][3]
	fmt.Println("  Received:", <-ch) // Gets 2, buffer now has [3]
	fmt.Println("  Received:", <-ch) // Gets 3, buffer is now empty
}

// ============================================================================
// Example 3: Closing Channels and Range
// ============================================================================
// close(ch) signals that no more values will be sent.
// range ch reads values one by one until the channel is closed.
//
// Important rules:
//   - Only the SENDER should close a channel
//   - Closing an already-closed channel causes a PANIC (runtime crash)
//   - Sending on a closed channel causes a PANIC
//   - Receiving from a closed channel returns the zero value (0, "", false, nil)
func closingAndRange() {
	ch := make(chan int, 5)

	// Launch a goroutine that sends 5 numbers then closes the channel.
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i // Send each number
		}
		// close(ch) is essential — without it, `range ch` below would block
		// forever after receiving all 5 values, waiting for more that never come.
		close(ch)
	}()

	// range ch automatically:
	//   1. Receives a value from ch and assigns it to `val`
	//   2. Executes the loop body
	//   3. Repeats until ch is closed, then exits
	fmt.Print("  Values: ")
	for val := range ch {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}

// ============================================================================
// Example 4: Directional Channels
// ============================================================================
// Go lets you restrict a channel to only send OR only receive in function signatures.
// This is a compile-time safety feature — it prevents accidental misuse.
//
//   chan<- int  means "you can only SEND to this channel" (arrow points away from chan)
//   <-chan int  means "you can only RECEIVE from this channel" (arrow points toward chan)
//
// When you create a channel with make(chan int), it's bidirectional (both send and receive).
// Go automatically converts bidirectional channels to directional ones when passed to functions.
func producer(ch chan<- int) {
	// ch is send-only here. Trying to do `val := <-ch` would cause a compile error.
	for i := 1; i <= 3; i++ {
		ch <- i
		fmt.Println("  Produced:", i)
	}
	close(ch) // The producer closes the channel — it owns the sending side.
}

func consumer(ch <-chan int) {
	// ch is receive-only here. Trying to do `ch <- 42` would cause a compile error.
	// This prevents the consumer from accidentally interfering with the channel.
	for val := range ch {
		fmt.Println("  Consumed:", val)
	}
}

func directionalChannels() {
	ch := make(chan int, 3) // Bidirectional channel
	go producer(ch)         // Go converts it to chan<- int automatically
	consumer(ch)            // Go converts it to <-chan int automatically
}

// ============================================================================
// Example 5: Select with Timeout
// ============================================================================
// `select` waits on multiple channel operations and executes the first one that's ready.
// `time.After` is the standard way to add a timeout to any channel operation.
//
// Why use timeouts?
// Without a timeout, if the other end of a channel never responds, your goroutine
// waits forever — this is a goroutine leak. Timeouts prevent this.
func selectWithTimeout() {
	ch := make(chan string)

	// This goroutine takes 2 seconds to send — too slow for our patience.
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "This is too slow!"
	}()

	// select waits for ONE of these cases to be ready:
	//   Case 1: Receive from ch (will be ready in 2 seconds)
	//   Case 2: Receive from time.After (will be ready in 500ms)
	//
	// Since 500ms < 2s, the timeout case wins and executes.
	select {
	case msg := <-ch:
		fmt.Println("  Got:", msg)
	case <-time.After(500 * time.Millisecond):
		// time.After returns a channel that sends the current time after the duration.
		// We use <- to receive from it (and discard the time value).
		// This case runs after 500ms of waiting with no message on ch.
		fmt.Println("  Timeout: operation took too long")
	}
}

// ============================================================================
// Example 6: Non-Blocking Send/Receive with Default
// ============================================================================
// Adding a `default` case to `select` makes it non-blocking.
// If no channel operation is ready, `default` executes immediately instead of waiting.
//
// This is useful when you want to check a channel without blocking.
func nonBlockingChannel() {
	ch := make(chan int, 1) // Buffer of 1

	// Try to receive from an empty channel.
	// Without `default`, this would block forever (nobody is sending).
	// With `default`, it immediately executes the default case.
	select {
	case val := <-ch:
		fmt.Println("  Received:", val)
	default:
		// This runs immediately because ch is empty — no blocking!
		fmt.Println("  No value ready to receive (non-blocking)")
	}

	// Try to send to the channel. It has buffer space (0/1 used), so it succeeds.
	select {
	case ch <- 42:
		fmt.Println("  Sent 42 to channel")
	default:
		// This would only run if the channel buffer was full.
		fmt.Println("  Channel full, cannot send")
	}
}

// ============================================================================
// Example 7: Fan-In Pattern
// ============================================================================
// Fan-in: Multiple goroutines send to a SINGLE channel.
// One receiver collects all the messages.
//
// Think of it like many people talking into one megaphone —
// the messages get combined into a single stream.
//
// Use case: Merging results from multiple parallel tasks.
func fanIn() {
	ch := make(chan string)

	// say() is a helper that launches a goroutine to send a message.
	say := func(msg string) {
		go func() { ch <- msg }()
	}

	// Launch 3 goroutines that each send one message.
	say("Hello")
	say("from")
	say("Go!")

	// We know exactly 3 messages are coming, so we receive exactly 3 times.
	// We CANNOT use `range ch` here because we don't close the channel
	// (the goroutines sending are internal to this function and finish at different times).
	fmt.Print("  Fan-in: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%s ", <-ch) // Blocks until each message arrives
	}
	fmt.Println()
}

// ============================================================================
// Example 8: Fan-Out Pattern
// ============================================================================
// Fan-out: ONE channel, MULTIPLE receivers.
// Each receiver gets a different value (not copies).
//
// Think of it like dealing cards — each worker gets a subset of the values.
//
// Use case: Distributing work across multiple workers for parallel processing.
func fanOut() {
	ch := make(chan int, 5)       // Channel with work items
	done := make(chan bool, 3)    // Channel to signal when each worker finishes

	// Put 5 values into the channel.
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch) // No more values — tell workers they can stop after draining.

	// Launch 3 workers (receivers). They all read from the SAME channel.
	// Each value is received by exactly ONE worker (not all of them).
	for w := 1; w <= 3; w++ {
		go func(id int) {
			// range ch reads until the channel is closed AND drained.
			for val := range ch {
				fmt.Printf("  Worker %d got: %d\n", id, val)
			}
			// Signal that this worker is done.
			// done is buffered (capacity 3), so these sends won't block.
			done <- true
		}(w)
	}

	// Wait for all 3 workers to finish.
	// Each worker sends `true` to done when it's finished.
	for i := 0; i < 3; i++ {
		<-done
	}
}

// ============================================================================
// Example 9: Pipeline Pattern
// ============================================================================
// A pipeline chains channels together: stage1 → stage2 → stage3 → output
// Each stage reads from one channel, processes the data, and writes to the next.
//
// Think of it like an assembly line:
//   Raw materials → Cut → Weld → Paint → Finished product
//
// Use case: Breaking complex processing into reusable, concurrent stages.
func pipeline() {
	// Stage 1 output / Stage 2 input
	numbers := make(chan int)
	// Stage 2 output / Final consumer input
	squares := make(chan int)

	// Stage 1: Generate numbers 1-5 and send them into `numbers`.
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers) // Stage 1 done — signals Stage 2 to finish after reading all.
	}()

	// Stage 2: Read from `numbers`, square each value, send to `squares`.
	go func() {
		for n := range numbers { // Reads until `numbers` is closed and drained
			squares <- n * n
		}
		close(squares) // Stage 2 done — signals the consumer below to finish.
	}()

	// Consumer: Read final results from `squares`.
	fmt.Print("  Pipeline (squares): ")
	for s := range squares { // Reads until `squares` is closed and drained
		fmt.Printf("%d ", s)
	}
	fmt.Println()
}

// main runs all 9 channel examples in sequence.
func main() {
	fmt.Println("=== 1. Unbuffered Channel ===")
	unbufferedChannel()

	fmt.Println("\n=== 2. Buffered Channel ===")
	bufferedChannel()

	fmt.Println("\n=== 3. Closing Channels and Range ===")
	closingAndRange()

	fmt.Println("\n=== 4. Directional Channels ===")
	directionalChannels()

	fmt.Println("\n=== 5. Select with Timeout ===")
	selectWithTimeout()

	fmt.Println("\n=== 6. Non-Blocking Send/Receive ===")
	nonBlockingChannel()

	fmt.Println("\n=== 7. Fan-In Pattern ===")
	fanIn()

	fmt.Println("\n=== 8. Fan-Out Pattern ===")
	fanOut()

	fmt.Println("\n=== 9. Pipeline Pattern ===")
	pipeline()
}
