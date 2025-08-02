package sort

// heap is on top of sort

// add push and pop to the PersonSlice
// push and pop are write operations, so we need to use a pointer receiver
func (ps *PersonSlice) Push(x any) {

	// maki: we need to cast the any to Person
	// the heap interface of push and pop uses any
	*ps = append(*ps, x.(Person))
}

func (ps *PersonSlice) Pop() any {
	old := *ps
	n := len(old)
	x := old[n-1]
	*ps = old[:n-1]
	return x
}
