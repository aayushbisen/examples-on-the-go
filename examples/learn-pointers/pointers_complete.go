package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// GO POINTERS - COMPLETE GUIDE (8 sections in one runnable file)
// Run: go run pointers_complete.go
// ============================================================================

// --- Types for Section 3 ---
type Person struct {
	Name string
	Age  int
}

// --- Types for Section 4 ---
type Counter struct {
	count int
}

// Value receiver: gets a COPY
func (c Counter) IncrementValue() {
	c.count++
	fmt.Printf("  IncrementValue: count = %d (local copy)\n", c.count)
}

// Pointer receiver: gets the ADDRESS
func (c *Counter) IncrementPointer() {
	c.count++
	fmt.Printf("  IncrementPointer: count = %d\n", c.count)
}

func (c *Counter) Reset() {
	c.count = 0
	fmt.Printf("  Reset: count set to 0\n")
}

func (c Counter) GetValue() int {
	return c.count
}

// --- Types for Section 5 ---
type DataStore struct {
	Users   []string
	Configs map[string]string
}

// --- Types for Section 6 ---
type HeapUser struct {
	Name  string
	Score int
}

// --- Types for Section 7 ---
type ServiceConfig struct {
	Debug      bool
	MaxRetries int
	Timeout    int
}

type SvcUser struct {
	ID    string
	Name  string
	Email string
}

type UserService struct {
	config *ServiceConfig
	mu     sync.RWMutex
	users  map[string]*SvcUser
}

func NewUserService(cfg *ServiceConfig) *UserService {
	return &UserService{
		config: cfg,
		users:  make(map[string]*SvcUser),
	}
}

func (s *UserService) CreateUser(id, name, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[id] = &SvcUser{ID: id, Name: name, Email: email}
	if s.config.Debug {
		fmt.Printf("  [DEBUG] Created user: %s\n", name)
	}
}

func (s *UserService) GetUser(id string) *SvcUser {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *UserService) UpdateEmail(id, newEmail string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[id]
	if !exists {
		return false
	}
	user.Email = newEmail
	if s.config.Debug {
		fmt.Printf("  [DEBUG] Updated email for %s\n", user.Name)
	}
	return true
}

func (s *UserService) ListUsers() []*SvcUser {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*SvcUser, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result
}

// --- Types for Section 8 ---
type Account struct {
	Balance float64
	Owner   string
}

type BigStruct struct {
	Data []byte
}

type SharedState struct {
	Value int
}

func main() {

	// ========================================================================
	// SECTION 1: BASICS OF POINTERS
	// & gets the address, * reads/writes the value at that address
	// ========================================================================
	fmt.Println("=== SECTION 1: BASICS ===")

	age := 25
	fmt.Printf("Value of age: %d\n", age)

	address := &age
	fmt.Printf("Memory address of age: %p\n", address)

	var ptr *int = &age
	fmt.Printf("Type of ptr: %T\n", ptr)

	fmt.Printf("Value at address ptr points to: %d\n", *ptr)

	*ptr = 30
	fmt.Printf("Age after modification through pointer: %d\n", age)

	name := "Alice"
	namePtr := &name
	fmt.Printf("String pointer - Value: %s, Address: %p\n", *namePtr, namePtr)

	var emptyPtr *int
	fmt.Printf("Uninitialized pointer value: %v\n", emptyPtr)
	// fmt.Println(*emptyPtr) // panic: nil pointer dereference

	fmt.Println()

	// ========================================================================
	// SECTION 2: POINTERS VS VALUES
	// Pass by value = copy, Pass by pointer = original is modified
	// ========================================================================
	fmt.Println("=== SECTION 2: VALUES VS POINTERS ===")

	x := 10
	fmt.Printf("Before modifyValue: x = %d\n", x)
	modifyValue(x)
	fmt.Printf("After modifyValue:  x = %d\n", x)

	fmt.Println()

	y := 10
	fmt.Printf("Before modifyPointer: y = %d\n", y)
	modifyPointer(&y)
	fmt.Printf("After modifyPointer:  y = %d\n", y)

	fmt.Println()

	a, b := 1, 2
	fmt.Printf("Before swap: a = %d, b = %d\n", a, b)
	swap(&a, &b)
	fmt.Printf("After swap:  a = %d, b = %d\n", a, b)

	fmt.Println()

	// ========================================================================
	// SECTION 3: POINTERS WITH STRUCTS
	// Structs are value types - pointers avoid copies and allow modification
	// ========================================================================
	fmt.Println("=== SECTION 3: STRUCTS ===")

	p := Person{Name: "Alice", Age: 25}
	fmt.Printf("Original: %+v\n", p)

	modifyPersonCopy(p)
	fmt.Printf("After modifyPersonCopy: %+v\n", p)

	fmt.Println()

	modifyPersonPtr(&p)
	fmt.Printf("After modifyPersonPtr:  %+v\n", p)

	fmt.Println()

	ptr2 := &p
	fmt.Printf("Direct field access through pointer: %s\n", ptr2.Name)

	fmt.Println()

	p2 := new(Person)
	p2.Name = "Charlie"
	p2.Age = 30
	fmt.Printf("Struct created with new(): %+v\n", *p2)

	fmt.Println()

	p3 := &Person{Name: "Diana", Age: 28}
	fmt.Printf("Pointer to struct literal: %+v\n", *p3)

	fmt.Println()

	// ========================================================================
	// SECTION 4: METHODS AND RECEIVERS
	// Value receivers work on copies, pointer receivers modify the original
	// ========================================================================
	fmt.Println("=== SECTION 4: RECEIVERS ===")

	fmt.Println("--- Value Receiver (copy) ---")
	c1 := Counter{count: 0}
	c1.IncrementValue()
	c1.IncrementValue()
	c1.IncrementValue()
	fmt.Printf("Final count after value receiver calls: %d\n", c1.GetValue())

	fmt.Println()
	fmt.Println("--- Pointer Receiver (original) ---")
	c2 := Counter{count: 0}
	c2.IncrementPointer()
	c2.IncrementPointer()
	c2.IncrementPointer()
	fmt.Printf("Final count after pointer receiver calls: %d\n", c2.GetValue())

	fmt.Println()
	fmt.Println("--- Reset (requires pointer) ---")
	c4 := Counter{count: 42}
	fmt.Printf("Before reset: %d\n", c4.count)
	c4.Reset()
	fmt.Printf("After reset: %d\n", c4.count)

	fmt.Println()

	// ========================================================================
	// SECTION 5: SLICES, MAPS, AND POINTERS
	// Slices/maps are already reference-like - pointers rarely needed
	// ========================================================================
	fmt.Println("=== SECTION 5: SLICES & MAPS ===")

	fmt.Println("--- Slice behavior ---")
	numbers := []int{1, 2, 3}
	fmt.Printf("Original slice: %v\n", numbers)
	modifySlice(numbers)
	fmt.Printf("After modifySlice: %v\n", numbers)

	fmt.Println()

	fmt.Println("--- Slice append caveat ---")
	names := []string{"Alice", "Bob"}
	fmt.Printf("Original names: %v (len=%d, cap=%d)\n", names, len(names), cap(names))
	addToSlice(names)
	fmt.Printf("After addToSlice: %v\n", names)

	fmt.Println()

	fmt.Println("--- Map behavior ---")
	scores := map[string]int{"Alice": 90, "Bob": 85}
	fmt.Printf("Original map: %v\n", scores)
	modifyMap(scores)
	fmt.Printf("After modifyMap: %v\n", scores)

	fmt.Println()

	fmt.Println("--- Pointer to slice (when needed) ---")
	data := DataStore{
		Users:   []string{"Alice"},
		Configs: map[string]string{"theme": "dark"},
	}
	fmt.Printf("Before: %+v\n", data)
	addUser(&data, "Bob")
	addUser(&data, "Charlie")
	fmt.Printf("After:  %+v\n", data)

	fmt.Println()

	// ========================================================================
	// SECTION 6: MEMORY BEHAVIOR - STACK VS HEAP
	// Returning pointers is safe - Go's escape analysis handles allocation
	// ========================================================================
	fmt.Println("=== SECTION 6: MEMORY ===")

	u1 := createHeapUser("Alice")
	fmt.Printf("User from createUser: %+v\n", *u1)

	fmt.Println()

	u2 := createUserValue("Bob")
	fmt.Printf("User from createUserValue: %+v\n", u2)

	fmt.Println()

	localPtrDemo()

	fmt.Println()

	users := []*HeapUser{}
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("User%d", i+1)
		users = append(users, createHeapUser(name))
	}
	for _, u := range users {
		fmt.Printf("  %s has score %d\n", u.Name, u.Score)
	}

	fmt.Println()

	// ========================================================================
	// SECTION 7: REAL-WORLD USE CASE - SERVICE LAYER
	// Shared config, dependency injection, thread-safe state
	// ========================================================================
	fmt.Println("=== SECTION 7: REAL-WORLD SERVICE ===")

	cfg := &ServiceConfig{
		Debug:      true,
		MaxRetries: 3,
		Timeout:    30,
	}

	service := NewUserService(cfg)

	fmt.Println("Creating users:")
	service.CreateUser("1", "Alice", "alice@example.com")
	service.CreateUser("2", "Bob", "bob@example.com")

	fmt.Println()

	fmt.Println("Updating user email:")
	service.UpdateEmail("1", "alice.smith@example.com")

	bob := service.GetUser("2")
	if bob != nil {
		bob.Email = "bob.jones@example.com"
		fmt.Printf("  Modified Bob directly: %s\n", bob.Email)
	}

	fmt.Println()

	fmt.Println("All users:")
	for _, u := range service.ListUsers() {
		fmt.Printf("  %s <%s>\n", u.Name, u.Email)
	}

	fmt.Println()

	fmt.Println("Changing shared config:")
	cfg.Debug = false
	fmt.Println("  Debug disabled")
	service.CreateUser("3", "Charlie", "charlie@example.com")

	fmt.Println()

	// ========================================================================
	// SECTION 8: COMMON MISTAKES AND BEST PRACTICES
	// ========================================================================
	fmt.Println("=== SECTION 8: MISTAKES & BEST PRACTICES ===")

	// MISTAKE 1: Nil pointer dereference
	fmt.Println("--- Mistake 1: Nil pointer dereference ---")
	var account *Account
	if account == nil {
		fmt.Println("  Account is nil - cannot access fields")
	}
	account = &Account{Balance: 100, Owner: "Alice"}
	fmt.Printf("  Safe access: %s has $%.2f\n", account.Owner, account.Balance)

	fmt.Println()

	// MISTAKE 2: Unexpected mutations
	fmt.Println("--- Mistake 2: Unexpected mutations ---")
	original := &Account{Balance: 500, Owner: "Bob"}
	fmt.Printf("  Original balance: $%.2f\n", original.Balance)
	applyFee(original)
	fmt.Printf("  After applyFee: $%.2f\n", original.Balance)

	fmt.Println()

	// MISTAKE 3: Pointer to loop variable
	fmt.Println("--- Mistake 3: Pointer to loop variable ---")
	nums := []int{1, 2, 3}
	var ptrs []*int
	for i := range nums {
		ptrs = append(ptrs, &nums[i])
	}
	for _, p := range ptrs {
		fmt.Printf("  Value: %d\n", *p)
	}

	fmt.Println()

	// MISTAKE 4: Overusing pointers
	fmt.Println("--- Mistake 4: Overusing pointers ---")
	val := new(int)
	*val = 42
	fmt.Printf("  Unnecessary pointer: %d\n", *val)
	str := "Alice"
	strPtr := &str
	fmt.Printf("  Pointer to read-only string: %s\n", *strPtr)

	fmt.Println()

	// BEST PRACTICES
	fmt.Println("--- When to use pointers ---")
	a2 := 10
	double(&a2)
	fmt.Printf("  After double: %d\n", a2)

	var maybeCfg *ServiceConfig
	if maybeCfg == nil {
		fmt.Println("  Nil pointer = no config = use defaults")
	}

	large := &BigStruct{Data: make([]byte, 10000)}
	processLarge(large)

	shared := &SharedState{Value: 0}
	go increment(shared)
	go increment(shared)

	fmt.Println()
	fmt.Println("=== SUMMARY ===")
	fmt.Println("1. Pointers hold memory addresses, not values")
	fmt.Println("2. Use & to get an address, * to read/write through it")
	fmt.Println("3. Pass pointers to functions that need to modify the original")
	fmt.Println("4. Use pointer receivers when methods modify the struct")
	fmt.Println("5. Slices/maps are already reference-like - pointers rarely needed")
	fmt.Println("6. Returning pointers is safe - Go handles memory automatically")
	fmt.Println("7. In real apps, pointers enable shared state and dependency injection")
	fmt.Println("8. Avoid pointers for small types, read-only data, and when nil is meaningless")
	fmt.Println("")
	fmt.Println("=== PRACTICE EXERCISES ===")
	fmt.Println("1. Write a function that takes two *int pointers and swaps their values")
	fmt.Println("2. Create a LinkedList struct with *Node pointers and implement Append/Print")
	fmt.Println("3. Build a simple cache (map[string]*Data) with Set/Get/Delete methods")
	fmt.Println("   that uses pointer receivers and handles nil keys gracefully")
}

// --- Section 2 helpers ---
func modifyValue(n int) {
	n = 100
	fmt.Printf("  Inside modifyValue: n = %d\n", n)
}

func modifyPointer(n *int) {
	*n = 100
	fmt.Printf("  Inside modifyPointer: *n = %d\n", *n)
}

func swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

// --- Section 3 helpers ---
func modifyPersonCopy(p Person) {
	p.Name = "Bob"
	fmt.Printf("  Inside modifyPersonCopy: %+v\n", p)
}

func modifyPersonPtr(p *Person) {
	p.Name = "Bob"
	fmt.Printf("  Inside modifyPersonPtr:  %+v\n", *p)
}

// --- Section 5 helpers ---
func modifySlice(s []int) {
	if len(s) > 0 {
		s[0] = 999
	}
	fmt.Printf("  Inside modifySlice: %v\n", s)
}

func addToSlice(s []string) {
	s = append(s, "Charlie")
	fmt.Printf("  Inside addToSlice: %v\n", s)
}

func modifyMap(m map[string]int) {
	m["Charlie"] = 95
	m["Alice"] = 100
	delete(m, "Bob")
	fmt.Printf("  Inside modifyMap: %v\n", m)
}

func addUser(ds *DataStore, name string) {
	ds.Users = append(ds.Users, name)
}

// --- Section 6 helpers ---
func createHeapUser(name string) *HeapUser {
	return &HeapUser{Name: name, Score: 100}
}

func createUserValue(name string) HeapUser {
	return HeapUser{Name: name, Score: 100}
}

func localPtrDemo() {
	x := 42
	ptr := &x
	fmt.Printf("  Local pointer value: %d\n", *ptr)
}

// --- Section 8 helpers ---
func applyFee(a *Account) {
	a.Balance -= 10
}

func double(n *int) {
	*n *= 2
}

func processLarge(l *BigStruct) {
	_ = len(l.Data)
}

func increment(s *SharedState) {
	s.Value++
}
