package gogymnastics

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
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // clean up on unexpected exit
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("stopping on context cancellation")
			return
		case t := <-ticker.C:
			fmt.Println("Heartbeat action performed at ", t)
			time.Sleep(100 * time.Millisecond) // Simulate async rpc time
		case msg, ok := <-inbound:
			if !ok {
				fmt.Println("stopping on inbound channel closed")
				return
			}
			ticker.Stop()
			fmt.Println("Connection maintained by Message: ", msg, " at ", time.Now())
			time.Sleep(100 * time.Millisecond) // Simulate async rpc time
			ticker.Reset(1 * time.Second)      // Reset the ticker after processing a message
		}
	}
}
