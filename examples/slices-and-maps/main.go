package main

import "fmt"

func main() {
	fmt.Println("Learn slices and maps in golang!")

	// arrays in golang
	// arrays are var a [num]T where T is type

	var ages [3]int

	fmt.Print(ages) // auto assigned all elements to 0 by default

	// array initiliasation
	// var a [n]T = [n]T{V1, V2, ... Vn}

	var roll_numbers = [9]int{1, 2, 3, 4} // other elements are 0

	fmt.Println(roll_numbers)

	// shorthand

	salaries := [3]int{10000, 20000, 300000}

	fmt.Println(salaries)

	// access using [element index]
	//
	// multiple iteration ways
	// for loop classic

	for i := 0; i < len(salaries); i++ {
		fmt.Printf("Index %d, element %d\n", i, salaries[i])
	}

	// range fn

	for i, e := range salaries {
		fmt.Printf("Range fn: Index %d, element %d\n", i, e)
	}

	/**
	 * for i, e := range arr {} // Normal usage of range

		for _, e := range arr {} // Omit index with _ and use element

		for i := range arr {} // Use index only

		for range arr {} // Simply loop over the array
	*/

	// compiler infered length of an arr

	asks := [...]int{1, 2, 3, 32, 4, 5, 6, 23, 7}

	fmt.Println("The len is == ", len(asks))

	// the arr len is part of its type
	// cannot resize the array because of this

	// var aa = [4]int{1, 2, 3, 4}
	// var bb [2]int = aa // Error, cannot use a (type [4]int) as type [2]int in assignment

	// array in go are value types unlike reference in other langs like py and js

	var a = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var b = a // Copy of a is assigned to b

	b[0] = "Monday"

	fmt.Println(a)
	fmt.Println(b)

	// =======================================================================================
	// Slices
	// Or what we can call segment of an array or flexible array

	// 3 things in a slice
	// 1. pointer to an underlying array
	// 2. the length of the segrement of array
	// 3. the capacity that is max the segment can grow

	// declaration is like var s []T same as array

	var name []string

	fmt.Println(name) // will be []
	// the zero value of an array is also nil we can verify
	fmt.Println(name == nil)

	// you can use make function to init a slice
	//
	var boxes = make([]string, 0, 0)

	fmt.Print(boxes)
}
