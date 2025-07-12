package concurrency

import (
	"context"
)

// pattern: or channels
// use case, one return is needed from all channels
func OrChn(ctx context.Context, chans ...<-chan any) <-chan any {

	if len(chans) == 0 {
		return nil
	}
	if len(chans) == 1 {
		return chans[0]
	}
	orChan := make(chan any)

	go func() {
		defer close(orChan)
		recurChan := OrChn(ctx, chans[1:]...)
		select {
		case <-ctx.Done():
			return
		case value, ok := <-chans[0]:
			if !ok {
				return
			}
			orChan <- value
			return
		case value, ok := <-recurChan:
			if !ok {
				return
			}
			orChan <- value
			return
		}
	}()
	return orChan
}
