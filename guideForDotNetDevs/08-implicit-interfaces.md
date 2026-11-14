# 08 - Implicit interfaces (how Go knows `Employee` is a `Person`)

**C#:** a type opts in by name: `class Employee : IPerson`. The compiler checks the members once, at that declaration.

**Go:** no declaration; the compiler checks, where a value is converted to the interface type, whether its type has the interface's methods. Expands the row in [01 - Go is not object-oriented](01-go-is-not-oop.md). What an interface value holds, embedding and conversions continue in [09 - Interface values and embedding](09-interface-values-and-embedding.md).

## Does Go search the methods in the current package?

Three rules, no search:

1. **A type implements an interface if its method set contains the interface's methods.** "A Go type implements an interface by implementing the methods of that interface, nothing more." (Go FAQ).
2. **The check happens where you use the value as the interface.** If it fails, that line does not compile. Effective Go: "passing an `*os.File` to a function expecting an `io.Reader` will not compile unless `*os.File` implements the `io.Reader` interface."
3. **A type's method set is fixed by its own package.** A method's receiver base type "must be declared in the same package as the method" (Go spec, Method declarations).

Different packages (verified): `a` declares `Person` and `Greet(p Person) string`; `b` declares `Employee`; `main` calls `a.Greet(&b.Employee{...})`. Printed `Hello, A B`; neither `a` nor `b` imports anything (`go list`).

## Pointer-receiver gotcha

```go
type Person interface {
	FullName() string
}

type Employee struct { /* ... */ }

func (e *Employee) FullName() string { return e.FirstName + " " + e.LastName }
func (e *Employee) Deactivate()      { e.IsActive = false }
```

Only `*Employee` is a `Person`. Verified on Go 1.27.3:

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

With a value receiver both compile (verified); see [05 - Value vs. pointer](05-value-vs-pointer.md). A missing method gives:

```
cannot use (*Employee)(nil) (value of type *Employee) as Person value in variable declaration: *Employee does not implement Person (missing method FullName)
```

## Compile-time assertion: `var _ Person = (*Employee)(nil)`

Go FAQ ("How can I guarantee my type satisfies an interface?"):

```go
var _ Person = (*Employee)(nil) // Verify that *Employee implements Person.
```

Misspelling `FullName` gives the "missing method" error (verified). Effective Go ("Interface checks"): "such declarations are only used when there are no static conversions already present in the code, which is a rare event."

## Where to define the interface

- Code Review Comments: "Go interfaces generally belong in the package that uses values of the interface type." and "Do not define interfaces before they are used".
- Google Go Style Guide: "The consumer of the interface should define it (not the package implementing the interface)", and "Functions should take interfaces as arguments but return concrete types".

So `NewEmployee` returns `*Employee`, not `Person`.

## Type assertions and the nil gotcha

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

- `x.(T)` panics on failure; `v, ok` does not (Go spec, Type assertions). A type switch picked the right cases.
- An interface holding a nil `*Employee` is not `nil` (FAQ: nil only if both type and value are unset; a nil `*int` gives `(T=*int, V=nil)`).

## C# to Go

| C# | Go |
|---|---|
| Nominal typing, checked at declaration | Structural typing (the FAQ's word), checked at conversion |
| Interface defined with the implementation | Define it in the consumer package |
| Extension methods: static methods called with instance syntax | Cannot add methods to a type from another package |
| Explicit implementation (`double IMetric.GetDistance()`); default interface methods | None; interfaces hold only method signatures |
| `p as IWorker` / `p is IWorker` | `w, ok := p.(Worker)` |

## Sources

- Go spec: <https://go.dev/ref/spec#Method_declarations>, <https://go.dev/ref/spec#Type_assertions>
- Go FAQ: <https://go.dev/doc/faq>
- Effective Go: <https://go.dev/doc/effective_go#interfaces_and_types>, <https://go.dev/doc/effective_go#blank_implements>
- Code Review Comments: <https://go.dev/wiki/CodeReviewComments>
- Google Style Guide: <https://google.github.io/styleguide/go/decisions.html>
- C# Interfaces: <https://learn.microsoft.com/en-us/dotnet/csharp/fundamentals/types/interfaces>
- C# Extension members: <https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/classes-and-structs/extension-methods>

Fact-checked 2026-11-11 on Go 1.27.3 (windows/amd64); all output came from running the examples.

**Unverified:**
- Spec, Effective Go and Style Guide were fetched via a tool reading only the first 100000 characters (spec about 264000, style guide about 124000).
- C# default interface members: only the Interfaces overview page was checked; behavior not read.
- The C# extension-members page does not say extension methods cannot make a type implement an interface. Go side verified: `cannot define new methods on non-local type time.Time`.
- C# "duck typing" (for example `foreach`) was not checked, so not claimed.
- The "missing method" error showed no "have ..." hint here; compiler wording is not a language guarantee.
