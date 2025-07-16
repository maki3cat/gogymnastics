package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// my question is: it this a good pattern that we don't call the cancel function and only relies on the timeout
// to ensure no goroutine leakage?
func onlyTimeout(ctx context.Context, timeout time.Duration) {
	ctxWithTimeout, _ := context.WithTimeout(ctx, timeout)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go someRPC(ctxWithTimeout, &waitGroup)
	waitGroup.Wait()
}

func someRPC(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	<-ctx.Done()
	fmt.Println("context done")
}
