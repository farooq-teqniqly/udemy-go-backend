# 19 - Producer/consumer with channels

**C#:** a bounded queue between producer(s) and consumer(s) is `System.Threading.Channels.Channel<T>` —
`Writer.Complete()` signals "no more items," `Reader.ReadAllAsync` drains and stops.

**Go:** the same shape with two channels and `close`. Example comes from `worker_pool_example`, which fans a
job channel out to a pool of workers and fans their results back into a second channel — see
[18 - Concurrency](18-concurrency-comparison.md) for the underlying goroutine/`WaitGroup` primitives.

```go
unprocessed := make(chan job, maxJobs)
processed := make(chan job, maxJobs)

var wg sync.WaitGroup
for range maxJobs {
	wg.Go(func() { processJob(unprocessed, processed) }) // worker: for j := range unprocessed { ... }
}

for id := range maxJobs {
	unprocessed <- id // producer
}
close(unprocessed) // no more jobs; workers' range loop drains the buffer then exits

go func() {
	wg.Wait()        // all workers returned
	close(processed) // safe to close now — no sender left
}()

for j := range processed { // consumer; exits once closed and drained
	// ...
}
```

## 1. `range` over a closed channel drains, then stops

The spec says a receive on a closed channel "can always proceed immediately, yielding the element type's
zero value after any previously sent values have been received" (<https://go.dev/ref/spec#Receive_operator>).
So each worker's `for j := range unprocessed` exits cleanly once `close(unprocessed)` runs and the buffer
empties — no sentinel value needed, and no reader-side change required to notice the channel is done.

## 2. Only the sender closes, and only once

"Sending to or closing a closed channel causes a run-time panic. Closing the nil channel also causes a
run-time panic" (<https://go.dev/ref/spec#Close>). That rule shapes the code above: `processed` is closed
from a dedicated goroutine gated on `wg.Wait()`, not from `main` directly before its own `range processed`.
Calling `wg.Wait()` inline in `main` would block it there, before it ever starts draining `processed`, and a
full buffered channel with no reader can stall every worker trying to send into it. Omit the close of either
channel and both `range` loops block forever: the runtime detects no goroutine can still make progress and
exits with `fatal error: all goroutines are asleep - deadlock!` instead of hanging silently — see
[18 - Concurrency](18-concurrency-comparison.md) for the same deadlock mechanism, and `section07_31`-style
channel examples for the single-channel case.

## 3. C# analogue: `Channel<T>`

`writer.Complete()` plays the role of `close(ch)`: "Mark the channel as being complete, meaning no more
items will be written to it" (<https://learn.microsoft.com/en-us/dotnet/api/system.threading.channels.channelwriter-1.complete>).
A reader drains the remaining items before observing completion — the same drain-then-stop order as `range`
over a closed Go channel. The difference is enforcement: nothing stops a second, unrelated goroutine from
calling `close` on a Go channel it doesn't own, since ownership is a convention, not a type-level restriction
— whereas `Channel<T>` splits `Writer`/`Reader` into separate types, so only code holding the `ChannelWriter<T>`
reference can call `Complete()`.

## Sources

- Go spec, Receive operator (closed-channel receive yields zero value): <https://go.dev/ref/spec#Receive_operator>
- Go spec, Close (panic on closing a closed or nil channel): <https://go.dev/ref/spec#Close>
- `sync.WaitGroup.Go` (added in Go 1.25): <https://pkg.go.dev/sync#WaitGroup>
- MS Learn, `ChannelWriter<T>.Complete`: <https://learn.microsoft.com/en-us/dotnet/api/system.threading.channels.channelwriter-1.complete>

Fact-checked 2026-11-26 against the Go spec (Receive operator, Close) and `ChannelWriter<T>.Complete`,
fetched directly. **Unverified:** exact Learn wording for `ChannelReader<T>`'s drain-then-stop behavior (not
fetched — the completion/draining order is stated from general `Channel<T>` documentation, not a quoted
`ChannelReader<T>` page); the Go Tour's own "Range and Close" wording (that page is a client-rendered SPA
that could not be fetched directly — the same substance is confirmed via the spec quotes above instead).
