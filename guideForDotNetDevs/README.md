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

## Conventions

- Every claim links to an official source (Go spec, Effective Go, Go FAQ, Code Review Comments,
  Google Go Style Guide, pkg.go.dev, Go issues).
- Claims that could not be confirmed from a source are marked **Unverified**.
- Check the Go spec first for language-behavior claims.
- Last fact-checked on 2026-11-11 against Go 1.27.3: sources were fetched and the code examples were
  compiled or run. Remaining gaps are marked **Unverified** in each doc.
