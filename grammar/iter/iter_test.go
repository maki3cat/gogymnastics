package grammar

import (
	"fmt"
	"testing"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	for k := range Keys(m) {
		fmt.Println("key: ", k)
	}
}

func TestKeys2(t *testing.T) {
	count := 0
	keys := Keys(map[string]int{"a": 1, "b": 2, "c": 3})
	keys(func(k string) bool {
		count++
		if count == 2 {
			return false
		}
		return true
		// fmt.Println("Got:", k)
		// if k == "b" {
		// 	return false
		// }
		// return true // <- this is the yield function deciding to keep going
	})
}
