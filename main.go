// Package main defines the entry point of the Go application.
// Every executable Go program MUST have a package named "main".
// This tells the Go compiler that this code should produce a runnable binary,
// not a library that other packages import.
package main

// import brings in external packages so we can use their functionality.
// "fmt" is the standard formatting package from Go's standard library.
// It provides functions like Println, Printf, and Sprintf for outputting text.
// Think of it like `#include <stdio.h>` in C or `import sys` in Python.
import "fmt"

// func main() is the entry point of the program.
// When you run `go run main.go`, Go looks for this function and starts here.
// Every Go executable program must have exactly ONE main() in the main package.
// The opening brace { MUST be on the same line as func main() — this is Go syntax, not a style choice.
func main() {
	// fmt.Println() prints the given text followed by a newline to the terminal.
	// It automatically adds spaces between multiple arguments and a newline at the end.
	// Example: fmt.Println("Hello", "World") prints "Hello World\n"
	fmt.Println("Hello, World!")
}
