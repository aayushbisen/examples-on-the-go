// Package main provides a comprehensive guide to interfaces in Go.
// This file covers everything from basic concepts to advanced patterns,
// with runnable examples and detailed explanations.
//
// ============================================================================
// WHAT IS AN INTERFACE?
// ============================================================================
// An interface in Go is a CONTRACT that defines WHAT a type can do, not WHAT it is.
// It's a collection of method signatures. If a type has those methods, it
// satisfies the interface.
//
// Think of it like a USB port: the port doesn't care what device you plug in
// (mouse, keyboard, drive), as long as it speaks "USB."
//
// ============================================================================
// HOW GO INTERFACES DIFFER FROM JAVA/C++
// ============================================================================
// | Feature          | Java / C++                          | Go                          |
// |------------------|-------------------------------------|-----------------------------|
// | Declaration      | class Dog implements Animal         | No keyword — automatic      |
// | Implementation   | Explicit (implements keyword)       | IMPLICIT (if methods match) |
// | Inheritance      | Classes inherit from other classes  | No inheritance, composition |
// | Purpose          | Define type hierarchies             | Define BEHAVIOR             |
//
// In Java, you must DECLARE that a class implements an interface.
// In Go, if your type has the methods, it AUTOMATICALLY satisfies the interface.
// This is called IMPLICIT IMPLEMENTATION.
//
// ============================================================================
// IMPLICIT IMPLEMENTATION — "Duck Typing" of Go
// ============================================================================
// "If it walks like a duck and quacks like a duck, it's a duck."
//
// You never write `implements`. If your struct has the methods the interface
// requires, Go accepts it — period.
//
// type Speaker interface { Speak() string }
// type Dog struct{}
// func (d Dog) Speak() string { return "Woof!" }
// // Dog now implements Speaker automatically — no declaration needed!
//
// ============================================================================
// HOW TO RUN
// ============================================================================
// go run examples/interfaces/interfaces.go
package main

import (
	"fmt"
	"math"
)

// ============================================================================
// SECTION 1: Basic Interface Example — Shapes
// ============================================================================

// Shape is an interface that defines what a shape can do.
// Any type with Area() and Perimeter() methods satisfies this interface.
//
// Note: Interfaces define BEHAVIOR, not data. Shape doesn't store radius or
// width — it just says "anything that is a Shape must be able to calculate
// its area and perimeter."
type Shape interface {
	Area() float64      // Returns the area of the shape
	Perimeter() float64 // Returns the perimeter of the shape
}

// Circle is a struct representing a circle.
// It doesn't know about the Shape interface — and it doesn't need to.
// Go will automatically recognize that Circle satisfies Shape.
type Circle struct {
	Radius float64
}

// Area calculates the circle's area.
// By defining this method, Circle AUTOMATICALLY satisfies the Shape interface.
// We use a VALUE receiver (c Circle) because Circle is small and we don't
// modify it. This also means copies of Circle are safe to use.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circle's circumference (2 * π * r).
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle is another shape.
type Rectangle struct {
	Width, Height float64
}

// Area calculates the rectangle's area.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the rectangle's perimeter.
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// printShape accepts ANY type that satisfies the Shape interface.
// This is polymorphism — one function, many types.
//
// The parameter "s Shape" means: "I don't care what concrete type this is,
// as long as it has Area() and Perimeter() methods."
//
// Under the hood, Go stores two things in the interface variable:
//   1. The concrete type (Circle or Rectangle)
//   2. The actual value
// This is called an INTERFACE VALUE: (type, value)
func printShape(s Shape) {
	fmt.Printf("  Area: %.2f | Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func basicShapesExample() {
	fmt.Println("=== Basic Shapes Example ===")

	// Create concrete types.
	circle := Circle{Radius: 5}
	rect := Rectangle{Width: 4, Height: 6}

	// Pass them to printShape. Go automatically treats them as Shape.
	fmt.Println("Circle:")
	printShape(circle) // Circle satisfies Shape implicitly

	fmt.Println("Rectangle:")
	printShape(rect) // Rectangle satisfies Shape implicitly

	// You can also store them in a slice of interfaces.
	// []Shape means "a slice of anything that satisfies Shape."
	shapes := []Shape{circle, rect}
	fmt.Println("\nAll shapes:")
	for i, s := range shapes {
		fmt.Printf("  Shape %d: Area = %.2f\n", i+1, s.Area())
	}
	fmt.Println()
}

// ============================================================================
// SECTION 2: Empty Interface (`any`)
// ============================================================================
//
// The empty interface `interface{}` (or its alias `any` since Go 1.18)
// has ZERO methods. This means EVERY type satisfies it.
//
// func printAnything(v any) {
//     fmt.Println(v) // Accepts int, string, struct, slice — anything.
// }
//
// This is Go's equivalent to `Object` in Java or `void*` in C.
// Use it when you truly don't know the type (e.g., JSON parsing).
// But avoid it when possible — it sacrifices type safety.

func emptyInterfaceExample() {
	fmt.Println("=== Empty Interface (`any`) Example ===")

	// `any` is an alias for `interface{}` — they are identical.
	// Since empty interface has zero methods, EVERY type satisfies it.
	values := []any{
		42,            // int
		"hello",       // string
		3.14,          // float64
		true,          // bool
		[]int{1, 2, 3}, // slice
		Circle{5},     // struct
	}

	fmt.Println("Values stored in []any:")
	for i, v := range values {
		// fmt.Printf("%v prints the value. We don't know the type at compile time.
		fmt.Printf("  [%d] %v\n", i, v)
	}
	fmt.Println()
}

// ============================================================================
// SECTION 3: Type Assertions
// ============================================================================
//
// When you have a value of type `any` (or any interface), you can "assert"
// that it's a specific concrete type.
//
// v.(string) — "I believe v is a string. Give me the string."
//
// There are two forms:
//   1. s := v.(string)          — Panics if v is not a string!
//   2. s, ok := v.(string)     — Safe: ok is true if assertion succeeded.
//
// ALWAYS prefer form #2 unless you're certain of the type.

func typeAssertionExample() {
	fmt.Println("=== Type Assertion Example ===")

	var v any = "hello world" // Store a string in an empty interface.

	// SAFE form: returns (value, bool). If v isn't a string, ok is false.
	s, ok := v.(string)
	if ok {
		fmt.Printf("  It's a string: %q (length: %d)\n", s, len(s))
	} else {
		fmt.Println("  Not a string")
	}

	// Try asserting a different type.
	n, ok := v.(int)
	if ok {
		fmt.Printf("  It's an int: %d\n", n)
	} else {
		fmt.Println("  Not an int (this is correct — v is a string)")
	}

	// UNSAFE form (would panic):
	// n := v.(int) // PANIC: interface conversion: any is string, not int
	// NEVER use this form unless you are 100% sure of the type.

	// Practical example: extracting values from a map[string]any
	// (common in JSON parsing or config files).
	config := map[string]any{
		"host": "localhost",
		"port": 8080,
		"debug": true,
	}

	if host, ok := config["host"].(string); ok {
		fmt.Printf("  Config host: %s\n", host)
	}
	if port, ok := config["port"].(int); ok {
		fmt.Printf("  Config port: %d\n", port)
	}
	fmt.Println()
}

// ============================================================================
// SECTION 4: Type Switches
// ============================================================================
//
// A type switch is the clean way to handle multiple possible types.
// It's like a switch statement, but it checks the TYPE of a value, not its value.
//
// switch t := v.(type) {
// case string:
//     fmt.Println("It's a string:", t) // t is of type string here
// case int:
//     fmt.Println("It's an int:", t)   // t is of type int here
// }

func typeSwitchExample() {
	fmt.Println("=== Type Switch Example ===")

	testValues := []any{42, "hello", 3.14, true, Circle{5}}

	for _, v := range testValues {
		// v.(type) is special syntax only allowed inside a switch.
		// It extracts the concrete type and assigns it to `t`.
		switch t := v.(type) {
		case string:
			// Inside this block, `t` is of type string.
			fmt.Printf("  String: %q (length: %d)\n", t, len(t))
		case int:
			// Inside this block, `t` is of type int.
			fmt.Printf("  Integer: %d (is even: %v)\n", t, t%2 == 0)
		case float64:
			fmt.Printf("  Float: %.2f\n", t)
		case bool:
			fmt.Printf("  Boolean: %v\n", t)
		case Circle:
			// We can even match our own types!
			fmt.Printf("  Circle: radius=%.1f, area=%.2f\n", t.Radius, t.Area())
		default:
			fmt.Printf("  Unknown type: %T\n", v) // %T prints the type name
		}
	}
	fmt.Println()
}

// ============================================================================
// SECTION 5: Real-World Use Case — Payment System
// ============================================================================
//
// This demonstrates how interfaces make your code EXTENSIBLE.
// The Checkout function doesn't know about Stripe or PayPal — it only
// knows about the PaymentGateway interface. Adding a new payment provider
// requires ZERO changes to Checkout.

// PaymentGateway defines the contract for any payment provider.
// Stripe, PayPal, Razorpay — as long as they implement these methods,
// our checkout system works with ALL of them.
//
// This is the "Dependency Inversion Principle": high-level modules (Checkout)
// depend on abstractions (PaymentGateway), not concrete types (Stripe).
type PaymentGateway interface {
	Charge(amount float64) (string, error) // Process a payment, return transaction ID
	Refund(transactionID string) error     // Refund a payment
	Name() string                          // Return the provider name
}

// Stripe is a concrete implementation of PaymentGateway.
// Notice: Stripe has NO idea it implements PaymentGateway.
// It just has the methods. Go does the rest.
type Stripe struct {
	APIKey string
}

// Charge simulates calling Stripe's API.
// We return (transactionID, error) — the standard Go error pattern.
func (s Stripe) Charge(amount float64) (string, error) {
	// In real code, this would make an HTTP call to Stripe's API.
	fmt.Printf("    [Stripe API] Charging $%.2f...\n", amount)
	return "stripe_txn_123", nil
}

func (s Stripe) Refund(transactionID string) error {
	fmt.Printf("    [Stripe API] Refunding %s...\n", transactionID)
	return nil
}

func (s Stripe) Name() string {
	return "Stripe"
}

// PayPal is another implementation of PaymentGateway.
type PayPal struct {
	ClientID string
}

func (p PayPal) Charge(amount float64) (string, error) {
	fmt.Printf("    [PayPal API] Charging $%.2f...\n", amount)
	return "paypal_txn_456", nil
}

func (p PayPal) Refund(transactionID string) error {
	fmt.Printf("    [PayPal API] Refunding %s...\n", transactionID)
	return nil
}

func (p PayPal) Name() string {
	return "PayPal"
}

// Checkout is our core business logic.
// It accepts ANY PaymentGateway — we don't care which one.
// This makes the system EXTENSIBLE: add a new provider, zero changes here.
//
// This is the #1 benefit of interfaces: write code that works with ANY
// type that satisfies the contract, not just specific types.
func Checkout(gateway PaymentGateway, amount float64) {
	fmt.Printf("  Processing $%.2f via %s...\n", amount, gateway.Name())

	txnID, err := gateway.Charge(amount)
	if err != nil {
		fmt.Println("    Payment failed!")
		return
	}

	fmt.Printf("    Success! Transaction ID: %s\n", txnID)
}

func paymentSystemExample() {
	fmt.Println("=== Payment System Example ===")

	// Create different payment gateways.
	stripe := Stripe{APIKey: "sk_test_..."}
	paypal := PayPal{ClientID: "client_..."}

	// The SAME Checkout function works with ANY gateway.
	Checkout(stripe, 99.99)
	fmt.Println()
	Checkout(paypal, 49.50)
	fmt.Println()

	// Store gateways in a slice — polymorphism in action.
	gateways := []PaymentGateway{stripe, paypal}
	fmt.Println("  All registered gateways:")
	for _, g := range gateways {
		fmt.Printf("    - %s\n", g.Name())
	}
	fmt.Println()
}

// ============================================================================
// SECTION 6: Interface Composition
// ============================================================================
//
// Go interfaces can be composed by embedding other interfaces.
// This is how the standard library builds powerful interfaces from tiny ones.
//
// Example from Go's standard library:
//   type ReadWriter interface {
//       Reader  // embeds Read()
//       Writer  // embeds Write()
//   }
//
// This means a ReadWriter must have BOTH Read() and Write() methods.

// Reader requires a Read method.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer requires a Write method.
type Writer interface {
	Write(p []byte) (n int, err error)
}

// ReadWriter is a COMPOSED interface. It embeds both Reader and Writer.
// Any type that satisfies ReadWriter must have BOTH Read() and Write().
//
// This is equivalent to:
// type ReadWriter interface {
//     Read(p []byte) (n int, err error)
//     Write(p []byte) (n int, err error)
// }
type ReadWriter interface {
	Reader
	Writer
}

// fileHandler demonstrates a type that satisfies the composed interface.
type fileHandler struct {
	name string
}

func (f fileHandler) Read(p []byte) (int, error) {
	fmt.Printf("    Reading from %s\n", f.name)
	return 0, nil
}

func (f fileHandler) Write(p []byte) (int, error) {
	fmt.Printf("    Writing to %s\n", f.name)
	return len(p), nil
}

func interfaceCompositionExample() {
	fmt.Println("=== Interface Composition Example ===")

	fh := fileHandler{name: "data.txt"}

	// fh satisfies Reader.
	var r Reader = fh
	r.Read(nil)

	// fh satisfies Writer.
	var w Writer = fh
	w.Write([]byte("hello"))

	// fh satisfies ReadWriter (the composed interface).
	var rw ReadWriter = fh
	rw.Read(nil)
	rw.Write([]byte("world"))

	fmt.Println()
}

// ============================================================================
// SECTION 7: Best Practices & Common Mistakes
// ============================================================================

// CustomError demonstrates the nil interface gotcha.
// This is a concrete error type with a pointer receiver.
type CustomError struct{}

func (e *CustomError) Error() string { return "custom error" }

// demonstrateNilInterfaceGotcha shows the most confusing aspect of Go interfaces.
//
// An interface value is nil ONLY when both its TYPE and VALUE are nil.
// If you store a nil pointer in an interface, the interface is NOT nil.
//
// var err error = nil          // nil interface — safe
// var e *CustomError = nil
// var err2 error = e           // NOT nil! (type=*CustomError, value=nil)
//
// This is a VERY common source of bugs. Always check the concrete type
// before checking for nil in interfaces.
func demonstrateNilInterfaceGotcha() {
	fmt.Println("=== Nil Interface Gotcha ===")

	// This is a nil interface — both type and value are nil.
	var nilErr error = nil
	fmt.Printf("  nilErr == nil: %v\n", nilErr == nil)

	// Create a nil pointer to our error type.
	var nilCustom *CustomError = nil

	// Assign it to an error interface. The interface now holds:
	//   type = *CustomError
	//   value = nil
	// This is NOT a nil interface!
	var err error = nilCustom
	fmt.Printf("  err == nil: %v\n", err == nil) // FALSE!
	fmt.Printf("  err's type: %T\n", err)        // *CustomError
	fmt.Printf("  err's value: %v\n", err)       // <nil>

	// The fix: always check the concrete value, not the interface.
	if nilCustom == nil {
		fmt.Println("  Correct way: check the concrete pointer, not the interface")
	}
	fmt.Println()
}

// ============================================================================
// SECTION 8: Interface Values Under the Hood
// ============================================================================
//
// Every interface value is a pair: (type, value).
//
// When you do:
//   var s Shape = Circle{5}
//
// The interface `s` internally stores:
//   type  = Circle
//   value = Circle{Radius: 5}
//
// When you call s.Area(), Go:
//   1. Looks up Circle's method table
//   2. Finds the Area() method
//   3. Calls it with the stored value
//
// This is why interface calls are slightly slower than direct calls —
// there's a small lookup cost. In 99% of cases, this doesn't matter.

func interfaceInternalsExample() {
	fmt.Println("=== Interface Values Under the Hood ===")

	var s Shape = Circle{Radius: 5}

	// %T prints the concrete TYPE stored in the interface.
	fmt.Printf("  Interface holds type: %T\n", s)

	// %v prints the VALUE stored in the interface.
	fmt.Printf("  Interface holds value: %+v\n", s)

	// When s is reassigned, the internal (type, value) pair changes.
	s = Rectangle{Width: 3, Height: 4}
	fmt.Printf("  After reassignment, type: %T\n", s)
	fmt.Printf("  After reassignment, value: %+v\n", s)
	fmt.Println()
}

// ============================================================================
// SECTION 9: When to Use Interfaces (and When Not To)
// ============================================================================
//
// RULE OF THUMB: Don't create an interface until you have at least TWO
// implementations OR you need to break a circular dependency.
//
// Good reasons to use interfaces:
//   - You have multiple implementations (Stripe, PayPal, Cash)
//   - You want to mock dependencies for testing
//   - You need to break import cycles (package A imports B, B imports A)
//   - The standard library expects it (io.Reader, http.Handler)
//
// Bad reasons to use interfaces:
//   - "It's good OOP practice" — Go is not Java
//   - "I might need it later" — YAGNI (You Aren't Gonna Need It)
//   - To avoid passing concrete types — interfaces add indirection cost

func whenToUseInterfaces() {
	fmt.Println("=== When to Use Interfaces ===")
	fmt.Println("  DO:")
	fmt.Println("    - Keep interfaces small (1-3 methods)")
	fmt.Println("    - Accept interfaces, return structs")
	fmt.Println("    - Define interfaces where they're USED, not where implemented")
	fmt.Println("    - Use pointer receivers when methods mutate state")
	fmt.Println()
	fmt.Println("  DON'T:")
	fmt.Println("    - Create large interfaces")
	fmt.Println("    - Use `any` unnecessarily")
	fmt.Println("    - Over-engineer with a single implementation")
	fmt.Println("    - Export interfaces with mutating methods without thought")
	fmt.Println()
}

// ============================================================================
// main runs all interface examples.
// ============================================================================
func main() {
	basicShapesExample()
	emptyInterfaceExample()
	typeAssertionExample()
	typeSwitchExample()
	paymentSystemExample()
	interfaceCompositionExample()
	demonstrateNilInterfaceGotcha()
	interfaceInternalsExample()
	whenToUseInterfaces()

	// ============================================================================
	// SUMMARY
	// ============================================================================
	// | Concept              | Explanation                                       |
	// |----------------------|---------------------------------------------------|
	// | Interface            | A set of method signatures — a contract for behavior |
	// | Implicit impl.       | No `implements` keyword — if methods match, done  |
	// | any / interface{}    | Accepts any type — use sparingly                  |
	// | Type assertion       | Extract concrete type from an interface           |
	// | Type switch          | Cleanly handle multiple possible types            |
	// | Small interfaces     | Many small > one large                            |
	// | nil interface gotcha | interface is nil only when type AND value are nil |
	//
	// ============================================================================
	// PRACTICE QUESTIONS
	// ============================================================================
	// 1. Create a `Notifier` interface with `Send(message string) error`.
	//    Implement it for `Email` and `SMS` structs. Write a `Notify`
	//    function that accepts any `Notifier` and sends a message.
	//
	// 2. What's the difference between:
	//      var err error = nil
	//    and
	//      var e *CustomError = nil; var err error = e
	//    Why does the second one NOT equal nil?
	//
	// 3. Why does Go favor "accept interfaces, return structs"? What
	//    problems does this pattern solve in larger codebases?
}
