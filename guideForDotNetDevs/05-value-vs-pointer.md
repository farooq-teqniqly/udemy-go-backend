# 05 - Value vs. pointer (`Employee` or `*Employee`)

**C#:** you decide once, at the type: `class` is a reference type, `struct` is a value type.

**Go:** there are only structs. A struct value is copied on assignment, like a C# `struct`. A `*T`
behaves like a C# `class` reference: callers share one instance. You choose at each use, through the
receiver type and the return type.

## Decision criteria

| Question | Points to |
|---|---|
| Will any method mutate it (`Deactivate()`, `GiveRaise()`)? | Pointer |
| Does it contain a `sync.Mutex` or other no-copy field? | Pointer |
| Is it small, immutable and value-like, like `time.Time`? | Value |
| Do you want `nil` to mean "absent"? | Pointer |
| Are you unsure? | Pointer |

- **Mutation:** "If the method needs to mutate the receiver, the receiver must be a pointer." (Code
  Review Comments, Receiver Type). The FAQ calls it "first, and most important".
- **Size:** for small structs a value receiver is "efficient and clear" (FAQ). Code Review Comments
  (Receiver Type) also says not to choose a value receiver just to reduce garbage without profiling first.
- **No-copy fields:** "you should generally take and return pointer types if your structs contain
  fields that should not be copied" (Google decisions, Copying).
- **Consistency:** "Don't mix receiver types" (Code Review Comments). If one method needs a pointer, the
  rest should use one too (FAQ).
- **When in doubt:** "If in doubt, use a pointer" (Code Review Comments, Receiver Type). Google
  decisions words it "When in doubt, use a pointer receiver" (Receiver type).
- **Method sets:** `T` has only value-receiver methods; `*T` has both. A pointer-receiver method means
  only `*T` satisfies an interface (Go spec, Method sets).
- **Method sets with embedding:** for a struct `S` embedding `T`, "the method sets of `S` and `*S` both
  include promoted methods with receiver `T`. The method set of `*S` also includes promoted methods with
  receiver `*T`" (Go spec, Struct types). So with `SavingsAccount` embedding `Account` and `Deposit`
  declared on `*Account`, only `*SavingsAccount` has `Deposit` in its method set. The spec does not say
  outright that `S` alone lacks it; that is inferred.

  ```
  var _ Depositor = SavingsAccount{}
  -> SavingsAccount does not implement Depositor (method Deposit has pointer receiver)
  var _ Depositor = &s          // compiles
  s.Deposit(10)                 // compiles: s is addressable, shorthand for (&s).Deposit(10)
  ```

  The last line is the spec's rule that for an addressable `x`, `x.m()` is shorthand for `(&x).m()`
  (Calls).
- **`fmt` and pointer-receiver `String()`:** with `String()` on `*Account`, observed on Go 1.27.3:
  `fmt.Println(&s)` printed `acct 10.00`, `fmt.Println(s)` printed `{{10} 0}`. The value's method set has
  no `String`, which is the likely reason `fmt` skips it. That `fmt` checks the method set this way was
  not confirmed against the fmt docs.

## Heap vs. stack is not the deciding factor

Returning a pointer does not automatically heap-allocate. The FAQ: "From a correctness standpoint, you
don't need to know." The compiler keeps a variable on the stack unless it cannot prove the variable is
not referenced after the function returns. Don't pick the return type for allocation reasons.

Observed with `go build -gcflags=-m` on Go 1.27.3 (windows/amd64): `NewEmployee` returning a value
reported no heap allocation for the struct. The pointer version reported `&Employee{...} escapes to
heap` when compiled as a standalone function, but `&Employee{...} does not escape` once it was inlined
into a caller that did not leak the pointer. When the caller stored the pointer in a package-level
variable, it escaped to heap. Escape analysis is a compiler detail that can change between versions.

## Example

Value version, fine while `Employee` stays immutable:

```go
func NewEmployee(firstName, lastName string, salary float64, isActive bool, joinedAt time.Time) Employee {
	return Employee{ID: uuid.New(), FirstName: firstName, LastName: lastName,
		Salary: salary, IsActive: isActive, JoinedAt: joinedAt}
}

func (e Employee) FullName() string { return e.FirstName + " " + e.LastName }
```

Pointer version, once you add mutators. Change every receiver at the same time:

```go
func NewEmployee(firstName, lastName string, salary float64, isActive bool, joinedAt time.Time) *Employee {
	return &Employee{ID: uuid.New(), FirstName: firstName, LastName: lastName,
		Salary: salary, IsActive: isActive, JoinedAt: joinedAt}
}

func (e *Employee) FullName() string { return e.FirstName + " " + e.LastName }
func (e *Employee) Deactivate()      { e.IsActive = false }
```

Taking the address of a composite literal (`&Employee{...}`) is fine; the storage survives the return
(Effective Go, Composite literals).

## Recommendation for `Employee`

88 bytes (`unsafe.Sizeof` on Go 1.27.3, windows/amd64, with a 16-byte array standing in for the UUID
and `time.Time` at 24 bytes), no mutex, and `FullName()` only reads. If it stays immutable, return
`Employee`. This is judgment, not a quoted rule. For a learning repo, "if in doubt, pointer" makes
`*Employee` the safer default.

## Sources

- Go Code Review Comments (Receiver Type, Pass Values): <https://go.dev/wiki/CodeReviewComments>
- Go FAQ (methods on values or pointers; heap or stack): <https://go.dev/doc/faq>
- Google Go Style Guide, decisions: <https://google.github.io/styleguide/go/decisions.html>
- Go spec, Method sets: <https://go.dev/ref/spec#Method_sets>
- Effective Go, pointers vs. values: <https://go.dev/doc/effective_go#pointers_vs_values>
- Effective Go, composite literals: <https://go.dev/doc/effective_go#composite_literals>

**Verified:** the quotes above against the fetched pages (the fetch tool summarizes pages, so quotes
were taken from its excerpts); method sets, `x.m()` as `(&x).m()` for addressable `x`, and `&T{...}`
(Go spec); returning the address of a local or composite literal (Effective Go); the C# comparison
(Microsoft Learn, "Value types" and "Reference types": `class` is a reference type, `struct` a value
type); both code examples compile and run on Go 1.27.3 with a stub `uuid` package; `Employee` size.
Google decisions has no rule on returning `T` vs. `*T` (whole page read; its only related text is the
Copying and Receiver type sections and "Pass values").

**Unverified:**
- That a nil pointer return signals "not found" is common practice. No fetched text confirms it.
- The spec's wording on struct assignment copying was not seen directly; the Representation of values
  section says structs are "self-contained" with a "complete copy of all its data". The Assignment
  statements and Zero value sections were only partly readable. The nil zero value of pointers is
  confirmed in Pointer types.
- Heap behavior was observed on one compiler version only, not documented as a guarantee.
- Added 2026-11-13: the method-set and `x.m()` quotes are from fetched spec pages (the fetch tool labelled
  them verbatim); the compiler message and `fmt` output are real runs. Why `fmt` skips the value's
  `String()` is unconfirmed (see above).
