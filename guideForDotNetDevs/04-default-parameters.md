# 04 - Default parameters

**C#:** optional parameters (`int salary = 0`), named arguments, and method overloads.

**Go:** none of these. A function has exactly one signature and every parameter is required. The Go FAQ
explains the related decision on overloading: having "a variety of methods with the same name but different
signatures was occasionally useful but that it could also be confusing and fragile in practice". Variadic
(`...T`) is the only form of optional arguments in the language itself, and it is just a final parameter
that accepts zero or more values.

## Idioms, using `NewEmployee`

### 1. Plain arguments

Fine when the list is short and every value is required.

### 2. Options struct

```go
type EmployeeOptions struct {
	Salary   float64
	IsActive bool
	JoinedAt time.Time // zero value means "not set"
}

func NewEmployee(firstName, lastName string, opts EmployeeOptions) Employee {
	if opts.JoinedAt.IsZero() {
		opts.JoinedAt = time.Now()
	}
	return Employee{ID: uuid.New(), FirstName: firstName, LastName: lastName,
		Salary: opts.Salary, IsActive: opts.IsActive, JoinedAt: opts.JoinedAt}
}

e := NewEmployee("Farooq", "Mahmud", EmployeeOptions{Salary: 150000})
```

Caveat: an unset field and an explicit zero look the same (`IsActive: false`). Google's guide says this form
is often preferred when some of these apply: all callers specify one or more options, many callers provide
many options, or the options are shared between multiple functions. Its benefits include self-documenting
field names, omitting defaults, and a struct that can grow without breaking call sites.

### 3. Functional options (closest to C# optional arguments)

```go
type EmployeeOption func(*Employee)

func WithSalary(s float64) EmployeeOption { return func(e *Employee) { e.Salary = s } }

func NewEmployee(firstName, lastName string, opts ...EmployeeOption) Employee {
	e := Employee{
		ID: uuid.New(), FirstName: firstName, LastName: lastName,
		IsActive: true,       // defaults live here
		JoinedAt: time.Now(), // defaults live here
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

e := NewEmployee("Farooq", "Mahmud", WithSalary(150000))
```

Google: prefer this when many of these apply: most callers need no options, most options are used
infrequently, there are many options, options take arguments, options could fail, or users/other packages
can provide custom options. It "requires a substantial amount of additional code", so use it only when the
advantages outweigh the overhead. Uber's guide recommends it for optional arguments in constructors and
public APIs you expect to expand, especially with three or more arguments.

### 4. Variadic (`...T`)

Only for optional values of one type. Not a general default mechanism.

## Which one?

There is no universal rule. Judge by how the call looks across all expected uses. Use an options struct
when most callers set the same fields; use functional options when most callers take the defaults.

## Sources

- Go FAQ, "Why does Go not support overloading of methods and operators?": <https://go.dev/doc/faq#overloading>
- Go spec, Function types / Parameters (variadic `...`; the parameter grammar has no default value):
  <https://go.dev/ref/spec#Function_types>
- Uber Go Style Guide, Functional Options: <https://github.com/uber-go/guide/blob/master/style.md>
- Google Go Style Guide, best practices (option structure, functional options):
  <https://google.github.io/styleguide/go/best-practices.html>
- Dave Cheney, Functional options for friendly APIs:
  <https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis>
- Rob Pike, Self-referential functions and the design of options (his post is
  about options that return an undo function, `type option func(f *Foo) option`; Cheney credits it as the
  source of the idea):
  <https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html>

Fact-checked against these pages on 2026-11-11. All three code samples compile (with a stub `uuid`, since the
installed toolchain's `uuid` is not the usual third-party package).

**Unverified:** the Google page was read in two chunks, and the second (final 38k characters) was only
summarized by the fetch tool as containing no further option-related guidance, not read verbatim.
