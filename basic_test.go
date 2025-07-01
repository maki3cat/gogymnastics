package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func fanInOut(n int) {
	res := make(chan struct{}, n)
	for i := range n {
		msg := fmt.Sprintf("this is number %d", i)
		go basicUnit(msg, res)
	}
	for range n {
		<-res
	}
	fmt.Println("finish all units")
}

func fanInOutV2(n int) {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := fmt.Sprintf("this is number %d", i)
			basicUnit(msg)
		}(i)
	}

	wg.Wait() // this is your "fan-in" — blocks until all are done
	fmt.Println("finished all units")
}

func basicUnit(msg string, res chan struct{}) {
	time.Sleep(time.Millisecond * 100)
	fmt.Println(msg)
	res <- struct{}{}
}

func TestGetFanInOut(t *testing.T) {
	fanInOut(10)
}

func TestGetFanInOutV2(t *testing.T) {
	fanInOutV2(10)
}
