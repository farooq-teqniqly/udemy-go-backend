# 18 - Concurrency: goroutines vs. Tasks/Threads

**C#:** concurrency is `Thread` (an OS-backed managed thread) or, almost always in practice, `Task`/`Task<T>`
run on the thread pool via `Task.Run`, coordinated with `async`/`await`, `Task.WaitAll`, or `CountdownEvent`.

**Go:** concurrency is a `go` statement starting a goroutine, coordinated with `sync.WaitGroup` or channels.
Examples come from `section12_53`, which launches 10 goroutines and waits for them with a `WaitGroup`.

## 1. Goroutine vs. `Thread` vs. `Task`

A `go` statement "starts the execution of a function call as an independent concurrent thread of control,
or goroutine, within the same address space... program execution does not wait for the invoked function to
complete." If the function has return values, "they are discarded when the function completes."

Goroutines are not OS threads. The Go FAQ: the runtime "multiplex[es] independently executing
functions... onto a set of threads," and a goroutine costs "little more than the allocation of stack
space" — a few kilobytes, growing and shrinking as needed, so "hundreds of thousands" can exist in one
process. This is Go's M:N scheduler: many goroutines over few OS threads.

A C# `Thread` is "the basic unit to which an operating system allocates processor time" — a real OS thread,
1:1. `Task` (via `Task.Run` or the thread pool) reuses pooled OS threads rather than creating one per task,
but each running `Task` still occupies a real OS thread while active — there is no stack-shrinking
equivalent to a goroutine.

```go
go printRandomNumber(5*time.Second, &wg) // fire-and-forget; no handle, no return value
```

```csharp
Task.Run(() => PrintRandomNumber(TimeSpan.FromSeconds(5))); // returns a Task — a handle you can await
```

## 2. `sync.WaitGroup`: `Add`/`Done`/`Wait`

`Add adds delta... to the WaitGroup task counter`. `Wait` "blocks until the WaitGroup task counter is
zero." The documented footgun: "calls with a positive delta that occur when the counter is zero must happen
before a Wait... Typically this means the calls to Add should execute before the statement creating the
goroutine" — i.e. call `Add` in the loop that starts the goroutines, not inside the goroutine itself, or
`Wait` can race past a counter that is still zero. `section12_53` does this correctly:

```go
wg.Add(workerCount) // before the loop
for range workerCount {
	go printRandomNumber(5*time.Second, &wg)
}
wg.Wait()
```

Each goroutine calls `defer wg.Done()` (equivalent to `Add(-1)`) once its work finishes.

The nearest C# type is `CountdownEvent`: `Signal()`/`AddCount` adjust `CurrentCount`, `Wait()` blocks until
it reaches zero — structurally the same Add/Done/Wait shape. In practice, though, C# code fanning out a
batch of `Task`s uses `Task.WaitAll(tasks)` instead, which blocks until every `Task` in the array completes
and wraps any exceptions in `AggregateException`.

## 3. Why Go needs `WaitGroup` but C# usually doesn't

A goroutine returns nothing — the spec confirms discarded return values — so there is no handle to wait on
unless you build one (`WaitGroup`, or a channel). A C# `Task`/`Task<T>` *is* that handle: it's awaitable and
(for `Task<T>`) carries the result, so `Task.WaitAll`/`Task.WhenAll` can aggregate a collection of existing
handles without any separate counter object. `WaitGroup` exists in Go precisely to fill the gap left by `go`
returning nothing.

## 4. No implicit join: `main` returning kills goroutines

Go spec, program execution: "When that function invocation returns, the program exits. It does not wait
for other (non-main) goroutines to complete." Drop `wg.Wait()` from `section12_53` and `main` can print
"Exiting..." and the process can end before any `printRandomNumber` goroutine wakes from its 5-second sleep
— nothing forces them to finish.

C# differs by thread kind, not by `Task` vs. `Thread`: "Background threads are identical to foreground
threads, except that background threads do not prevent a process from terminating. Once all foreground
threads belonging to a process have terminated, the... runtime ends the process. Any remaining background
threads are stopped and do not complete." Thread-pool threads (so most `Task` work) are background threads
by default, so an un-awaited `Task` can be cut off the same way a goroutine can. But an explicitly
constructed `Thread` defaults to `IsBackground = false` (foreground) and *will* keep the process alive until
it finishes — a case with no Go equivalent, since no goroutine can keep `main` from returning.

## 5. Channels, briefly

Channels are Go's other coordination primitive: "a typed conduit through which you can send and receive
values... By default, sends and receives block until the other side is ready. This allows goroutines to
synchronize without explicit locks or condition variables." The nearest C# types are `System.Threading.
Channels.Channel<T>` (async producer/consumer) and the older `BlockingCollection<T>` (blocking, no async
support).

## 6. `async`/`await` vs. `go`

C#'s `async`/`await` is a compiler transform: the Task-based asynchronous programming docs describe the
compiler handling "the hard parts," tracking what resumes a suspended method at each `await`. This creates
two kinds of methods — `async` and not — that must be called differently ("function coloring").

Go has no such split. A `go` statement is not `Task.Run` (a thread-pool work item) and not an async method
call — Go has one execution model, and the runtime transparently multiplexes blocking and non-blocking code
onto OS threads underneath every goroutine. There is no `await` keyword and no second, differently-typed
version of a function for calling it concurrently: `printRandomNumber` is called the same way whether
behind a `go` or not.

## Sources

- Go spec, Go statements (goroutine creation, discarded return values): <https://go.dev/ref/spec#Go_statements>
- Go spec, Program execution (`main` returning exits without waiting): <https://go.dev/ref/spec#Program_execution>
- Effective Go, Goroutines: <https://go.dev/doc/effective_go#goroutines>
- Go FAQ, why goroutines instead of threads (M:N scheduling, stack size): <https://go.dev/doc/faq#goroutines>
- `sync.WaitGroup` (`Add`, `Done`, `Wait`, the Add-before-go rule): <https://pkg.go.dev/sync#WaitGroup>
- Go Tour, channels: <https://go.dev/tour/concurrency/2>
- MS Learn, Threads and threading (`Thread`, thread pool): <https://learn.microsoft.com/en-us/dotnet/standard/threading/threads-and-threading>
- MS Learn, Foreground and background threads: <https://learn.microsoft.com/en-us/dotnet/standard/threading/foreground-and-background-threads>
- MS Learn, `Task.WaitAll`: <https://learn.microsoft.com/en-us/dotnet/api/system.threading.tasks.task.waitall>
- MS Learn, `CountdownEvent`: <https://learn.microsoft.com/en-us/dotnet/standard/threading/countdownevent>
- MS Learn, Task-based asynchronous programming: <https://learn.microsoft.com/dotnet/csharp/asynchronous-programming/task-asynchronous-programming-model>

Fact-checked 2026-11-25 against the sources above (Go spec sections, `sync` package doc, Go FAQ, Go Tour,
and the MS Learn pages fetched directly). **Unverified:** the exact wording of `Task.WaitAll` blocking the
calling thread (confirmed in Learn Q&A/overview pages, not the API reference page itself); that the C#
compiler generates an `IAsyncStateMachine` struct (sourced from a Microsoft engineering blog post, not an
official Learn/spec page); exact Learn wording for `System.Threading.Channels`/`BlockingCollection` (types
and behavior are correct, but no Learn quote was fetched for them).
