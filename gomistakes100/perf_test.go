package gomistakes100

import "testing"

const n = 1_000_000

var global map[int]struct{}

func BenchmarkMapWithoutSize(bench *testing.B) {
	var local map[int]struct{}
	for i := 0; i < bench.N; i++ {
		m := make(map[int]struct{})
		for j := 0; j < n; j++ {
			m[j] = struct{}{}
		}
	}
	global = local
}

func BenchmarkMapWithSize(bench *testing.B) {
	var local map[int]struct{}
	for i := 0; i < bench.N; i++ {
		m := make(map[int]struct{}, n)
		for j := 0; j < n; j++ {
			m[j] = struct{}{}
		}
	}
	global = local
}
