package grammar

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ps = []Person{
	{Name: "a", Age: 10, Group: 1, Addr: ""},
	{Name: "b", Age: 5, Group: 2, Addr: ""},
	{Name: "c", Age: 7, Group: 3, Addr: ""},
}

// sort.Slice with a less function which is easier
// more of a once-off method
func TestSortSlice2(t *testing.T) {
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Group < ps[j].Group
	})
	assert.True(t, ps[0].Group == 1)
	fmt.Println(ps)
}

// sort.Sort with the Sort interface (less/len/swap)
// reusable method
// sort.Sort is generic method where sort.Slice only works on a slice
func TestSortSlice(t *testing.T) {
	// warning:
	// we need to change the slice of person to PersonSlice in order to use the sort.Sort function
	sort.Sort(PersonSlice(ps))
	assert.True(t, ps[0].Age == 5)
	fmt.Println(ps)
}
