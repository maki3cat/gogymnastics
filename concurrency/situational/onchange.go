package situational

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// What we are doing
// general problem: when something is changed, some worker should come to work;
// the concrete problem: when raft commitIdx is changed, newly arrived logs should be applied to the state machine

// Solution 1: using conditional variable
// the problem here is that we cannot just use context with conditional variable wait
// in the same with as we do with channel
type RaftState struct {
	committedIdx int // inclusive, data index starting from 1
	appliedIdx   int // inclusive, data index starting from 1
	stopped      *atomic.Bool
	cv           *sync.Cond
}

func NewRaftState() *RaftState {
	return &RaftState{
		committedIdx: 0,
		appliedIdx:   0,
		stopped:      &atomic.Bool{},
		cv:           sync.NewCond(&sync.Mutex{}),
	}
}

func (s *RaftState) GetAppliedIdx() int {
	s.cv.L.Lock()
	defer s.cv.L.Unlock()
	return s.appliedIdx
}

func (s *RaftState) GetCommittedIdx() int {
	s.cv.L.Lock()
	defer s.cv.L.Unlock()
	return s.committedIdx
}

func (s *RaftState) SetCommitIdx(idx int) {
	s.cv.L.Lock()
	defer s.cv.L.Unlock()
	if idx <= s.committedIdx {
		return
	}
	s.committedIdx = idx
	s.cv.Broadcast()
}

func (s *RaftState) mockApply(startingIdx int, toIdx int) {
	fmt.Printf("mockApply: %d -> %d\n", startingIdx, toIdx)
}

func (s *RaftState) startApplyWorker() {
	for {
		s.cv.L.Lock()
		for s.appliedIdx == s.committedIdx && !s.stopped.Load() {
			s.cv.Wait()
		}
		if s.stopped.Load() {
			s.cv.L.Unlock()
			break
		} else {
			s.mockApply(s.appliedIdx, s.committedIdx)
			s.appliedIdx = s.committedIdx
			s.cv.L.Unlock()
		}
	}
}

func (s *RaftState) stopApplyWorker(ctx context.Context) {
	<-ctx.Done()
	s.stopped.Store(true)
	s.cv.Broadcast()
}

func (s *RaftState) Start(ctx context.Context) {
	go s.startApplyWorker()
	go s.stopApplyWorker(ctx)
}
