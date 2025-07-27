package grammar

import (
	"context"
)

// check if the function returns because of ctx cancel
func someFunc(ctx context.Context, intput chan bool) bool {
	select {
	case <-ctx.Done():
		return false
	case val := <-intput:
		return val
	}
}

func isCtxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
