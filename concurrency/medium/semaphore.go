package medium

type Sem struct {
	n int
	c chan struct{}
}

func NewBinarySemaphore() *Sem {
	sem := &Sem{
		n: 0,
		c: make(chan struct{}, 1),
	}
	sem.c <- struct{}{}
	return sem
}

func NewSem(n int) *Sem {
	sem := &Sem{
		n: n,
		c: make(chan struct{}, n),
	}
	for range n {
		sem.c <- struct{}{}
	}
	return sem
}

func (n *Sem) Acquire() {
	<-n.c
}

func (n *Sem) Release() {
	n.c <- struct{}{}
}
