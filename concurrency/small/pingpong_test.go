package small

import (
	"testing"
	"time"
)

func TestPingPong(t *testing.T) {
	table := make(chan *Ball)
	go ppplayer("ping", table)
	time.Sleep(10 * time.Millisecond)
	go ppplayer("pong", table)

	// table <- new(Ball) // game on; toss the ball
	time.Sleep(1 * time.Second)
	<-table // game over; grab the ball
	close(table)
}
