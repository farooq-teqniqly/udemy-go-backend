# 07 - Implicit interfaces (how Go knows `Employee` is a `Person`)

**C#:** a type opts in by name: `class Employee : IPerson`. The compiler checks the members once, at that declaration.

**Go:** there is no declaration. Nothing in `Employee` mentions `Person`. The compiler checks, at each place where a value is converted to the interface type, whether the value's type has the interface's methods.

This expands the one-line row in [01 - Go is not object-oriented](01-go-is-not-oop.md).

## The question: "does Go search the methods in the current package?"

No. There is no search and no registry of "implementers". Three rules together give the behavior:

1. **A type implements an interface if its method set contains the interface's methods.** "A Go type implements an interface by implementing the methods of that interface, nothing more." (Go FAQ). The spec definition: `T` implements `I` if `T` is not an interface and is an element of the type set of `I`.
2. **The check happens where you use the value as the interface.** Assignment, passing an argument, returning a value, or a type assertion. If it fails, that line does not compile. Effective Go: "passing an `*os.File` to a function expecting an `io.Reader` will not compile unless `*os.File` implements the `io.Reader` interface."
3. **A type's method set is fixed by the type's own package.** A method's receiver base type "must be declared in the same package as the method" (Go spec, Method declarations). You cannot add methods to `time.Time` or to a type from another package, so nobody can change what `Employee` implements from the outside.

So the interface and the type can live in different packages, and neither needs to import the other. Only the code that converts one to the other (the consumer) needs to see both.

## The user's example

```go
type Person interface {
	FullName() string
}

type Employee struct { /* ... */ }

func (e *Employee) FullName() string { return e.FirstName + " " + e.LastName }
func (e *Employee) Deactivate()      { e.IsActive = false }
```

Method sets (Go spec, Method sets):

| Type | Method set |
|---|---|
| `Employee` | methods declared with receiver `Employee` (here: none) |
| `*Employee` | methods declared with receiver `*Employee` or `Employee` (here: `FullName`, `Deactivate`) |

So `*Employee` implements `Person`, and `Employee` does not.

## Pointer-receiver gotcha

Because `FullName` has a pointer receiver, only `*Employee` is a `Person`. Verified on Go 1.27.3:

```go
employee := &Employee{"A", "B"}
var p Person = employee    // compiles: *Employee implements Person
var q Person = *employee   // does not compile
var r Person = Employee{}  // does not compile
```

```
cannot use *employee (variable of struct type Employee) as Person value in variable declaration: Employee does not implement Person (method FullName has pointer receiver)
cannot use Employee{} (value of struct type Employee) as Person value in variable declaration: Employee does not implement Person (method FullName has pointer receiver)
```

With a value receiver (`func (e Employee) FullName()`), both `Employee` and `*Employee` compile as `Person` (verified: the same program with the receiver changed compiled and ran). Which receiver to pick is covered in [05 - Value vs. pointer](05-value-vs-pointer.md). `NewEmployee` already returns `*Employee`, so the user's file is fine as is.

## Missing method

If the method is absent (or misspelled), the error names the missing method:

```
cannot use (*Employee)(nil) (value of type *Employee) as Person value in variable declaration: *Employee does not implement Person (missing method FullName)
```

## Compile-time assertion: `var _ Person = (*Employee)(nil)`

If nothing in your package converts `*Employee` to `Person`, a typo in a method name goes unnoticed until someone else's code fails. The Go FAQ ("How can I guarantee my type satisfies an interface?") gives the idiom:

```go
var _ Person = (*Employee)(nil) // Verify that *Employee implements Person.
```

The blank identifier means no variable is created; the line exists only for the type check. With `FullName` removed or misspelled, the build fails with the "missing method" error above (verified). Effective Go ("Interface checks") says not to do this for every type: "such declarations are only used when there are no static conversions already present in the code, which is a rare event."

## Different packages

Verified with a three-package module: package `a` declares `Person` and `func Greet(p Person) string`; package `b` declares `Employee` with a pointer-receiver `FullName`; `main` imports both and calls `a.Greet(&b.Employee{...})`. It printed `Hello, A B`. `go list` showed neither `a` nor `b` imports anything.

## Where to define the interface

Define it where it is consumed, not where it is implemented:

- Code Review Comments: "Go interfaces generally belong in the package that uses values of the interface type." and "Do not define interfaces before they are used".
- Google Go Style Guide: "The consumer of the interface should define it (not the package implementing the interface)", and "Functions should take interfaces as arguments but return concrete types".

That is why `NewEmployee` returns `*Employee`, not `Person`. In C# the interface usually lives with the implementation or in a shared contracts assembly; in Go, small interfaces next to the consumer are the norm.

## Type assertions and the nil gotcha

Verified output (Go 1.27.3):

```go
var p Person = &Employee{"A", "B"}
e, ok := p.(*Employee)          // A true
_, ok = p.(Robot)               // false, no panic
_ = p.(Robot)                   // panics
var np *Employee
var q Person = np
fmt.Println(q == nil)           // false
```

```
recovered: interface conversion: main.Person is *main.Employee, not main.Robot
```

- Type assertion `x.(T)`: the single-value form panics on failure; the `v, ok` form does not (Go spec, Type assertions). If `T` is a concrete type it must implement the interface type of `x`, otherwise the assertion is a compile error.
- A type switch (`switch v := p.(type)`) selected `*Employee` and `Robot` cases correctly.
- An interface holding a nil `*Employee` is not `nil`. The FAQ: an interface value is nil only if both type and value are unset; storing a nil `*int` gives `(T=*int, V=nil)`, which is non-nil. Relevant if a function returns `Person` and you `return np`.

## C# to Go

| C# | Go |
|---|---|
| `class Employee : IPerson` declares the contract | No declaration. Having the methods is enough. |
| Compiler checks at the class declaration | Compiler checks where a value is converted to the interface. |
| Nominal typing: `IPerson` must be named in the base list | Structural typing (the FAQ's own word). |
| Interface usually defined with the implementation | Define it in the consumer package. |
| Extension methods are static methods called with instance syntax | Cannot add methods to a type from another package at all. |
| Explicit interface implementation (`double IMetric.GetDistance()`) | No equivalent. One method name satisfies every interface that wants it. |
| Default interface methods | Interfaces hold only method signatures; no bodies. |
| `IPerson p = employee;` always works for a class that lists it | `var p Person = employee` depends on the receiver type (`*Employee` vs. `Employee`). |
| Reflection or `is`/`as` at runtime | Type assertion / type switch, with `v, ok` form. |

## How the runtime represents it (implementation detail)

An interface value is a (type, value) pair. Russ Cox's article says implicit conversions are checked at compile time, and that for dynamic conversions the runtime computes (or finds in a cache) an itable by looking up each interface method in the concrete type's method table. This is an implementation detail and may change; the language guarantee is the spec's.

## How can a `*Employee` be assigned to a `Person`? (what an interface variable holds)

The question: `employee` is a `*Employee`, so how can it be assigned to `Person`, which is not a pointer type?

`Person` is declared as `type Person interface { ... }`, so it is an interface type, not a pointer type and not a struct type. A variable of interface type does not hold "a `Person` struct". The spec: "A variable of interface type can store a value of any type that is in the type set of the interface." (Go spec, Interface types). At run time it holds a (dynamic type, dynamic value) pair. The spec: "Variables of interface type also have a distinct _dynamic type_, which is the (non-interface) type of the value assigned to the variable at run time" (Go spec, Variables), and its example `x = v // x has value (*T)(nil) and dynamic type *T` shows the stored value can be a pointer. The spec has no separate definition of "dynamic value"; it only uses the phrase "has value" in that example.

**Assignability.** Go spec, Assignability: "A value `x` of type `V` is _assignable_ to a variable of type `T` ... if one of the following conditions applies", and one condition is "`T` is an interface type, but not a type parameter, and `x` implements `T`." So the only question is whether `*Employee` implements `Person`.

**Implements depends on the method set.** Go spec, Method sets:

- "The method set of a defined type `T` consists of all methods declared with receiver type `T`."
- "The method set of a pointer to a defined type `T` ... is the set of all methods declared with receiver `*T` or `T`."

`FullName` is declared with receiver `*Employee`, so it is in the method set of `*Employee` and not of `Employee`. Therefore `var p Person = employee` compiles and `var p Person = *employee` does not (the error is shown in [Pointer-receiver gotcha](#pointer-receiver-gotcha), and the receiver choice is covered in [05 - Value vs. pointer](05-value-vs-pointer.md)).

**What `p` holds.** After `var p Person = employee`, `p` is (dynamic type `*Employee`, dynamic value: the pointer). The struct is not copied. Verified on Go 1.27.3 (windows/amd64) with the `Employee` and `Person` from `section08_37` (the `uuid` field and `NewEmployee` left out):

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

**C# contrast.** In C#, `IPerson p = employee;` works because the class lists the interface after a colon in its declaration ("A class or struct lists the interfaces it implements after a colon in its declaration", Microsoft Learn, Interfaces). The existing C# to Go table already has a row for this, so no row was added. In Go there is no such declaration. The compiler checks the method set at the assignment. `p` is like an `IPerson` variable holding a reference: copying the interface value copies the (type, pointer) pair, not the struct.

How the runtime stores the pair (itable) is an implementation detail; see [How the runtime represents it](#how-the-runtime-represents-it-implementation-detail).

## Interfaces can embed other interfaces

Yes. Listing an interface name inside another interface is called *embedding*. Go spec, Embedded interfaces: "an interface `T` may use a (possibly qualified) interface type name `E` as an interface element. This is called _embedding_ interface `E` in `T`" (added in Go 1.14). Effective Go says the same in plain words: "A `ReadWriter` can do what a `Reader` does _and_ what a `Writer` does; it is a union of the embedded interfaces. Only interfaces can be embedded within interfaces."

```go
type Person interface {
	FullName() string
}

type Worker interface {
	Person        // embeds Person: Worker's method set includes FullName
	Deactivate()
}
```

`Worker` has exactly two methods: `FullName() string` and `Deactivate()`. It is the same as writing both out by hand. `*Employee` has both methods, so it implements `Worker` (and `Person`) with no declaration.

**Composition of method sets, not inheritance.** There is no base type, no `base.FullName()` call and no override. The spec describes it with type sets: the type set of `Worker` is "the _intersection_ of the type sets defined by `Worker`'s explicitly declared methods and the type sets of `Worker`'s embedded interfaces", i.e. the types that have all the methods. Effective Go calls the method list a "union". These are the same thing seen from two sides: the method list is a union, the set of qualifying types is an intersection.

**Interface embedding vs. struct embedding.** Both use the same syntax (a bare type name) but are different features. Embedding in a struct adds a field and promotes its methods to the struct; embedding in an interface only adds method requirements. This section covers the interface kind only.

### What converts to what

Verified on Go 1.27.3 (windows/amd64) with `Employee` as in `section08_37` (the `uuid` field replaced by a `[16]byte` type):

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

- `Worker` to `Person` works without a cast. Spec, Implementing an interface: a type `T` implements `I` if "`T` is an interface and the type set of `T` is a subset of the type set of `I`", and Assignability allows assigning a value to an interface type it implements. `Worker` has every method `Person` has, so `Worker` implements `Person`.
- `Person` to `Worker` does not compile, because `Person` lacks `Deactivate`:

```go
var p Person = &Employee{}
var w Worker = p
var x Worker = Worker(p)
```

```
cannot use p (variable of interface type Person) as Worker value in variable declaration: Person does not implement Worker (missing method Deactivate)
cannot convert q (variable of interface type Person) to type Worker: Person does not implement Worker (missing method Deactivate)
```

(The second line came from a variable named `q` in the test program.)

- The compile-time check looks only at the static type `Person`. The type assertion `p2.(Worker)` instead asks at run time whether the value inside the interface (a `*Employee`) has the methods. It succeeded above. Same idea with the standard library: an `io.Reader` holding a `*strings.Reader` asserted to `io.Seeker` gave `true`, to `io.ReadWriter` gave `false` (no panic with the `v, ok` form).

### Overlapping and conflicting methods

```go
type A interface{ M() }
type B interface{ M() }
type C interface {
	A
	B
}
```

This compiles on Go 1.27.3. A `T` with `func (T) M()` was assigned to `C` and then to `A`, and `M` ran on each call. The spec: "When embedding interfaces, methods with the same names must have identical signatures." The spec marks embedding as `[Go 1.14]`; I did not test with an older toolchain.

Same name with a different signature does not compile:

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

The compiler reported this once for `E` and once for `F`. The spec's own example of this case: "illegal: signatures of Reader.Close and Close are different". Declaring the same method twice explicitly in one interface is also illegal: "The name of each explicitly specified method must be unique".

### Embedding a concrete type

Embedding a struct type such as `Employee` in an interface compiles, but it turns the interface into a constraint, not a normal interface:

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

The spec: interfaces that are not basic "may only be used as type constraints, or as elements of other interfaces used as constraints. They cannot be the types of values or variables". So for ordinary interface use, only interfaces can be embedded, which matches Effective Go. Type constraints belong to generics and are out of scope here.

### Standard library example: `io.ReadWriter`

pkg.go.dev/io defines:

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

Verified: `var rw io.ReadWriter = &bytes.Buffer{}` compiled, `var r io.Reader = rw` compiled with no cast, and reading from `r` returned what was written through `rw`.

### Small interfaces

Both `io.ReadWriter` and `Worker` are built from small pieces. The Google Go Style Guide (decisions, Interfaces) says "Design interfaces to be small for easier implementation and composition" (from the fetched text). That is the motivation for embedding: define small interfaces and compose them when a consumer needs more.

### C# to Go for interface composition

| C# | Go |
|---|---|
| `interface IWorker : IPerson` (interfaces "can inherit from one or more other interfaces") | `type Worker interface { Person; Deactivate() }` |
| Multiple interface inheritance: `interface IX : IA, IB` | List several: `interface { A; B }` |
| A class implementing `IShape` "can be implicitly converted to `IDrawable`, because `IShape` inherits from it" | A `Worker` value is assigned to a `Person` variable with no cast, because `Worker` implements `Person` |
| Base to derived: `p as IWorker` / `p is IWorker` | `p.(Worker)` / `w, ok := p.(Worker)` |
| Two interfaces with the same member: explicit implementation when the implementations must differ (`double IMetric.GetDistance()`) | Identical signatures merge into one method that satisfies both. Different signatures: compile error. No explicit implementation. |
| Default interface members: an interface can provide a method body | Not possible: interface elements are method signatures (or constraint terms) |
| `class Employee : IWorker` must name the interface | Nothing in `Employee` names `Worker`; having the methods is enough |

The Go difference: no inheritance. `Worker` is not declared as a subtype of `Person`; it is a list of methods (`Person`'s methods plus its own), and a type implements it by having all of them.

## Sources

- Go spec, Interface types / Embedded interfaces / Implementing an interface: <https://go.dev/ref/spec#Interface_types>, <https://go.dev/ref/spec#General_interfaces>
- Go spec, Method sets: <https://go.dev/ref/spec#Method_sets>
- Go spec, Method declarations (receiver base type in same package): <https://go.dev/ref/spec#Method_declarations>
- Go spec, Assignability: <https://go.dev/ref/spec#Assignability>
- Go spec, Variables (dynamic type): <https://go.dev/ref/spec#Variables>
- Go spec, Type assertions: <https://go.dev/ref/spec#Type_assertions>
- Go FAQ, "Why doesn't Go have 'implements' declarations?", "How can I guarantee my type satisfies an interface?", "Why is my nil error value not equal to nil?": <https://go.dev/doc/faq>
- Effective Go, Embedding: <https://go.dev/doc/effective_go#embedding>
- Effective Go, Interfaces and other types; Interface checks: <https://go.dev/doc/effective_go#interfaces_and_types>, <https://go.dev/doc/effective_go#blank_implements>
- Go Code Review Comments, Interfaces: <https://go.dev/wiki/CodeReviewComments>
- Google Go Style Guide, decisions (Interfaces): <https://google.github.io/styleguide/go/decisions.html>
- Package io (ReadWriter, ReadCloser, ReadWriteCloser): <https://pkg.go.dev/io>
- Russ Cox, Go Data Structures: Interfaces: <https://research.swtch.com/interfaces>
- C# Interfaces (implementing, explicit implementation, default members): <https://learn.microsoft.com/en-us/dotnet/csharp/fundamentals/types/interfaces>
- C# Extension members: <https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/classes-and-structs/extension-methods>

Fact-checked on 2026-11-11 on Go 1.27.3 (windows/amd64). All compiler errors and program output above were produced by building and running the examples with that toolchain.

**Unverified:**
- Interface embedding: the spec, Effective Go and Google Style Guide pages were fetched through a summarizing tool that reads only the first 100000 characters (the spec is about 264000, the style guide about 124000). The Interface types section and the Effective Go Embedding passage were in range; the style guide sentence "Design interfaces to be small for easier implementation and composition" came from a summary and its exact heading and wording were not checked. The Go FAQ and Code Review Comments were not consulted for composition.
- That overlapping methods were an error before Go 1.14 is not stated here as fact; the spec only marks embedding as `[Go 1.14]`, and no older toolchain was run.
- C# default interface members were checked only against the Interfaces overview page, which says they exist and links to a separate article; their behavior was not read.
- The spec pages are long and the fetch tool summarizes them and reads only the first 100000 of about 264000 characters. The Assignability, Method sets, Interface types and Variables quotes in the last section came from those excerpts. The Assignability and Method sets bullets were found in that range, but the quotes were not checked against the full page text.
- That `copying an interface value copies the (type, pointer) pair, not the struct` rests on the run above and the spec's "has value (*T)(nil)" example, not on a spec sentence about copying.
- The spec's Conversions text was fetched only as "x is assignable to T"; the claim that `Person(x)` follows the same implements rule rests on that bullet.
- The C# extension-members page does not state that extension methods cannot make a type implement an interface; the table row only says what the page states (static methods called with instance syntax). The Go side (no methods on foreign types) was verified with the compiler: `cannot define new methods on non-local type time.Time`.
- "Duck typing" in C# (for example the `foreach` pattern) was not checked against Microsoft Learn, so it is not claimed here.
- The exact itable construction timing for static conversions is not stated in the Cox article.
- The "missing method" error for a misspelled method name showed no "have ..." hint in this toolchain; other versions may differ. Compiler message wording is not a language guarantee.
