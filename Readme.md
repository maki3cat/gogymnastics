

This repo including this readme is still Growing :)

Still, hope you like it and star it :star:


## Concurrency Patterns Summarized

please read the code in 
- concurrency (haven't sized up yet)
- concurrency_small
- concurrency_medium


## Principles than Patterns Summarized

1. goroutines should be preemptable/interruptable/able to exit
   - either waiting on `ctx.Done` with only small units in between
   - either the code is non-blocking and exits almost instantenously
2. most of the time when reading from a `chan`, read the ok flag, `data, ok := <- someChan`
3. we can use buffered `chan` to avoid goroutine leakage 
4. if there is only one channel to read from which is sure to be closed, consider `for range` than `for select` in the first place; but if the channel is not closed, this for range may block forever without a chance to read from `ctx.Done`;