# 06 - Encapsulation (`private set` in Go)

**C#:** `public bool IsActive { get; private set; }` lets everyone read the property, but only code
inside the class can write it. Related tools: `init`, `required`, `readonly` fields and `record`s.

**Go:** there are no properties and no `private` keyword. Visibility is decided by the **first letter of
the name**: `IsActive` is exported (visible to other packages), `isActive` is not. The usual substitute
for `private set` is an unexported field, an exported getter method and, if the value can change, an
exported mutation method such as `Deactivate()`.

The big difference: C# `private` is per **type**; Go "unexported" is per **package**.

## The idiom

Because encapsulation is per package, this only works with the type in its own package (for example
`employee/`), not in a single `package main` file:

```go
package employee

type Employee struct {
	id        uuid.UUID // unexported: read-only to outsiders
	FirstName string    // exported: freely read and written
	LastName  string
	salary    float64
	isActive  bool
	joinedAt  time.Time
}

func New(first, last string, salary float64, joinedAt time.Time) *Employee {
	return &Employee{id: uuid.New(), FirstName: first, LastName: last,
		salary: salary, isActive: true, joinedAt: joinedAt}
}

func (e *Employee) ID() uuid.UUID       { return e.id }       // getter: no Get prefix
func (e *Employee) IsActive() bool      { return e.isActive }
func (e *Employee) Salary() float64     { return e.salary }
func (e *Employee) JoinedAt() time.Time { return e.joinedAt }
func (e *Employee) Deactivate()         { e.isActive = false } // controlled mutation
```

Caller, in another package:

```go
e := employee.New("A", "B", 100, time.Now())
fmt.Println(e.IsActive()) // true
e.Deactivate()            // the only way to change it
```

- **Getter naming:** a field `owner` gets a getter `Owner`, not `GetOwner`; a setter, if needed, is
  `SetOwner` (Effective Go, Getters). Google decisions agrees: "Function and method names should not use
  a `Get` or `get` prefix" (Naming, Getters).
- **Why a method named like the field works:** "The use of upper-case names for export provides the hook
  to discriminate the field from the method" (Effective Go, Getters).
- **Setters are optional.** Prefer a method that expresses intent (`Deactivate()`, `GiveRaise(pct)`) over
  a bare `SetIsActive(bool)`; that is the equivalent of `ChangeAddress` in Microsoft's `private set`
  example. This is a design suggestion, not a quoted rule.

## Why the field must be unexported for a same-named getter

A C# property is one member (`Balance`). Go splits it into a field and a method, and the two cannot share
a name. The spec: "If the base type is a struct type, the non-blank method and field names must be
distinct" (Method declarations). So an exported field `Balance` plus a getter `Balance()` does not compile.
Real output from Go 1.27.3 for `type Account struct { Balance float64 }` with `func (a *Account) Balance() float64`:

```
main.go:67:19: field and method with the same name Balance
	main.go:11:2: other declaration of Balance
main.go:48:2: invalid operation: a.Balance += amount (mismatched types func() float64 and float64)
```

The follow-on errors appear because `a.Balance` now resolves to the method (type `func() float64`), not the
field. The fix is on the field, as in the `Employee` idiom above: unexported `balance`, exported `Balance()`.

```go
type Account struct{ balance float64 }

func (a *Account) Balance() float64 { return a.balance }
```

This is exactly what Effective Go describes: "If you have a field called `owner` (lower case, unexported),
the getter method should be called `Owner` (upper case, exported), not `GetOwner`" (Getters). Note the
side effect: the exported field is gone, so only code in the package can set the balance, which is the
`private set` behavior (see the limits under "Where it is NOT airtight").

## C# to Go table

| C# | Go |
|---|---|
| `{ get; private set; }` | unexported field + exported getter + exported mutation method |
| `{ get; }` (read-only, set in ctor) | unexported field + exported getter, no mutator; set in `New...` |
| `{ get; init; }` | no equivalent. Exported field (anyone can write it) or the getter-only idiom |
| `required` | no equivalent. A constructor with parameters is the only enforcement |
| `private readonly` field | unexported field (the compiler has no `readonly`; just never assign it) |
| `private` (per type) | unexported name (per **package**) |
| `public` | exported name (capital letter) |
| `internal` | no direct equivalent. `internal/` directories restrict imports (not covered here) |
| `record` / `with` | none. Copy the struct value and change the copy |
| `private` constructor + factory | unexported fields + `New...` function (nothing stops `var e Employee`, see below) |

## Compile errors an outside package gets

Real output from Go 1.27.3, with `employee` as a separate package and the caller in `cmd/errs`:

```
e.isActive = false
  -> e.isActive undefined (type *employee.Employee has no field or method isActive, but does have method IsActive)

e.notes = "x"          // a field with no same-named method
  -> e.notes undefined (cannot refer to unexported field notes)

employee.Employee{isActive: true}
  -> cannot refer to unexported field isActive in struct literal of type employee.Employee

employee.Employee{"a", "b"}      // unkeyed literal
  -> implicit assignment to unexported field id in struct literal of type employee.Employee
  -> too few values in struct literal of type employee.Employee (type has unexported fields - use key:value pairs)

var v employee.Employee; v.salary = 5
  -> v.salary undefined (type employee.Employee has no field or method salary, but does have method Salary)
```

## Where it is NOT airtight

All of these were run (`cmd/ext`, `cmd/wipe` in the scratch module):

1. **The zero value is always constructible.** `var z employee.Employee` compiles in any package.
   Result: `IsActive()` false, `Salary()` 0, `ID()` all zeros. A constructor cannot force its use
   (C# can, with a `private` constructor).
2. **Value copies are the caller's own.** `v := *p2; v.Deactivate()` changes only `v`; `p2` stays active.
   The copy includes the unexported fields (the caller just cannot name them). With `*Employee`
   (`q := p`), both variables share one instance, so `q.Deactivate()` changes `p` too. See doc 05.
3. **Whole-struct assignment overwrites hidden fields.** `*p = employee.Employee{FirstName: "X"}` compiles
   outside the package and resets `isActive`, `salary` and `id` to zero values. Unexported fields stop
   being *named*, not being *replaced*.
4. **Anything in the same package can touch them.** `employee.Tamper(e)` (a plain function in package
   `employee`) changed `isActive` and `salary`. Go has no per-type privacy, so every file in
   `package main` can reach every unexported field. That is why `section08_36` (a single `package main`)
   cannot demonstrate encapsulation. Google decisions says it too: "Code within a package can access
   unexported identifiers in the package" (best-practices, Package size).
5. **Reflection** can read unexported fields but not set them. `CanSet()` was `false` for `isActive` and
   `true` for the exported `FirstName`; `SetBool` panicked with `reflect: reflect.Value.SetBool using
   value obtained using unexported field`. `Value.Interface()` also panics for such fields
   (pkg.go.dev/reflect). `unsafe` can bypass this (not tested).
6. **`encoding/json` skips unexported fields.** `json.Marshal(p2)` printed
   `{"FirstName":"C","LastName":"D"}`: `id`, `salary`, `isActive`, `joinedAt` are dropped. Unmarshalling
   cannot fill them either. If a value must be serialised, export it or write `MarshalJSON`.

7. **Embedding adds no `protected`.** A C# derived class cannot touch a private base field. In Go the
   rule is the package. With `SavingsAccount` embedding `Account` in `package account`, a method on
   `SavingsAccount` can read `sa.balance` directly (it compiled and ran). From another package, Go 1.27.3
   reports `sa.balance undefined (type account.SavingsAccount has no field or method balance, but does have
   method Balance)`, while `sa.Balance()` works. So moving `SavingsAccount` to its own package would force it
   through the getter.

## When not to bother with getters and setters

- Effective Go: "There's nothing wrong with providing getters and setters yourself, and it's often
  appropriate to do so" but they are not automatic. Plain exported fields are normal Go: `FirstName` and
  `LastName` above need no accessors.
- Keep a field unexported only when an invariant needs protecting (an ID that must not change, an
  `isActive` that should only flip via `Deactivate()`).
- Google style on getters: only naming guidance was found (no `Get` prefix). See **Unverified**.

## Applying it to `section08_36`

Move `Employee` into its own package directory (for example `section08_36/employee/employee.go`,
`package employee`), keep `main.go` as the caller, and import it as `teqniqly.com/section08_36/employee`.
Rename `NewEmployee` to `employee.New` (the package name already says `employee`, Effective Go, Package
names). If `Employee` stays in `package main`, lowercase fields still compile, but nothing is protected.

## Sources

- Go spec, Exported identifiers; Struct types; Type identity: <https://go.dev/ref/spec#Exported_identifiers>
- Go spec, Composite literals: <https://go.dev/ref/spec#Composite_literals>
- Go spec, Method declarations (field and method names must be distinct); Struct types (field names unique):
  <https://go.dev/ref/spec#Method_declarations>
- Effective Go, Getters; Names; Package names: <https://go.dev/doc/effective_go#Getters>
- Google Go Style Guide, decisions (Naming, Getters; Copying; Literal formatting):
  <https://google.github.io/styleguide/go/decisions.html>
- Google Go Style Guide, best practices (Package size): <https://google.github.io/styleguide/go/best-practices.html>
- reflect (`Value.CanSet`): <https://pkg.go.dev/reflect#Value.CanSet>
- encoding/json (`Marshal`): <https://pkg.go.dev/encoding/json#Marshal>
- C# `init` (includes the `private set` / read-only / `init` comparison):
  <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/init>
- C# `required`: <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/required>
- C# `readonly`: <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/readonly>
- C# records: <https://learn.microsoft.com/dotnet/csharp/language-reference/builtin-types/record>
- C# lightweight immutable class (`private set`): <https://learn.microsoft.com/dotnet/csharp/programming-guide/classes-and-structs/how-to-implement-a-lightweight-class-with-auto-implemented-properties>

Fact-checked on 2026-11-11 on Go 1.27.3 (windows/amd64): quotes taken from fetched pages (the fetch tool
summarises, so quotes are its excerpts); every behavior and compiler message above comes from a program
run in a scratch module with a separate `employee` package.

**Verified:** identifier export rule (spec); `GetX` naming and setter naming (Effective Go, Google decisions);
`CanSet` wording; all six behaviors in "Where it is NOT airtight"; the five compiler errors; the C# side
(`private set` modifiable after construction only inside the class, `init` only during initialisation,
`required`, `readonly` fields assignable only in declaration/constructor, positional records are
init-only).

**Added 2026-11-13:** the "must be distinct" rule (spec, Method declarations; quote from the fetch tool's
verbatim excerpt) and the Effective Go `owner`/`Owner` sentence were confirmed; the compiler output is
from building `bankAccount/main.go` on Go 1.27.3. The corrected `balance`/`Balance()` snippet was
compiled and run in a scratch module (printed `5`). Item 7 under "Where it is NOT airtight" (embedding and
same-package access) was also compiled in a scratch module with a separate `account` package.

**Unverified:**
- How a *promoted* field or method (from an embedded type) interacts with the distinct-names rule. The
  spec sentences cover the struct's own fields and the methods bound to it; the Selectors
  shallowest-depth rule was not read for this addition.
- The Go spec text for "unexported fields cannot be named in a literal or selector from another package"
  was not seen (the fetched excerpts had no such sentence). The behavior is confirmed by the compiler
  errors, not by spec wording. The Selectors and Assignment sections gave no rule on exportedness or
  struct copying.
- No Google Style Guide text on avoiding trivial getters/setters, or on exported vs. unexported fields,
  was found in the pages read (decisions: first 100,000 characters, best-practices: all, but fetched in
  parts through a summarising tool).
- Go Tour "Exported names" page could not be fetched (page content did not load); it is not cited.
- `encoding/json` doc says "Each exported struct field becomes a member of the object"; the explicit
  "unexported fields are ignored" wording was not seen. The behavior was confirmed by running it.
- That `unsafe` can set unexported fields was not tested or sourced.
- `internal/` directory semantics (mentioned in the table) were not fetched.
