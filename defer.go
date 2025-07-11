package main

import "fmt"

func demoV2() (result int) {
	defer func() {
		result = 42 // overwrites the named return value
		fmt.Println("defer runs")
		return // legal, but doesn't exit the function again
	}()
	return 7
}

func demo() int {
	defer func() int {
		fmt.Println("defer runs")
		return 100 // legal, but doesn't exit the function again
	}()
	return 7
}

// func main() {
// 	fmt.Println(demo())
// 	fmt.Println(demoV2())
// }
