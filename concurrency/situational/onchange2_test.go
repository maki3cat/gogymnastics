package situational

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnChange2(t *testing.T) {
	raftState := NewRaftState2()
	ctx, cancel := context.WithCancel(context.Background())
	go raftState.StartApplyWorker(ctx)
	raftState.SetCommitIdx(10)
	raftState.SetCommitIdx(20)
	time.Sleep(1 * time.Second)
	assert.Equal(t, raftState.GetAppliedIdx(), 20)
	cancel()
	raftState.SetCommitIdx(30)
	assert.Equal(t, raftState.GetAppliedIdx(), 20)
	assert.Equal(t, raftState.GetCommittedIdx(), 30)
}

func TestOnChange2_MultipleCommits(t *testing.T) {
	raftState := NewRaftState2()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go raftState.StartApplyWorker(ctx)

	raftState.SetCommitIdx(5)
	time.Sleep(100 * time.Millisecond)
	raftState.SetCommitIdx(8)
	time.Sleep(100 * time.Millisecond)
	raftState.SetCommitIdx(12)
	time.Sleep(600 * time.Millisecond) // ensure apply worker has time to process

	assert.Equal(t, 12, raftState.GetAppliedIdx())
	assert.Equal(t, 12, raftState.GetCommittedIdx())
}

func TestOnChange2_NoCommit(t *testing.T) {
	raftState := NewRaftState2()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go raftState.StartApplyWorker(ctx)

	time.Sleep(600 * time.Millisecond)
	assert.Equal(t, 0, raftState.GetAppliedIdx())
	assert.Equal(t, 0, raftState.GetCommittedIdx())
}

func TestOnChange2_CommitAfterStop(t *testing.T) {
	raftState := NewRaftState2()
	ctx, cancel := context.WithCancel(context.Background())
	go raftState.StartApplyWorker(ctx)

	raftState.SetCommitIdx(3)
	time.Sleep(600 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
	raftState.SetCommitIdx(7)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 3, raftState.GetAppliedIdx())
	assert.Equal(t, 7, raftState.GetCommittedIdx())
}

func TestOnChange2_BatchedSignal(t *testing.T) {
	raftState := NewRaftState2()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go raftState.StartApplyWorker(ctx)

	// Rapidly set commit index multiple times before worker wakes up
	raftState.SetCommitIdx(1)
	raftState.SetCommitIdx(2)
	raftState.SetCommitIdx(3)
	raftState.SetCommitIdx(4)
	raftState.SetCommitIdx(5)
	time.Sleep(600 * time.Millisecond)

	assert.Equal(t, 5, raftState.GetAppliedIdx())
	assert.Equal(t, 5, raftState.GetCommittedIdx())
}

