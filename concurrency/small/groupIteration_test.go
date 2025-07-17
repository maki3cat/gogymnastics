package small

import (
	"context"
	"testing"
	"time"
)

func TestGroupIteration_Basic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	k := 5
	A := make(chan struct{}, k)
	B := make(chan struct{}, k)
	C := make(chan struct{}, k)

	// Pre-fill channels to trigger each case
	for range k {
		A <- struct{}{}
		B <- struct{}{}
		C <- struct{}{}
	}

	done := make(chan struct{})
	go func() {
		groupIteration(ctx, A, B, C)
		close(done)
	}()

	// Let it run a bit
	time.Sleep(100 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("groupIteration did not exit after context cancel")
	}
}
