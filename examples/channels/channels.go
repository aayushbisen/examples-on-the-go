// Package main for the worker pool example using channels.
// This example demonstrates the classic Go concurrency pattern:
// sending work through channels to goroutines that process it concurrently.
package main

// Import multiple packages using a grouped import block.
// This is the idiomatic way to import multiple packages in Go.
import (
	"fmt"   // fmt provides formatted I/O functions (printing, formatting)
	"time"  // time provides functions for measuring and displaying time
)

// worker is a function that processes jobs from a channel and sends results to another.
//
// Parameters:
//   - id: identifies which worker this is (useful for logging/debugging)
//   - jobs <-chan int: a RECEIVE-ONLY channel. The arrow points toward chan, meaning
//     this function can ONLY read from this channel. It cannot send to it.
//     This is a safety feature — the compiler will reject any attempt to send on jobs.
//   - results chan<- int: a SEND-ONLY channel. The arrow points away from chan,
//     meaning this function can ONLY write to this channel. It cannot read from it.
//
// Why directional channels? They prevent bugs. If worker accidentally tries to
// close the jobs channel or read from results, the compiler catches it at build time.
func worker(id int, jobs <-chan int, results chan<- int) {
	// range jobs reads from the channel one value at a time.
	// It BLOCKS (pauses) until a value is available or the channel is closed.
	// When the channel is closed, the loop exits automatically.
	for job := range jobs {
		// fmt.Printf formats and prints text.
		// %d is a placeholder for a decimal integer (like Python's f"{id}").
		// \n is a newline character.
		fmt.Printf("Worker %d started job %d\n", id, job)

		// time.Sleep pauses this goroutine (NOT the entire program).
		// Other goroutines keep running during this sleep.
		// This simulates real work like making an API call or processing data.
		time.Sleep(time.Second)

		fmt.Printf("Worker %d finished job %d\n", id, job)

		// Send the result (job multiplied by 2) to the results channel.
		// The <- operator sends a value INTO a channel.
		// This blocks if the channel buffer is full until a receiver is ready.
		results <- job * 2
	}
}

// main is the entry point. It sets up the channels, launches workers, and coordinates everything.
func main() {
	// make(chan Type, capacity) creates a channel.
	//
	// The second argument (5) creates a BUFFERED channel.
	// A buffered channel can hold up to 5 values before the sender blocks.
	// Think of it like a pipe with a small holding tank — you can push 5 items
	// in without anyone pulling them out yet.
	//
	// Without the capacity (just make(chan int)), it would be UNBUFFERED.
	// Unbuffered channels block the sender until a receiver is ready (and vice versa).
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// := is the short variable declaration operator.
	// Go infers the type automatically. This is equivalent to:
	// var jobs chan int = make(chan int, 5)

	// Start 3 worker goroutines.
	// The `go` keyword before a function call launches it as a goroutine.
	// Goroutines are lightweight threads managed by Go's runtime (not the OS).
	// You can have thousands of goroutines — they only use ~2KB of stack each.
	//
	// All 3 workers run CONCURRENCY (at the same time), sharing the jobs channel.
	// They compete for jobs — whoever is ready grabs the next one.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 5 jobs into the jobs channel.
	// The <- operator sends a value INTO a channel.
	// Since the channel has a buffer of 5, all 5 sends succeed immediately without blocking.
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	// close(jobs) tells all workers: "No more jobs are coming."
	// Workers that are using `for job := range jobs` will finish their current job
	// and then exit the loop when they see the channel is closed.
	//
	// Important: Only the SENDER should close a channel. Closing from a receiver
	// or closing an already-closed channel will cause a panic (runtime crash).
	close(jobs)

	// Collect all 5 results from the results channel.
	// The <- operator receives a value FROM a channel.
	// We use a blank identifier _ or just discard the value with <-results
	// since we only care that results were produced, not what they were.
	for a := 1; a <= 5; a++ {
		<-results
	}

	// If we didn't collect all 5 results, the program might exit before workers finish.
	// Go does NOT wait for goroutines to finish — when main() returns, the program exits
	// and ALL goroutines are killed immediately.

	fmt.Println("All jobs done!")
}
