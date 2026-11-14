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
| [07 - Implicit interfaces](07-implicit-interfaces.md) | `class Employee : IPerson` | No declaration; the compiler checks the method set where the value is used as the interface. |
| [08 - Time formatting layouts](08-time-formatting.md) | `"yyyy-MM-dd"` token format | Fixed reference instant `Mon Jan 2 15:04:05 MST 2006`; layout is built from its literal substrings. |
| [09 - Equality and `GetHashCode`](09-equality-and-hashcode.md) | `Equals`/`GetHashCode` overrides, `IEquatable<T>`, `record` | `==` is fixed/structural, not overridable; no hash override; custom equality is a plain `Equal` method by convention. |
| [10 - Generics and type constraints](10-generics-and-type-constraints.md) | `where T : struct`, `where T : IComparable<T>` | Constraints are interfaces; can be a union of type terms (`int \| float32 \| float64`), constraint-only once non-basic. |
| [11 - DDD in Go primer](11-ddd-in-go-primer.md) | entities, value objects, aggregates, domain exceptions, repositories | Package is the boundary; `New` returns `(T, error)`; errors as domain language; embedding is composition, used sparingly. |
| [12 - Functions as parameters](12-functions-as-parameters.md) | `Func<>`, `Action<>`, delegates, lambdas | Function types are first-class; literals are closures; method values replace method groups; funcs compare only to `nil`. |
| [15 - Where encapsulation is not airtight](15-encapsulation-limits.md) | `private` enforced against everything outside the class | Unexported is per package; zero values, copies, reflection, `encoding/json` and embedding reach around it. |
| [14 - Interface values and embedding](14-interface-values-and-embedding.md) | `IPerson p = employee;`, `IWorker : IPerson` | An interface value is a (type, value) pair; interfaces embed other interfaces; conversions go up implicitly, down by assertion. |
| [13 - Documenting code](13-documenting-code.md) | `///` XML comments, `<summary>`, `<param>` | Plain `//` doc comments above declarations, full sentences starting with the name; read by `go doc` and pkg.go.dev. |

## Conventions

- Every claim links to an official source (Go spec, Effective Go, Go FAQ, Code Review Comments,
  Google Go Style Guide, pkg.go.dev, Go issues).
- Claims that could not be confirmed from a source are marked **Unverified**.
- Check the Go spec first for language-behavior claims.
- Last fact-checked on 2026-11-13 against Go 1.27.3 (doc 11 and the 2026-11-13 additions to docs 01, 05
  and 06; the rest on 2026-11-11): sources were
  fetched and the code examples were compiled or run. Remaining gaps are marked **Unverified** in each doc.
