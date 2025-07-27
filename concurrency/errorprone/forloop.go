package main

import (
	"fmt"
	"sync"
	"time"
)

// processNumber simulates some work with a delay
func processNumber(num int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
	fmt.Printf("Processing number (incorrect): %d\n", num)
	time.Sleep(100 * time.Millisecond) // Simulate work
	fmt.Printf("Finished processing number (incorrect): %d\n", num)
}

func main() {
	fmt.Println("--- Running Experiment: Incorrectly Capturing Loop Variable ---")

	numbers := []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1) // Increment WaitGroup for each goroutine
		// ❌ Problematic: 'n' is captured by reference.
		// By the time the goroutines run, 'n' will likely be its final value (5).
		go processNumber(n, &wg)
	}

	wg.Wait() // Wait for all goroutines to complete
	fmt.Println("--- Experiment: Incorrectly Capturing Loop Variable Finished ---")
}
