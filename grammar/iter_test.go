package grammar

import (
	"fmt"
	"testing"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	for k := range Keys(m) {
		fmt.Println(k)
	}
}
