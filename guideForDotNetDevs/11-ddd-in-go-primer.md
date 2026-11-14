# 11 - DDD in Go: a primer for .NET developers

**C#:** DDD usually means entities as classes with a `private set` Id, value objects as `record`s,
aggregates guarded by private constructors, domain exceptions, repository interfaces, and one assembly
per bounded context.

**Go:** the language has no DDD features, and the official sources cited here do not discuss DDD. The
*language facts* below are sourced; the *DDD mapping* is convention and judgment. Examples use `Address`,
`Customer` and `Company` from `section09_41`.

## Mapping

| DDD / C# | Go |
|---|---|
| Value object (`record`) | Small struct passed by value; `==` compares fields |
| Entity (class with `Id`) | Struct with an identity field, usually used via pointer |
| Aggregate root | Exported type in its own package |
| Validating ctor / factory | `New(...) (T, error)` function |
| `private set`, invariants | Unexported fields (per **package**, doc 06) |
| Domain exception | Returned `error`; sentinel or typed errors |
| Repository interface | Small interface, satisfied implicitly (doc 07) |
| Bounded context / assembly | Package (directory) |
| Entity base class | Embedding, used sparingly |

## 1. The package is the boundary

Exported means a capitalized first letter (Go spec, Exported identifiers), so invariants are protected
only once the type lives in its own package. In `section09_41`'s `package main`, `Address{}` and
`Customer{Email: "x"}` skip the constructor. See docs 06 and 15.

```go
package customer

type ID int // distinct type instead of a bare int

type Customer struct {
	id   ID
	name string
}

var ErrNameRequired = errors.New("name not provided")

func New(id ID, name string) (*Customer, error) {
	if name = strings.TrimSpace(name); name == "" {
		return nil, ErrNameRequired
	}
	return &Customer{id: id, name: name}, nil
}

func (c *Customer) ID() ID { return c.id } // no Get prefix
```

From another package, `z.id` fails with `z.id undefined (type customer.Customer has no field or method id)`
and `customer.Customer{id: 1}` with `cannot refer to unexported field id in struct literal`. But
`var z customer.Customer` always compiles, so a constructor cannot be *forced* (doc 15).

## 2. Value objects and entities

`Address` is a value object: no identity, and structs are comparable when all fields are (Go spec,
Comparison operators), so `a == b` is value equality. There is no `with`; copy and change the copy.
`Customer` is an entity. Once it gains behavior like `ChangeEmail`, it needs pointer receivers and
constructors usually return `*T` (doc 05; design suggestion).

## 3. Constructors and naming

A constructor is a package-level function, not a method (a method needs a value to call it on). Names are
read with the package name: `customer.New`, not `customer.NewCustomer` (Effective Go, Package names;
Google decisions: `widget.NewWidget` becomes `widget.New`).

## 4. Errors are the domain language

Go returns an `error` instead of throwing. Match on identity, not message text: sentinels with
`errors.Is`, typed errors with `errors.As`, several failures via `errors.Join` (Go 1.20+)
(pkg.go.dev/errors). Use `fmt.Errorf` when formatting or wrapping with `%w`; otherwise `errors.New`.
When wrapping, `%w` keeps the original reachable; `%v` only copies its text. Run: `errors.Is` was `false`
for `%v` and `true` for `%w`.

## 5. Repositories and interfaces

The consumer declares the small interface it needs, such as `Save(*Customer) error`. Code Review Comments
(Interfaces): "Go interfaces generally belong in the package that uses values of the interface type, not
the package that implements those values."

## 6. Embedding: composition, not inheritance

Effective Go (Embedding): promoted methods are invoked with "the inner type, not the outer one" as the
receiver, so there is no virtual dispatch.

- **Fits:** a base type shared by aggregates, as in `get-eventually/go-eventually`, whose README says
  to embed `aggregate.BaseRoot` ("you MUST embed"). That it records events was seen only in a page
  summary: **Unverified**.
- **Fits:** `SavingsAccount` embeds `Account` (`bankAccount`): every `Account` method belongs on it.
  `NewSavingsAccount` calls `NewAccount` and stores `Account: *a` (no `base(...)` in Go). It is still not
  an `Account` (not assignable, doc 01), only `*SavingsAccount` gets `*Account`'s methods (doc 05), and
  same-package code can read the unexported `balance` (doc 15).
- **Fits:** embedding an interface to decorate one method, as `sort.Reverse` does with `sort.Interface`.
- **Poor fit:** plain has-a. `Company` has an `Address`: use a named field. Embedding promotes every
  exported member, and `Address.String()` is hidden only because `Company` declares its own (spec,
  Selectors: shallowest depth wins).
- **Never embed a mutex** (Uber guide, Zero-value Mutexes).

Uber: embedded types in public structs "leak implementation details, inhibit type evolution, and obscure
documentation"; embedding "should provide tangible benefit" (Avoid Embedding Types in Public Structs;
Embedding in Structs).

## 7. What to leave behind

Skip inheritance hierarchies, DI containers and ORM-shaped entities. Start with plain structs, one `New`
per aggregate and small interfaces (judgment).

## Sources

- Go spec: <https://go.dev/ref/spec>
- Effective Go: <https://go.dev/doc/effective_go>
- Google decisions: <https://google.github.io/styleguide/go/decisions>
- Uber guide: <https://github.com/uber-go/guide/blob/master/style.md>
- Code Review Comments: <https://go.dev/wiki/CodeReviewComments#interfaces>
- errors: <https://pkg.go.dev/errors>
- `sort.Reverse`: <https://github.com/golang/go/blob/master/src/sort/sort.go>
- go-eventually: <https://github.com/get-eventually/go-eventually>

Fact-checked 2026-11-13 on Go 1.27.3 (fetch tool summarises; some quotes from raw sources; the snippet
compiled and the compiler messages above are real output). **Verified:** all quotes and spec rules
cited. **Unverified / judgment:** the DDD mapping (sections 1 to 3, 5, 7, `type ID int`, pointer
receivers); "match on identity" and `errors.New` advice; `customer.ID` as a stutter example (a
principle in the guide, not a literal one); go-eventually recording events.
