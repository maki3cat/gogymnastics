package shortcircuit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTarget(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	res := target(ctx, []Point{{0, 0}, {1, 1}, {2, 2}})
	t.Logf("res: %v", res)
	assert.True(t, res)
}

func TestTarget_Fail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	res := target(ctx, []Point{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}, {10, 10}})
	t.Logf("res: %v", res)
	assert.False(t, res)
}

func TestTarget_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	res := target(ctx, []Point{{0, 0}, {1, 1}, {2, 2}})
	t.Logf("res: %v", res)
	assert.False(t, res)
	assert.Equal(t, ctx.Err(), context.DeadlineExceeded)
}

func TestTarget_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res := target(ctx, []Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}, {10, 10}})
	t.Logf("res: %v", res)
	assert.False(t, res)
	assert.Equal(t, ctx.Err(), context.Canceled)
}
