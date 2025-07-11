package gogymnastics

import (
	"fmt"
	"sync"
	"testing"
)

func lock_example() {
	var mu sync.Mutex
	counter := 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	for range 10 {
		go func() {
			mu.Lock()
			counter += 1
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

func TestLock(t *testing.T) {
	lock_example()
}