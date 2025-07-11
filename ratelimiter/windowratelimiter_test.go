package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

type SlidingWindow struct {
	capacity int
	rate     int
	window   time.Duration
	requests []time.Time
	mu       sync.Mutex
}

func NewSlidingWindow(capacity int, rate int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		capacity: capacity,
		rate:     rate,
		window:   window,
		requests: make([]time.Time, 0),
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	validRequests := sw.requests[:0] // reuse of the slice
	for _, t := range sw.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	sw.requests = validRequests

	// Check if we're at capacity
	if len(sw.requests) >= sw.capacity {
		return false
	}

	// Add new request
	sw.requests = append(sw.requests, now)
	return true
}

func TestSlidingWindow(t *testing.T) {
	// Create a sliding window with capacity 3, rate 3, and 1 second window
	sw := NewSlidingWindow(3, 3, time.Second)

	// Test initial requests within capacity
	for i := 0; i < 3; i++ {
		if !sw.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// Test request exceeding capacity
	if sw.Allow() {
		t.Error("Expected request to be denied when at capacity")
	}

	// Wait for window to pass and test again
	time.Sleep(time.Second)
	if !sw.Allow() {
		t.Error("Expected request to be allowed after window passed")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	sw = NewSlidingWindow(5, 5, time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sw.Allow()
		}()
	}
	wg.Wait()

	// Verify the number of requests after concurrent access
	sw.mu.Lock()
	if len(sw.requests) > 5 {
		t.Errorf("Expected at most 5 requests, got %d", len(sw.requests))
	}
	sw.mu.Unlock()
}
