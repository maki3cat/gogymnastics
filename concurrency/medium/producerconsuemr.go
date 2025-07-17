package medium

import (
	"fmt"
	"sync"
)

type ProducerConsumer struct {
	cv      *sync.Cond // this is a struct than interface
	remains int
}

func NewProducerConsumer(initial int) *ProducerConsumer {
	return &ProducerConsumer{
		cv:      &sync.Cond{L: &sync.Mutex{}}, // key line
		remains: initial,
	}
}

func (pc *ProducerConsumer) Produce(resources int) {
	cv := pc.cv // key line

	cv.L.Lock()         // key line
	defer cv.L.Unlock() // key line
	pc.remains += resources
	cv.Broadcast() // key line
	fmt.Println("produced", resources, "resources, remaining", pc.remains)
}

func (pc *ProducerConsumer) Consumer(resources int) {
	cv := pc.cv
	cv.L.Lock()
	defer cv.L.Unlock()
	for pc.remains < resources { // key line
		fmt.Println("should wait:", pc.remains, "<", resources)
		cv.Wait()
	}
	pc.remains -= resources
	fmt.Println("has consumed", resources, "resources, remaining", pc.remains)
}
