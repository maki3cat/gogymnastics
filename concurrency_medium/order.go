package concurrency_medium

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type OrderedThreads struct {
	names       []string
	currentName string
	cond        *sync.Cond
	lock        sync.Locker
}

func NewOrderedThreads(names []string) *OrderedThreads {
	lock := &sync.Mutex{}
	return &OrderedThreads{
		names:       names,
		currentName: names[0],
		lock:        lock,
		cond:        sync.NewCond(lock),
	}
}
func (ot *OrderedThreads) Start(ctx context.Context) {
	for idx, name := range ot.names {
		go ot.startWorker(ctx, name, idx)
		fmt.Println("starting the log processor of ", name)
	}
}

func (ot *OrderedThreads) getNextName(currentIdx int) string {
	return ot.names[(currentIdx+1)%len(ot.names)]
}

func (ot *OrderedThreads) startWorker(ctx context.Context, name string, idx int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ot.lock.Lock()
			for ot.currentName != name {
				ot.cond.Wait()
			}
			fmt.Println("handling the log of ", name)
			time.Sleep(1 * time.Second)
			ot.currentName = ot.getNextName(idx)
			ot.cond.Broadcast()
			ot.lock.Unlock()
		}
	}
}

// func main() {
// 	names := []string{"A", "B", "C"}
// 	ot := NewOrderedThreads(names)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ot.Start(ctx)
// 	time.Sleep(10 * time.Second)
// 	cancel()
// }
