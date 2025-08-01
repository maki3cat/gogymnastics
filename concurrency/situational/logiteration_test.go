package logiteration

import (
	"testing"
	"time"
)

func TestLogIteration(t *testing.T) {
	lp := NewLogPrinter()
	startFun := func(num int) {
		for range 3 {
			switch num {
			case 0:
				go lp.printFirst("hello world")
			case 1:
				go lp.printSecond("hello world")
			case 2:
				go lp.printThird("hello world")
			}
		}
	}
	go startFun(2)
	go startFun(0)
	go startFun(1)

	time.Sleep(1 * time.Second)
}
