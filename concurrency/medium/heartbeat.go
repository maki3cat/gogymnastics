package medium

import (
	"context"
	"fmt"
	"time"
)

// Reading:
// Basics Understanding of Ticker
// https://gobyexample.com/tickers
// read the ticker source code, the buffer channel is 1, and
// if reading is slow other ticks will be dropped
// len(ch) == 1, drop until the channel is empty

// The Pattern:
// The Heartbeat Pattern utilizes a ticker to perform a heartbeat action at regular, fixed intervals.
// However, when a request is received,
// the heartbeat is not needed, as the incoming request itself can be used to maintain the connection.
func heartbeatServer(ctx context.Context, inbound <-chan string) {
	tickerGap := 100 * time.Millisecond
	ticker := time.NewTicker(tickerGap)
	defer ticker.Stop() // clean up on unexpected exit
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("stopping on context cancellation")
			return
		default:
			{
				select {
				case <-ctx.Done():
					ticker.Stop()
					fmt.Println("stopping on context cancellation")
					return
				case <-ticker.C:
					fmt.Println("Heartbeat action performed at ", time.Now().UnixNano()/1000_000)
				case msg, ok := <-inbound:
					if !ok {
						fmt.Println("stopping on inbound channel closed")
						return
					}
					ticker.Reset(tickerGap) // Reset the ticker after processing a message
					go func() {
						fmt.Println("async handle some request: ", msg, " at ", time.Now().UnixNano()/1000_000)
					}()
				}
			}
		}
	}
}
