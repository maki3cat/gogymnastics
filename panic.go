package main

import (
	"context"
	"fmt"
	"time"
)

func panic_in_subgoroutine(ctx context.Context) {
	go func() {
		fmt.Println("sub gorotine starts")
		panic("sub goroutine panicks")
	}()
}

func main() {
	panic_in_subgoroutine(context.Background())
	time.Sleep(time.Second)
	fmt.Println("main is still running")
}
