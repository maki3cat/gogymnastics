package medium

import (
	"context"
	"fmt"
	"testing"
)

func TestPipeline(t *testing.T) {
	inbound := make(chan int)
	go func() {
		defer close(inbound)
		for i := range 10 {
			inbound <- i
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for v := range pipeline(ctx, inbound) {
		fmt.Println(v)
	}
}

func TestPipelineWithCancel(t *testing.T) {
	inbound := make(chan int)
	go func() {
		defer close(inbound)
		for i := range 5 {
			inbound <- i
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for v := range pipeline(ctx, inbound) {
		fmt.Println(v)
	}
}
