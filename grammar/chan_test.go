package grammar

import (
	"context"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	var nilChan chan int
	test := func(ctx context.Context, ch chan int) {
		defer func() {
			t.Log("test finished")
			t.Logf("in the test function, is the channel still nil: %v", ch == nil)
		}()
		select {
		case <-ch:
			t.Log("ch is not nil")
		case <-ctx.Done():
			t.Log("ctx is done")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go test(ctx, nilChan)
	nilChan = make(chan int, 1)
	t.Logf("in the outer scope, nilChan is nil: %v", nilChan == nil)
	nilChan <- 1
	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}
