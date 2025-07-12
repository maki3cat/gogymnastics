package grammar

import (
	"fmt"
	"iter"
)

func Keys[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[V] {
	return func(yield func(V) bool) {
		// print("yield: ", yield, "\n")
		for k, v := range m {
			fmt.Println("k: ", k, "v: ", v)
			if !yield(v) {
				return
			}
		}
	}
}
