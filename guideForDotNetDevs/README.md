# Go Guide for .NET Developers

Short, sourced notes on where Go differs from C#/.NET. Each doc starts from a C# habit, shows the Go
equivalent (or the lack of one), and cites the source. Examples use the `Employee` type from
`section08_36`.

## Docs

| Doc | C# concept | Go answer in one line |
|---|---|---|
| [01 - Go is not object-oriented](01-go-is-not-oop.md) | classes, inheritance, `this` | Types + methods + interfaces; no classes. |
| [02 - UUIDs and `Guid`](02-uuid-and-guid.md) | `Guid.NewGuid()` | `uuid.New()`; a `[16]byte` value type. |
| [03 - Date-only types](03-date-only-types.md) | `DateOnly` | None in stdlib; use `time.Time` at midnight UTC. |
| [04 - Default parameters](04-default-parameters.md) | optional params, overloads | None; use options structs or functional options. |
| [05 - Value vs. pointer](05-value-vs-pointer.md) | `class` vs `struct` | You choose per type via receiver and return type. |
| [06 - Encapsulation](06-encapsulation-and-private-set.md) | `private set`, `init`, `readonly` | Unexported field + exported getter; privacy is per package. |
| [07 - Where encapsulation is not airtight](07-encapsulation-limits.md) | `private` enforced against everything outside the class | Unexported is per package; zero values, copies, reflection, `encoding/json` and embedding reach around it. |
| [08 - Implicit interfaces](08-implicit-interfaces.md) | `class Employee : IPerson` | No declaration; the compiler checks the method set where the value is used as the interface. |
| [09 - Interface values and embedding](09-interface-values-and-embedding.md) | `IPerson p = employee;`, `IWorker : IPerson` | An interface value is a (type, value) pair; interfaces embed other interfaces; conversions go up implicitly, down by assertion. |
| [10 - Time formatting layouts](10-time-formatting.md) | `"yyyy-MM-dd"` token format | Fixed reference instant `Mon Jan 2 15:04:05 MST 2006`; layout is built from its literal substrings. |
| [11 - Equality and `GetHashCode`](11-equality-and-hashcode.md) | `Equals`/`GetHashCode` overrides, `IEquatable<T>`, `record` | `==` is fixed/structural, not overridable; no hash override; custom equality is a plain `Equal` method by convention. |
| [12 - Generics and type constraints](12-generics-and-type-constraints.md) | `where T : struct`, `where T : IComparable<T>` | Constraints are interfaces; can be a union of type terms (`int \| float32 \| float64`), constraint-only once non-basic. |
| [13 - DDD in Go primer](13-ddd-in-go-primer.md) | entities, value objects, aggregates, domain exceptions, repositories | Package is the boundary; `New` returns `(T, error)`; errors as domain language; embedding is composition, used sparingly. |
| [14 - Functions as parameters](14-functions-as-parameters.md) | `Func<>`, `Action<>`, delegates, lambdas | Function types are first-class; literals are closures; method values replace method groups; funcs compare only to `nil`. |
| [15 - Documenting code](15-documenting-code.md) | `///` XML comments, `<summary>`, `<param>` | Plain `//` doc comments above declarations, full sentences starting with the name; read by `go doc` and pkg.go.dev. |
| [16 - Running benchmarks](16-running-benchmarks.md) | BenchmarkDotNet, `[Benchmark]`, `[Params]`, `[MemoryDiagnoser]` | Built in: `BenchmarkXxx(b *testing.B)` with `b.Loop()`, `go test -bench`, `b.Run` for sizes, `benchstat` to compare runs. |
| [17 - Stepping through code in VS Code](17-debugging-in-vscode.md) | Visual Studio debugger, F9/F10/F11, Locals/Watch/Call Stack | Same keys, but the debugger is Delve driven by the Go extension; Call Stack lists goroutines; `launch.json` `program` points at the package. |
| [18 - Concurrency: goroutines vs. Tasks/Threads](18-concurrency-comparison.md) | `Thread`, `Task`/`Task.Run`, `async`/`await`, `Task.WaitAll` | M:N-scheduled goroutines via `go`; no return value or handle, so `sync.WaitGroup`/channels fill the gap; `main` returning does not wait for goroutines. |
| [19 - Producer/consumer with channels](19-producer-consumer-with-channels.md) | `System.Threading.Channels.Channel<T>`, `Writer.Complete()` | Worker pool: `close()` signals done, `range` drains then stops; only the sender closes, or it's a panic; unclosed channels deadlock. |

## Conventions

- Every claim links to an official source (Go spec, Effective Go, Go FAQ, Code Review Comments,
  Google Go Style Guide, pkg.go.dev, Go issues).
- Claims that could not be confirmed from a source are marked **Unverified**.
- Check the Go spec first for language-behavior claims.
- Last fact-checked against Go 1.27.3 on 2026-11-26 for doc 19 and the `WaitGroup.Go` addition to doc 18
  (Go spec, `sync` package doc, MS Learn `ChannelWriter<T>.Complete`, fetched directly), on 2026-11-25 for
  the rest of doc 18 (Go spec, `sync` package doc, Go FAQ, Go Tour, and MS Learn pages fetched directly), on
  2026-11-15 for docs 16 and 17 (doc 17 from fetched pages only; its `launch.json` was not run), and on
  2026-11-13 for docs 13 to 15 and the 2026-11-13 additions to docs 01, 05 and 06; the rest on 2026-11-11.
  Docs 07, 09, and 19 were split out of 06, 08 (formerly 15 and 14), and 18 respectively without
  re-checking the parent doc's remaining content. Sources were fetched and the code examples were compiled
  or run. Remaining gaps are marked **Unverified** in each doc.
