package grammar

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
			if val := recover(); val != nil {
				fmt.Println(fmt.Sprintf("recovered and got value %v", val))
			}
		}()
		time.Sleep(time.Microsecond * 100)
		fmt.Println("sub gorotine starts")
		panic(19983420)
	}()
}

// func main() {
// 	panic_with_recover(context.Background())
// 	time.Sleep(time.Second)
// 	fmt.Println("main is still running")
// 	fmt.Println("main exits")
// }
