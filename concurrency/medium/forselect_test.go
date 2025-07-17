package medium

import (
	"context"
	"testing"
	"time"
)

func TestForSelect(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		ch := make(chan any)
		ctx, cancel := context.WithCancel(context.Background())
		
		go forSelect(ctx, ch)
		
		// Give goroutine time to start
		time.Sleep(10 * time.Millisecond)
		
		cancel()
		// Give time for goroutine to exit
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("channel close", func(t *testing.T) {
		ch := make(chan any)
		ctx := context.Background()
		
		go forSelect(ctx, ch)
		
		ch <- "test"
		close(ch)
		
		// Give time for goroutine to process and exit
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("receive value", func(t *testing.T) {
		ch := make(chan any)
		ctx := context.Background()
		
		go forSelect(ctx, ch)
		
		ch <- "test value"
		
		// Give time for goroutine to process
		time.Sleep(10 * time.Millisecond)
	})
}