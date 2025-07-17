package medium

import (
	"context"
	"fmt"
)

// use case like wait on exit signal and request chan
func forSelect(ctx context.Context, ch <-chan any) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return
		default:
			select {
			// maki: actually, the for select is a case that this ch can avoid close
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok {
					fmt.Println("channel closed")
					return
				}
				// as if handle this request
				fmt.Println(v)
			}
		}
	}
}
