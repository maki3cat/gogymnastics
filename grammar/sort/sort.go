package sort

type Person struct {
	Name  string
	Age   int
	Group int
	Addr  string
}

// add the len/swap/less to the slice of the datastructure
type PersonSlice []Person

func (ps PersonSlice) Len() int {
	return len(ps)
}

func (ps PersonSlice) Swap(i int, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PersonSlice) Less(i int, j int) bool {
	return ps[i].Age < ps[j].Age
}
