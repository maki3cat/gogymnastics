package main

import (
	"fmt"
)

var example1 = [1024 * 1024]byte{}
var example2 = [1024 * 1024]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func main() {
	size := 100 // hand-change it
	fmt.Printf("size: %d\n", size)

	localData1 := []byte{}
	for range size {
		localData1 = append(localData1, byte(255))
	}
	localData2 := [1024]byte{}

	fmt.Printf("localData1 addr: %p\n", &localData1)
	fmt.Printf("localData2 addr: %p\n", &localData2)
	fmt.Printf("example1 addr: %p\n", &example1)
	fmt.Printf("example2 addr: %p\n", &example2)
}

//ASLR impacts virtual address
// === RUN   TestStaticAlloc
// size: 30
// localData1 addr: 0x140000a4120
// localData2 addr: 0x140000aa400
// globalData addr: 0x100366880
// === RUN   TestStaticAlloc
// size: 10
// localData1 addr: 0x1400000c138
// localData2 addr: 0x1400007e400
// globalData addr: 0x10079a880
