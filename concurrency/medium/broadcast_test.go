package medium

import (
	"context"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	broadcast(context.TODO(), 10)
}

func TestBroadcast2(t *testing.T) {
	inputCh := make(chan string)
	go broadcast2(context.TODO(), 3, inputCh)
	inputCh <- "hello"
	time.Sleep(time.Millisecond * 100)
	inputCh <- "world"
	time.Sleep(time.Millisecond * 500)
	close(inputCh)
}
