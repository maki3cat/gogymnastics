package medium

import (
	"context"
	"testing"
	"time"
)

func TestOrder(t *testing.T) {

	names := []string{"A", "B", "C"}
	ot := NewOrderedThreads(names)
	go func() {
		for range 10 {
			ot.Signal(1)
			time.Sleep(2 * time.Second)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	ot.Start(ctx)
	time.Sleep(10 * time.Second)
	cancel()

}
