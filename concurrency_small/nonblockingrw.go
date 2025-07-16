package concurrency_small

import (
	"context"
	"fmt"
	"time"
)

// how to send to a channel only when it doesn't block
func dosomething_once_nonblocking(c chan<- struct{}) {
	select {
	case c <- struct{}{}:
		fmt.Println("sent 1")
	default:
		fmt.Println("channel is full, exit")
		return
	}
}

// for example sending a heartbeat only when the previous is not consumed
func startHeartbeat(ctx context.Context, bufferSize int, interval time.Duration) <-chan struct{} {
	heartbeatChan := make(chan struct{}, bufferSize)
	go func(ctx context.Context) {
		defer close(heartbeatChan)
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done, exit")
				return
			case <-ticker.C:
				select {
				case heartbeatChan <- struct{}{}:
					fmt.Println("sent heartbeat")
				default:
					fmt.Println("channel is full, skip")
				}
			}
		}
	}(ctx)
	return heartbeatChan
}
