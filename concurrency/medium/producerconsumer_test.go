package medium

import (
	"sync"
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {
	pc := NewProducerConsumer(1)

	wg := sync.WaitGroup{}
	wg.Add(4)

	// should be able to consume immediately
	go func() {
		pc.Consumer(1)
		wg.Done()
	}()
	// should block
	go func() {
		pc.Consumer(5)
		wg.Done()
	}()
	go func() {
		pc.Consumer(2)
		wg.Done()
	}()

	go func() {
		for range 3 {
			time.Sleep(1 * time.Second)
			pc.Produce(3)
		}
		wg.Done()
	}()
	wg.Wait()
}
