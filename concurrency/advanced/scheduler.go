package advanced

import (
	"container/heap"
	"context"
	"fmt"
	"time"
)

type Task struct {
	runAt    time.Time
	callback chan struct{}
}

// taskHeap implements heap.Interface and holds tasks.
type taskHeap []Task

func (h taskHeap) Len() int           { return len(h) }
func (h taskHeap) Less(i, j int) bool { return h[i].runAt.Before(h[j].runAt) }
func (h taskHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *taskHeap) Push(x interface{}) {
	*h = append(*h, x.(Task))
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TimebasedScheduler uses a min-heap to schedule tasks based on time.
type TimebasedScheduler struct {
	tasks        taskHeap
	sleepTrigger chan struct{} // buffer 1
}

func NewTimebasedScheduler() *TimebasedScheduler {
	s := &TimebasedScheduler{}
	s.sleepTrigger = make(chan struct{}, 1)
	heap.Init(&s.tasks)
	return s
}

func (s *TimebasedScheduler) triggerImmediately(task Task) {
	go func() {
		duration := time.Until(task.runAt)
		if duration > 0 {
			time.Sleep(duration)
			task.callback <- struct{}{}
		} else {
			task.callback <- struct{}{}
		}
	}()
}

func (s *TimebasedScheduler) AddTask(task Task) {
	// if task time < 1 s, we need to trigger immediately
	if task.runAt.Before(time.Now().Add(1 * time.Second)) {
		s.triggerImmediately(task)
		return
	}

	// enqueue
	heap.Push(&s.tasks, task)
	select {
	case s.sleepTrigger <- struct{}{}:
	default:
		// if the channel is full, we don't need to do anything
	}
}

func (s *TimebasedScheduler) sleepWithTrigger(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("context done")
		return
	case <-time.After(duration):
		fmt.Println("wake up from time.After", duration)
		return
	case <-s.sleepTrigger:
		fmt.Println("wake up from addTask trigger")
		return
	}
}

func (s *TimebasedScheduler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return
		default:
			if s.tasks.Len() == 0 {
				// default trigger 1 minute
				s.sleepWithTrigger(ctx, time.Minute)
			} else {
				now := time.Now()
				nextTask := s.tasks[0]
				duration := nextTask.runAt.Sub(now)
				if duration <= 0 {
					nextTask.callback <- struct{}{}
					heap.Pop(&s.tasks)
				}
				s.sleepWithTrigger(ctx, duration)
			}
		}
	}
}
