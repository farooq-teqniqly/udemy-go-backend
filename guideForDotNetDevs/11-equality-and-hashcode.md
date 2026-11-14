# 11 - Equality and `GetHashCode`

**C#:** every type inherits `object.Equals()`/`GetHashCode()`, both overridable. `record` types get
compiler-generated structural versions. `IEquatable<T>` lets a type opt into typed equality used by
generic collections.

**Go:** none of this exists as a per-type virtual/override mechanism. Equality is fixed by the
language; hashing is internal to `map` and not exposed per type.

## `==` is structural and fixed, not overridable

Go has no operator overloading (Go FAQ, "Why does Go not support overloading of methods and
operators?": "things are simpler without it"). For structs and arrays, `==` is defined purely
structurally by the spec (go.dev/ref/spec#Comparison_operators):

> Two struct values are equal if their corresponding non-blank field values are equal.
> Two array values are equal if their corresponding element values are equal.

No type can change what `==` does for itself — there's no `IEquatable<T>` hook.

```go
type Employee struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Salary    float64
	IsActive  bool
	JoinedAt  time.Time
}

e1 := NewEmployee("Farooq", "Mahmud", 150000, true, time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC))
e2 := NewEmployee("Farooq", "Mahmud", 150000, true, time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC))

e1 == e2    // false: e1, e2 are *Employee -- compares pointer addresses, not contents
*e1 == *e2  // false: dereferenced struct compare is field-by-field, and IDs differ (uuid.New() per call)
```

**Caveat vs. C# records:** struct `==` compares *unexported* fields too — the spec says "all
non-blank fields," no visibility exception. `record`'s generated `Equals` does the same (all
fields), so this part matches C#; just know Go gives you no way to carve out a field the way you'd
override `Equals` to skip one.

## Not every type supports `==`

Slices, maps, and function types are not comparable (spec, same section). A struct or array
containing one of those fields becomes non-comparable itself. Using `==` on a type that's
statically known to be non-comparable is a compile error; doing it through an `interface{}` value
whose dynamic type turns out non-comparable is a runtime panic instead.

## Custom equality: a plain `Equal` method, not an interface

When equality should be something other than "all fields match" -- e.g. "same business entity,"
identified by ID -- the idiom is a method, by convention named `Equal`, same as `time.Time.Equal`:

```go
func (e *Employee) Equal(other *Employee) bool {
	return e.ID == other.ID
}

e1.Equal(e2) // true once e1.ID == e2.ID
```

This is **not** backed by any standard interface -- there is no stdlib `Equal`-method interface
like `IEquatable<T>`. `slices.Equal` uses the built-in `comparable` constraint (tied to `==`), and
`slices.EqualFunc` just takes a plain `func(E1, E2) bool` parameter; neither declares an
`Equal`-bearing interface. The closest thing to a recognized convention is the third-party
`google/go-cmp` package, which reflectively looks for an `Equal` method -- that's `go-cmp`'s own
convention, not something the language or stdlib enforces.

`time.Time` is itself a good reason to prefer a method over `==`: its docs say to "prefer
t.Equal(u) to t == u" because `==` also compares the location and monotonic-clock reading, so two
values representing the same instant can compare unequal with `==`.

## `GetHashCode`: no per-type override; `map` hashing is internal

Go's built-in `map` needs its key type to be `comparable`, but the hash function it uses internally
is not exposed and not overridable per type -- the spec only requires `==`/`!=` to be "fully
defined for operands of the key type."

The closest built-in to a user-facing hash function is `hash/maphash`:

```go
seed := maphash.MakeSeed()
h1 := maphash.Comparable(seed, e1.ID)
h2 := maphash.Comparable(seed, e2.ID)
```

`maphash.Comparable[T comparable](seed Seed, v T) uint64` guarantees equal values hash equal, but
it works only for `comparable` types (same restriction as map keys), is not cryptographically
secure, and a `Seed` "is local to a single process and cannot be serialized or otherwise recreated
in a different process" -- so hash values aren't stable across runs or machines. There's no
convention equivalent to overriding `GetHashCode`; types just need to be `comparable` to be used as
map keys, full stop.

## Summary table

| C# | Go |
|---|---|
| `object.Equals()` / `record` structural `Equals` | `==` on structs/arrays, fixed by the spec, not overridable |
| `IEquatable<T>` | No equivalent interface; a plain `Equal(other T) bool` method by convention |
| `object.GetHashCode()` (overridable) | No per-type hook; `map` hashing is internal/unexported |
| `HashCode.Combine(...)` / custom hash in `GetHashCode` | `hash/maphash.Comparable` for a one-off hash of a comparable value (not identity-based, not cross-process stable) |

## Sources

- Comparison operators: <https://go.dev/ref/spec#Comparison_operators>
- Go FAQ, operator overloading: <https://go.dev/doc/faq#overloading>
- `time.Time.Equal`: <https://pkg.go.dev/time#Time.Equal>
- `hash/maphash`: <https://pkg.go.dev/hash/maphash>
- `slices` package (`Equal`, `EqualFunc`): <https://pkg.go.dev/slices>
- `go/types.Comparable`: <https://pkg.go.dev/go/types>
- `google/go-cmp`'s `Equal`-method convention: <https://pkg.go.dev/github.com/google/go-cmp/cmp>

**Unverified:** whether a conventional `Equal(other T) bool` method is prescribed anywhere in
Effective Go or the Go blog as a named idiom -- it was only confirmed as stdlib convention
(`time.Time.Equal`), not as a documented rule.
