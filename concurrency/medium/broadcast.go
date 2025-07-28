package medium

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// solution1: broadcast using "close a chan", this is same with ctx.Done()
// but this broadcast is one-off solution
// it is cannot be used to deliver content (can be modified to ctx like structure for simple message, but still not ready for complex message)
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

// solution2: broadcast using "different channel"
func worker2(ctx context.Context, inputCh chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inputCh:
			// always need to return on closing input channel
			if !ok {
				fmt.Println("broadcast closed", msg)
				return
			}
			fmt.Println("worker received: ", msg)
		}
	}
}

// suppose we don't need to cleanup
func broadcast2(ctx context.Context, workerCnt int, inputCh chan string) {
	// give each worker a separate channel to receive the message
	broadcastCh := make([]chan string, workerCnt)
	for i := range workerCnt {
		broadcastCh[i] = make(chan string)
		go worker2(ctx, broadcastCh[i])
	}

	// cleanup the dedicated channel for each worker
	defer func() {
		for _, ch := range broadcastCh {
			close(ch)
		}
	}()

	// first round of select is to select the input with ctx.Done
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inputCh:
			if !ok {
				fmt.Println("input channel closed")
				return
			}
			// try to broadcast to all,
			// still need to wait for ctx,
			// so has a select should be inside this for-loop
			for _, ch := range broadcastCh {
				select {
				case <-ctx.Done():
					return
				case ch <- msg:
				}
			}
		}
	}
}
