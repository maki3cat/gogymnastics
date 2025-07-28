package medium

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// solution1: broadcast using "close a chan"
func worker(ctx context.Context, broadcastCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case <-broadcastCh:
		fmt.Println("worker received broadcast")
		return
	}
}

func broadcast(ctx context.Context, workerCnt int) {
	broadcastCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	for range workerCnt {
		wg.Add(1)
		go worker(ctx, broadcastCh, wg)
	}
	time.Sleep(time.Second * 1)
	close(broadcastCh)
	wg.Wait()
	fmt.Println("broadcast done")
}
