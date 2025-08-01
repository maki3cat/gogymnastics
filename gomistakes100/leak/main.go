package main

// copied from https://github.com/teivah/100-go-mistakes/blob/master/src/03-data-types/26-slice-memory-leak/slice-pointers/main.go
// this is not my work, and I just want to do some experiments

import (
	"fmt"
	"runtime"
)

type Foo struct {
	v []byte
}

func main() {
	foos := make([]Foo, 1_000)
	printAlloc()

	for i := 0; i < len(foos); i++ {
		foos[i] = Foo{
			v: make([]byte, 1024*1024),
		}
	}
	printAlloc()

	two := keepFirstTwoElementsOnly(foos)
	runtime.GC()
	printAlloc()
	runtime.KeepAlive(two)

	two = keepFirstTwoElementsOnlyMarkNil(foos)
	runtime.GC()
	printAlloc()
	runtime.KeepAlive(two)

	// two = keepFirstTwoElementsOnlyCopy(foos)
	// runtime.GC()
}

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
	return foos[:2]
}

func keepFirstTwoElementsOnlyCopy(foos []Foo) []Foo {
	res := make([]Foo, 2)
	copy(res, foos)
	return res
}

func keepFirstTwoElementsOnlyMarkNil(foos []Foo) []Foo {
	for i := 2; i < len(foos); i++ {
		foos[i].v = nil
	}
	return foos[:2]
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}
