package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestOnlyTimeout(t *testing.T) {
	onlyTimeout(context.Background(), 1*time.Second)
}
