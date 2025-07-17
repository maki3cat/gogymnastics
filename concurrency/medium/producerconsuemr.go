package medium

import "sync"

type ProducerConsumer struct {
	cv      *sync.Cond // this is a struct than interface
	remains int
}

func NewProducerConsumer(initial int) *ProducerConsumer {
	return &ProducerConsumer{
		cv:      &sync.Cond{L: &sync.Mutex{}},
		remains: initial,
	}
}

func (pc *ProducerConsumer) produce(resources int) {
	cv := pc.cv

	cv.L.Lock()
	defer cv.L.Unlock()
	pc.remains += resources
	cv.Broadcast()
}

func (pc *ProducerConsumer) consumer(resources int) {
	cv := pc.cv

	cv.L.Lock()
	defer cv.L.Unlock()
	for pc.remains <= resources {
		cv.Wait()
	}
	pc.remains -= resources
}
