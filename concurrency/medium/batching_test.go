package medium

import (
	"context"
	"testing"
	"time"
)

func TestBatching(t *testing.T) {
	t.Run("basic batching", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		b := NewBatching(ctx, 3)
		b.Start(ctx)

		// Send 5 items which should create 2 batches
		b.Send(ctx, 1)
		b.Send(ctx, 2)
		b.Send(ctx, 3) // Should trigger first batch [1,2,3]
		b.Send(ctx, 4)
		b.Send(ctx, 5) // Should trigger second batch [4,5]
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		b := NewBatching(ctx, 3)
		b.Start(ctx)

		b.Send(ctx, 1)
		b.Send(ctx, 2)

		// Cancel before batch is full
		time.Sleep(10 * time.Millisecond)
		cancel()
	})

	t.Run("drain partial batch", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		b := NewBatching(ctx, 4)
		b.Start(ctx)

		// Send 2 items - should create partial batch
		b.Send(ctx, 1)
		b.Send(ctx, 2)
		// Give time for batch processing
		time.Sleep(10 * time.Millisecond)
	})
}
