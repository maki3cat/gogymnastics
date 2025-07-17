

This repo including this readme is still Growing :)

Still, hope you like it and star it :star:

## Concurrency

### Basic Primitives

- chan: rendezvous chan, buffered chan; (share by communication)
- sync: cond, lock, pool/once, etc; (communicate by sharing)
- sync/atomic: relatively lowerlevel primitives with atomic Swap

## Patterns Summarized

please read the code in
- concurrency (haven't sized up yet)
- concurrency_small
- concurrency_medium


## Principles

1. whenever using goroutins: check goroutine leakage / goroutines should be preemptable/interruptable/able to exit
   - either waiting on `ctx.Done` with only small units in between
   - either the code is non-blocking and exits almost instantenously
   - (*usually channels that are not closed don't cause the memory leak of the chan, but may trigger goroutine leakage if there is some goroutine only exiting on the channel's close signal)

2. most of the time when reading from a `chan`, read the ok flag, `data, ok := <- someChan`
3. we can use buffered `chan` to avoid goroutine leakage 
4. if there is only one channel to read from which is sure to be closed, consider `for range` than `for select` in the first place; but if the channel is not closed, this for range may block forever without a chance to read from `ctx.Done`;

5. `cancel()` should "always" be called when the caller generated the `cxt, cancel` exits
6. very very rare problem solving requires the leverage of `tryLock/tryUnlock`
