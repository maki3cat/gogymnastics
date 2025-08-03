package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// a possible habit learned from this example
// 1. we can set a default size for the slice as something (0, 64/128/256)
// 2. we defer close the channer in the sender
// 3. the receiver always handle the case when the channel is closed
// 4. check edgey cases: zero/nil slice, unbounded growth/scalability
// 5. if the input of slice can be empty, use a local chan of nil when empty/real chan of not nil when not empty so that the select can walk around this one

// unbounded buffer
func startUnboundedBuffer(ctx context.Context, in chan any) chan any {
	out := make(chan any)
	go func() {
		buf := make([]any, 0, 128) // we can set a default size for the buffer
		defer close(out)           // we can let the sender to defer closing the channel
		var localOut chan any
		var outElement any
		for {
			if len(buf) == 0 {
				localOut = nil
			} else {
				localOut = out
				outElement = buf[0]
			}
			select {
			case <-ctx.Done():
				fmt.Println("exits buffer")
				return
			case localOut <- outElement:
				// maki: this should be here instead of the forloop,
				// because the outElement may not be used in this round of loop
				buf = buf[1:]
			case inElement := <-in:
				buf = append(buf, inElement)
				// fmt.Println("adding", inElement, "currentBuffer", buf)
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan any, 0)
	out := startUnboundedBuffer(ctx, in)

	wg := sync.WaitGroup{}
	wg.Add(1)
	// sender
	go func() {
		defer wg.Done()
		for i := range 10 {
			in <- i
		}
		fmt.Println("sender finished")
	}()
	wg.Wait()

	wg.Add(1)
	// receiver
	go func() {
		defer wg.Done()
		for v := range out {
			fmt.Println("out:", v)
		}
	}()
	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
}
