package tasks

import "fmt"

func PrintOptions() int {
	var option int
	fmt.Println("Choose option")
	fmt.Println("1. Add task")
	fmt.Println("2. List taks")
	fmt.Println("3. Exit")
	fmt.Scan(&option)
	return option
}

func AddTask(list []string, newT string) []string {
	return append(list, newT)
}
