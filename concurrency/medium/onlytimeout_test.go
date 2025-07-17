package medium

import (
	"context"
	"testing"
	"time"
)

func TestOnlyTimeout(t *testing.T) {
	waitOnOnlyOne(context.Background(), 1*time.Second)
}
