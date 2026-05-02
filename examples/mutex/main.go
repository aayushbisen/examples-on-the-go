package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
	===========================================================================
	MUTEX IN GO - COMPLETE GUIDE
	===========================================================================

	WHAT IS A MUTEX?
	----------------
	Mutex stands for "Mutual Exclusion". It ensures that only one goroutine
	can access a critical section of code at a time.

	REAL-WORLD ANALOGY:
	Think of a mutex like a key to a single-occupancy bathroom:
	- You take the key and lock the door (Lock)
	- Others must wait until you're done
	- You unlock and return the key (Unlock)
	- Next person can now enter

	WHY DO WE NEED MUTEX?
	---------------------
	Go allows multiple goroutines to run concurrently. When they access the
	same variable AND at least one is writing, you get a "race condition".

	Race Condition: Two or more goroutines access shared data simultaneously,
	and at least one is writing. The result depends on timing (non-deterministic).

	QUICK DECISION GUIDE:
	---------------------
	- Simple counter?     → Use sync/atomic (Example 2)
	- Multiple operations? → Use sync.Mutex (Example 3)
	- Many reads, few writes? → Use sync.RWMutex (Examples 5, 6)
	- Read-only data?     → No synchronization needed!

	GOLDEN RULES:
	------------
	1. A mutex protects STATE, not code. Name it after what it protects.
	2. Never copy a struct that contains a mutex (pass by pointer).
	3. Keep locked sections as short as possible.
	4. Don't sleep or do I/O inside a lock.
	5. Always use 'go run -race' during development.
*/

// ============================================================================
// EXAMPLE 1: RACE CONDITION (WHAT NOT TO DO)
// ============================================================================
//
// This example demonstrates what happens when multiple goroutines
// modify the same variable without any synchronization.
//
// The counter++ operation is NOT atomic. It involves:
// 1. Read current value
// 2. Add 1
// 3. Write new value
//
// If two goroutines do this simultaneously, they might both read the same
// value (say 100), both add 1 (to 101), and both write 101. We lost one increment!
//
// Run with: go run -race main.go
// The race detector will flag this as a data race.

func raceConditionExample() {
	fmt.Println("=== Example 1: Race Condition (Without Mutex) ===")
	fmt.Println("Notice: the result is often less than 10000!\n")

	counter := 0
	var wg sync.WaitGroup

	// Launch 10 goroutines, each incrementing 1000 times
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				// RACE CONDITION: This is NOT atomic!
				// Multiple goroutines can interleave here.
				counter++
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter: %d (expected 10000, likely wrong!)\n\n", counter)
}

// ============================================================================
// EXAMPLE 2: SYNC/ATOMIC (BEST FOR SIMPLE COUNTERS)
// ============================================================================
//
// For simple operations like counters, atomic operations are:
// - Faster than mutex (no lock overhead)
// - Cleaner (no Lock/Unlock pairs)
// - Impossible to forget to unlock
//
// atomic.Int64 provides:
// - Add(delta) - atomically add and return new value
// - Load() - atomically read the value
// - Store(val) - atomically write the value
// - CompareAndSwap(old, new) - conditional update
//
// USE ATOMIC WHEN: You have a single value (counter, flag, pointer) that
// needs concurrent access.

func atomicCounterExample() {
	fmt.Println("=== Example 2: Atomic Counter (Preferred for Counters) ===")
	fmt.Println("sync/atomic is faster and simpler than mutex for counters\n")

	// atomic.Int64 is a zero-value usable atomic counter
	var counter atomic.Int64
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				// Add is atomic - safe for concurrent use
				// No lock needed, no race condition possible
				counter.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter: %d (always correct!)\n\n", counter.Load())
}

// ============================================================================
// EXAMPLE 3: MUTEX FOR RELATED OPERATIONS
// ============================================================================
//
// A mutex is needed when you have MULTIPLE related operations that must
// happen as a single unit (atomically).
//
// In this example, we read the old value, compute a new value based on it,
// and write it back. These three steps must not be interrupted by other
// goroutines, otherwise they'll read stale data.
//
// KEY POINT: Mutex protects the ENTIRE sequence of operations, not just
// the write. This is called a "critical section".

func mutexBasicExample() {
	fmt.Println("=== Example 3: Mutex for Related Operations ===")
	fmt.Println("Mutex protects multiple operations that must happen together\n")

	// SafeCounter demonstrates the classic mutex pattern:
	// - mu protects 'value'
	// - Always Lock before accessing value
	// - Always Unlock after (use defer for safety)
	type SafeCounter struct {
		mu    sync.Mutex // Protects 'value' field
		value int
	}

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Lock: Enter critical section
			counter.mu.Lock()

			// Critical section: Multiple related operations
			// These must happen atomically (as one unit)
			old := counter.value
			// Simulate some computation based on current value
			counter.value = old + (id + 1)

			fmt.Printf("Goroutine %d: %d -> %d\n", id, old, counter.value)

			// Unlock: Leave critical section
			counter.mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final value: %d\n\n", counter.value)
}

// ============================================================================
// EXAMPLE 4: BANK ACCOUNT (REAL-WORLD WITH PROPER LOCK HYGIENE)
// ============================================================================
//
// This example shows PROPER mutex usage patterns:
//
// 1. VALIDATE INPUT BEFORE LOCKING
//    Don't acquire a lock just to find out the input is invalid.
//    Validation is single-threaded, do it outside the lock.
//
// 2. UNLOCK BEFORE I/O (fmt.Printf, network, file operations)
//    I/O is slow. Holding a lock during I/O blocks all other goroutines.
//    Read the data, unlock, then print.
//
// 3. HANDLE EARLY RETURNS / ERRORS
//    If you have multiple return paths, use 'defer' to ensure unlock.
//    Or manually unlock on each path (shown in Withdraw).
//
// 4. METHODS THAT READ STATE
//    GetBalance() also needs the lock because another goroutine
//    might be writing at the same time.

type BankAccount struct {
	mu      sync.Mutex // Protects: balance
	balance int
}

// Deposit adds money to the account.
// Best practice: Validate before locking, unlock before I/O.
func (a *BankAccount) Deposit(amount int) error {
	// VALIDATION: Do this BEFORE acquiring the lock
	// No need to block other goroutines for invalid input
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive: %d", amount)
	}

	// LOCK: Enter critical section
	a.mu.Lock()

	// Read state while holding lock
	oldBalance := a.balance

	// Modify state
	a.balance += amount
	newBalance := a.balance

	// UNLOCK: Before I/O (printing)
	// This allows other goroutines to proceed while we print
	a.mu.Unlock()

	fmt.Printf("Deposited %d, balance: %d -> %d\n", amount, oldBalance, newBalance)
	return nil
}

// Withdraw removes money from the account if sufficient funds exist.
// This shows manual unlock on different code paths (no defer).
func (a *BankAccount) Withdraw(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive: %d", amount)
	}

	a.mu.Lock()

	// CHECK: Insufficient funds?
	if a.balance < amount {
		// MANUAL UNLOCK: Must unlock before returning!
		// Forgetting this causes a DEADLOCK.
		a.mu.Unlock()
		fmt.Printf("Failed to withdraw %d, balance: %d (insufficient funds)\n", amount, a.balance)
		return fmt.Errorf("insufficient funds")
	}

	// Update state
	a.balance -= amount
	newBalance := a.balance

	// MANUAL UNLOCK: After successful withdrawal
	a.mu.Unlock()

	fmt.Printf("Withdrew %d, new balance: %d\n", amount, newBalance)
	return nil
}

// GetBalance returns the current balance.
// Even reads need the lock if writes can happen concurrently.
func (a *BankAccount) GetBalance() int {
	a.mu.Lock()
	defer a.mu.Unlock() // defer is safe here since no early returns
	return a.balance
}

func bankAccountExample() {
	fmt.Println("=== Example 4: Real-World - Bank Account ===")
	fmt.Println("Shows: validate before lock, unlock before I/O\n")

	account := &BankAccount{balance: 1000}
	var wg sync.WaitGroup

	operations := []struct {
		name   string
		amount int
		op     func(int) error
	}{
		{"Deposit", 5000, account.Deposit},
		{"Deposit", 2000, account.Deposit},
		{"Withdraw", 1000, account.Withdraw},
		{"Withdraw", 3000, account.Withdraw},
		{"Withdraw", 5000, account.Withdraw}, // Will fail - insufficient funds
		{"Deposit", -100, account.Deposit},   // Will fail - invalid amount
	}

	for _, op := range operations {
		wg.Add(1)
		go func(name string, amount int, fn func(int) error) {
			defer wg.Done()
			if err := fn(amount); err != nil {
				fmt.Printf("%s %d failed: %v\n", name, amount, err)
			}
		}(op.name, op.amount, op.op)
	}

	wg.Wait()
	fmt.Printf("Final balance: %d\n\n", account.GetBalance())
}

// ============================================================================
// EXAMPLE 5: TICKET BOOKING (RW MUTEX FOR READ-HEAVY WORKLOADS)
// ============================================================================
//
// sync.RWMutex is a reader/writer mutex:
// - RLock() / RUnlock() - for READING (multiple goroutines can read at once)
// - Lock() / Unlock()   - for WRITING (only one goroutine, no readers either)
//
// USE RW MUTEX WHEN:
// - You have many more reads than writes
// - Reads are frequent, writes are rare
//
// In this example:
// - Book() uses Lock() because it modifies state
// - Available() uses RLock() because it only reads
// - Multiple people can check availability simultaneously

type TicketBooking struct {
	mu        sync.RWMutex // Protects: available, totalSold
	available int
	totalSold int
}

// Book attempts to book tickets. Uses exclusive lock (Lock).
func (t *TicketBooking) Book(user string, count int) bool {
	if count <= 0 {
		fmt.Printf("%s: invalid ticket count %d\n", user, count)
		return false
	}

	// EXCLUSIVE LOCK: No other readers or writers can proceed
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check availability
	if t.available < count {
		fmt.Printf("%s failed: only %d tickets left\n", user, t.available)
		return false
	}

	// Process booking
	t.available -= count
	t.totalSold += count
	fmt.Printf("%s booked %d ticket(s), %d remaining\n", user, count, t.available)
	return true
}

// Available returns remaining tickets. Uses shared read lock (RLock).
func (t *TicketBooking) Available() int {
	// SHARED LOCK: Multiple goroutines can call this concurrently
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.available
}

func ticketBookingExample() {
	fmt.Println("=== Example 5: Ticket Booking (RWMutex) ===")
	fmt.Println("Writes are exclusive, reads can happen concurrently\n")

	booking := &TicketBooking{available: 10}

	users := []struct {
		name  string
		count int
	}{
		{"Alice", 3}, {"Bob", 4}, {"Charlie", 5}, {"Diana", 2},
	}

	var wg sync.WaitGroup
	for _, u := range users {
		wg.Add(1)
		go func(user string, count int) {
			defer wg.Done()
			booking.Book(user, count)
		}(u.name, u.count)
	}

	wg.Wait()
	fmt.Printf("Sold: %d, Remaining: %d\n\n", booking.totalSold, booking.Available())
}

// ============================================================================
// EXAMPLE 6: SAFE CACHE (RW MUTEX WITH MAP)
// ============================================================================
//
// This shows a concurrent-safe cache using RWMutex with a map.
//
// WHY RWMutex HERE?
// - Set() is called rarely (writes)
// - Get() is called frequently (reads)
// - Multiple Get() calls can run simultaneously
//
// IMPORTANT: Returning a value from inside the lock is SAFE because:
// - string is a value type (copied on return)
// - If we returned a pointer or slice, we'd need to be more careful
//
// NOTE: For production, consider sync.Map which is optimized for this pattern.

type SafeCache struct {
	mu   sync.RWMutex       // Protects: data map
	data map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{data: make(map[string]string)}
}

// Set adds or updates a key-value pair. Exclusive lock needed.
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	fmt.Printf("SET %s = %s\n", key, value)
}

// Get retrieves a value by key. Shared read lock.
func (c *SafeCache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Safe to return: string is a value type (copied)
	return c.data[key]
}

func cacheExample() {
	fmt.Println("=== Example 6: RWMutex Cache ===")
	fmt.Println("Multiple readers concurrent, exclusive writer\n")

	cache := NewSafeCache()
	cache.Set("user:1", "Alice")
	cache.Set("user:2", "Bob")

	var wg sync.WaitGroup

	// Many concurrent readers - they all run at the same time!
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val := cache.Get("user:1")
			fmt.Printf("Reader %d: %s\n", id, val)
		}(i)
	}

	// One writer - blocks readers while active
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.Set("user:1", "Alice Updated")
	}()

	wg.Wait()
	fmt.Println()
}

// ============================================================================
// MAIN: RUN ALL EXAMPLES
// ============================================================================

func main() {
	fmt.Println("========================================")
	fmt.Println("   Understanding Mutex in Go")
	fmt.Println("========================================\n")

	// Uncomment to see race condition in action:
	// raceConditionExample()

	atomicCounterExample()
	mutexBasicExample()
	bankAccountExample()
	ticketBookingExample()
	cacheExample()

	fmt.Println("========================================")
	fmt.Println("   End of Mutex Examples")
	fmt.Println("========================================")
	fmt.Println("\nTIP: Run with 'go run -race main.go' to detect race conditions!")
}
