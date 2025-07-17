package small

import (
	"context"
	"fmt"
)

func groupIteration(ctx context.Context, A <-chan struct{}, B <-chan struct{}, C <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done, exit")
			return
		default:
			tmpA, tmpB, tmpC := A, B, C
			fmt.Println("------------------")
			for range 3 {
				select {
				case <-ctx.Done():
					fmt.Println("context done, exit")
					return
				case <-tmpA:
					fmt.Println("A")
					tmpA = nil
				case <-tmpB:
					fmt.Println("B")
					tmpB = nil
				case <-tmpC:
					fmt.Println("C")
					tmpC = nil
				}
			}
		}
	}
}
