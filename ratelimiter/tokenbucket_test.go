package ratelimiter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type TokenBucket struct {
	capacity int
	rate     int
	bucket   chan struct{}
}

func NewTokenBucket(capacity int, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		rate:     rate,
		bucket:   make(chan struct{}, capacity),
	}
}

func (tb *TokenBucket) Start(ctx context.Context) {
	go tb.refill(ctx)
}

func (tb *TokenBucket) getToken() bool {
	select {
	case <-tb.bucket:
		return true
	default:
		return false
	}
}

func (tb *TokenBucket) refill(ctx context.Context) {
	ticker := time.NewTicker(time.Second / time.Duration(tb.rate))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case tb.bucket <- struct{}{}:
					fmt.Println("refill one token")
				default:
					// channel is full, drop the token
				}
			}
		}
	}
}

func TestTokenBucket(t *testing.T) {
	// Create a token bucket with capacity 5 and rate 2 tokens/second
	tb := NewTokenBucket(5, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the token bucket
	tb.Start(ctx)
	time.Sleep(1 * time.Second)

	// Test initial capacity
	tokensReceived := 0
	for range 5 {
		if tb.getToken() {
			tokensReceived++
		}
	}
	if tokensReceived != 5 {
		t.Errorf("Expected 5 initial tokens, got %d", tokensReceived)
	}

	// Test cancellation
	cancel()
	time.Sleep(100 * time.Millisecond) // Give some time for goroutine to stop
}
