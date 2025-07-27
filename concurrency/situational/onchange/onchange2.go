package onchange

import (
	"context"
	"sync"
	"time"
)

// in this version, we use channel to notify the apply worker to wake up
type RaftState2 struct {
	appliedIdx    int // inclusive, data index starting from 1
	committedIdx  int // inclusive, data index starting from 1
	lock          *sync.Mutex
	updatedSignal chan int
}

func NewRaftState2() *RaftState2 {
	return &RaftState2{
		appliedIdx:    0,
		committedIdx:  0,
		lock:          &sync.Mutex{},
		updatedSignal: make(chan int, 1), // we only need to buffer one signal, and batch update
	}
}

func (s *RaftState2) GetAppliedIdx() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.appliedIdx
}

func (s *RaftState2) GetCommittedIdx() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.committedIdx
}

func (s *RaftState2) SetCommitIdx(idx int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.committedIdx = idx
	// non-blockingly send a signal, especially useful when we only need one signal and do batching
	// use the signal only for notificaiton than the real idx
	select {
	case s.updatedSignal <- idx:
	default:
	}
}

// should be called in a separate goroutine
func (s *RaftState2) StartApplyWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// todo: cleanup
			return
		default:
			select {
			case <-ctx.Done():
				// todo: cleanup
				return
			case <-s.updatedSignal:
				// simpliy the apply as commited idx
				time.Sleep(500 * time.Millisecond)
				s.appliedIdx = s.committedIdx
			}
		}
	}
}
