package logiteration

import (
	"fmt"
	"sync"
)

// Each subsystem is represented by a function in your code: printFirst(), printSecond(), and printThird().
// When these functions are called, they will print a log line. These functions may be called concurrently.
// You must ensure that the log lines are always processed (i.e., printed) in the order: first A, then B, and then C.
func NewLogPrinter() *LogPrinter {
	return &LogPrinter{
		cv:    sync.NewCond(&sync.Mutex{}),
		order: 0,
	}
}

type LogPrinter struct {
	cv    *sync.Cond
	order int // 0: first, 1: second, 2: third
}

// the 3 system interfaces
// todo: add error handling
// todo: add context timeout
func (l *LogPrinter) printFirst(msg string) {
	l.cv.L.Lock()
	defer l.cv.L.Unlock()

	for l.order != 0 {
		l.cv.Wait()
	}
	fmt.Printf("A: %s\n", msg)
	l.order = 1
	l.cv.Broadcast()
}

func (l *LogPrinter) printSecond(msg string) {
	l.cv.L.Lock()
	defer l.cv.L.Unlock()

	for l.order != 1 {
		l.cv.Wait()
	}
	fmt.Printf("B: %s\n", msg)
	l.order = 2
	l.cv.Broadcast()
}

func (l *LogPrinter) printThird(msg string) {
	l.cv.L.Lock()
	defer l.cv.L.Unlock()

	for l.order != 2 {
		l.cv.Wait()
	}
	fmt.Printf("C: %s\n", msg)
	l.order = 0
	l.cv.Broadcast()
}
