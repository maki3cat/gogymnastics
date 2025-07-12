package main

import "context"

func TeeChn(ctx context.Context, ch <-chan any) (<-chan any, <-chan any) {
	out1, out2 := make(chan any), make(chan any)
	go func() {
		defer close(out1)
		defer close(out2)
		for v := range ch {
			tmp1, tmp2 := out1, out2
			for range 2 {
				select {
				case <-ctx.Done():
					return
				default:
					select {
					case <-ctx.Done():
						return
					case tmp1 <- v:
						tmp1 = nil
					case tmp2 <- v:
						tmp2 = nil
					}
				}
			}
		}
	}()
	return out1, out2
}
