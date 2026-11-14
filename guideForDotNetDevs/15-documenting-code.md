# 15 - Documenting code (`///` XML comments vs. Go doc comments)

**C#:** you write `///` comments with XML tags (`<summary>`, `<param>`, `<returns>`). The compiler can
produce an XML documentation file from them ("The compiler finds all comment fields with XML tags in the
source code and creates an XML documentation file", Microsoft Learn).

**Go:** there is one standard convention, the **doc comment**. It is an ordinary `//` comment with no tags.
Examples use `bankAccount`.

| C# | Go |
|---|---|
| `/// <summary>` | A `//` comment directly above the declaration |
| `<param>`, `<returns>` | No tags: say it in the sentence |
| `<see cref="X"/>` | `[X]` doc link (`[pkg.Name]` across packages) |
| `<example>` | An `Example` function with an `// Output:` comment |
| `[Obsolete]` | A paragraph starting `Deprecated:` |
| XML file from the compiler | `go doc`, pkg.go.dev, and `gopls` hover read the comments directly |

## The rule

"Doc comments are comments that appear immediately before top-level package, const, func, type, and var
declarations with no intervening newlines. Every exported (capitalized) name should have a doc comment" (Go
Doc Comments). Code Review Comments (Doc Comments) adds "as should non-trivial unexported type or function
declarations".

Write full sentences that begin with the name: "Comments should begin with the name of the thing being
described and end in a period" (Code Review Comments, Comment Sentences). Google decisions agrees and
allows a leading article ("a", "an", "the").

```go
// Deposit adds amount to the balance. It returns an error, and leaves the balance
// unchanged, if amount is negative.
func (a *Account) Deposit(amount float64) error {
```

Describe what the code **does**, including failure cases, since no `<exception>` or `<returns>` tag
carries them. Google best practices: "Not every parameter must be enumerated" and
document "the error-prone or non-obvious fields and parameters by saying why they are interesting".

## What the tools show

`go doc` "looks up and prints the doc comment for a given package or symbol" (Go Doc Comments). Real
output from `go doc Account.Deposit` on Go 1.27.3:

```
func (a *Account) Deposit(amount float64) error
    Deposit adds amount to the balance. It returns an error, and leaves the
    balance unchanged, if amount is negative.
```

`go doc -all .` on `bankAccount` listed the exported fields of `Account` and replaced the unexported
`balance` with `// Has unexported fields.`. It did not list `Deposit` and `Balance` again under
`SavingsAccount`, so say in the type comment that they are promoted (doc 01).

## Viewing the docs as HTML

`go help doc` on Go 1.27.3 lists `go doc -http`: "Serve HTML documentation over HTTP for the current
package." Run it from the package directory. It serves on a random local port, and the page looks like
pkg.go.dev. Observed in Chrome with the `bankAccount` code:

- **As a library package** (`package bankaccount`, same code and comments, no `main()`): the page showed the
  package comment under Overview, an Index, and every type, function and method with its comment. The
  unexported `balance` field was replaced by `// contains filtered or unexported fields`. `SavingsAccount`
  listed only its own members (`NewSavingsAccount`, `AddInterest`, `String`), not the promoted `Deposit`,
  `Withdraw` and `Balance`.
- **As `package main`** (the real `bankAccount`): the page was labelled "command" and showed only the
  package comment and the source file list. No types or methods appeared.

So API docs only render in HTML once the code lives in an importable package, which fits the package split
in doc 13. Why the viewer hides a command's API is **Unverified**: the behavior was reproduced, but no
documentation was found that explains it.

## Package, field, and example comments

- **Package comment:** "Every package should have a package comment introducing the package", its first
  sentence begins with `Package `, and in a multi-file package it should be in one file (Go Doc Comments).
- **Struct fields:** "either the doc comment or per-field comments should explain the meaning of each
  exported field".
- **Examples:** "Godoc examples are snippets of Go code that are displayed as package documentation and
  that are verified by running them as tests" (Go blog, Testable Examples). An `Example` function ending in
  `// Output:` is compared with its standard output when tests run (testing package).

## Syntax

There is no HTML or XML. "Go doc comments are written in a simple syntax that supports paragraphs,
headings, links, lists, and preformatted code blocks ... Markdown aficionados can view the syntax as a
simplified subset of Markdown." `gofmt` rewrites doc comments to a canonical form. "Paragraphs starting with
`Deprecated:` are treated as deprecation notices ... pkg.go.dev will hide their docs by default."

## Checking that comments exist

`go vet` has no doc-comment check (its analyzer list has none). `revive`'s `exported` rule "warns on
undocumented exported functions and methods", and staticcheck has `ST1000` (package comment) and
`ST1020` to `ST1022` (exported names), all non-default. `golint` is "deprecated and frozen".

## Sources

- Go Doc Comments: <https://go.dev/doc/comment>
- Effective Go, Commentary (now defers to Go Doc Comments): <https://go.dev/doc/effective_go#commentary>
- Go Code Review Comments (Doc Comments, Comment Sentences): <https://go.dev/wiki/CodeReviewComments>
- Google Go Style Guide, decisions and best practices:
  <https://google.github.io/styleguide/go/decisions>, <https://google.github.io/styleguide/go/best-practices>
- `go doc`: <https://pkg.go.dev/cmd/go#hdr-Show_documentation_for_package_or_symbol>
- Testable examples: <https://go.dev/blog/examples>, <https://pkg.go.dev/testing#hdr-Examples>
- staticcheck checks: <https://staticcheck.dev/docs/checks/>; revive rules:
  <https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md>
- C# XML documentation comments: <https://learn.microsoft.com/dotnet/csharp/language-reference/xmldoc/>

Fact-checked on 2026-11-13 on Go 1.27.3: the quotes were taken from the raw page text (fetched with `curl`,
not the fetch tool's summary), and the `go doc` output is real.

**Verified:** everything quoted above, `go doc` behavior on `bankAccount`, `go vet` having no doc-comment
check, the `go doc -http` flag text, and the HTML rendering results (observed once, in Chrome).

**Unverified:**
- Why `go doc -http` shows only an overview for a `package main` command.
- Whether a `package main` command should use the `Package main ...` form; Go Doc Comments' wording for
  commands was not checked. The `bankAccount` comment uses the `Package ` form from the quoted rule.
- An explicit rule that interface methods need comments.
- Whether `Example` functions must live in `_test.go` files (not stated in the text read).
- Whether staticcheck's `ST1020` to `ST1022` fire on a comment that is missing entirely, and which of
  these linters `golangci-lint` enables by default.
- The standalone `godoc` command's current status. Only `go doc` was checked.
- Microsoft Learn's `<param>` tag text was not extracted (`<returns>` was).
