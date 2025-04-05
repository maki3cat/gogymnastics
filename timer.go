package main

import (
	"context"
	"fmt"
	"time"
)

func heartbeat(ctx context.Context, inbound <-chan string) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("Heartbeat action performed at ", time.Now())
		case msg := <-inbound:
			fmt.Println("Received message: ", msg, " at ", time.Now())
			time.Sleep(500 * time.Millisecond) // Simulate processing time
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inbound := make(chan string)

	go heartbeat(ctx, inbound)

	// Simulate some work
	go func() {
		for i := 0; i < 5; i++ {
			inbound <- fmt.Sprintf("Message %d", i)
			time.Sleep(2 * time.Second)
		}
		// the sender closes the channel,
		// closed channel can be read from but not written to
		close(inbound)
	}()
	time.Sleep(1 * time.Minute)

	// Cancel the context to stop the heartbeat
	cancel()

	// Wait a moment to ensure the heartbeat stops
	time.Sleep(2 * time.Second)
	fmt.Println("Main function exiting")
}
