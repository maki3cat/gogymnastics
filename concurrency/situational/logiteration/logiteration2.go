package logiteration

import (
	"context"
	"fmt"
	"time"
)

type LogPrinter2 struct {
	logCh1   chan string
	logCh2   chan string
	logCh3   chan string
	waitTime time.Duration // if there is no log after waitTime, we move to next system
}

func NewLogPrinter2(bufferSize int) *LogPrinter2 {
	return &LogPrinter2{
		logCh1:   make(chan string, bufferSize),
		logCh2:   make(chan string, bufferSize),
		logCh3:   make(chan string, bufferSize),
		waitTime: 500 * time.Millisecond,
	}
}

func (l *LogPrinter2) printFirst(msg string) {
	l.logCh1 <- msg
}

func (l *LogPrinter2) printSecond(msg string) {
	l.logCh2 <- msg
}

func (l *LogPrinter2) printThird(msg string) {
	l.logCh3 <- msg
}

// should be call asynchronously
func (l *LogPrinter2) start(ctx context.Context) {
	ticker := time.NewTicker(l.waitTime)
	defer ticker.Stop()

	for {
		// first system
		select {
		case <-ctx.Done():
			fmt.Println("context done, stop log iteration")
			return
		case msg := <-l.logCh1:
			fmt.Println("1st system message:", msg)
		case <-ticker.C:
			fmt.Println("1st system timeout, move to 2nd system")
		}
		ticker.Reset(l.waitTime)

		// second system
		select {
		case <-ctx.Done():
			fmt.Println("context done, stop log iteration")
			return
		case msg := <-l.logCh2:
			fmt.Println("2nd system message:", msg)
		case <-ticker.C:
			fmt.Println("2nd system timeout, move to 3rd system")
		}
		ticker.Reset(l.waitTime)

		// third system
		select {
		case <-ctx.Done():
			fmt.Println("context done, stop log iteration")
			return
		case msg := <-l.logCh3:
			fmt.Println("3rd system message:", msg)
		case <-ticker.C:
			fmt.Println("3rd system timeout, move to 1st system")
		}
		ticker.Reset(l.waitTime)
	}
}
