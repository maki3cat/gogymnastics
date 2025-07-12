package concurrency_medium

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inbound := make(chan string, 10)

	go heartbeatServer(ctx, inbound)

	// Simulate inbound requests
	go func() {
		for i := range 5 {
			inbound <- fmt.Sprintf("Message %d", i)
			gap := rand.Intn(500)
			time.Sleep(time.Duration(gap) * time.Millisecond)
		}
		// maki:
		// the Channel should not close the channel to make
		// the reader never block on the Channel
		// in Go, reading a closed channel will return the zero value, and ok=False of the channel type
		// close(inbound) -> cannot do this
	}()

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("Main function exiting")
}
