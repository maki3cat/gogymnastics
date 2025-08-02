package grammar

import (
	"fmt"
	"iter"
)

func Keys[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[K] {
	// maki:
	// the yield generic type K should be the same as the type of the iter.Seq[K]
	return func(yield func(K) bool) {
		// print("yield: ", yield, "\n")
		for k, v := range m {
			fmt.Println("k: ", k, "v: ", v)
			if !yield(k) {
				return
			}
		}
	}
}
