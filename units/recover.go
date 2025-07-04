package main

import (
	"context"
	"fmt"
	"time"
)

// DOESN'T WORK !!
// func panic_with_recover(ctx context.Context) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("recovered")
// 		}
// 	}()
//
// 	go func() {
// 		fmt.Println("sub gorotine starts")
// 		panic("sub goroutine panicks")
// 	}()
// }

// WORKS !!
func panic_with_recover(ctx context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered")
			}
		}()
		time.Sleep(time.Microsecond * 100)
		fmt.Println("sub gorotine starts")
		panic("sub goroutine panicks")
	}()
}

func main() {
	panic_with_recover(context.Background())
	time.Sleep(time.Second)
	fmt.Println("main is still running")
	fmt.Println("main exits")
}
