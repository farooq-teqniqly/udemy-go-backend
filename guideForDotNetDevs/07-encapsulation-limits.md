# 07 - Where Go encapsulation is not airtight

**C#:** a `private` field or `private` constructor is enforced by the compiler against everything outside the class, including derived classes.

**Go:** "unexported" is enforced per **package**, and several language features reach around it. Continues [06 - Encapsulation](06-encapsulation-and-private-set.md); `Employee` and `employee.New` are from there.

## Where it is NOT airtight

All run (`cmd/ext`, `cmd/wipe`):

1. **The zero value is always constructible.** `var z employee.Employee` compiles anywhere: `IsActive()`
   false, `Salary()` 0, `ID()` all zeros. A constructor cannot force its use (C# can).
2. **Value copies are the caller's own.** `v := *p2; v.Deactivate()` changes only `v`; `p2` stays active
   (the copy includes unexported fields); with `q := p`, `q.Deactivate()` changes `p`. See doc 05.
3. **Whole-struct assignment overwrites hidden fields.** `*p = employee.Employee{FirstName: "X"}` compiles
   outside the package and zeroes `isActive`, `salary` and `id`.
4. **Anything in the same package can touch them.** `employee.Tamper(e)` changed `isActive` and `salary`.
   So `section08_36` (a single `package main`) cannot demonstrate encapsulation. Google decisions: "Code
   within a package can access unexported identifiers in the package" (best-practices, Package size).
5. **Reflection** can read unexported fields but not set them. `CanSet()` was `false` for `isActive`, `true`
   for `FirstName`; `SetBool` panicked with `reflect: reflect.Value.SetBool using
   value obtained using unexported field`. `Value.Interface()` also panics (pkg.go.dev/reflect). `unsafe`
   can bypass this (not tested).
6. **`encoding/json` skips unexported fields.** `json.Marshal(p2)` printed
   `{"FirstName":"C","LastName":"D"}` (`id`, `salary`, `isActive`, `joinedAt` dropped; unmarshalling cannot
   fill them). To serialise one, export it or write `MarshalJSON`.
7. **Embedding adds no `protected`.** With `SavingsAccount` embedding `Account` in `package account`, its
   methods can read `sa.balance` (compiled and ran). From another package,
   Go 1.27.3 reports `sa.balance undefined (type account.SavingsAccount has no field or method balance, but does have
   method Balance)`, while `sa.Balance()` works.

## Sources

- Google Go Style Guide, best practices (Package size): <https://google.github.io/styleguide/go/best-practices.html>
- reflect (`Value.CanSet`): <https://pkg.go.dev/reflect#Value.CanSet>
- encoding/json (`Marshal`): <https://pkg.go.dev/encoding/json#Marshal>

Fact-checked on 2026-11-11 on Go 1.27.3 (windows/amd64): the Google quote is a summarising-fetch-tool excerpt; every
behavior and compiler message comes from a run in a scratch module with a separate `employee` (items 1 to 6) or `account` (item 7) package.

**Verified:** all six behaviors in items 1 to 6; item 7 (added 2026-11-13).

**Unverified:**
- That `unsafe` can set unexported fields was not tested or sourced.
- `encoding/json` doc says "Each exported struct field becomes a member of the object"; the explicit "unexported fields are ignored" wording was not seen. Behavior confirmed by running it.
- The spec's wording on struct copying and exportedness was not seen: the Selectors and Assignment sections gave no rule on either. Item 2 rests on the run.
