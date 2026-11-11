# 02 - UUIDs and `Guid`

**C#:** `Guid.NewGuid()` returns a `Guid` struct. `Guid.Empty` is the zero value.

**Go:** a `UUID` is a plain `[16]byte` value type. Since Go 1.27 there is a standard library `uuid`
package (import path is just `uuid`; Go 1.27 release notes, "New uuid package"). Before that, `github.com/google/uuid` was the standard choice.

## Generating

```go
import "uuid" // Go 1.27+ standard library

id := uuid.New()   // random (v4); currently equivalent to uuid.NewV4()
v7 := uuid.NewV7() // time-ordered (v7); no error return
```

Other stdlib helpers: `uuid.Nil()` (a function, not a variable), `uuid.Max()`,
`uuid.Parse(s) (UUID, error)`, `uuid.MustParse(s)`, `u.Compare(v)`, `u.String()`.

Third-party equivalent (older Go, or a codebase that already uses it):

```go
import "github.com/google/uuid"

id := uuid.New()        // v4; panics on failure ("New creates a new random UUID or panics")
s := uuid.NewString()   // same as uuid.New().String(); also panics
v7, err := uuid.NewV7() // returns Nil and an error instead of panicking
```

## As a struct field

```go
type Employee struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
}
```

| C# | Go |
|---|---|
| `Guid?` | `*uuid.UUID` (`nil` = absent) |
| `Guid.Empty` | zero value, all zeros (stdlib also has `uuid.Nil()`) |
| `id == Guid.Empty` | `id == uuid.UUID{}` (arrays are comparable with `==`) |
| `Dictionary<Guid, T>` | `map[uuid.UUID]T` (arrays are valid map keys) |

- `encoding/json` writes and reads the field as a string, because `UUID` has `MarshalText` and
  `UnmarshalText`. Verified by running on Go 1.27.3: `{"id":"f81d4fae-7dec-11d0-a765-00a0c91e6bf6"}`
  round-trips, and `id == uuid.UUID{}` and `id == uuid.Nil()` are both `true` for a zero value.
- Pass and return it by value. At 16 bytes a pointer buys nothing.

## v4 vs. v7 for database keys

v7 is time-ordered, which usually gives better index locality than random v4 (RFC 9562, section 6.11:
"Time-ordered monotonic UUIDs benefit from greater database-index locality"). Rough .NET analogies:
`Guid.NewGuid()` is v4, `Guid.CreateVersion7()` is v7 (available in .NET 9 and later; Microsoft Learn,
`Guid.CreateVersion7`). Stdlib `NewV7` sorts in increasing order unless the system clock moves backwards.

## Sources

- Proposal and accepted API (issue title changed from `crypto/uuid` to `uuid` on 2026-03-12; accepted
  2026-04-08; closed as completed 2026-04-21, milestone Go1.27): <https://github.com/golang/go/issues/62026>
- `Nil` to `Zero` rename (closed not planned 2026-04-22): <https://github.com/golang/go/issues/78612>
- Go 1.27 release notes: <https://go.dev/doc/go1.27>
- RFC 9562: <https://www.rfc-editor.org/rfc/rfc9562>
- `google/uuid` v1.6.0: <https://pkg.go.dev/github.com/google/uuid>
- Spec, Comparison operators (arrays comparable): <https://go.dev/ref/spec#Comparison_operators>
- Spec, Struct types (tags): <https://go.dev/ref/spec#Struct_types>

**Verified:** the stdlib API above matches `go doc -all uuid` on Go 1.27.3 (`Max`, `MustParse`, `New`, `NewV4`,
`NewV7`, `Nil`, `Parse`, and methods `AppendText`, `Compare`, `MarshalText`, `String`, `UnmarshalText`).
Issue #78612 proposed renaming `Nil` to `Zero`; it was closed as not planned, so `Nil()` shipped.

**Unverified:** the exact date the import path moved from `crypto/uuid` to `uuid` (the issue title was
renamed on 2026-03-12; the accepted design comment of 2026-04-08 uses `uuid`). Spec wording for array
comparison was not fetched in full (the page was truncated), but `go doc uuid` states UUIDs are comparable
with `==`.
