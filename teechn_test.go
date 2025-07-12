package main

import (
	"context"
	"testing"
)

func TestTeeChn(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		ctx := context.Background()
		in := make(chan any)
		out1, out2 := TeeChn(ctx, in)

		// Test sending values
		go func() {
			in <- "test1"
			in <- 42
			close(in)
		}()

		// Verify both channels receive all values in any order
		values1 := make([]any, 0, 2)
		values2 := make([]any, 0, 2)
		for i := 0; i < 2; i++ {
			values1 = append(values1, <-out1)
			values2 = append(values2, <-out2)
		}

		// Verify channels are closed
		if _, ok := <-out1; ok {
			t.Error("out1 should be closed")
		}
		if _, ok := <-out2; ok {
			t.Error("out2 should be closed")
		}

		// Verify all values were received
		expected := []any{"test1", 42}
		if len(values1) != len(expected) || len(values2) != len(expected) {
			t.Errorf("Expected %d values, got %d and %d", len(expected), len(values1), len(values2))
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		in := make(chan any)
		out1, out2 := TeeChn(ctx, in)

		// Cancel context before sending
		cancel()

		// Try to send value
		go func() {
			in <- "test"
			close(in)
		}()

		// Verify channels are closed
		if _, ok := <-out1; ok {
			t.Error("out1 should be closed after context cancellation")
		}
		if _, ok := <-out2; ok {
			t.Error("out2 should be closed after context cancellation")
		}
	})

	t.Run("empty input channel", func(t *testing.T) {
		ctx := context.Background()
		in := make(chan any)
		out1, out2 := TeeChn(ctx, in)

		// Close input immediately
		close(in)

		// Verify output channels are closed
		if _, ok := <-out1; ok {
			t.Error("out1 should be closed")
		}
		if _, ok := <-out2; ok {
			t.Error("out2 should be closed")
		}
	})

}
