package sort

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	ps := PersonSlice{
		{Name: "Alice", Age: 20},
		{Name: "Bob", Age: 21},
		{Name: "Charlie", Age: 12},
	}

	heap.Init(&ps)
	heap.Push(&ps, Person{Name: "David", Age: 13})

	for ps.Len() > 0 {
		fmt.Println(heap.Pop(&ps).(Person))
	}

	// pop from an empty heap will panic
	assert.Panics(t, func() {
		heap.Pop(&ps)
	})
}
