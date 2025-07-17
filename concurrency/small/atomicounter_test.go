package small

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	ans := AtomicCounter(context.TODO(), 10, 100)
	assert.Equal(t, ans, 10*100, "result should be")
}
