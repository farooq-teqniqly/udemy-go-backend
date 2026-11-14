# 06 - Encapsulation (`private set` in Go)

**C#:** `public bool IsActive { get; private set; }` lets everyone read the property, but only code
inside the class can write it. Related tools: `init`, `required`, `readonly` fields and `record`s.

**Go:** no properties, no `private` keyword. Visibility is the **first letter of the name**: `IsActive` is
exported, `isActive` is not. The substitute for `private set` is an unexported field, an exported getter
and a mutation method such as `Deactivate()`. C# `private` is per **type**; Go "unexported" is per **package**.

## The idiom

Needs its own package (`employee/`), not `package main`:

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

- **Getter naming:** a setter, if needed, is `SetOwner` (Effective Go, Getters). Google decisions: "Function
  and method names should not use a `Get` or `get` prefix" (Naming, Getters).
- **Why it works:** "The use of upper-case names for export provides the hook
  to discriminate the field from the method" (Effective Go, Getters).
- **Setters are optional.** Prefer `Deactivate()` or `GiveRaise(pct)` over `SetIsActive(bool)` (like
  `ChangeAddress` in Microsoft's example); a design suggestion, not a quoted rule.
- **Getters are not automatic.** Effective Go: "There's nothing wrong with providing getters and setters
  yourself, and it's often appropriate to do so". Plain exported fields (`FirstName`) are normal Go.

## Why the field must be unexported for a same-named getter

A C# property is one member; Go splits it into a field and a method that cannot share a name. Spec: "If the base type is a struct type, the non-blank method and field names must be distinct" (Method
declarations). Real output from Go 1.27.3 for `type Account struct { Balance float64 }` with
`func (a *Account) Balance() float64`:

```
main.go:67:19: field and method with the same name Balance
	main.go:11:2: other declaration of Balance
main.go:48:2: invalid operation: a.Balance += amount (mismatched types func() float64 and float64)
```

(`a.Balance` now resolves to the method, type `func() float64`.) Fix: unexported `balance`.

```go
type Account struct{ balance float64 }

func (a *Account) Balance() float64 { return a.balance }
```

Effective Go: "If you have a field called `owner` (lower case, unexported), the getter method should be
called `Owner` (upper case, exported), not `GetOwner`" (Getters). Now only package code can set the balance.

## C# to Go table

| C# | Go |
|---|---|
| `{ get; private set; }` | unexported field + exported getter + exported mutation method |
| `{ get; }` | unexported field + getter, no mutator; set in `New...` |
| `{ get; init; }` | none. Exported field or getter-only idiom |
| `required` | none. Only constructor parameters |
| `private readonly` | unexported field (no `readonly`; never assign it) |
| `internal` | no direct equivalent (`internal/` directories restrict imports) |
| `record` / `with` | none. Copy the struct value, change the copy |
| `private` constructor | unexported fields + `New...` (nothing stops `var e Employee`) |

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

The zero value, value copies, same-package access, reflection, `encoding/json` and embedding all reach around
it. They are in [15 - Where Go encapsulation is not airtight](15-encapsulation-limits.md).

## Applying it to `section08_36`

Move `Employee` to `section08_36/employee/employee.go` (import `teqniqly.com/section08_36/employee`) and
rename `NewEmployee` to `employee.New` (Effective Go, Package names).

## Sources

- Go spec, Exported identifiers; Struct types; Type identity: <https://go.dev/ref/spec#Exported_identifiers>
- Go spec, Composite literals: <https://go.dev/ref/spec#Composite_literals>
- Go spec, Method declarations; Struct types (field names unique): <https://go.dev/ref/spec#Method_declarations>
- Effective Go, Getters; Names; Package names: <https://go.dev/doc/effective_go#Getters>
- Google Go Style Guide, decisions (Naming, Getters; Copying; Literal formatting): <https://google.github.io/styleguide/go/decisions.html>
- C# `init`: <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/init>
- C# `required`: <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/required>
- C# `readonly`: <https://learn.microsoft.com/dotnet/csharp/language-reference/keywords/readonly>
- C# records: <https://learn.microsoft.com/dotnet/csharp/language-reference/builtin-types/record>
- C# lightweight class (`private set`): <https://learn.microsoft.com/dotnet/csharp/programming-guide/classes-and-structs/how-to-implement-a-lightweight-class-with-auto-implemented-properties>

Fact-checked on 2026-11-11 on Go 1.27.3 (windows/amd64): quotes are summarising-fetch-tool excerpts; every
behavior and compiler message comes from a run in a scratch module with a separate `employee` package.

**Verified:** identifier export rule (spec); `GetX` and setter naming (Effective Go, Google decisions);
the five compiler errors; the behaviors in doc 15; the C# side
(`private set` modifiable after construction only inside the class, `init` only during initialisation,
`required`, `readonly` fields assignable only in declaration/constructor, positional records init-only).

**Added 2026-11-13:** the "must be distinct" rule (spec, Method declarations; the fetch tool's verbatim
excerpt) and the Effective Go `owner`/`Owner` sentence were confirmed; the compiler output is from building
`bankAccount/main.go` on Go 1.27.3. The `balance`/`Balance()` snippet was compiled and run (printed `5`).

**Unverified:**
- How a *promoted* field or method (from an embedded type) interacts with the distinct-names rule. The
  spec sentences cover the struct's own fields and bound methods; the Selectors shallowest-depth rule was
  not read for this addition.
- The spec text for "unexported fields cannot be named in a literal or selector from another package" was
  not seen (excerpts had no such sentence); the behavior is confirmed by compiler errors, not spec wording.
  The Selectors and Assignment sections gave no rule on exportedness or struct copying.
- No Google Style Guide text on avoiding trivial getters/setters, or on exported vs. unexported fields,
  was found (decisions: first 100,000 characters, best-practices: all, both fetched in parts through a
  summarising tool).
- Go Tour "Exported names" page could not be fetched (content did not load); it is not cited.
- `internal/` directory semantics (mentioned in the table) were not fetched.
