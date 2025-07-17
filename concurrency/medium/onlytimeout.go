package medium

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// my question is: it this a good pattern that we don't call the cancel function and only relies on the timeout
// to ensure no goroutine leakage?
func waitOnOnlyOne(ctx context.Context, timeout time.Duration) {
	ctxWithTimeout, _ := context.WithTimeout(ctx, timeout)
	returns := make(chan bool, 2)
	go someRPC(ctxWithTimeout, returns, "A")
	go someRPC(ctxWithTimeout, returns, "B")
	for res := range returns {
		if res {
			fmt.Println("one caller returns true")
			return
		}
	}
	fmt.Println("no caller returns true")
}

func someRPC(ctx context.Context, returns chan bool, name string) {
	select {
	case <-ctx.Done():
		fmt.Println(name, "exits on context done")
		returns <- false
	case <-time.After(time.Duration(rand.Intn(1100)) * time.Millisecond):
		fmt.Println(name, "returns true")
		returns <- true
	}
}
