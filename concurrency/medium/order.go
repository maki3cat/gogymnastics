package medium

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type OrderedThreads struct {
	names   []string        // can be a,b,c,d...
	signals []chan struct{} // signals that there are logs

	currentName string // subsystem's name
	cond        *sync.Cond
	lock        sync.Locker
}

func NewOrderedThreads(names []string) *OrderedThreads {
	lock := &sync.Mutex{}
	signals := make([]chan struct{}, len(names))
	for idx := range names {
		signals[idx] = make(chan struct{}, 1)
	}
	return &OrderedThreads{
		names:       names,
		currentName: names[0],
		lock:        lock,
		cond:        sync.NewCond(lock),
		signals:     signals,
	}
}

func (ot *OrderedThreads) Signal(idx int) {
	select {
	case ot.signals[idx] <- struct{}{}:
	default:
		return

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
	defer func() {
		fmt.Println("graceful shutdown, open files/buffered logs cleanning up")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ot.lock.Lock()
			for ot.currentName != name {
				ot.cond.Wait()
			}
			select {
			case <-ot.signals[idx]:
				fmt.Println("mock batching handling the log of ", name)
			default:
				fmt.Println("this subsystem has no logs input", name)
			}
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
