package concurrency

import (
	"context"
	"testing"
)

func TestFanIn(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		ch1 := make(chan any)
		ch2 := make(chan any)
		ch3 := make(chan any)

		// Send values on input channels
		go func() {
			ch1 <- 1
			ch2 <- 2
			ch3 <- 3
			close(ch1)
			close(ch2)
			close(ch3)
		}()

		// Fan in the channels
		merged := FanIn(context.Background(), ch1, ch2, ch3)

		// Collect results
		results := make([]int, 0, 3)
		for val := range merged {
			if v, ok := val.(int); ok {
				results = append(results, v)
			} else {
				t.Errorf("Expected int, got %T", val)
			}
		}

		// Verify we got all values
		if len(results) != 3 {
			t.Errorf("Expected 3 values, got %d", len(results))
		}

		// Verify we got 1, 2, and 3 in any order
		seen := make(map[int]bool)
		for _, v := range results {
			seen[v] = true
		}
		for _, expected := range []int{1, 2, 3} {
			if !seen[expected] {
				t.Errorf("Missing expected value: %d", expected)
			}
		}
	})

	t.Run("empty channels list", func(t *testing.T) {
		result := FanIn(context.Background())
		if result != nil {
			t.Error("Expected nil channel for empty input")
		}
	})

	t.Run("single channel", func(t *testing.T) {
		ch := make(chan any)
		go func() {
			ch <- 42
			close(ch)
		}()

		result := FanIn(context.Background(), ch)
		val := <-result
		if v, ok := val.(int); !ok || v != 42 {
			t.Errorf("Expected 42, got %v", val)
		}
	})
}
