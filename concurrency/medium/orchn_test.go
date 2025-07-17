package medium

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestOrChn(t *testing.T) {
	c1 := make(chan any, 1)
	c2 := make(chan any, 1)
	c3 := make(chan any, 1)
	c4 := make(chan any, 1)
	c5 := make(chan any, 1)
	c4 <- 4
	orChn := OrChn(context.Background(), c1, c2, c3, c4, c5)
	if orChn != nil {
		v, ok := <-orChn
		if ok {
			fmt.Println(v) // Changed to fmt.Println to properly display the value
		}
	}
}

func TestOrChn_WithContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c1 := make(chan any, 1)
	c2 := make(chan any, 1)
	c3 := make(chan any, 1)
	c4 := make(chan any, 1)
	c5 := make(chan any, 1)
	orChn := OrChn(ctx, c1, c2, c3, c4, c5)
	cancel()
	if orChn != nil {
		v, ok := <-orChn
		if ok {
			t.Errorf("orChn should be nil, but got %v", v)
		}
	}
}

func signal(duration time.Duration) <-chan any {
	ch := make(chan any)
	go func() {
		defer close(ch)
		time.Sleep(duration)
		ch <- struct{}{}
	}()
	return ch
}
func TestOrChn_Timeout(t *testing.T) {
	fmt.Println("start at", time.Now().Format(time.RFC3339Nano))
	<-OrChn(
		context.Background(), signal(time.Second), signal(time.Minute), signal(10*time.Millisecond), signal(time.Minute*2))
	fmt.Println("end at", time.Now().Format(time.RFC3339Nano))
}
