// =============================================================================
// Go Structs & Composition — A Practical, Interview-Ready Guide
//
// Run:  go run examples/structs-and-composition/structs_and_composition.go
// Build: go build -o demo structs_and_composition.go
//
// This file demonstrates structs, methods, embedding, interfaces, and
// real-world composition patterns in idiomatic Go.
// =============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

// =============================================================================
// 1. STRUCT DEFINITIONS — Exported vs Unexported Fields
// =============================================================================

// Person is our base struct.
// Exported fields (capitalized) are visible outside this package.
// Unexported fields (lowercase) are package-private — useful for encapsulation.
type Person struct {
	Name      string    // Exported — visible to other packages
	Email     string    // Exported
	Age       int       // Exported
	createdAt time.Time // Unexported — internal tracking, not part of public API
}

// =============================================================================
// 2. STRUCT INITIALIZATION — Multiple Ways
// =============================================================================

// newPerson is a factory function (Go's closest equivalent to a constructor).
// Factory functions let you validate input, set defaults, and hide implementation
// details — a core Go idiom.
func newPerson(name, email string, age int) (*Person, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if age < 0 {
		return nil, fmt.Errorf("age cannot be negative")
	}
	return &Person{
		Name:      name,
		Email:     email,
		Age:       age,
		createdAt: time.Now(),
	}, nil
}

// =============================================================================
// 3. METHODS — Value Receivers vs Pointer Receivers
// =============================================================================

// Greet uses a VALUE receiver (p Person).
//
//   - Go passes a COPY of the struct to this method.
//   - Mutations inside this method do NOT affect the original struct.
//   - Use value receivers when the method only READS data or when the struct
//     is small and cheap to copy (time.Duration, sync.Mutex are NOT cheap to copy).
//
// Common mistake: mutating fields in a value-receiver method and expecting
// the caller to see the change. They won't.
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s (%d)", p.Name, p.Age)
}

// HaveBirthday uses a POINTER receiver (p *Person).
//
//   - Go passes a pointer to the original struct.
//   - Mutations DO affect the original.
//   - Use pointer receivers when the method MODIFIES the receiver or when
//     the struct is large and copying is expensive.
//
// Rule of thumb: if ANY method on a struct needs a pointer receiver, make
// ALL methods on that struct use pointer receivers for consistency.
func (p *Person) HaveBirthday() {
	p.Age++
}

// Summary demonstrates that value receivers can still read unexported fields
// within the same package.
func (p Person) Summary() string {
	return fmt.Sprintf("%s, email: %s, member since %s",
		p.Name, p.Email, p.createdAt.Format("Jan 2006"))
}

// =============================================================================
// 4. COMPOSITION VIA EMBEDDING — Employee embeds Person
// =============================================================================

// Employee COMPOSES a Person by embedding it anonymously (no field name).
//
//   - This is NOT inheritance. Go has no inheritance.
//   - Employee IS-NOT-A Person; Employee HAS-A Person.
//   - The embedded Person's fields and methods are PROMOTED to Employee.
//     That means you can call e.Name, e.Greet() as if they belong to Employee.
//
// Why composition over inheritance?
//   - Inheritance creates tight coupling and fragile base-class problems.
//   - Composition is explicit, flexible, and avoids diamond problems.
//   - You can compose multiple unrelated structs (Go has no single-inheritance limit).
type Employee struct {
	Person             // Embedded — fields and methods are promoted
	EmployeeID   string
	Department   string
	Salary       float64
	Skills       []string
}

// Override demonstrates method shadowing.
// Employee defines its own Greet() which SHADOWS Person.Greet().
// When you call e.Greet(), this version runs — NOT Person.Greet().
//
// To call the embedded version explicitly: e.Person.Greet()
func (e Employee) Greet() string {
	return fmt.Sprintf("Hi, I'm %s from %s (ID: %s)", e.Name, e.Department, e.EmployeeID)
}

// GiveRaise uses a pointer receiver to modify the Salary field.
func (e *Employee) GiveRaise(percent float64) {
	e.Salary *= (1 + percent/100)
}

// AnnualCost returns the fully loaded annual cost (salary + 30% overhead).
func (e Employee) AnnualCost() float64 {
	return e.Salary * 1.3
}

// =============================================================================
// 5. MANAGER — Deeper Embedding Chain
// =============================================================================

// Manager embeds Employee, which embeds Person.
// Promoted fields and methods travel up the chain: m.Name, m.Greet(), m.HaveBirthday().
type Manager struct {
	Employee        // Embeds Employee (which embeds Person)
	TeamSize   int
	Bonus      float64
}

// Manager overrides Greet() again — further down the chain.
func (m Manager) Greet() string {
	return fmt.Sprintf("I'm %s, managing a team of %d in %s", m.Name, m.TeamSize, m.Department)
}

// =============================================================================
// 6. INTERFACES — Decoupled Behavior
// =============================================================================

// Greeter defines a single-method interface.
// Any type with a Greet() string method satisfies this interface implicitly.
// There is NO "implements" keyword in Go — satisfaction is duck-typed.
type Greeter interface {
	Greet() string
}

// Billable is another interface for cost calculation.
type Billable interface {
	AnnualCost() float64
}

// Multi-interface types: Manager satisfies BOTH Greeter and Billable because
// it has Greet() string and AnnualCost() float64 methods.

// Welcome is a function that accepts ANY Greeter.
// This is polymorphism without inheritance.
func Welcome(g Greeter) {
	fmt.Println("  →", g.Greet())
}

// PrintCost accepts ANY Billable and prints its annual cost.
func PrintCost(b Billable) {
	fmt.Printf("  → Annual cost: $%.2f\n", b.AnnualCost())
}

// =============================================================================
// 7. INTERFACE + COMPOSITION — A Realistic Pattern
// =============================================================================

// Notifier is an interface for sending notifications.
type Notifier interface {
	Notify(message string) error
}

// EmailNotifier implements Notifier using an embedded Person for the email field.
type EmailNotifier struct {
	Person           // Embed Person to access Email field
	Provider   string
}

// Notify uses the promoted Email field from Person.
func (e EmailNotifier) Notify(message string) error {
	fmt.Printf("  → Sending email to %s via %s: %q\n", e.Email, e.Provider, truncate(message, 40))
	return nil
}

// SlackNotifier implements Notifier but does NOT embed Person.
// It has its own Contact field instead.
//
// This demonstrates that you can satisfy the same interface with different
// internal structures — composition gives you flexibility.
type SlackNotifier struct {
	Contact string // Explicit field, not embedded
	Channel string
}

// Notify uses the explicit Contact field.
func (s SlackNotifier) Notify(message string) error {
	fmt.Printf("  → Slack to %s in #%s: %q\n", s.Contact, s.Channel, truncate(message, 40))
	return nil
}

// SendAlert accepts any Notifier — decoupled from concrete types.
func SendAlert(n Notifier, msg string) {
	_ = n.Notify(msg)
}

// =============================================================================
// 8. NESTED STRUCTS — Without Embedding
// =============================================================================

// Address is a standalone struct.
type Address struct {
	Street  string
	City    string
	Country string
}

// Contractor uses Address as a NAMED (not embedded) field.
//
// When to use explicit (named) fields vs embedding:
//   - Use EMBEDDING when you want promoted access (c.Name instead of c.Person.Name)
//     and when the relationship is truly "is implemented in terms of."
//   - Use NAMED FIELDS when you want to be explicit about the relationship
//     and avoid namespace collisions. This is safer and more readable.
//
// Common mistake: overusing embedding just to save a few keystrokes.
// If you find yourself qualifying with the type name often (c.Person.Name),
// a named field is probably clearer.
type Contractor struct {
	Person           // Embedded for promoted fields
	Address  Address // Named field — must access as c.Address.City
	Rate     float64
}

// AnnualCost satisfies Billable.
func (c Contractor) AnnualCost() float64 {
	return c.Rate * 2080 // 40hrs/week * 52 weeks
}

// =============================================================================
// 9. SHAPES — Another Composition Example (Geometry)
// =============================================================================

// Shape is an interface for geometric shapes.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle is a concrete shape.
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Rectangle is a concrete shape.
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// ColoredShape uses composition to ADD behavior to any Shape.
// This is the Decorator pattern via embedding.
//
// Note: Go embedding does NOT dispatch dynamically. When ColoredShape embeds
// Shape, the Area() and Perimeter() methods are promoted BUT they call the
// underlying type's methods directly — NOT through the interface.
// This is a subtle but important distinction from OOP inheritance.
type ColoredShape struct {
	Shape            // Embedded interface — promotes Area() and Perimeter()
	Color     string
	Material  string
}

// DescribeShape prints area and perimeter for any Shape.
func DescribeShape(s Shape) {
	fmt.Printf("  → Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// =============================================================================
// 10. JSON TAGS — Serialization
// =============================================================================

// Product demonstrates struct tags for JSON marshaling.
// Tags are string literals that encode metadata for the encoding/json package
// (and other packages like form, xml, validate, etc.).
//
//   - json:"name"           → use lowercase "name" in JSON
//   - json:"-"              → omit this field from JSON entirely
//   - json:"price,omitempty" → omit if zero value
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price,omitempty"`
	Description string  `json:"description,omitempty"`
	internalSKU string  `json:"-"` // Never serialized
}

// =============================================================================
// 11. HELPER FUNCTIONS
// =============================================================================

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func printSection(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  " + title)
	fmt.Println(strings.Repeat("=", 60))
}

// =============================================================================
// MAIN — Demonstrations
// =============================================================================

func main() {

	// -------------------------------------------------------------------------
	// SECTION 1: Zero Values
	// -------------------------------------------------------------------------
	printSection("1. ZERO VALUES")

	// Go initializes structs to their zero values automatically.
	//   string → ""
	//   int    → 0
	//   bool   → false
	//   slice  → nil
	//   struct → all fields zeroed
	var zeroPerson Person
	fmt.Printf("  Zero Person: Name=%q, Age=%d, Email=%q\n",
		zeroPerson.Name, zeroPerson.Age, zeroPerson.Email)
	fmt.Println("  → Zero values are safe to use but often not semantically valid.")
	fmt.Println("  → Use factory functions to enforce invariants.")

	// -------------------------------------------------------------------------
	// SECTION 2: Struct Initialization
	// -------------------------------------------------------------------------
	printSection("2. STRUCT INITIALIZATION")

	// Way 1: Factory function (idiomatic — allows validation)
	p1, err := newPerson("Alice", "alice@example.com", 30)
	if err != nil {
		fmt.Println("  Error:", err)
	} else {
		fmt.Println("  Factory:", p1.Summary())
	}

	// Way 2: Literal with named fields (order doesn't matter — preferred)
	p2 := Person{
		Name: "Bob",
		Age:  25,
	}
	fmt.Printf("  Named literal: %+v\n", p2)

	// Way 3: Literal with positional fields (order MUST match — fragile)
	p3 := Person{"Charlie", "charlie@example.com", 28, time.Now()}
	fmt.Printf("  Positional: %+v\n", p3)

	// Way 4: Zero-value then assign
	var p4 Person
	p4.Name = "Diana"
	p4.Age = 35
	fmt.Printf("  Zero + assign: %+v\n", p4)

	// -------------------------------------------------------------------------
	// SECTION 3: Value vs Pointer Receivers
	// -------------------------------------------------------------------------
	printSection("3. VALUE vs POINTER RECEIVERS")

	p, _ := newPerson("Eve", "eve@example.com", 29)
	fmt.Printf("  Before HaveBirthday: Age=%d\n", p.Age)

	// HaveBirthday has a pointer receiver — it modifies the original.
	p.HaveBirthday()
	fmt.Printf("  After HaveBirthday:  Age=%d (mutated through pointer receiver)\n", p.Age)

	// Greet has a value receiver — it works on a copy.
	// But since it only reads, the distinction doesn't matter to the caller.
	fmt.Println("  Greet:", p.Greet())

	// -------------------------------------------------------------------------
	// SECTION 4: Composition & Promoted Fields
	// -------------------------------------------------------------------------
	printSection("4. COMPOSITION — EMBEDDING & PROMOTED FIELDS")

	emp := Employee{
		Person: Person{
			Name:  "Frank",
			Email: "frank@company.com",
			Age:   34,
		},
		EmployeeID: "EMP-001",
		Department: "Engineering",
		Salary:     95000,
		Skills:     []string{"Go", "Kubernetes", "gRPC"},
	}

	// Promoted fields — accessing Person.Name as if it were Employee.Name
	fmt.Printf("  Promoted field: emp.Name = %q\n", emp.Name)
	fmt.Printf("  Promoted field: emp.Email = %q\n", emp.Email)
	fmt.Printf("  Direct field:   emp.EmployeeID = %q\n", emp.EmployeeID)

	// Promoted method — but Employee overrides Greet()
	fmt.Println("  Promoted (overridden) method: emp.Greet():")
	fmt.Println("    →", emp.Greet())

	// Explicit access to the embedded Person's Greet()
	fmt.Println("  Embedded method (explicit): emp.Person.Greet():")
	fmt.Println("    →", emp.Person.Greet())

	// HaveBirthday is NOT overridden — it's promoted from Person
	emp.HaveBirthday()
	fmt.Printf("  Promoted method: emp.HaveBirthday() → Age is now %d\n", emp.Age)

	// -------------------------------------------------------------------------
	// SECTION 5: Method Shadowing / Overriding
	// -------------------------------------------------------------------------
	printSection("5. METHOD SHADOWING (NOT TRUE OVERRIDING)")

	mgr := Manager{
		Employee: Employee{
			Person: Person{
				Name:  "Grace",
				Email: "grace@company.com",
				Age:   42,
			},
			EmployeeID: "MGR-001",
			Department: "Platform",
			Salary:     140000,
		},
		TeamSize: 8,
		Bonus:    15000,
	}

	// Manager has its own Greet() — shadows both Employee.Greet() and Person.Greet()
	fmt.Println("  mgr.Greet() (Manager's version):")
	fmt.Println("    →", mgr.Greet())

	// To reach Employee.Greet():
	fmt.Println("  mgr.Employee.Greet() (Employee's version):")
	fmt.Println("    →", mgr.Employee.Greet())

	// To reach Person.Greet():
	fmt.Println("  mgr.Employee.Person.Greet() (Person's version):")
	fmt.Println("    →", mgr.Employee.Person.Greet())

	// HaveBirthday is NOT shadowed anywhere — promoted all the way from Person
	mgr.HaveBirthday()
	fmt.Printf("  mgr.HaveBirthday() → Age is now %d (promoted from Person)\n", mgr.Age)

	// -------------------------------------------------------------------------
	// SECTION 6: Interfaces + Composition
	// -------------------------------------------------------------------------
	printSection("6. INTERFACES + COMPOSITION")

	fmt.Println("  Passing Employee (satisfies Greeter):")
	Welcome(emp)

	fmt.Println("  Passing Manager (satisfies Greeter + Billable):")
	Welcome(mgr)
	PrintCost(mgr)

	fmt.Println("  Passing Contractor (satisfies Billable):")
	contractor := Contractor{
		Person: Person{
			Name:  "Hank",
			Email: "hank@freelance.io",
			Age:   38,
		},
		Address: Address{
			Street:  "123 Remote St",
			City:    "Austin",
			Country: "US",
		},
		Rate: 85,
	}
	PrintCost(contractor)
	fmt.Printf("  Named field access: contractor.Address.City = %q\n", contractor.Address.City)

	// -------------------------------------------------------------------------
	// SECTION 7: Notification System — Same Interface, Different Internals
	// -------------------------------------------------------------------------
	printSection("7. INTERFACE DECOUPLING — NOTIFICATION SYSTEM")

	// EmailNotifier embeds Person — uses promoted Email field
	emailNotif := EmailNotifier{
		Person: Person{
			Email: "team@company.com",
		},
		Provider: "SendGrid",
	}

	// SlackNotifier has explicit fields — no embedding
	slackNotif := SlackNotifier{
		Contact: "@dev-team",
		Channel: "alerts",
	}

	// Both satisfy Notifier — SendAlert doesn't care about the concrete type
	fmt.Println("  SendAlert via EmailNotifier:")
	SendAlert(emailNotif, "Server CPU at 95% — investigate immediately")

	fmt.Println("  SendAlert via SlackNotifier:")
	SendAlert(slackNotif, "Server CPU at 95% — investigate immediately")

	fmt.Println("  → Same interface, completely different internal structures.")
	fmt.Println("  → Composition lets you pick the right structure per use case.")

	// -------------------------------------------------------------------------
	// SECTION 8: Shapes — Composition with Interface Embedding
	// -------------------------------------------------------------------------
	printSection("8. GEOMETRY — INTERFACE EMBEDDING (DECORATOR PATTERN)")

	circle := Circle{Radius: 5}
	fmt.Println("  Circle:")
	DescribeShape(circle)

	rect := Rectangle{Width: 4, Height: 6}
	fmt.Println("  Rectangle:")
	DescribeShape(rect)

	// ColoredShape embeds the Shape interface — adds metadata without changing behavior
	coloredCircle := ColoredShape{
		Shape:    circle,
		Color:    "Blue",
		Material: "Metal",
	}
	fmt.Printf("  ColoredCircle (%s, %s):\n", coloredCircle.Color, coloredCircle.Material)
	DescribeShape(coloredCircle)
	// Note: Area() is promoted from the embedded Circle, not from the interface
	// at runtime. This is static dispatch — a key Go difference from OOP.

	// -------------------------------------------------------------------------
	// SECTION 9: JSON Serialization
	// -------------------------------------------------------------------------
	printSection("9. JSON SERIALIZATION WITH TAGS")

	product := Product{
		ID:          "PRD-42",
		Name:        "Widget Pro",
		Price:       29.99,
		Description: "A professional-grade widget for everyday use.",
		internalSKU: "WG-PRO-2024",
	}

	data, _ := json.MarshalIndent(product, "  ", "  ")
	fmt.Println("  Marshaled JSON:")
	fmt.Println("  ", string(data))
	fmt.Println("  → Note: internalSKU is omitted (json:\"-\")")
	fmt.Println("  → Fields follow camelCase in Go, snake_case in JSON via tags")

	// Demonstrate omitempty with zero price
	zeroPriceProduct := Product{
		ID:   "PRD-99",
		Name: "Free Sample",
	}
	data2, _ := json.MarshalIndent(zeroPriceProduct, "  ", "  ")
	fmt.Println("  Product with zero Price (omitempty):")
	fmt.Println("  ", string(data2))
	fmt.Println("  → Price field omitted because it's zero and has omitempty")

	// -------------------------------------------------------------------------
	// SECTION 10: Anonymous Structs
	// -------------------------------------------------------------------------
	printSection("10. ANONYMOUS STRUCTS")

	// Anonymous struct — useful for one-off data grouping, test fixtures,
	// or JSON unmarshaling without polluting the package namespace.
	config := struct {
		Host string
		Port int
		Debug bool
	}{
		Host:  "localhost",
		Port:  8080,
		Debug: true,
	}
	fmt.Printf("  Anonymous struct: %+v\n", config)

	// Slice of anonymous structs — handy for table-driven tests
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"even", 4, true},
		{"odd", 7, false},
		{"zero", 0, true},
	}
	for _, tt := range tests {
		fmt.Printf("  Test %q: input=%d, expected=%v\n", tt.name, tt.input, tt.expected)
	}

	// -------------------------------------------------------------------------
	// SECTION 11: BEST PRACTICES SUMMARY
	// -------------------------------------------------------------------------
	printSection("11. BEST PRACTICES — QUICK REFERENCE")

	fmt.Println("  EMBEDDING vs EXPLICIT FIELDS:")
	fmt.Println("    • Embed when you want promoted access and behavioral reuse")
	fmt.Println("    • Use explicit fields when clarity and decoupling matter more")
	fmt.Println("    • Avoid embedding just to save typing — readability > brevity")
	fmt.Println()
	fmt.Println("  RECEIVERS:")
	fmt.Println("    • Pointer receiver when mutating or for large structs")
	fmt.Println("    • Value receiver for small, immutable structs (time.Duration)")
	fmt.Println("    • Be consistent — don't mix on the same type")
	fmt.Println()
	fmt.Println("  COMPOSITION vs INHERITANCE:")
	fmt.Println("    • Go has NO inheritance — only composition")
	fmt.Println("    • Embedding IS composition, not inheritance")
	fmt.Println("    • Promoted methods call the embedded type's methods directly")
	fmt.Println("    • There is no dynamic/virtual dispatch through embedding")
	fmt.Println()
	fmt.Println("  INTERFACES:")
	fmt.Println("    • Keep interfaces small (1-3 methods) — Interface Segregation")
	fmt.Println("    • Accept interfaces, return concrete types (Rob Pike)")
	fmt.Println("    • Define interfaces where they are USED, not where they are satisfied")
	fmt.Println()
	fmt.Println("  FACTORY FUNCTIONS:")
	fmt.Println("    • Use factories (newFoo) to validate and enforce invariants")
	fmt.Println("    • Return (*T, error) to allow callers to handle construction failures")
	fmt.Println()
	fmt.Println("  ZERO VALUES:")
	fmt.Println("    • Design types so their zero value is useful (e.g., sync.Mutex)")
	fmt.Println("    • When zero value is not useful, use a factory and unexport the struct")

	// -------------------------------------------------------------------------
	// FINAL DEMO: Polymorphic slice
	// -------------------------------------------------------------------------
	printSection("BONUS: POLYMORPHIC SLICE OF SHAPES")

	shapes := []Shape{
		Circle{Radius: 3},
		Rectangle{Width: 5, Height: 10},
		Circle{Radius: 7},
	}

	var totalArea float64
	for i, s := range shapes {
		area := s.Area()
		totalArea += area
		fmt.Printf("  Shape %d: Area = %.2f\n", i+1, area)
	}
	fmt.Printf("  Total area of all shapes: %.2f\n", totalArea)

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  END OF DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 60))
}
