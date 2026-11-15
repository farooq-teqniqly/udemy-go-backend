# 16 - Running benchmarks in Go

**C#:** you add the BenchmarkDotNet package, mark methods with `[Benchmark]`, add `[MemoryDiagnoser]` for
allocations and `[Params(10, 100, 1000)]` for sizes, build in Release, and call `BenchmarkRunner.Run<T>()`.
It refuses to run unless built in Release.

**Go:** benchmarking is built into the `testing` package and `go test`. There is nothing to install, no
attributes and no separate runner. The `testing` page documents no warmup phase and no statistics; you
compare runs yourself (`benchstat`, below). That absence is **Unverified**. Examples come from
`section10_44`, which compares ways to build `"1,2,3,...,n"`.

## 1. Writing one

A benchmark is a function in a `_test.go` file named `BenchmarkXxx` that takes `*testing.B`. Since Go 1.24
the body is a `b.Loop()` loop:

```go
var sink string // package-level, so the result is "used"

func BenchmarkStringsJoin(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = stringsJoin(1000)
	}
}
```

- `b.Loop()` returns `true` "as long as the benchmark should continue running", so you do not pick the
  iteration count. Older code writes `for i := 0; i < b.N; i++`, where the framework picks `b.N`. Use one
  form, not both. The docs say the condition must be written exactly `b.Loop()` for the keep-alive below.
- For `b.Loop`, "arguments to and results from function calls and assigned variables within the loop are
  kept alive, preventing the compiler from fully optimizing away the loop body". The `sink` variable is the
  same trick for the old `b.N` form.
- Setup written before the loop is not timed: `b.Loop` resets the timer on its first call and stops it when
  it returns `false`, so cleanup is not timed either. With `b.N`, call `b.ResetTimer()` after setup.
  `b.StopTimer()` and `b.StartTimer()` exclude work inside the loop.

## 2. Running it

```
go -C section10_44 test -run '^$' -bench . -benchmem
```

| Flag | Effect |
|---|---|
| `-bench regexp` | Run benchmarks whose name matches; `.` is all. Without it none run. |
| `-run '^$'` | A regexp that matches no test names, so only benchmarks run. It is the idiom on benchstat's page. |
| `-benchmem` | Add `B/op` and `allocs/op` (like `[MemoryDiagnoser]`); `b.ReportAllocs()` does it per benchmark. |
| `-benchtime 5s` or `100x` | Run for a duration or an exact iteration count. Default is 1s. |
| `-count 10` | Repeat each benchmark, giving several samples to compare. |
| `-cpu 1,4` | Run with each `GOMAXPROCS` value. |
| `-cpuprofile`, `-memprofile` | Write profiles for `go tool pprof`. |

## 3. Reading the output

```
BenchmarkStringsJoin/n=1000-20    55624    21666 ns/op    6981 B/op    902 allocs/op
```

From left: name with a `-20` suffix, iterations run, time per operation, bytes allocated per operation,
allocations per operation. The `testing` page's sample line (`BenchmarkRandInt-8 68453040 17.8 ns/op`)
explains iterations and `ns/op`. That the suffix is `GOMAXPROCS` (20 logical CPUs here) and the exact
meaning of `B/op` and `allocs/op` are **Unverified**; `-benchmem` is documented only as "memory allocation
statistics". Compare `allocs/op` as well as time: for strings, allocations were the story.

## 4. Sizes: `b.Run` is `[Params]`

`b.Run(name, func(b *testing.B))` runs a sub-benchmark. Its name is appended to the parent's after a `/`
(the testing docs only imply this), and `-bench` splits its regexp on unbracketed slashes, with each part
matching one level. Parents that might match run with `b.N=1` to discover their sub-benchmarks:

```go
func BenchmarkJoinNumbers(b *testing.B) {
	for _, v := range variants {
		b.Run(v.name, func(b *testing.B) {
			for _, size := range sizes {
				b.Run("n="+strconv.Itoa(size), func(b *testing.B) {
					for b.Loop() {
						sink = v.fn(size)
					}
				})
			}
		})
	}
}
```

```
go -C section10_44 test -run '^$' -bench 'JoinNumbers/Builder.*/n=1000$' -benchmem
```

## 5. Comparing runs

One run is noisy; in `section10_44` the spread reached 30 to 60% at some sizes. Run `-count 10`, save the
output to a file, change the code, save again, then compare with `benchstat` (a separate tool, installed
with `go install golang.org/x/perf/cmd/benchstat@latest`; that command is **Unverified**, since its page
does not show it):

```
go -C section10_44 test -run '^$' -bench . -benchmem -count 10 > old.txt
benchstat old.txt new.txt
```

It "computes statistical summaries and A/B comparisons", shows a percentage change per benchmark, and
prints `~` when the difference is not statistically significant.

## 6. Pitfalls

- **A benchmark can mislead through compiler behavior.** In `section10_44`, turning the item count from a
  constant into a function parameter made `strings.Join` allocate more (2 allocs instead of 1 at n=10) and
  take about 60% longer. The likely cause is the slice moving from stack to heap (escape analysis); this was
  **Unverified**. Benchmark with realistic, non-constant inputs.
- **Keep the baseline.** Adding the naive `s +=` version showed it 700 times slower than a `Builder`
  at 10,000 items. Without the baseline the other differences looked large.
- **Do not trust one machine's numbers as universal**; they depend on CPU, OS and Go version.
- **Closures capture the loop variable.** `for _, v := range variants` with `v.fn` in a closure is safe
  from Go 1.22, where "each iteration of the loop creates new variables". It applies to modules whose
  `go.mod` says `go 1.22` or later. That `range` loops are covered was not checked against the spec:
  **Unverified**.

## Sources

- `testing` package (Benchmarks, `B.Loop`, `B.Run`, `B.ReportAllocs`, `B.ResetTimer`): <https://pkg.go.dev/testing>
- `go test` flags (`go help testflag`): <https://pkg.go.dev/cmd/go#hdr-Testing_flags>
- More predictable benchmarking with `testing.B.Loop`: <https://go.dev/blog/testing-b-loop>
- benchstat: <https://pkg.go.dev/golang.org/x/perf/cmd/benchstat>
- Go 1.22 release notes (loop variables): <https://go.dev/doc/go1.22>
- BenchmarkDotNet: <https://benchmarkdotnet.org/articles/overview.html>

Fact-checked 2026-11-15 on Go 1.27.3 (flag text from `go help testflag` and `go doc testing`; web pages via
a fetch tool that summarises). **Verified:** `BenchmarkXxx`/`-bench`, `b.Loop` behavior and timer rules,
`ResetTimer`/`StopTimer`, the flags in the table, `-bench` slash splitting, `ReportAllocs`, benchstat's
purpose, the Go 1.22 loop change, and the BenchmarkDotNet attributes and Release requirement.
**Unverified:** the `-20` suffix meaning, `B/op` and `allocs/op` definitions, the `benchstat` install
command, Go having no warmup or statistics, `range` loops under the 1.22 change, and the `/` name joining
(implied only). Section 6's numbers and the `strings.Join` allocation cause come from `section10_44`, not
from a source.
