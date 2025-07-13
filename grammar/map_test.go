package grammar

import (
	"fmt"
	"testing"
)

func Map() {
	m := make(map[string]int)
	m["a"] = 0
	m["b"] = 2
	fmt.Println(m)
	if _, ok := m["d"]; !ok {
		fmt.Println("d not found")
	}
	if _, ok := m["a"]; !ok { // zero value doesn't mean not found
		fmt.Println("a not found")
	}

	// read non-existent key will not panic, but zero value is returned
	fmt.Println(m["d"])
	fmt.Println(m["e"])
}

func TestMap(t *testing.T) {
	Map()
}
