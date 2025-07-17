package medium

import (
	"context"
	"testing"
	"time"
)

func TestCancelWithSameContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	CancelWithSameContext(ctx)
}

func TestWithLayers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	CancelWithLayers(ctx)
}

func TestTimeoutCancellation(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	CancelWithSameContext(ctx)
}
