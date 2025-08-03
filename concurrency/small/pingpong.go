package small

import (
	"fmt"
	"time"
)

type Ball struct{ hits int }

func ppplayer(name string, table chan *Ball) {
	for {
		ball, ok := <-table
		if !ok {
			fmt.Println("table closed")
			return
		}
		ball.hits += 1
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
