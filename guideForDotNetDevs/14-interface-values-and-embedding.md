# 14 - Interface values and embedding (what an interface holds, what converts to what)

**C#:** `IPerson p = employee;` stores a reference to the object. Interfaces can inherit from other interfaces (`IWorker : IPerson`).

**Go:** an interface value holds a (dynamic type, value) pair, and interfaces can embed other interfaces. Continues [07 - Implicit interfaces](07-implicit-interfaces.md); `Person` and `Employee` are the types from there.

## What an interface variable holds

Spec, Interface types: "A variable of interface type can store a value of any type that is in the type set of the interface." Variables: "Variables of interface type also have a distinct _dynamic type_, which is the (non-interface) type of the value assigned to the variable at run time"; example `x = v // x has value (*T)(nil) and dynamic type *T`. Assignability: "A value `x` of type `V` is _assignable_ to a variable of type `T` ... if one of the following conditions applies", one being "`T` is an interface type, but not a type parameter, and `x` implements `T`." Method sets:

- "The method set of a defined type `T` consists of all methods declared with receiver type `T`."
- "The method set of a pointer to a defined type `T` ... is the set of all methods declared with receiver `*T` or `T`."

`FullName` has receiver `*Employee`, so only `*Employee` has it. `p` holds (dynamic type `*Employee`, the pointer); the struct is not copied. Verified on Go 1.27.3 with `section08_37`'s types:

```go
employee := &Employee{FirstName: "Farooq", LastName: "Mahmud", IsActive: true}
var p Person = employee
fmt.Printf("%T\n", p)                  // *main.Employee
fmt.Println(p.(*Employee) == employee) // true: same pointer
p.(*Employee).Deactivate()
fmt.Println(employee.IsActive)         // false: change made through p seen through employee
employee.LastName = "Changed"
fmt.Println(p.FullName())              // Farooq Changed: change made through employee seen through p
```

Output: `*main.Employee`, `true`, `false`, `Farooq Changed`.

C#: `IPerson p = employee;` works because "A class or struct lists the interfaces it implements after a colon in its declaration" (Microsoft Learn). Copying `p` copies the (type, pointer) pair. Runtime (Russ Cox): implicit conversions are checked at compile time; dynamic ones compute or cache an itable by looking up each interface method in the concrete type's method table. An implementation detail that may change.

## Interfaces can embed other interfaces

Spec, Embedded interfaces: "an interface `T` may use a (possibly qualified) interface type name `E` as an interface element. This is called _embedding_ interface `E` in `T`" (Go 1.14). Effective Go: "A `ReadWriter` can do what a `Reader` does _and_ what a `Writer` does; it is a union of the embedded interfaces. Only interfaces can be embedded within interfaces."

```go
type Person interface {
	FullName() string
}

type Worker interface {
	Person        // embeds Person: Worker's method set includes FullName
	Deactivate()
}
```

`*Employee` implements `Worker` and `Person` with no declaration: composition of method sets, not inheritance. The type set of `Worker` is "the _intersection_ of the type sets defined by `Worker`'s explicitly declared methods and the type sets of `Worker`'s embedded interfaces".

### What converts to what

Verified on Go 1.27.3:

```go
e := &Employee{FirstName: "Farooq", LastName: "Mahmud", IsActive: true}
var w Worker = e              // compiles: *Employee has FullName and Deactivate
var p Person = w              // compiles: implicit, no cast needed
fmt.Println(p.FullName())     // Farooq Mahmud

var p2 Person = e
w2, ok := p2.(Worker)         // runtime type assertion
fmt.Println(ok)               // true
w2.Deactivate()
fmt.Println(e.IsActive)       // false
```

- `Worker` to `Person` needs no cast (spec, Implementing an interface): `T` implements `I` if "`T` is an interface and the type set of `T` is a subset of the type set of `I`".
- `Person` to `Worker` does not compile:

```go
var p Person = &Employee{}
var w Worker = p
var x Worker = Worker(p)
```

```
cannot use p (variable of interface type Person) as Worker value in variable declaration: Person does not implement Worker (missing method Deactivate)
cannot convert q (variable of interface type Person) to type Worker: Person does not implement Worker (missing method Deactivate)
```

(Second line: variable named `q` in the test program.)

- Compile time sees the static type; `p2.(Worker)` checks the held `*Employee` at run time. An `io.Reader` holding a `*strings.Reader` asserted to `io.Seeker` gave `true`, to `io.ReadWriter` `false` (no panic with `v, ok`).

### Overlapping and conflicting methods

```go
type A interface{ M() }
type B interface{ M() }
type C interface {
	A
	B
}
```

Compiles on Go 1.27.3; a `T` with `func (T) M()` was assigned to `C` then `A`, and `M` ran each time. Spec: "When embedding interfaces, methods with the same names must have identical signatures." A different signature does not compile:

```go
type D interface{ M() int }
type E interface {
	A   // M()
	D   // M() int
}
type F interface {
	A
	M() int
}
```

```
duplicate method M
	other declaration of method M
```

Reported for `E` and for `F`. Spec example: "illegal: signatures of Reader.Close and Close are different". Declaring a method twice explicitly is also illegal: "The name of each explicitly specified method must be unique".

### Embedding a concrete type

Embedding a struct compiles but makes the interface a constraint:

```go
type W1 interface {
	Employee
	Deactivate()
}
var w W1
```

```
cannot use type W1 outside a type constraint: interface contains type constraints
```

Spec: such interfaces "may only be used as type constraints, or as elements of other interfaces used as constraints. They cannot be the types of values or variables".

### Standard library: `io.ReadWriter`

```go
type ReadWriter interface {
	Reader
	Writer
}

type ReadCloser interface {
	Reader
	Closer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}
```

Verified: `var rw io.ReadWriter = &bytes.Buffer{}` and `var r io.Reader = rw` compiled (no cast); `r` returned what was written through `rw`. Google Go Style Guide (Interfaces): "Design interfaces to be small for easier implementation and composition".

## C# to Go

| C# | Go |
|---|---|
| `interface IWorker : IPerson` (interfaces "can inherit from one or more other interfaces"); `IX : IA, IB` | `interface { Person; Deactivate() }`; `interface { A; B }` |
| A class implementing `IShape` "can be implicitly converted to `IDrawable`, because `IShape` inherits from it" | `Worker` assigns to `Person`, no cast |
| Same member in two interfaces: explicit implementation if they differ | Identical signatures merge; different ones fail to compile |

## Sources

- Go spec: <https://go.dev/ref/spec#Interface_types>, <https://go.dev/ref/spec#General_interfaces>, <https://go.dev/ref/spec#Method_sets>, <https://go.dev/ref/spec#Assignability>, <https://go.dev/ref/spec#Variables>
- Effective Go: <https://go.dev/doc/effective_go#embedding>
- Google Style Guide: <https://google.github.io/styleguide/go/decisions.html>
- Package io: <https://pkg.go.dev/io>
- Russ Cox, Go Data Structures: Interfaces: <https://research.swtch.com/interfaces>
- C# Interfaces: <https://learn.microsoft.com/en-us/dotnet/csharp/fundamentals/types/interfaces>

Fact-checked 2026-11-11 on Go 1.27.3 (windows/amd64); all output came from running the examples.

**Unverified:**
- Spec, Effective Go and Style Guide were fetched via a tool reading only the first 100000 characters (spec about 264000, style guide about 124000). The style guide sentence "Design interfaces to be small for easier implementation and composition" came from a summary; exact heading and wording unchecked. The Go FAQ and Code Review Comments were not consulted for composition.
- The Assignability, Method sets, Interface types and Variables quotes came from those excerpts (bullets were in range) but were not checked against the full page.
- Pre-Go-1.14 overlapping-method behavior is not claimed; the spec only marks embedding `[Go 1.14]`, no older toolchain was run.
- That copying an interface value copies the (type, pointer) pair rests on the run above and the spec's "has value (*T)(nil)" example, not a spec sentence.
- That `Person(x)` follows the same implements rule rests only on the Conversions text "x is assignable to T".
- The itable construction timing for static conversions is not stated in the Cox article.
