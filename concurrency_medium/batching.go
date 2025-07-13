package concurrency_medium

import (
	"context"
	"fmt"
)

// a very typical use-case of buffered channel
// and batching is a very useful pattern in system engineering
type Batching struct {
	buffer chan any
	size   int
}

func NewBatching(ctx context.Context, size int) *Batching {
	return &Batching{
		buffer: make(chan any, size),
		size:   size,
	}
}

func (b *Batching) drainBatch() []any {
	batch := make([]any, 0, b.size-1)
	count := 0
	for {
		select {
		case v := <-b.buffer:
			batch = append(batch, v)
			count++
			if count == b.size-1 {
				return batch
			}
		default:
			return batch
		}
	}
}

func (b *Batching) Start(ctx context.Context) {
	go func() {
		defer func() {
			close(b.buffer)
			fmt.Println("close batching")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-b.buffer:
				batch := b.drainBatch()
				batch = append(batch, v)
				fmt.Println("handle batch", batch)
			}
		}
	}()
	fmt.Println("start batching")
}

func (b *Batching) Send(ctx context.Context, v any) {
	select {
	case <-ctx.Done():
		return
	case b.buffer <- v:
	}
}
