package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// ============================================================
// GO ERROR HANDLING - COMPLETE LEARNING GUIDE
// ============================================================
// This single file covers every major error handling concept in Go.
// Run it with: go run main.go
//
// Topics covered:
//   1. Basic errors (errors.New, fmt.Errorf, if err != nil)
//   2. Custom error types (struct + Error() method)
//   3. Error wrapping (fmt.Errorf %w, errors.Is, errors.As)
//   4. Sentinel errors (package-level Err variables)
//   5. Panic and recover (when to use, defer + recover pattern)
//   6. Defer usage (cleanup, LIFO, timing)
//   7. Multiple errors (early return, error collectors)
//   8. Real-world file I/O with proper error handling
//   9. HTTP error handling (status codes, middleware recovery)
//  10. Common mistakes and how to avoid them
// ============================================================

// ============================================================
// SENTINEL ERRORS (used across multiple sections)
// ============================================================

var (
	ErrNotFound        = errors.New("record not found")
	ErrUserNotFound    = errors.New("users: user not found")
	ErrUserExists      = errors.New("users: user already exists")
	ErrInvalidInput    = errors.New("users: invalid user input")
	ErrDatabase        = errors.New("database connection failed")
	ErrInvalidPayment  = errors.New("payment processing failed")
)

// ============================================================
// CUSTOM ERROR TYPES
// ============================================================

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s (got: %v)",
		e.Field, e.Message, e.Value)
}

type DatabaseError struct {
	Code      int
	Table     string
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	base := fmt.Sprintf("database error [%d]: %s on table '%s'",
		e.Code, e.Operation, e.Table)
	if e.Err != nil {
		return fmt.Sprintf("%s - underlying: %v", base, e.Err)
	}
	return base
}

type AuthError struct {
	Reason    string
	Retryable bool
}

func (e *AuthError) Error() string {
	if e.Retryable {
		return fmt.Sprintf("auth error (retryable): %s", e.Reason)
	}
	return fmt.Sprintf("auth error (permanent): %s", e.Reason)
}

type PaymentError struct {
	CardNumber string
	Amount     float64
	Reason     string
}

func (e *PaymentError) Error() string {
	return fmt.Sprintf("payment failed for card %s ($%.2f): %s",
		e.CardNumber, e.Amount, e.Reason)
}

type ValidationErrorCollector struct {
	errors []string
}

func (c *ValidationErrorCollector) Add(field, message string) {
	c.errors = append(c.errors, fmt.Sprintf("%s: %s", field, message))
}

func (c *ValidationErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

func (c *ValidationErrorCollector) Error() string {
	if len(c.errors) == 0 {
		return ""
	}
	result := fmt.Sprintf("%d validation error(s):", len(c.errors))
	for _, e := range c.errors {
		result += "\n  - " + e
	}
	return result
}

type User struct {
	ID    int
	Name  string
	Email string
}

type UserService struct {
	users  map[int]User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
		},
		nextID: 2,
	}
}

func (s *UserService) GetUser(id int) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidInput
	}
	user, exists := s.users[id]
	if !exists {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(name, email string) (User, error) {
	if name == "" || email == "" {
		return User{}, ErrInvalidInput
	}
	for _, u := range s.users {
		if u.Email == email {
			return User{}, ErrUserExists
		}
	}
	user := User{ID: s.nextID, Name: name, Email: email}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

type Resource struct {
	name string
}

func (r *Resource) Acquire() {
	fmt.Printf("   Resource '%s' acquired\n", r.name)
}

func (r *Resource) Release() {
	fmt.Printf("   Resource '%s' released\n", r.name)
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

func parseAge(age int) (string, error) {
	if age < 0 {
		return "", fmt.Errorf("invalid age %d: age cannot be negative", age)
	}
	if age > 150 {
		return "", fmt.Errorf("invalid age %d: exceeds maximum human lifespan", age)
	}
	return fmt.Sprintf("Age %d is valid", age), nil
}

func processUser(name string, age int) {
	msg, err := parseAge(age)
	if err != nil {
		fmt.Printf("  User '%s': FAILED - %v\n", name, err)
		return
	}
	fmt.Printf("  User '%s': %s\n", name, msg)
}

func validateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Value: email, Message: "email is required"}
	}
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email", Value: email, Message: "email must contain '@'"}
	}
	return nil
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Value: age, Message: "age must be between 0 and 150"}
	}
	return nil
}

func handleError(err error) {
	switch e := err.(type) {
	case *ValidationError:
		fmt.Printf("   [VALIDATION] Field '%s': %s\n", e.Field, e.Message)
	case *AuthError:
		if e.Retryable {
			fmt.Printf("   [AUTH] Retryable: %s\n", e.Reason)
		} else {
			fmt.Printf("   [AUTH] Permanent: %s\n", e.Reason)
		}
	default:
		fmt.Printf("   [UNKNOWN] %v\n", e)
	}
}

func findUser(id int) error {
	err := ErrNotFound
	err = fmt.Errorf("user repository: %w", err)
	err = fmt.Errorf("get user %d: %w", id, err)
	return err
}

func processPayment(card string, amount float64) error {
	paymentErr := &PaymentError{
		CardNumber: card,
		Amount:     amount,
		Reason:     "insufficient funds",
	}
	return fmt.Errorf("checkout service: %w", paymentErr)
}

func runWithRecover(fn func() string) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered from: %v", r)
		}
	}()
	return fn()
}

func demonstrateInitPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   Caught init panic: %v\n", r)
		}
	}()
	configLoaded := false
	if !configLoaded {
		panic("required configuration file not found")
	}
}

func demonstrateInvariantPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   Caught invariant panic: %v\n", r)
		}
	}()
	status := "unknown"
	if status == "unknown" {
		panic("invariant violated: status should never be 'unknown' here")
	}
}

func demonstrateAssertion() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   Caught assertion: %v\n", r)
		}
	}()
	value := -5
	if value < 0 {
		panic(fmt.Sprintf("assertion failed: expected positive value, got %d", value))
	}
}

func safeDivide(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   safeDivide(%d, %d) recovered: %v\n", a, b, r)
		}
	}()
	if b == 0 {
		panic("division by zero")
	}
	fmt.Printf("   safeDivide(%d, %d) = %d\n", a, b, a/b)
}

func simulateHTTPHandler(path string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   [RECOVER] Handler panicked: %v\n", r)
		}
	}()
	if path == "/api/crash" {
		panic("nil pointer dereference in handler")
	}
	return "200 OK"
}

func demonstrateLIFO() {
	fmt.Println("   Registering defers 1, 2, 3...")
	defer fmt.Println("   Defer 1 runs LAST")
	defer fmt.Println("   Defer 2 runs SECOND")
	defer fmt.Println("   Defer 3 runs FIRST")
	fmt.Println("   Function body ends, defers execute:")
}

func readFileWithDefer() error {
	f, err := os.CreateTemp("", "defer-example-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString("Hello from defer!\n")
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	fmt.Printf("   Wrote to temp file: %s\n", f.Name())
	fmt.Println("   File will be closed automatically by defer")
	os.Remove(f.Name())
	return nil
}

func demonstrateArgumentEvaluation() {
	value := 1
	defer fmt.Printf("   Deferred prints: %d (captured at registration)\n", value)
	value = 2
	fmt.Printf("   Current value: %d (changed after defer registration)\n", value)
}

func slowOperation() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("   slowOperation took %v\n", elapsed)
	}()
	for i := 0; i < 1000000; i++ {
		_ = i * 2
	}
}

func fastOperation() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("   fastOperation took %v\n", elapsed)
	}()
	fmt.Println("   Quick work done")
}

func demonstrateMethodDefer() {
	r := Resource{name: "database-connection"}
	r.Acquire()
	defer r.Release()
	r2 := Resource{name: "file-handle"}
	r2.Acquire()
	defer r2.Release()
	fmt.Println("   Doing work with resources...")
}

func processOrder(orderID string, quantity int) (string, error) {
	if orderID == "" {
		return "", fmt.Errorf("step 1 (validate): order ID is required")
	}
	fmt.Printf("   Step 1: Validated order %s\n", orderID)
	if quantity <= 0 {
		return "", fmt.Errorf("step 2 (inventory): item out of stock")
	}
	fmt.Printf("   Step 2: Checked inventory (%d items)\n", quantity)
	if quantity > 50 {
		return "", fmt.Errorf("step 3 (payment): payment declined for large order")
	}
	fmt.Println("   Step 3: Payment processed")
	if quantity < 5 {
		return "", fmt.Errorf("step 4 (shipping): minimum order of 5 for shipping")
	}
	fmt.Println("   Step 4: Shipping arranged")
	return fmt.Sprintf("Order %s processed (%d items)", orderID, quantity), nil
}

func validateAndReport(name, email string, age int) {
	fmt.Printf("   Validating: name=%q, email=%q, age=%d\n", name, email, age)
	collector := &ValidationErrorCollector{}
	if name == "" || len(name) < 2 {
		collector.Add("name", "must be at least 2 characters")
	}
	if email == "" || !containsAt(email) {
		collector.Add("email", "must be a valid email address")
	}
	if age < 0 || age > 150 {
		collector.Add("age", "must be between 0 and 150")
	}
	if collector.HasErrors() {
		fmt.Printf("   FAILED: %v\n", collector)
	} else {
		fmt.Printf("   PASSED: all fields valid\n")
	}
}

func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}

func provisionServer(name string) error {
	fmt.Printf("   Provisioning server: %s\n", name)
	fmt.Printf("   Step 1: Created VM for %s\n", name)
	fmt.Printf("   Step 2: Configured networking for %s\n", name)
	err := fmt.Errorf("step 3: software installation failed (timeout)")
	if err != nil {
		fmt.Printf("   Rollback: Removing networking for %s\n", name)
		fmt.Printf("   Rollback: Destroying VM for %s\n", name)
		return fmt.Errorf("provisioning %s failed: %w", name, err)
	}
	fmt.Printf("   Server %s provisioned successfully\n", name)
	return nil
}

func main() {
	// ========================================================
	// SECTION 1: BASIC ERRORS
	// ========================================================
	// The built-in error type is an interface: type error interface { Error() string }
	// Functions that can fail return (value, error). Always check with if err != nil.

	fmt.Println("============================================================")
	fmt.Println("SECTION 1: BASIC ERRORS")
	fmt.Println("============================================================")
	fmt.Println()

	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("divide(10, 0) failed: %v\n", err)
	}

	fmt.Println()

	msg, err := parseAge(25)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(msg)
	}

	msg, err = parseAge(-5)
	if err != nil {
		fmt.Printf("parseAge(-5) failed: %v\n", err)
	}

	fmt.Println()
	fmt.Println("Early return pattern:")
	processUser("Alice", 30)
	processUser("Bob", -1)

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Functions that can fail return (value, error)")
	fmt.Println("  2. Always check: if err != nil { handle it }")
	fmt.Println("  3. Use errors.New() for static errors")
	fmt.Println("  4. Use fmt.Errorf() for errors with dynamic values")
	fmt.Println("  5. Return early on error to avoid deep nesting")

	// ========================================================
	// SECTION 2: CUSTOM ERROR TYPES
	// ========================================================
	// Create a struct with an Error() string method to satisfy the error interface.
	// This lets you attach structured data (fields, codes, metadata) to errors.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 2: CUSTOM ERROR TYPES")
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Println("1. Validation errors with structured data:")

	err = validateEmail("")
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	err = validateEmail("not-an-email")
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	err = validateAge(-5)
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	fmt.Println()
	fmt.Println("2. Type assertion to access custom fields:")

	err = validateEmail("bad-email")
	if err != nil {
		if vErr, ok := err.(*ValidationError); ok {
			fmt.Printf("   Field: %s\n", vErr.Field)
			fmt.Printf("   Value: %v\n", vErr.Value)
			fmt.Printf("   Message: %s\n", vErr.Message)
		}
	}

	fmt.Println()
	fmt.Println("3. Custom errors wrapping other errors:")

	dbErr := &DatabaseError{
		Code:      1045,
		Table:     "users",
		Operation: "INSERT",
		Err:       fmt.Errorf("connection refused"),
	}
	fmt.Printf("   %v\n", dbErr)

	fmt.Println()
	fmt.Println("4. Type switch for error-specific handling:")

	errs := []error{
		&ValidationError{Field: "name", Value: "", Message: "required"},
		&AuthError{Reason: "token expired", Retryable: true},
		&AuthError{Reason: "invalid credentials", Retryable: false},
		fmt.Errorf("unexpected error"),
	}

	for _, e := range errs {
		handleError(e)
	}

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Custom errors = struct + Error() string method")
	fmt.Println("  2. They carry structured data, not just messages")
	fmt.Println("  3. Use type assertion (err.(*Type)) to access fields")
	fmt.Println("  4. Use type switch to handle different error types")

	// ========================================================
	// SECTION 3: ERROR WRAPPING
	// ========================================================
	// fmt.Errorf("context: %w", err) wraps errors while preserving the original.
	// errors.Is() checks if any layer matches. errors.As() extracts typed errors.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 3: ERROR WRAPPING")
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Println("1. Wrapping errors with fmt.Errorf and %w:")

	err = findUser(123)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	}

	fmt.Println()
	fmt.Println("2. errors.Is - checking wrapped errors:")

	err = findUser(123)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("   errors.Is found ErrNotFound in the chain")
		}
	}

	err = findUser(123)
	if err != nil {
		if err == ErrNotFound {
			fmt.Println("   This will NOT print - == does not unwrap")
		} else {
			fmt.Println("   err == ErrNotFound is false (== does not unwrap)")
		}
	}

	fmt.Println()
	fmt.Println("3. errors.As - extracting specific error types:")

	err = processPayment("card_4242", 100)
	if err != nil {
		var pErr *PaymentError
		if errors.As(err, &pErr) {
			fmt.Printf("   Extracted PaymentError: Card=%s, Amount=%.2f\n",
				pErr.CardNumber, pErr.Amount)
			fmt.Printf("   Full error: %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("4. errors.Unwrap - manual unwrapping:")

	err = findUser(123)
	fmt.Printf("   Full wrapped error:    %v\n", err)

	unwrapped := errors.Unwrap(err)
	fmt.Printf("   After first unwrap:    %v\n", unwrapped)

	unwrapped = errors.Unwrap(unwrapped)
	fmt.Printf("   After second unwrap:   %v\n", unwrapped)

	unwrapped = errors.Unwrap(unwrapped)
	fmt.Printf("   After third unwrap:    %v\n", unwrapped)

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Use %w in fmt.Errorf to wrap errors (not %v or %s)")
	fmt.Println("  2. Use errors.Is() instead of == for error comparison")
	fmt.Println("  3. Use errors.As() to extract typed errors from chains")
	fmt.Println("  4. Wrap at every boundary: repo -> service -> handler")

	// ========================================================
	// SECTION 4: SENTINEL ERRORS
	// ========================================================
	// Sentinel errors are predefined package-level error variables.
	// Convention: Err prefix (e.g., ErrNotFound). Check with errors.Is().

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 4: SENTINEL ERRORS")
	fmt.Println("============================================================")
	fmt.Println()

	svc := NewUserService()

	fmt.Println("1. Checking sentinel errors with errors.Is:")

	_, err = svc.GetUser(999)
	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("   User 999 not found (matched ErrUserNotFound)")
	}

	_, err = svc.GetUser(-1)
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("   User -1 rejected (matched ErrInvalidInput)")
	}

	fmt.Println()
	fmt.Println("2. Sentinel errors work through wrapping:")

	_, err = svc.GetUser(999)
	wrappedErr := fmt.Errorf("HTTP handler: %w", err)
	if errors.Is(wrappedErr, ErrUserNotFound) {
		fmt.Println("   errors.Is found ErrNotFound through the wrapper")
	}

	fmt.Println()
	fmt.Println("3. Successful operations:")

	user, err := svc.GetUser(1)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found user: %s (%s)\n", user.Name, user.Email)
	}

	newUser, err := svc.CreateUser("Bob", "bob@example.com")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Created user: %s (ID: %d)\n", newUser.Name, newUser.ID)
	}

	fmt.Println()
	fmt.Println("4. Decision tree based on sentinel errors:")

	testCases := []struct {
		name  string
		email string
	}{
		{"Carol", "carol@example.com"},
		{"", "empty@example.com"},
		{"Dave", "bob@example.com"},
	}

	for _, tc := range testCases {
		_, err := svc.CreateUser(tc.name, tc.email)
		switch {
		case errors.Is(err, ErrInvalidInput):
			fmt.Printf("   '%s': rejected - missing required fields\n", tc.name)
		case errors.Is(err, ErrUserExists):
			fmt.Printf("   '%s': rejected - email already registered\n", tc.name)
		case err != nil:
			fmt.Printf("   '%s': unexpected error - %v\n", tc.name, err)
		default:
			fmt.Printf("   '%s': created successfully\n", tc.name)
		}
	}

	fmt.Println()
	fmt.Println("Sentinel vs Custom Type decision guide:")
	fmt.Println("  Use SENTINEL ERRORS when:")
	fmt.Println("    - Error is a simple, well-known condition")
	fmt.Println("    - No extra data needed")
	fmt.Println("  Use CUSTOM ERROR TYPES when:")
	fmt.Println("    - You need to attach data (IDs, fields, codes)")
	fmt.Println("    - Multiple variants of the same error category")

	// ========================================================
	// SECTION 5: PANIC AND RECOVER
	// ========================================================
	// panic stops execution. recover (only inside defer) catches it.
	// Never use panic for normal errors. Only for truly unrecoverable situations.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 5: PANIC AND RECOVER")
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Println("1. Panic behavior (caught by recover):")

	result2 := runWithRecover(func() string {
		fmt.Println("   Starting operation...")
		fmt.Println("   About to panic!")
		panic("something went terribly wrong")
	})
	fmt.Printf("   Result after recover: %v\n", result2)

	fmt.Println()
	fmt.Println("2. Appropriate uses of panic:")

	fmt.Println("   a) init() failures (program cannot start)")
	demonstrateInitPanic()

	fmt.Println("   b) Programmer errors (violated invariants)")
	demonstrateInvariantPanic()

	fmt.Println("   c) Development assertions")
	demonstrateAssertion()

	fmt.Println()
	fmt.Println("3. Defer + Recover pattern:")

	safeDivide(10, 2)
	safeDivide(10, 0)
	safeDivide(20, 4)

	fmt.Println()
	fmt.Println("4. Simulated HTTP handler recovery:")

	for _, path := range []string{"/api/ok", "/api/crash", "/api/also-ok"} {
		res := simulateHTTPHandler(path)
		fmt.Printf("   GET %s -> %s\n", path, res)
	}

	fmt.Println()
	fmt.Println("Anti-patterns (what NOT to do):")
	fmt.Println("  BAD: Using panic for validation errors")
	fmt.Println("  BAD: Using panic for expected failures (file not found)")
	fmt.Println("  BAD: Using panic for user input errors")
	fmt.Println("  GOOD: Use (value, error) returns for all expected failures")

	// ========================================================
	// SECTION 6: DEFER USAGE
	// ========================================================
	// defer schedules a function to run when the surrounding function returns.
	// Multiple defers execute in LIFO order. Arguments are evaluated at defer time.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 6: DEFER USAGE")
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Println("1. Basic defer:")
	fmt.Println("   Start of main")
	defer fmt.Println("   End of main (deferred - runs last)")
	fmt.Println("   Middle of main")

	fmt.Println()
	fmt.Println("2. LIFO order (last-in, first-out):")
	demonstrateLIFO()

	fmt.Println()
	fmt.Println("3. File cleanup with defer:")
	_ = readFileWithDefer()

	fmt.Println()
	fmt.Println("4. Arguments evaluated at defer registration:")
	demonstrateArgumentEvaluation()

	fmt.Println()
	fmt.Println("5. Timing function execution:")
	slowOperation()
	fastOperation()

	fmt.Println()
	fmt.Println("6. Deferred method calls:")
	demonstrateMethodDefer()

	// ========================================================
	// SECTION 7: MULTIPLE ERRORS
	// ========================================================
	// Check errors immediately after each operation. Return early on first error.
	// Use error collectors when you need ALL failures (e.g., form validation).

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 7: MULTIPLE ERRORS")
	fmt.Println("============================================================")
	fmt.Println()

	fmt.Println("1. Early return pattern (multi-step workflow):")
	orderResult, err := processOrder("ORD-001", 100)
	if err != nil {
		fmt.Printf("   Order failed: %v\n", err)
	} else {
		fmt.Printf("   Order succeeded: %s\n", orderResult)
	}

	fmt.Println()
	fmt.Println("2. Errors at different steps:")

	_, err = processOrder("", 50)
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	_, err = processOrder("ORD-002", 0)
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	_, err = processOrder("ORD-003", 100)
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	_, err = processOrder("ORD-004", 1)
	if err != nil {
		fmt.Printf("   %v\n", err)
	}

	fmt.Println()
	fmt.Println("3. Collecting multiple errors (form validation):")

	validateAndReport("  ", "not-an-email", -5)
	fmt.Println()
	validateAndReport("Alice", "alice@example.com", 25)

	fmt.Println()
	fmt.Println("4. Multi-step with rollback on failure:")

	err = provisionServer("web-01")
	if err != nil {
		fmt.Printf("   Provisioning failed: %v\n", err)
	}

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Check every error immediately, return early")
	fmt.Println("  2. Happy path stays at base indentation (no else blocks)")
	fmt.Println("  3. Use error collectors when you need ALL failures")
	fmt.Println("  4. Plan rollback/cleanup for multi-step operations")

	// ========================================================
	// SECTION 8: REAL-WORLD FILE I/O
	// ========================================================
	// File operations are a common source of errors. Always:
	//   - Check errors on Open/Create
	//   - defer Close immediately after
	//   - Check errors on Read/Write
	//   - Handle partial reads/writes gracefully

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 8: REAL-WORLD FILE I/O")
	fmt.Println("============================================================")
	fmt.Println()

	// Pattern: Write to file with proper error handling
	err = writeToFile("example.txt", []byte("Hello, Go Error Handling!\nThis is a test file.\n"))
	if err != nil {
		fmt.Printf("   Write failed: %v\n", err)
	} else {
		fmt.Println("   File written successfully")
	}

	// Pattern: Read from file with proper error handling
	data, err := readFromFile("example.txt")
	if err != nil {
		fmt.Printf("   Read failed: %v\n", err)
	} else {
		fmt.Printf("   Read %d bytes from file\n", len(data))
		fmt.Printf("   Content:\n   %s", string(data))
	}

	// Pattern: Handle missing file
	_, err = readFromFile("nonexistent.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("   File does not exist (handled gracefully)")
		} else {
			fmt.Printf("   Read failed: %v\n", err)
		}
	}

	// Pattern: Copy file with error handling at each step
	err = copyFile("example.txt", "example-copy.txt")
	if err != nil {
		fmt.Printf("   Copy failed: %v\n", err)
	} else {
		fmt.Println("   File copied successfully")
	}

	// Cleanup demo files
	os.Remove("example.txt")
	os.Remove("example-copy.txt")

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Always check errors on Open/Create/Read/Write")
	fmt.Println("  2. defer f.Close() right after checking open error")
	fmt.Println("  3. Use errors.Is(err, os.ErrNotExist) for missing files")
	fmt.Println("  4. Handle partial reads/writes in production code")
	fmt.Println("  5. Clean up temp files even on error")

	// ========================================================
	// SECTION 9: HTTP ERROR HANDLING
	// ========================================================
	// HTTP handlers must return proper status codes and handle failures.
	// Use middleware to recover from panics so one bad request doesn't crash the server.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 9: HTTP ERROR HANDLING")
	fmt.Println("============================================================")
	fmt.Println()

	// Simulate HTTP handler responses
	fmt.Println("Simulated HTTP handler responses:")

	handler := httpSimulator{}
	handler.handleGetUser(1)
	handler.handleGetUser(999)
	handler.handleGetUser(-1)
	handler.handleCreateUser("", "test@test.com")
	handler.handleCreateUser("Test", "bad-email")
	handler.handleCrashEndpoint()

	fmt.Println()
	fmt.Println("Key takeaways:")
	fmt.Println("  1. Return proper HTTP status codes (400, 404, 500)")
	fmt.Println("  2. Use errors.Is to map errors to status codes")
	fmt.Println("  3. Recover from panics in middleware")
	fmt.Println("  4. Log errors server-side, return generic messages to clients")
	fmt.Println("  5. Never expose internal error details to clients")

	// ========================================================
	// SECTION 10: COMMON MISTAKES
	// ========================================================
	// The most frequent error handling mistakes Go developers make and how to fix them.

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SECTION 10: COMMON MISTAKES")
	fmt.Println("============================================================")
	fmt.Println()

	// Mistake 1: Ignoring errors
	fmt.Println("1. Ignoring errors:")
	fmt.Println("   BAD:  os.Remove(\"file.txt\")         // error ignored!")
	fmt.Println("   GOOD: if err := os.Remove(\"file.txt\"); err != nil { ... }")
	demonstrateIgnoredError()

	fmt.Println()

	// Mistake 2: Overusing panic
	fmt.Println("2. Overusing panic:")
	fmt.Println("   BAD:  panic(\"user not found\")       // should return error")
	fmt.Println("   GOOD: return nil, errors.New(\"user not found\")")

	fmt.Println()

	// Mistake 3: Losing error context
	fmt.Println("3. Losing error context:")
	fmt.Println("   BAD:  return nil, errors.New(\"something went wrong\")")
	fmt.Println("   GOOD: return nil, fmt.Errorf(\"processing user %d: %w\", id, err)")
	demonstrateLostContext()

	fmt.Println()

	// Mistake 4: Nil error pitfalls
	fmt.Println("4. Nil error pitfalls:")
	fmt.Println("   BAD:  returning a nil pointer to error interface")
	fmt.Println("   GOOD: return nil, nil (or just nil for single return)")
	demonstrateNilErrorPitfall()

	fmt.Println()

	// Mistake 5: String comparison instead of errors.Is
	fmt.Println("5. String comparison instead of errors.Is:")
	fmt.Println("   BAD:  if err.Error() == \"not found\" { ... }")
	fmt.Println("   GOOD: if errors.Is(err, ErrNotFound) { ... }")

	// ========================================================
	// SUMMARY AND PRACTICE EXERCISES
	// ========================================================

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("SUMMARY: KEY LESSONS")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("  1. ERRORS ARE VALUES - Return them, check them, don't ignore them")
	fmt.Println("  2. if err != nil IS THE FOUNDATION - Check immediately after every call")
	fmt.Println("  3. EARLY RETURN PATTERN - Handle error and return, keep happy path flat")
	fmt.Println("  4. CUSTOM ERROR TYPES - Use structs with Error() for rich error data")
	fmt.Println("  5. ERROR WRAPPING - Use %w to preserve error chains across layers")
	fmt.Println("  6. errors.Is / errors.As - Check and extract from wrapped error chains")
	fmt.Println("  7. SENTINEL ERRORS - Use Err prefix for simple, well-known conditions")
	fmt.Println("  8. DEFER FOR CLEANUP - Always defer Close() right after successful Open()")
	fmt.Println("  9. PANIC RARELY - Only for truly unrecoverable situations or bugs")
	fmt.Println(" 10. NEVER IGNORE ERRORS - Even _ = someFunc() should be deliberate")
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("PRACTICE EXERCISES")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("  Exercise 1: Write a function that reads a JSON file and parses it")
	fmt.Println("    into a struct. Handle all possible errors (file not found,")
	fmt.Println("    permission denied, invalid JSON) with appropriate error types.")
	fmt.Println()
	fmt.Println("  Exercise 2: Create a custom MultiError type that implements the")
	fmt.Println("    error interface and can hold multiple errors. Add a method")
	fmt.Println("    that returns all individual errors for inspection.")
	fmt.Println()
	fmt.Println("  Exercise 3: Build a retry mechanism that attempts an operation")
	fmt.Println("    up to N times. It should only retry on specific errors (e.g.,")
	fmt.Println("    network timeouts) and give up immediately on others.")
	fmt.Println()
	fmt.Println("  Exercise 4: Write an HTTP middleware that logs all errors,")
	fmt.Println("    recovers from panics, and maps custom error types to the")
	fmt.Println("    correct HTTP status codes (400, 401, 404, 500).")
	fmt.Println()
	fmt.Println("  Exercise 5: Create a database transaction helper that uses")
	fmt.Println("    defer to automatically rollback on error and commit on")
	fmt.Println("    success. Handle the case where rollback itself fails.")
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("END OF GO ERROR HANDLING GUIDE")
	fmt.Println("============================================================")
}

// ============================================================
// SECTION 8: FILE I/O HELPER FUNCTIONS
// ============================================================

func writeToFile(path string, data []byte) error {
	// Open file with create/truncate flags, read/write permissions
	f, err := os.Create(path)
	if err != nil {
		// Wrap with context about what operation failed
		return fmt.Errorf("creating file %s: %w", path, err)
	}
	// Defer close immediately - ensures cleanup on all paths
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		// Error during write - defer will still close the file
		return fmt.Errorf("writing to file %s: %w", path, err)
	}

	// Sync ensures data is flushed to disk
	if err = f.Sync(); err != nil {
		return fmt.Errorf("syncing file %s: %w", path, err)
	}

	return nil
}

func readFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}
	defer f.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	return data, nil
}

func copyFile(src, dst string) error {
	// Open source
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading source %s: %w", src, err)
	}

	// Create destination
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating destination %s: %w", dst, err)
	}
	defer f.Close()

	// Write data
	_, err = f.Write(data)
	if err != nil {
		// Try to clean up the incomplete destination file
		os.Remove(dst)
		return fmt.Errorf("writing to destination %s: %w", dst, err)
	}

	return nil
}

// ============================================================
// SECTION 9: HTTP SIMULATOR
// ============================================================

type httpSimulator struct {
	svc *UserService
}

func (h *httpSimulator) init() {
	h.svc = NewUserService()
}

func (h *httpSimulator) handleGetUser(id int) {
	if h.svc == nil {
		h.init()
	}
	_, err := h.svc.GetUser(id)
	status, msg := mapErrorToHTTPResponse(err)
	fmt.Printf("   GET /users/%d -> %d %s\n", id, status, msg)
}

func (h *httpSimulator) handleCreateUser(name, email string) {
	if h.svc == nil {
		h.init()
	}
	_, err := h.svc.CreateUser(name, email)
	status, msg := mapErrorToHTTPResponse(err)
	fmt.Printf("   POST /users (name=%q) -> %d %s\n", name, status, msg)
}

func (h *httpSimulator) handleCrashEndpoint() {
	// Simulate panic recovery in HTTP middleware
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("   GET /crash -> 500 Internal Server Error (recovered from panic: %v)\n", r)
			}
		}()
		panic("unexpected nil pointer")
	}()
}

func mapErrorToHTTPResponse(err error) (int, string) {
	if err == nil {
		return 200, "OK"
	}
	switch {
	case errors.Is(err, ErrUserNotFound):
		return 404, "User not found"
	case errors.Is(err, ErrInvalidInput):
		return 400, "Invalid input"
	case errors.Is(err, ErrUserExists):
		return 409, "User already exists"
	default:
		// Never expose internal error details to the client
		// Log the real error server-side
		return 500, "Internal server error"
	}
}

// ============================================================
// SECTION 10: COMMON MISTAKES DEMONSTRATIONS
// ============================================================

// demonstrateIgnoredError shows why ignoring errors is dangerous.
func demonstrateIgnoredError() {
	// BAD example (commented out because it would silently fail):
	// os.Remove("nonexistent-file.txt") // error ignored!

	// GOOD example:
	err := os.Remove("nonexistent-file.txt")
	if err != nil {
		fmt.Printf("   Caught the error we would have missed: %v\n", err)
	}
}

// demonstrateLostContext shows how to preserve error context.
func demonstrateLostContext() {
	// BAD: losing the original error
	badErr := doRiskyOperation()
	if badErr != nil {
		// BAD: original error details are lost
		_ = fmt.Errorf("operation failed")
	}

	// GOOD: preserving the original error with context
	goodErr := doRiskyOperation()
	if goodErr != nil {
		// GOOD: original error is wrapped and preserved
		wrapped := fmt.Errorf("risky operation failed: %w", goodErr)
		fmt.Printf("   Wrapped error preserves context: %v\n", wrapped)
	}
}

func doRiskyOperation() error {
	return &DatabaseError{
		Code:      500,
		Operation: "QUERY",
		Table:     "orders",
		Err:       fmt.Errorf("connection timeout after 30s"),
	}
}

// demonstrateNilErrorPitfall shows the nil interface gotcha.
// This is a subtle but common Go mistake.
func demonstrateNilErrorPitfall() {
	// Scenario: A function returns a pointer to a custom error type.
	// When there's no error, returning nil directly is correct.
	// But returning a nil pointer as an error interface is NOT nil!

	_, err := mightFail(true)
	fmt.Printf("   mightFail(true)  -> err == nil: %v\n", err == nil)

	_, err = mightFail(false)
	fmt.Printf("   mightFail(false) -> err == nil: %v (surprising!)\n", err == nil)
	fmt.Println("   The error interface holds a nil *MyError, which is NOT a nil interface!")
	fmt.Println("   Fix: return nil directly, not a nil pointer typed as error")
}

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

// mightFail demonstrates the nil interface pitfall.
// When success=true, it returns nil, nil (correct).
// When success=false, it returns a nil *MyError as error (BUG).
func mightFail(success bool) (result string, err error) {
	if success {
		return "ok", nil
	}
	// BUG: returning a nil pointer as error interface
	// The interface value is (type=*MyError, value=nil), which != nil
	var myErr *MyError
	return "", myErr // myErr is nil, but error interface is NOT nil!
}
