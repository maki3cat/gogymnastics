package shortcircuit

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Point struct {
	X int
	Y int
}

func shoot(ctx context.Context, point Point) bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return point.X == 0
}

const workers = 10

func target(ctx context.Context, coordinate []Point) bool {

	tasks := make(chan Point)
	results := make(chan bool)
	finalResult := make(chan bool)

	// fanout the task
	go func() {
		defer close(tasks)
		for _, point := range coordinate {
			select {
			case <-ctx.Done():
				return
			case tasks <- point:
			}
		}
	}()
	// do the task
	worker := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case point, ok := <-tasks:
				if !ok {
					return
				}
				res := shoot(ctx, point)
				results <- res
			}
		}
	}
	wg := sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}
	// wait for all the workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// fan in the result
	go func() {
		for {
			select {
			case <-ctx.Done():
				finalResult <- false
				return
			case res, ok := <-results:
				if !ok {
					finalResult <- false
					return
				}
				if res {
					finalResult <- true
					return // short circuit
				}
			}
		}
	}()

	return <-finalResult
}
