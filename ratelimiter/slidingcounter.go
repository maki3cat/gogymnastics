package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

type SlidingCounter struct {
	rate       int           // requests per second
	bucketSize time.Duration // size of each bucket
	buckets    []int         // array of counters for each bucket
	bucketTime time.Time     // start time of first bucket
	mu         sync.Mutex
}

func NewSlidingCounter(rate int) *SlidingCounter {
	return &SlidingCounter{
		rate:       rate,
		bucketSize: time.Second / 10, // 100ms per bucket
		buckets:    make([]int, 10),  // 10 buckets for 1s window
		bucketTime: time.Now(),
	}
}

func (sc *SlidingCounter) Allow() bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(sc.bucketTime)
	bucketCount := len(sc.buckets)

	// Calculate how many buckets have passed
	passedBuckets := int(elapsed / sc.bucketSize)
	if passedBuckets > 0 {
		// Reset old buckets
		if passedBuckets >= bucketCount {
			// If more than window has passed, reset all buckets
			for i := range sc.buckets {
				sc.buckets[i] = 0
			}
			sc.bucketTime = now
		} else {
			// Shift buckets left and zero out the rest
			copy(sc.buckets, sc.buckets[passedBuckets:])
			for i := bucketCount - passedBuckets; i < bucketCount; i++ {
				sc.buckets[i] = 0
			}
			sc.bucketTime = now.Add(-time.Second) // Set to 1 second before now
		}
	}

	// Calculate current total requests in window
	total := 0
	for _, count := range sc.buckets {
		total += count
	}

	// Check if we're at capacity
	if total >= sc.rate {
		return false
	}

	// Add request to current bucket
	currentBucket := int(elapsed / sc.bucketSize)
	if currentBucket >= bucketCount {
		currentBucket = bucketCount - 1
	}
	sc.buckets[currentBucket]++
	return true
}

func TestSlidingCounter(t *testing.T) {
	// Create sliding counter with 10 requests/sec split into 10 buckets
	sc := NewSlidingCounter(10)

	// Test initial requests within capacity
	for i := 0; i < 10; i++ {
		if !sc.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// Test request exceeding capacity
	if sc.Allow() {
		t.Error("Expected request to be denied when at capacity")
	}

	// Wait for half window to pass and test again
	time.Sleep(500 * time.Millisecond)
	if !sc.Allow() {
		t.Error("Expected request to be allowed after half window passed")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	sc = NewSlidingCounter(20)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Allow()
		}()
	}
	wg.Wait()

	// Verify the total count after concurrent access
	sc.mu.Lock()
	total := 0
	for _, count := range sc.buckets {
		total += count
	}
	if total > 20 {
		t.Errorf("Expected at most 20 requests, got %d", total)
	}
	sc.mu.Unlock()
}
