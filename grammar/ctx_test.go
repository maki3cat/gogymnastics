package grammar

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// not cancelled
func TestSomeFunc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	intput := make(chan bool)
	go func() {
		intput <- false
	}()
	// time.Sleep(100 * time.Millisecond)
	result := someFunc(ctx, intput)
	assert.False(t, result)
	assert.Equal(t, ctx.Err(), nil)
	cancel()
	assert.Equal(t, ctx.Err(), context.Canceled)
}

func TestSomeFunc_Cancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	intput := make(chan bool)
	cancel()
	result := someFunc(ctx, intput)
	assert.False(t, result)
	assert.Equal(t, ctx.Err(), context.Canceled)
}

func TestSomeFunc_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	intput := make(chan bool)
	go func() {
		intput <- true
	}()
	result := someFunc(ctx, intput)
	assert.True(t, result)
	fmt.Println(ctx.Err(), isCtxDone(ctx))
	time.Sleep(200 * time.Millisecond)
	fmt.Println(ctx.Err(), isCtxDone(ctx))
	assert.Equal(t, ctx.Err(), context.DeadlineExceeded)
	//  you call its cancel() function again, nothing changes about the context's state or its error.
	cancel() // when deadline exceeded, cancel will not work
	fmt.Println(ctx.Err(), isCtxDone(ctx))
	assert.Equal(t, ctx.Err(), context.DeadlineExceeded)
}

func TestSomeFunc_Timeout_Cancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	intput := make(chan bool)
	go func() {
		intput <- true
	}()
	result := someFunc(ctx, intput)
	assert.True(t, result)
	fmt.Println(ctx.Err(), isCtxDone(ctx))
	cancel()
	assert.Equal(t, ctx.Err(), context.Canceled)
	fmt.Println(ctx.Err(), isCtxDone(ctx))
	time.Sleep(200 * time.Millisecond)
	// when context is cancelled, deadline doesn't work anymore
	assert.Equal(t, ctx.Err(), context.Canceled)
}
