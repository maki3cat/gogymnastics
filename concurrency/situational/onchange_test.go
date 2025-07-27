package situational

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnChange(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	raftState := NewRaftState()
	raftState.Start(ctx)
	raftState.SetCommitIdx(10)
	time.Sleep(10 * time.Millisecond)
	cancel()
	raftState.SetCommitIdx(20)
	assert.Equal(t, raftState.GetAppliedIdx(), 10)
	assert.Equal(t, raftState.GetCommittedIdx(), 20)
}


func TestOnChange_MultipleCommits(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	raftState := NewRaftState()
	raftState.Start(ctx)

	raftState.SetCommitIdx(5)
	time.Sleep(5 * time.Millisecond)
	raftState.SetCommitIdx(8)
	time.Sleep(5 * time.Millisecond)
	raftState.SetCommitIdx(12)
	time.Sleep(10 * time.Millisecond)

	cancel()
	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, raftState.GetAppliedIdx(), 12)
	assert.Equal(t, raftState.GetCommittedIdx(), 12)
}

func TestOnChange_NoCommit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	raftState := NewRaftState()
	raftState.Start(ctx)

	time.Sleep(10 * time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, raftState.GetAppliedIdx(), 0)
	assert.Equal(t, raftState.GetCommittedIdx(), 0)
}

func TestOnChange_CommitAfterStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	raftState := NewRaftState()
	raftState.Start(ctx)

	raftState.SetCommitIdx(3)
	time.Sleep(5 * time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond)
	raftState.SetCommitIdx(7)
	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, raftState.GetAppliedIdx(), 3)
	assert.Equal(t, raftState.GetCommittedIdx(), 7)
}

