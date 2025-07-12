package concurrency

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
	pipeline := sub(context.Background(), mul(context.Background(), add(context.Background(), inbound)))
	for v := range pipeline {
		fmt.Println(v)
	}
}
