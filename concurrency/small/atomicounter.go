package small

import (
	"context"
	"sync"
	"sync/atomic"
)

func AtomicCounter(ctx context.Context, threadCnt int, perThread int) int {
	var ans atomic.Int32
	wg := sync.WaitGroup{}
	unit := func(k int) {
		defer wg.Done()
		for range k {
			ans.Add(1)
		}
	}
	wg.Add(threadCnt)
	for range threadCnt {
		unit(perThread)
	}
	wg.Wait()
	return int(ans.Load())
}
