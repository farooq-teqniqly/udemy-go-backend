# 03 - Date-only types

**C#:** `DateOnly` (and `TimeOnly`) since .NET 6.

**Go:** no date-only type in the standard library, including 1.27. `time.Time` is the only date/time
type.

## Options

### 1. `time.Time` at midnight UTC (simplest)

```go
type Employee struct {
	JoinedAt time.Time // convention: midnight UTC, date part only
}

joinedAt := time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)
joinedAt, err := time.Parse(time.DateOnly, "2024-03-15")
```

`time.DateOnly` is a layout string (`"2006-01-02"`), not a type. It was added in Go 1.20, along with
`time.DateTime` and `time.TimeOnly` (Go 1.20 release notes, "time"). Nothing stops a non-midnight value, so
the "date only" rule is a convention.

### 2. `civil.Date` from `cloud.google.com/go/civil`

```go
JoinedAt civil.Date

civil.Date{Year: 2024, Month: time.March, Day: 15}
civil.ParseDate("2024-03-15")
```

A struct of `Year int`, `Month time.Month`, `Day int` with `AddDays`, `Compare`, `In(loc)`,
`MarshalText`/`UnmarshalText`, and `Scan`/`Value` (so it works with `database/sql`). It lives in the
`cloud.google.com/go` module, which is third-party and pre-v1 (v0.123.0 when checked).

### 3. A string or your own named type

Not researched.

## Why it's not in the stdlib

The civil date/time proposal, <https://github.com/golang/go/issues/19700>, is still open (milestone
Proposal, last updated 2025-04-19, 79 comments). In April 2025 a Go team member (neild) said the open question
is whether the types go in `time` or a new package, and jba noted `cloud.google.com/go/civil` drags in
dependencies. No decision was made. The Go 1.27 release notes list no new date type.

## Sources

- `time` package: <https://pkg.go.dev/time>
- `civil` package: <https://pkg.go.dev/cloud.google.com/go/civil>
- Civil date proposal: <https://github.com/golang/go/issues/19700>
- ISO 8601 date constant request (closed as completed 2020-02-03; the constant arrived later as `time.DateOnly`): <https://github.com/golang/go/issues/36650>
- Go 1.27 release notes: <https://go.dev/doc/go1.27>

**Verified by running on Go 1.27.3** (scratch module, `cloud.google.com/go` v0.123.0 fetched fine):
- `time.Parse(time.DateOnly, "2024-03-15")` gives `2024-03-15 00:00:00 +0000 UTC`.
- A midnight-UTC `time.Time` marshals to JSON as the full timestamp `"2024-03-15T00:00:00Z"`.
- `civil.Date` marshals to JSON as `"2024-03-15"`.

**SQL `DATE` into `time.Time`:** `database/sql` allows scanning a `time.Time` source into `*time.Time`
(<https://pkg.go.dev/database/sql#Rows.Scan>). For go-sql-driver/mysql (the driver checked), `DATE` and
`DATETIME` come back as `[]byte` unless the DSN has `parseTime=true`, which returns `time.Time`; `loc`
(default UTC) sets its location. The result is still a `time.Time`, so the date-only rule stays a convention.

**Unverified:** pgx and lib/pq behavior for `DATE` (the pgx docs page was too large to read the `Date` scan
rules). The full text of issue #19700 comments was only skimmed.
