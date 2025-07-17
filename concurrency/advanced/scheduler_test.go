package advanced

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTimebasedScheduler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	scheduler := NewTimebasedScheduler()
	go scheduler.Run(ctx)

	task1 := Task{runAt: time.Now().Add(3 * time.Second), callback: make(chan struct{})}
	scheduler.AddTask(task1)
	task2 := Task{runAt: time.Now().Add(500 * time.Millisecond), callback: make(chan struct{})}
	scheduler.AddTask(task2)

	for {
		select {
		case <-task1.callback:
			fmt.Println("task executed at", time.Now().Unix())
		case <-task2.callback:
			fmt.Println("task2 executed at", time.Now().Unix())
		case <-ctx.Done():
			t.Fatal("context done")
		}

	}
}
