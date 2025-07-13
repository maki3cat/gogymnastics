package concurrency

import (
	"context"
	"sync"
)

// fan-in requires len(chns)+1 goroutines
func FanIn(ctx context.Context, chs ...<-chan any) <-chan any {
	// base cases
	if len(chs) == 0 {
		return nil
	}
	if len(chs) == 1 {
		return chs[0]
	}

	// cases that really need fanin
	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	fanInCh := make(chan any)

	// n workers to wait on faninCh
	fanInWorker := func(ch <-chan any) {
		defer wg.Done()
		for v := range ch {
			fanInCh <- v
		}
	}
	for _, ch := range chs {
		go fanInWorker(ch)
	}

	// the one to wait all fanin workers to finish
	go func() {
		wg.Wait()
		close(fanInCh)
	}()
	return fanInCh
}
