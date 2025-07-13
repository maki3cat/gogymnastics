package concurrency_medium

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func CancelWithSameContext(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	child := func(name string) {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				fmt.Println(name, "done")
				return
			case <-ticker.C:
				fmt.Println(name, " is working")
			}
		}
	}

	go child("a")
	go child("b")
	go child("c")

	wg.Wait()
}

func CancelWithLayers(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	child := func(ctx context.Context, name string) {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				fmt.Println(name, "done")
				return
			case <-ticker.C:
				fmt.Println(name, " is working")
			}
		}
	}
	go child(subCtx, "x")
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	ticker := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("the outer function is working")
		case <-ctx.Done():
			fmt.Println("the outer function is done")
			return
		}
	}
}
