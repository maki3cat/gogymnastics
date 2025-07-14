package concurrency_medium

import (
	"context"
	"testing"
	"time"
)

func TestOrder(t *testing.T) {

	names := []string{"A", "B", "C"}
	ot := NewOrderedThreads(names)
	ctx, cancel := context.WithCancel(context.Background())
	ot.Start(ctx)
	time.Sleep(10 * time.Second)
	cancel()

}
