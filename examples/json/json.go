package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Person is our main struct for demonstrating JSON marshaling.
// The json tags after each field control how the field appears in JSON output.
// For example, `json:"name"` means the JSON key will be "name", not "Name".
type Person struct {
	Name      string          `json:"name"`           // Basic field, always included
	Age       int             `json:"age,omitempty"`  // omitempty = skip if zero/empty
	Email     string          `json:"email,omitempty"`
	Address   Address         `json:"address,omitempty"`
	Nicknames []string        `json:"nicknames,omitempty"`
	RawData   json.RawMessage `json:"rawData,omitempty"` // Stores raw JSON bytes as-is
}

// Address is a nested struct to show how embedded structs are handled in JSON.
// Each field's json tag works the same way.
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

// CustomTime wraps time.Time to demonstrate custom JSON marshaling.
// Sometimes you want a different JSON format than the default.
// By default, time.Time marshals to something like "2006-01-02T15:04:05Z"
// We want just "2006-01-02" for our API, so we implement custom marshaling.
type CustomTime struct {
	time.Time
}

// MarshalJSON is called automatically when Go needs to convert CustomTime to JSON.
// It returns the date formatted as "YYYY-MM-DD" instead of the default ISO format.
// This lets us control the output format completely.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// time.Date.Format uses reference time "Mon Jan 2 15:04:05 MST 2006"
	// We format it as just the date portion: "2006-01-02"
	return json.Marshal(ct.Time.Format("2006-01-02"))
}

// UnmarshalJSON is the reverse of MarshalJSON - it converts JSON back to CustomTime.
// When we receive JSON like "2024-03-15", this method parses it back into a time.Time.
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// First, extract the string from the JSON bytes
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Then parse that string as our date format
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	// Store the parsed time back into our CustomTime
	ct.Time = t
	return nil
}

// Product demonstrates more complex JSON scenarios:
// - float64 (decimal numbers)
// - map[string]any (dynamic key-value pairs)
// - CustomTime (our own type with custom marshaling)
type Product struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Price     float64       `json:"price"`
	Tags      []string      `json:"tags,omitempty"`    // List of strings
	Metadata  map[string]any `json:"metadata,omitempty"` // Flexible key-value storage
	CreatedAt CustomTime    `json:"createdAt"`         // Uses our custom marshaling
}

func main() {
	fmt.Println("=== JSON Learning Examples ===")

	// We call each function to demonstrate a different JSON concept.
	// Each function is standalone so you can read them in any order.

	basics()           // Start here! This shows the fundamental Marshal/Unmarshal
	structTags()       // Learn how json tags control field names and inclusion
	omitemptyBehavior() // Learn how omitempty actually works (it's subtle!)
	customMarshaling() // Create your own JSON format for custom types
	rawMessageExample() // Store unparsed JSON for later processing
	streamingEncoderDecoder() // Stream large data without loading all into memory
	unknownStructure() // Decode JSON when you don't know the structure
	numberHandling()   // Why JSON numbers become float64 and how to handle it
	embeddedStructs()  // How struct embedding affects JSON output
	partialDecode()    // Decode only the fields you need
}

// basics demonstrates the two fundamental operations: Marshal and Unmarshal.
// Marshal = Go struct → JSON string (encoding)
// Unmarshal = JSON string → Go struct (decoding)
func basics() {
	fmt.Println("\n--- Basic Marshal/Unmarshal ---")

	// Create a Person struct (our source data)
	p := Person{
		Name:      "Alice",
		Age:       30,
		Email:     "alice@example.com",
		Nicknames: []string{"Ali", " lace"},
	}

	// json.Marshal converts a Go value to JSON bytes.
	// It returns ([]byte, error) - always check the error!
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("Marshaled: %s\n", string(data))
	// Output: {"name":"Alice","age":30,"email":"alice@example.com","nicknames":["Ali"," lace"]}
	// Notice: Address is omitted because it's an empty struct (omitempty applies to it too)

	// json.Unmarshal converts JSON bytes back to a Go struct.
	// The struct must be passed as a pointer (&p2) so Unmarshal can modify it.
	var p2 Person
	err = json.Unmarshal(data, &p2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("Unmarshaled: %+v\n", p2)
	// Output shows we got our data back!
}

// structTags show how json tags control JSON serialization behavior.
// A struct tag is a string that appears after a field declaration in backticks.
// The json:"..." tag tells the JSON encoder/decoder how to handle that field.
func structTags() {
	fmt.Println("\n--- Struct Tags ---")

	// This struct uses various json tag options:
	type User struct {
		ID       int    `json:"id"`                   // Key name in JSON will be "id" (lowercase)
		Name     string `json:"userName"`              // JSON key is "userName", not "Name"
		Password string `json:"-"`                     // "-" means NEVER include this field (good for secrets!)
		Internal string `json:"-"`                     // Same as above - completely excluded
		Notes    string `json:"notes,omitempty"`       // "omitempty" = skip if the value is empty/zero
	}

	u := User{
		ID:       1,
		Name:     "Bob",
		Password: "secret123",
		Internal: "hidden",
		Notes:    "", // Empty string - will be omitted due to omitempty
	}

	// json.MarshalIndent works like Marshal but adds indentation for readability.
	// Second parameter is a prefix (empty here), third is an indent string.
	data, _ := json.MarshalIndent(u, "", "  ")
	fmt.Printf("With tags:\n%s\n", string(data))

	// Key observations:
	// 1. "password" and "internal" don't appear at all (json:"-")
	// 2. "notes" doesn't appear because it's an empty string (omitempty)
	// 3. "id" appears as "id" (not "ID") and "userName" appears as "userName" (not "Name")

	fmt.Println("Note: Password and Internal are hidden due to json:\"-\"")
	fmt.Println("Note: Notes is omitted because it's empty (omitempty)")
}

// omitemptyBehavior demonstrates the nuanced behavior of omitempty.
// IMPORTANT: "empty" in Go means: nil, false, 0, empty string, empty slice, empty map, empty pointer.
// This is more comprehensive than you might expect!
func omitemptyBehavior() {
	fmt.Println("\n--- Omitempty Behavior ---")

	type Demo struct {
		Str    string           `json:"str,omitempty"`     // Empty string = omitted
		Int    int              `json:"int,omitempty"`     // 0 = omitted
		Slice  []string         `json:"slice,omitempty"`   // nil = omitted
		Map    map[string]int  `json:"map,omitempty"`     // nil = omitted
		Ptr    *int             `json:"ptr,omitempty"`     // nil = omitted
		Zero   int              `json:"zero"`              // NO omitempty - always included!
	}

	// Create Demo with ALL zero/empty/nil values
	d := Demo{
		Str:    "",    // empty string
		Int:    0,     // zero int
		Slice:  nil,   // nil slice
		Map:    nil,    // nil map
		Ptr:    nil,    // nil pointer
		Zero:   0,      // zero (but no omitempty!)
	}

	data, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("nil/empty values:\n%s\n", string(data))
	// Only "zero" appears! Everything with omitempty is hidden.

	fmt.Println("str, int, slice, map, ptr are OMITTED (omitempty)")
	fmt.Println("zero (int with zero value) is INCLUDED (no omitempty)")

	// Now create Demo with non-empty values
	d2 := Demo{
		Str:    "hello",           // non-empty string
		Int:    42,               // non-zero int
		Slice:  []string{},       // EMPTY slice (not nil!) - different from above!
		Map:    map[string]int{}, // EMPTY map (not nil!) - different from above!
		Ptr:    new(int),         // pointer to 0 (not nil!) - different from above!
		Zero:   0,               // still zero, but no omitempty
	}

	data2, _ := json.MarshalIndent(d2, "", "  ")
	fmt.Printf("\nNon-empty values:\n%s\n", string(data2))
	// Now ALL fields appear, even though some are "empty" in a sense.
	// Key insight: empty slices/maps and non-nil pointers are NOT omitted!

	fmt.Println("All fields now included (even empty slices/maps)")
	// The crucial distinction: nil vs empty ({} or [])
}

// customMarshaling shows how to create custom JSON formats for your own types.
// By implementing MarshalJSON and/or UnmarshalJSON, you control the JSON format completely.
func customMarshaling() {
	fmt.Println("\n--- Custom Marshaling (CustomTime) ---")

	prod := Product{
		ID:        1,
		Name:      "Widget",
		Price:     29.99,
		Tags:      []string{"electronics", "sale"},
		Metadata:  map[string]any{"color": "blue", "weight": 0.5}, // any = interface{}, accepts any type
		CreatedAt: CustomTime{time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)},
	}

	data, err := json.MarshalIndent(prod, "", "  ")
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("Custom time format:\n%s\n", string(data))
	// Notice "createdAt" shows as "2024-03-15" instead of "2024-03-15T00:00:00Z"
	// This is because our CustomTime.MarshalJSON returned a different format.

	fmt.Println("CreatedAt uses custom MarshalJSON to format as YYYY-MM-DD")

	// Unmarshal also works - our custom UnmarshalJSON parses "2024-03-15" back to time.Time
	var prod2 Product
	err = json.Unmarshal(data, &prod2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("Unmarshaled: %+v\n", prod2)
}

// rawMessageExample demonstrates json.RawMessage, which stores raw JSON bytes without parsing.
// This is useful when you want to:
// 1. Store JSON that you don't have a struct for yet
// 2. Pass through JSON without parsing and re-marshaling (preserves exact bytes)
// 3. Delay parsing until later when you know the structure
func rawMessageExample() {
	fmt.Println("\n--- json.RawMessage (Delayed Parsing) ---")

	// Here's some raw JSON as a string
	rawJSON := `{"name":"Test","metadata":{"nested":"value","numbers":[1,2,3]}}`

	// We can store this in a Person's RawData field directly
	var p Person
	p.RawData = json.RawMessage(rawJSON) // No parsing happens here - just stores the bytes

	// When we marshal Person, RawData appears as-is in the output
	data, _ := json.MarshalIndent(p, "", "  ")
	fmt.Printf("Raw JSON stored as-is:\n%s\n", string(data))
	// The raw JSON is embedded directly in the output - no escaping or reformatting

	fmt.Println("\nParsing RawData later:")
	// Later, when we're ready, we can parse just the RawData
	var metadata map[string]any
	json.Unmarshal(p.RawData, &metadata) // Now we parse it
	fmt.Printf("Parsed metadata: %+v\n", metadata)
	// metadata becomes: map[metadata:map[nested:value numbers:[1 2 3]] name:Test]
}

// streamingEncoderDecoder shows efficient ways to handle large JSON data streams.
// Instead of reading all JSON into memory at once, we stream it piece by piece.
// json.Encoder writes JSON directly to a writer (like a file or network connection).
// json.Decoder reads JSON from a reader without loading everything into memory.
func streamingEncoderDecoder() {
	fmt.Println("\n--- Streaming with Encoder/Decoder ---")

	// Use a buffer to simulate a stream (in real use, this could be a file or network)
	var buf bytes.Buffer

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	// json.NewEncoder writes JSON to a writer (we pass &buf, a *bytes.Buffer)
	// This is more efficient than Marshal + Write for multiple objects
	enc := json.NewEncoder(&buf)
	for _, p := range people {
		// Encode writes JSON directly to the buffer - no intermediate string
		if err := enc.Encode(p); err != nil {
			log.Fatalf("Encode error: %v", err)
		}
	}
	// Each Encode adds a newline, making it newline-delimited JSON (NDJSON)
	fmt.Printf("Encoded (streaming):\n%s", buf.String())

	// Now read it back with json.Decoder
	dec := json.NewDecoder(&buf)
	// dec.More() returns true if there's more JSON to decode (checks for remaining bytes)
	for dec.More() {
		var p Person
		// Decode reads one JSON object from the stream
		if err := dec.Decode(&p); err != nil {
			log.Fatalf("Decode error: %v", err)
		}
		fmt.Printf("Decoded: %+v\n", p)
	}
	// This is memory-efficient for large datasets - we only hold one object at a time
}

// unknownStructure demonstrates decoding JSON when you don't know its structure.
// By decoding into map[string]any, you can handle any JSON structure dynamically.
// Then use type switches to handle different value types.
func unknownStructure() {
	fmt.Println("\n--- Decoding Unknown Structure ---")

	// JSON with various value types
	jsonData := `{"name":"Test","count":42,"active":true,"items":["a","b"]}`

	// map[string]any is like a JavaScript object or Python dict
	// any is an alias for interface{} - can hold any value type
	var result map[string]any
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	fmt.Printf("Decoded map: %+v\n", result)
	// Output: map[active:true count:42 items:[a b] name:Test]
	// Notice: count is float64, not int (JSON has no integer type)

	// Use type switch to handle each value according to its actual type
	for k, v := range result {
		switch val := v.(type) {
		case string:
			// Handle string values
			fmt.Printf("  %s is string: %q\n", k, val)
		case float64:
			// Handle numbers (JSON numbers are always float64 by default!)
			fmt.Printf("  %s is float64: %f\n", k, val)
		case bool:
			// Handle boolean values
			fmt.Printf("  %s is bool: %t\n", k, val)
		case []any:
			// Handle arrays (elements are also any)
			fmt.Printf("  %s is array: %v\n", k, val)
		default:
			// Handle anything else (like nested objects)
			fmt.Printf("  %s is unknown type\n", k)
		}
	}
	fmt.Println("Note: JSON numbers become float64, not int")
	// This is because JSON doesn't distinguish between integer and floating-point
}

// numberHandling explores JSON's handling of large numbers and precision.
// JSON has only one number type (no int vs float distinction).
// By default, Go unmarshals JSON numbers as float64, which loses precision for large integers.
// json.Number preserves the exact string representation so you can parse it yourself.
func numberHandling() {
	fmt.Println("\n--- Number Handling ---")

	// Large integer and decimal number
	jsonWithNumber := `{"value": 18446744073709551615, "decimal": 3.14}`

	// Default behavior: numbers become float64
	var result map[string]any
	json.Unmarshal([]byte(jsonWithNumber), &result)
	// The large number loses precision: 1.8446744073709552e+19
	fmt.Printf("Default parsing: %+v\n", result)
	// Note the precision loss with large numbers!

	// Use json.Number to preserve the exact string representation
	type WithNumber struct {
		Value   json.Number `json:"value"`   // Keeps the number as a string internally
		Decimal json.Number `json:"decimal"` // Same here
	}

	var withNumber WithNumber
	json.Unmarshal([]byte(jsonWithNumber), &withNumber)

	fmt.Printf("With json.Number: %+v\n", withNumber)
	// Now Value and Decimal are json.Number types, preserving exact input

	// json.Number has methods to convert to specific types when needed
	i64, _ := withNumber.Value.Int64()
	fmt.Printf("Value as int64: %d\n", i64)
	// Careful: the number is too large for int64, so it overflows!

	f64, _ := withNumber.Decimal.Float64()
	fmt.Printf("Decimal as float64: %f\n", f64)
	// 3.14 converts cleanly to float64
}

// embeddedStructs shows how struct embedding affects JSON output.
// When you embed a struct (without a name), its fields appear at the parent level.
// This is called "field promotion" - embedded fields behave as if their fields were declared directly.
func embeddedStructs() {
	fmt.Println("\n--- Embedded Structs ---")

	type Inner struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	// Outer embeds Inner (no field name, just the type)
	// This is "composition" in Go - Inner is part of Outer
	type Outer struct {
		Inner        // Embedded - no field name means Inner's fields are promoted
		Name string `json:"name"`
	}

	o := Outer{
		Inner: Inner{X: 10, Y: 20}, // Must use field name when creating
		Name:  "Outer",
	}

	data, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf("Embedded struct (fields promoted):\n%s\n", string(data))
	// Notice: x and y appear at the TOP level, not nested under "Inner"
	// This is because Inner is embedded, so its fields are "promoted" to Outer

	fmt.Println("Note: X and Y appear at top level (promoted by embedding)")
	// If Inner were a named field (Inner Inner), output would be nested under "Inner"
}

// partialDecode demonstrates a common pattern: decode only what you need.
// Sometimes you have JSON with many fields, but you only care about a few.
// Option 1: Try to decode into a struct - fails if unknown fields exist (before Go 1.8)
// Option 2: Use map[string]any, then extract only what you need
func partialDecode() {
	fmt.Println("\n--- Partial Decode (Target Type Unknown) ---")

	jsonData := `{
		"id": 1,
		"name": "Product",
		"unknownField1": "value1",
		"unknownField2": 123
	}`

	// Try to decode into Person - this will fail because Person doesn't have
	// unknownField1 and unknownField2
	var p Person
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		fmt.Printf("Full decode error: %v\n", err)
		fmt.Println("This is expected - Person struct doesn't have all fields")
	}

	// Instead, decode into map[string]any - this accepts ANY JSON structure
	var generic map[string]any
	err = json.Unmarshal([]byte(jsonData), &generic)
	if err != nil {
		log.Fatalf("Generic decode error: %v", err)
	}
	fmt.Printf("Decoded to map[string]any: %+v\n", generic)
	// All fields are captured, even ones we don't have structs for

	// Now we can access only the fields we care about
	fmt.Println("Then access only the fields you need")
	// Example: if generic["id"] is all we need...
	if id, ok := generic["id"].(float64); ok { // Remember: numbers are float64!
		fmt.Printf("Found id: %.0f\n", id)
	}
}