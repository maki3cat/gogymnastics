package concurrency

import "context"

// maki:
// I think in the pipeline pattern,
// closing the channel is a good practice

func add(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// maki:
		// but if the in channel is blocked and forgets to close it, the goroutine will leak
		// ctx cancel cannot help it
		// so it seems closing a channel is a good practice
		for v := range in {
			select {
			// maki:
			// the key line is here to prevent possible goroutine leek
			case <-ctx.Done():
				return
			case out <- v + 1:
			}
		}
	}()
	return out
}

func mul(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			// maki: the key line is here to prevent possible goroutine leek, but what if the for range blocks
			case <-ctx.Done():
				return
			case out <- v * 2:
			}
		}
	}()
	return out
}

func pipeline(ctx context.Context, in <-chan int) <-chan int {
	return mul(ctx, add(ctx, in))
}
