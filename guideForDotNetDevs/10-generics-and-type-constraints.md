# 10 - Generics and type constraints

**C#:** generics since 1.0, constraints via `where T : struct`, `where T : IComparable<T>`,
`where T : SomeBaseClass`, `where T : new()` -- each constrains by a single type, interface, or
value-vs-reference-ness. There's no `where T : int | float | double`-style closed union of
concrete types.

**Go:** generics arrived in Go 1.18 (go.dev/doc/go1.18, "Generics"). Type parameters are
constrained by interfaces -- but an interface used purely as a constraint can list a **union of
type terms**, something no ordinary ("basic") interface or C# constraint can express.

```go
type Number interface {
	int | float32 | float64
}

func Product[T Number](numbers ...T) T {
	var total T = 1
	for _, num := range numbers {
		total *= num
	}
	return total
}

// equivalent, written inline instead of as a named interface:
func Sum[T int | float32 | float64](numbers ...T) T {
	var total T
	for _, num := range numbers {
		total += num
	}
	return total
}
```

`Number` and the inline `int | float32 | float64` constraint on `Sum` are the same thing: a
**type element**, defined by the spec as "a union of one or more type terms"
(go.dev/ref/spec#General_interfaces). The union's allowed types are its **type set** -- "The type
set of a union of terms t1|t2|...|tn is the union of the type sets of the terms."

## This makes the interface constraint-only

An interface containing type terms (not just method signatures) is **non-basic**. The spec:

> Interfaces that are not [basic] may only be used as type constraints, or as elements of other
> interfaces used as constraints. They cannot be the types of values or variables, or components
> of other, non-interface types.

So `var n Number` does not compile -- `Number` can only appear as `[T Number]`. This is the same
rule covered in [07 - Implicit interfaces](07-implicit-interfaces.md#embedding-a-concrete-type) for
embedding a concrete struct type in an interface; a union of type terms is the other way to end up
with a non-basic, constraint-only interface.

## `~T`: match by underlying type, not exact type

A plain term like `int` in a union matches only `int` itself. Prefixing with `~` -- the
**underlying type term** (`UnderlyingType = "~" Type`) -- widens that to "the set of all types
whose underlying type is T":

```go
type Number interface {
	~int | ~float32 | ~float64
}

type Celsius int // underlying type is int

func Product[T Number](numbers ...T) T { /* ... */ }

Product(Celsius(10), Celsius(20)) // compiles with ~int; would not with plain int
```

Rules on `~T`: `T`'s own underlying type must be `T` itself (so `~error` is illegal -- `error` is
an interface), and non-interface terms in a union must have pairwise-disjoint type sets (`~int |
MyInt` is illegal if `MyInt`'s underlying type is `int` -- the sets overlap).

## `comparable`

A predeclared interface, added with generics in Go 1.18: "denotes the set of all types which can
be compared using `==` or `!=`. It may only be used as (or embedded in) a type constraint." This is
what backs `maphash.Comparable` and map-key type parameters -- see
[09 - Equality and `GetHashCode`](09-equality-and-hashcode.md).

## `cmp.Ordered`

Go 1.21 added the `cmp` package with `cmp.Ordered`, a constraint for types supporting `<`, `<=`,
`>`, `>=` (go.dev/doc/go1.21, "New cmp package"). Before that, the experimental
`golang.org/x/exp/constraints` package (explicitly outside the Go 1 compatibility guarantee) was
the common source for an `Ordered` constraint; `cmp.Ordered` is its stabilized stdlib successor
(**Unverified**: no stdlib doc page states the succession explicitly in those words -- this is
inferred from the Go 1.18 and 1.21 release notes together, not a direct quote).

```go
import "cmp"

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
```

## C# to Go

| C# | Go |
|---|---|
| `where T : struct` | No direct equivalent; closest is listing concrete types in a union, or `comparable` |
| `where T : IComparable<T>` | `cmp.Ordered`, or a hand-written constraint interface |
| `where T : SomeBaseClass` | No base-class constraints (no classes); embed an interface instead |
| `where T : new()` | **Unverified** -- not researched for this doc |
| No union-of-concrete-types constraint syntax | `interface { int \| float32 \| float64 }`, optionally with `~` per term |

## Sources

- Go spec, Interface types / General interfaces: <https://go.dev/ref/spec#Interface_types>,
  <https://go.dev/ref/spec#General_interfaces>
- Go 1.18 release notes, Generics, `comparable`: <https://go.dev/doc/go1.18>
- Go 1.21 release notes, New `cmp` package: <https://go.dev/doc/go1.21>
- `cmp` package: <https://pkg.go.dev/cmp#Ordered>

**Unverified:**
- That `cmp.Ordered` is explicitly documented as the successor to
  `golang.org/x/exp/constraints.Ordered` -- inferred from release-note timing, not a direct quote.
- C# `where T : new()` was not researched for this doc.
- No MS Learn page on C# generic constraints was fetched; the C# column is general knowledge of
  `where T : ...` syntax, not a cited source.
