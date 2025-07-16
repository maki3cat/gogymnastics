package concurrency_small

import (
	"context"
	"testing"
	"time"
)

func TestNonblocking_Pass(t *testing.T) {
	ch := make(chan struct{})
	dosomething_once_nonblocking(ch)
}

func TestNonblocking_Done(t *testing.T) {
	ch := make(chan struct{}, 1)
	dosomething_once_nonblocking(ch)
}

func TestHeartbeat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()
	heartbeatChan := startHeartbeat(ctx, 1, time.Second)
	// maki: for and single select channel should be replaced with heartbeatChan
	for range heartbeatChan {
		t.Log("received heartbeat")
	}
}
