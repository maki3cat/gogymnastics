package main

import (
	"context"
	"sync"
)

func FanIn(ctx context.Context, chs ...<-chan any) <-chan any {

	if len(chs) == 0 {
		return nil
	}
	if len(chs) == 1 {
		return chs[0]
	}

	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	fanInCh := make(chan any)

	// n workers to wait on faninCh
	for _, ch := range chs {
		go func(ch <-chan any) {
			defer wg.Done()
			for v := range ch {
				fanInCh <- v
			}
		}(ch)
	}

	// another one to wait on wait and close fanInCh
	go func() {
		wg.Wait()
		close(fanInCh)
	}()
	return fanInCh
}
