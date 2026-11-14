# 14 - Functions as parameters (`Func<>`, `Action<>`, delegates)

**C#:** you pass behavior with `Func<T, TResult>`, `Action<T>` or a custom `delegate`, usually as a lambda
(`a => a.Balance > 400`) or a method group (`acct.Deposit`). Delegates can be multicast (`+=`) and compared.

**Go:** functions are first-class values with a *function type*. There is no `delegate` keyword and no
`Func`/`Action` split: the type is just the signature. Examples use the `Account` from `bankAccount`.

| C# | Go |
|---|---|
| `Func<Account, bool>` | `func(*Account) bool` |
| `Action<Account>` | `func(*Account)` (empty result list) |
| `delegate bool Pred(Account a);` | `type Pred func(*Account) bool` |
| Lambda `a => a.Balance > 400` | Function literal `func(a *Account) bool { return a.Balance() > 400 }` |
| Method group `acct.Deposit` | Method value `acct.Deposit` |
| No direct equivalent (closest: a lambda taking the instance as its first parameter) | Method expression `(*Account).Withdraw` |
| Multicast `d += other` | None: a func value holds one function |
| `d?.Invoke()`, `NullReferenceException` | Check `f != nil`; calling a nil func panics |
| Delegates can be compared (`==`, by invocation list) | A func value can only be compared to `nil` |

## Passing a function

```go
func ApplyAll(accounts []*Account, op func(*Account) error) error {
	for _, a := range accounts {
		if err := op(a); err != nil {
			return err
		}
	}
	return nil
}

fee := 25.0 // captured by the closure
err := ApplyAll(accts, func(a *Account) error { return a.Withdraw(fee) })
```

A callback reports failure by returning an `error`, like any Go function; the caller decides what to do.
The spec's *function literal* is an anonymous function. There is no `=>` short form, so Go lambdas are
longer than C# ones.

## Closures capture variables, not copies

Function literals are closures: "they may refer to variables declared in a surrounding function. Those
variables are then shared between the surrounding function and the function literal" (Go spec, Function
literals). Run on Go 1.27.3: after `fee = 10` the same closure withdrew 10, not 25.

Loop variables: since Go 1.22 each `for` iteration has its own variable (Go spec, For statements), so
three closures built in `for i := 0; i < 3; i++` returned `0 1 2` when called afterward. Older toolchains
shared one variable. Your `go.mod` uses 1.27.x.

## Method values and method expressions

Both were run on Go 1.27.3 (`%T` output shown):

```go
deposit := acct.Deposit          // func(float64) error   (receiver already bound to acct)
withdraw := (*Account).Withdraw  // func(*main.Account, float64) error  (receiver is parameter 1)
withdraw(other, 1000)
```

`acct.Deposit` is the C# method group: it carries `acct` with it. The method expression is the unbound form.

## nil and comparison

- A func variable's zero value is `nil`. Calling it panicked: `runtime error: invalid memory address or nil
  pointer dereference` (recovered, as in `section07_33`).
- `f == g` for two funcs does not compile: `invalid operation: f == g (func can only be compared to nil)`.

## Function parameter or interface?

Go lets you attach methods to a named function type, which is how `net/http` turns a plain function into a
handler (`http.HandlerFunc` is documented as an adapter). Standard library APIs that need one small piece
of behavior take a function, for example `slices.SortFunc`. A rule of thumb (judgment, not a quoted rule):
use a function parameter for one behavior, and an interface (doc 08) when the callee needs several related
methods or you want to name the concept in a domain model. Generics (doc 12) let one helper work for any
element type:

```go
func Filter[T any](xs []T, keep func(T) bool) []T { /* loop, append when keep(x) */ }
rich := Filter(accts, func(a *Account) bool { return a.Balance() > 400 })
```

## Sources

- Go spec, Function types; Function literals; Method values; Method expressions; Calls; Comparison
  operators; For statements; Method declarations: <https://go.dev/ref/spec>
- `net/http.HandlerFunc`: <https://pkg.go.dev/net/http#HandlerFunc>
- `slices.SortFunc` and the `slices` function list: <https://pkg.go.dev/slices>
- C# delegates (multicast, comparison, method group assigned to a delegate):
  <https://learn.microsoft.com/dotnet/csharp/programming-guide/delegates/using-delegates>
- C# method group conversions: <https://learn.microsoft.com/dotnet/csharp/language-reference/language-specification/conversions>

Fact-checked on 2026-11-13 on Go 1.27.3: the spec and pkg.go.dev quotes were checked against the raw page
text (some items via the fetch tool's summary), and every example and compiler message above was run.

**Verified:** function types and the optional result list; the closure, loop-variable, method value and
method expression text; func values compare only to `nil`; calling a nil func panics; methods on a defined
func type (`HandlerFunc`); the `slices` package has no `Filter` or `Map` (it offers `IndexFunc`,
`ContainsFunc`, `DeleteFunc`, `SortFunc` and others); C# multicast `+=`, delegate comparison and method
group conversion.

**Unverified / judgment:**
- The function-versus-interface guidance is judgment, not a quoted rule.
- The spec says only "run-time panic" for a nil call; the exact message `invalid memory address or nil
  pointer dereference` is what Go 1.27.3 printed.
- The Go FAQ claim that methods can be defined on any type was not fetched; the spec wording used here is
  "a defined type in the same package".
