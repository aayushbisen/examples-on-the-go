package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ---------------- DEFER BASICS ----------------
func deferBasics() {
	fmt.Println("--- Basic defer execution ---")
	basicDefer()
	fmt.Println()

	fmt.Println("--- Multiple defers (LIFO order) ---")
	multipleDefers()
	fmt.Println()

	fmt.Println("--- Defer with file/resource cleanup ---")
	fileCleanup()
	fmt.Println()

	fmt.Println("--- Defer with mutex unlock ---")
	mutexUnlock()
	fmt.Println()

	fmt.Println("--- Defer argument evaluation timing ---")
	deferArgEval()
	fmt.Println()

	fmt.Println("--- Defer modifying named return values ---")
	res := namedReturnExample()
	fmt.Println("Final returned value:", res)
	fmt.Println()
}

func basicDefer() {
	// Defer schedules a function call to run just before the surrounding function returns.
	// Key rule: Defer runs after all other code in the function, before the function exits to its caller.
	fmt.Println("Step 1: Before defer")
	defer fmt.Println("Step 3: Deferred print (runs when basicDefer returns)")
	fmt.Println("Step 2: After defer")
}

func multipleDefers() {
	// Defers execute in LIFO (Last In, First Out) order, like a stack.
	// Key rule: Last registered defer runs first when the function exits.
	defer fmt.Println("Defer 3 (last registered, first to run)")
	defer fmt.Println("Defer 2")
	defer fmt.Println("Defer 1 (first registered, last to run)")
	fmt.Println("Main logic done")
}

func fileCleanup() {
	// Defer is idiomatic for resource cleanup (files, connections). Ensures cleanup runs even on early return/panic.
	f, err := os.CreateTemp("", "example")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer func() {
		fmt.Println("Closing temp file:", f.Name())
		f.Close()
		os.Remove(f.Name())
	}()
	fmt.Println("Temp file created:", f.Name())
	_, err = f.WriteString("hello world")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Written to file successfully")
}

func mutexUnlock() {
	// Defer avoids deadlocks by ensuring mutex unlock even on early return/panic.
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Mutex locked, doing work...")
	fmt.Println("Work done, mutex will be unlocked by defer")
}

func deferArgEval() {
	// Defer arguments are evaluated immediately when the defer statement is executed, not when the deferred function runs.
	x := 10
	defer fmt.Println("Defer argument evaluated at defer time, x =", x)
	x = 20
	fmt.Println("After modifying x, x =", x)
}

func namedReturnExample() (ret int) {
	// Defer can modify named return values. Named returns are function-scoped variables, defer runs before return.
	defer func() {
		ret += 10
		fmt.Println("Defer modifying named return, new ret:", ret)
	}()
	ret = 5
	fmt.Println("Initial named return ret:", ret)
	return
}

// ---------------- PANIC BASICS ----------------
func panicBasics() {
	// NOTE: Section-level recover keeps the program running to demonstrate later sections.
	// Without this, unrecovered panics would crash the program.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("RECOVERED (section-level):", r, "(recover only to keep program running)")
		}
	}()

	fmt.Println("--- Manual panic ---")
	manualPanicDemo()
	fmt.Println()

	fmt.Println("--- Runtime panic (index out of range) ---")
	runtimePanicDemo()
	fmt.Println()

	fmt.Println("--- Panic propagation across functions ---")
	propagationDemo()
	fmt.Println()

	fmt.Println("--- Panic in recursion ---")
	recursionPanicDemo(3)
	fmt.Println()
}

func manualPanicDemo() {
	// Manual panic stops normal execution and starts stack unwinding.
	fmt.Println("Before manual panic")
	panic("manual panic triggered")
	fmt.Println("After panic (never runs)")
}

func runtimePanicDemo() {
	// Runtime panics are triggered by invalid operations (e.g., out-of-bounds index).
	arr := []int{1, 2, 3}
	fmt.Println("Accessing index 3 of arr (length 3)")
	_ = arr[3]
	fmt.Println("After index access (never runs)")
}

func propagationDemo() {
	// Panic propagates up the call stack, running defers in each frame until recovered.
	fmt.Println("Calling function A")
	funcA()
}

func funcA() {
	defer fmt.Println("Defer in funcA (runs when funcA unwinds)")
	fmt.Println("Calling function B")
	funcB()
}

func funcB() {
	defer fmt.Println("Defer in funcB (runs when funcB unwinds)")
	fmt.Println("Calling function C")
	funcC()
}

func funcC() {
	defer fmt.Println("Defer in funcC (runs when funcC unwinds)")
	fmt.Println("Panicking in funcC")
	panic("panic from funcC")
}

func recursionPanicDemo(n int) {
	// Panic in recursive functions unwinds all recursive frames.
	defer fmt.Println("Defer in recursion level", n)
	if n == 0 {
		fmt.Println("Panicking at recursion level 0")
		panic("recursion panic")
	}
	fmt.Println("Recursion level", n, "calling", n-1)
	recursionPanicDemo(n - 1)
}

// ---------------- RECOVER BASICS ----------------
func recoverBasics() {
	fmt.Println("--- Recover inside defer (correct usage) ---")
	correctRecoverDemo()
	fmt.Println()

	fmt.Println("--- Recover returning panic value ---")
	recoverWithValueDemo()
	fmt.Println()

	fmt.Println("--- Failed recover (wrong placement) ---")
	failedRecoverDemo()
	fmt.Println()

	fmt.Println("--- Recover outside defer does nothing ---")
	recoverOutsideDeferDemo()
	fmt.Println()
}

func correctRecoverDemo() {
	// Recover only works inside a deferred function. Stops panic propagation and returns panic value.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered correctly:", r)
		}
	}()
	fmt.Println("Before panic")
	panic("test panic")
}

func recoverWithValueDemo() {
	// Recover returns the value passed to panic (any type).
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered panic value: type=%T, value=%v\n", r, r)
		}
	}()
	panic(42)
}

func failedRecoverDemo() {
	// Recover outside a deferred function returns nil and does nothing.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered (section defer):", r)
			fmt.Println("Note: Earlier recover outside defer had no effect")
		}
	}()
	fmt.Println("Calling recover outside defer (wrong placement)")
	r := recover()
	fmt.Println("Recover outside defer returned:", r)
	panic("panic after failed recover")
}

func recoverOutsideDeferDemo() {
	// Recover outside defer with no active panic returns nil.
	fmt.Println("Calling recover outside defer (no active panic)")
	r := recover()
	fmt.Println("Recover returned:", r)
}

// ---------------- COMBINED FLOW ----------------
func combinedFlow() {
	fmt.Println("--- Full flow: panic → defer → recover ---")
	fullFlowDemo()
	fmt.Println()

	fmt.Println("--- Stack unwinding with multiple functions ---")
	stackUnwindingDemo()
	fmt.Println()

	fmt.Println("--- Logging style recovery (like middleware) ---")
	loggingRecoveryDemo()
	fmt.Println()
}

func fullFlowDemo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Full flow recovered:", r)
		}
	}()
	fmt.Println("Before panic")
	panic("full flow panic")
}

func stackUnwindingDemo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Stack unwinding recovered:", r)
		}
	}()
	defer fmt.Println("Defer in stackUnwindingDemo (first to run during unwind)")
	fmt.Println("Calling inner function")
	innerFunc()
}

func innerFunc() {
	defer fmt.Println("Defer in innerFunc (runs when innerFunc unwinds)")
	fmt.Println("Calling deeper function")
	deeperFunc()
}

func deeperFunc() {
	defer fmt.Println("Defer in deeperFunc (runs when deeperFunc unwinds)")
	fmt.Println("Panicking in deeperFunc")
	panic("stack unwinding panic")
}

func loggingRecoveryDemo() {
	// Mimics middleware-style recovery: log panics for debugging.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("LOG: panic recovered:", r)
			fmt.Println("LOG: add debug.Stack() in production for stack traces")
		}
	}()
	fmt.Println("Doing work...")
	panic("middleware-style panic")
}

// ---------------- ADVANCED CASES ----------------
func advancedCases() {
	fmt.Println("--- Recover inside goroutine (correct handling) ---")
	goroutineRecoverDemo()
	fmt.Println()

	fmt.Println("--- Panic in goroutine crashes program if not recovered ---")
	fmt.Println("The next demo will crash the program intentionally:")
	goroutineCrashDemo()
}

func goroutineRecoverDemo() {
	// Each goroutine has its own call stack. Recover only works for panics in the same goroutine.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Goroutine recovered:", r)
			}
		}()
		fmt.Println("Goroutine doing work")
		panic("goroutine panic")
	}()
	wg.Wait()
}

func goroutineCrashDemo() {
	// Unrecovered goroutine panics crash the entire program.
	// Recover in one goroutine cannot catch panics from another.
	go func() {
		fmt.Println("Goroutine panicking...")
		panic("unrecovered goroutine panic")
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("Main done (this will not print if the goroutine crashes first)")
}

// ---------------- MAIN ----------------
func main() {
	fmt.Println("===== STARTING DEFER, PANIC, RECOVER EXAMPLES =====")
	fmt.Println()

	fmt.Println("===== DEFER BASICS =====")
	deferBasics()
	fmt.Println()

	fmt.Println("===== PANIC BASICS =====")
	panicBasics()
	fmt.Println()

	fmt.Println("===== RECOVER BASICS =====")
	recoverBasics()
	fmt.Println()

	fmt.Println("===== COMBINED FLOW =====")
	combinedFlow()
	fmt.Println()

	fmt.Println("===== ADVANCED CASES =====")
	advancedCases()

	fmt.Println("===== ALL EXAMPLES COMPLETED =====")
}
