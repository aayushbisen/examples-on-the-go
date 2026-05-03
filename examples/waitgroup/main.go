package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Go WaitGroups Tutorial
=====================

## Concurrency in Go
Go supports concurrency using goroutines: lightweight, independently scheduled functions
created with the `go` keyword (e.g., `go myFunc()`). The main function runs in a goroutine,
and when main exits, all other running goroutines are terminated immediately.

## Why Synchronization is Needed
Goroutines run concurrently with the main goroutine. Without synchronization, main may exit
before spawned goroutines complete, leading to missing output or incomplete work.

## The Problem WaitGroups Solve
`sync.WaitGroup` coordinates multiple goroutines, allowing a parent goroutine to wait until
all spawned goroutines finish. This prevents premature exit of the main function.

## WaitGroup Methods
A WaitGroup maintains an internal counter tracking active goroutines:
- `Add(delta int)`: Increments the counter by `delta`. Must be called *before* starting tracked goroutines.
- `Done()`: Decrements the counter by 1 (equivalent to `Add(-1)`). Called when a goroutine finishes.
- `Wait()`: Blocks the caller until the counter reaches 0 (all goroutines have called `Done()`).

## Common Mistakes
1. Forgetting `Done()`: Counter never reaches 0, `Wait()` blocks forever (deadlock).
2. Calling `Add()` inside a goroutine: `Wait()` may execute before `Add()` is called, causing premature exit.
3. Incorrect `Add()` delta: Overcounting/undercounting leads to deadlocks or missed goroutines.
4. Copying WaitGroups: WaitGroups contain internal state and must be passed by pointer (`*sync.WaitGroup`).
5. Loop variable capture: Goroutines in loops may see updated loop values if variables are not captured correctly.

## Idiom: `defer wg.Done()`
Always use `defer wg.Done()` at the start of goroutines using a WaitGroup. This ensures `Done()` is
called even if the goroutine panics or returns early, preventing accidental deadlocks.
*/

// Demonstrates broken behavior without WaitGroup (commented to avoid runtime issues)
/*
func brokenNoWaitGroup() {
	fmt.Println("Broken Example: No WaitGroup")
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d starting\n", id)
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("Goroutine %d done\n", id)
		}(i)
	}
	// No Wait() here: main exits immediately, goroutines are killed
	fmt.Println("Broken Example: Main exiting (goroutines may not finish)")
}
*/

// Example 1: Basic WaitGroup with multiple goroutines
func basicWaitGroupExample() {
	fmt.Println("\n=== Example 1: Basic WaitGroup ===")

	var wg sync.WaitGroup
	numGoroutines := 3

	// Add count before starting goroutines
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		// Pass loop variable as argument to avoid capturing by reference
		go func(id int) {
			defer wg.Done() // Ensure counter is decremented on exit
			fmt.Printf("Basic Example: Goroutine %d starting\n", id)
			time.Sleep(300 * time.Millisecond) // Simulate work
			fmt.Printf("Basic Example: Goroutine %d done\n", id)
		}(i)
	}

	fmt.Println("Basic Example: Waiting for all goroutines...")
	wg.Wait() // Block until counter reaches 0
	fmt.Println("Basic Example: All goroutines finished!")
}

// Example 2: Realistic task processing simulation
func taskProcessingExample() {
	fmt.Println("\n=== Example 2: Task Processing (Realistic) ===")

	tasks := []string{"task-1", "task-2", "task-3", "task-4"}
	var wg sync.WaitGroup

	// Add count equal to number of tasks
	wg.Add(len(tasks))

	for _, task := range tasks {
		// Capture task variable to avoid loop variable race condition
		currentTask := task
		go func() {
			defer wg.Done()
			fmt.Printf("Processing: %s\n", currentTask)
			// Simulate variable work time based on task name length
			workTime := time.Duration(200+len(currentTask)*100) * time.Millisecond
			time.Sleep(workTime)
			fmt.Printf("Completed: %s\n", currentTask)
		}()
	}

	fmt.Println("Task Processing: Waiting for all tasks...")
	wg.Wait()
	fmt.Println("Task Processing: All tasks completed!")
}

/*
## WaitGroups vs Channels
Both tools handle synchronization but serve different use cases:
- Use WaitGroups when you only need to wait for a group of goroutines to finish (no data passing).
- Use channels when you need to send data between goroutines, signal completion via done channels,
  or select between multiple communication events.
- WaitGroups are simpler for pure coordination; channels are more flexible for data flow.
*/

/*
## Summary
- `sync.WaitGroup` coordinates goroutine completion.
- Call `Add()` with the correct count *before* starting goroutines.
- Use `defer wg.Done()` in goroutines to ensure counter decrement.
- Call `Wait()` to block until all goroutines finish.
- Avoid common pitfalls: missing `Done()`, incorrect `Add()` calls, loop variable capture issues.
*/

func main() {
	// Uncomment to see broken behavior (output will likely miss goroutine logs):
	// brokenNoWaitGroup()

	basicWaitGroupExample()
	taskProcessingExample()
}
