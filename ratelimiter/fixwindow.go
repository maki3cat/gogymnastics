package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindow struct {
	rate        int
	count       int
	windowStart time.Time
	mu          sync.Mutex
}

func NewFixedWindow(rate int) *FixedWindow {
	return &FixedWindow{
		rate:        rate,
		count:       0,
		windowStart: time.Now(),
	}
}

// (1)
// unlike sliding window, fix window doesn't need to save the request time
// the fix window only needs a counter, and (current time - windowStart time) < 1 second
// thus only needs O(1) space
// (2)
// only needs the Allow() method
// (3)
// the window counter shall need lock
func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	// If we're in a new window, reset the count
	if now.Sub(fw.windowStart) >= time.Second {
		fw.count = 0
		fw.windowStart = now
	}

	// Check if we're at capacity
	if fw.count >= fw.rate {
		return false
	}

	// Increment count and allow request
	fw.count++
	return true
}
