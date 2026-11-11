# 01 - Go is not object-oriented (the way C# is)

Many of the differences in the other docs trace back to one fact: Go has no classes. This is a mental
model summary; the details are in the linked docs.

| C# | Go |
|---|---|
| `class` (reference type) | No equivalent. Use a `struct` and pass `*T` when you need sharing. |
| `struct` (value type) | `struct`, copied on assignment. |
| Methods inside the class body | Functions with a receiver, declared outside the type. |
| `this` | Explicit receiver name, e.g. `func (e Employee) FullName()`. |
| Constructors | Ordinary functions by convention: `NewEmployee(...)`. |
| Inheritance | None (no type hierarchy). Use composition; struct embedding is "analogous—but not identical—to subclassing". |
| `interface` implemented explicitly (`: IFoo`) | Satisfied implicitly. No `implements` keyword. |
| `public` / `private` keywords | No keywords. A capitalized name is exported (`IsActive`), a lowercase one is not (`isActive`). Privacy is per package, not per type. See [06](06-encapsulation-and-private-set.md). |
| Method overloading, optional params | Not supported. See [04](04-default-parameters.md). |
| `null` on any reference | `nil` only for pointers, slices, maps, channels, funcs, interfaces. |

## Methods are attached to types, not declared inside them

```go
type Employee struct {
	FirstName string
	LastName  string
}

// Receiver (e Employee) plays the role of `this`.
func (e Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}
```

Whether the receiver is `Employee` or `*Employee` decides if the method works on a copy or the
original. See [05 - Value vs. pointer](05-value-vs-pointer.md).

## Constructors are a convention

Go has no constructor syntax. `NewEmployee` is just a function that returns a value. Effective Go advises
designing types so the zero value is usable without further initialization (`var e Employee`), and says an
initializing constructor is needed "sometimes" when the zero value isn't good enough, such as generating an
ID. This is a design convention, not a guarantee for every struct.

## Sources

- Go FAQ, "Is Go an object-oriented language?": <https://go.dev/doc/faq#Is_Go_an_object-oriented_language>
  (no type hierarchy; embedding is analogous but not identical to subclassing; methods can be defined on any
  type, not only structs)
- Go FAQ, "Why doesn't Go have 'implements' declarations?": <https://go.dev/doc/faq#implements_interface>
  ("A Go type implements an interface by implementing the methods of that interface, nothing more.")
- Go spec, Method declarations: <https://go.dev/ref/spec#Method_declarations> (a method is a function with a
  receiver; the receiver base type must be declared in the same package)
- Go spec, Interface types (Implementing an interface): <https://go.dev/ref/spec#Interface_types>
- Go spec, Struct types (embedded fields) and Zero value / Predeclared identifiers (`nil` is the zero value
  for types whose values can contain references): <https://go.dev/ref/spec>
- Effective Go, "Allocation with new" and "Constructors and composite literals":
  <https://go.dev/doc/effective_go>

Fact-checked against these pages on 2026-11-11.
