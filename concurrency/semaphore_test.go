package concurrency

import (
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	t.Run("binary semaphore", func(t *testing.T) {
		sem := NewBinarySemaphore()
		wg := sync.WaitGroup{}
		counter := 0
		mutex := sync.Mutex{}

		// Launch 5 goroutines that try to increment counter
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem.Acquire()
				defer sem.Release()

				mutex.Lock()
				counter++
				// Simulate some work
				time.Sleep(time.Millisecond * 10)
				mutex.Unlock()
			}()
		}

		wg.Wait()
		if counter != 5 {
			t.Errorf("Expected counter to be 5, got %d", counter)
		}
	})

	t.Run("n-semaphore", func(t *testing.T) {
		n := 3
		sem := NewSem(n)
		active := 0
		maxActive := 0
		mutex := sync.Mutex{}
		wg := sync.WaitGroup{}

		// Launch 10 goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem.Acquire()
				defer sem.Release()

				mutex.Lock()
				active++
				if active > maxActive {
					maxActive = active
				}
				mutex.Unlock()

				// Simulate work
				time.Sleep(time.Millisecond * 50)

				mutex.Lock()
				active--
				mutex.Unlock()
			}()
		}

		wg.Wait()
		if maxActive > n {
			t.Errorf("Expected max concurrent routines to be %d, got %d", n, maxActive)
		}
	})
}
