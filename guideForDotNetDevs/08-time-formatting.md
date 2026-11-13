# 08 - Time formatting layouts

**C#:** format strings are letter tokens — `"yyyy-MM-dd"`, `"HH:mm:ss"` — each letter is a
case-sensitive placeholder defined by the format-string spec. Tokens don't have to be a real date;
`y`, `M`, `d` are just symbols.

**Go:** `time.Format`/`time.Parse` take a *layout string* built from one specific reference
instant:

```
Mon Jan 2 15:04:05 MST 2006
```

This is not an example you can substitute — it's the only reference value the formatter
recognizes. `Format`/`Parse` scan the layout for its literal substrings (`"2006"`, `"01"`, `"02"`,
`"15"`, `"04"`, `"05"`, `"MST"`, `"Mon"`, `"Jan"`, etc.) and swap each for the real value in that
slot. Any other date or number in a layout string is just literal text, not a placeholder.

```go
joinedAt.Format("2006-01-02")          // "2025-11-17"
joinedAt.Format("Jan 2, 2006 3:04 PM") // "Nov 17, 2025 12:00 AM"
```

## The mnemonic

The reference instant's fields read as `1 2 3 4 5 6 -7` in order: month=1 (`01`), day=2 (`02`),
hour=3 as `15` (3 PM in 24h), minute=4 (`04`), second=5 (`05`), year=06 (`06` or `2006`), timezone
offset=-7 (`MST` / `-0700`). Written as a timestamp it's `01/02 03:04:05PM '06 -0700`.

## Source

Go's `time` package docs, "Reference time" (`pkg.go.dev/time#pkg-constants`): "The reference time
used in these layouts is the specific time... Jan 2 15:04:05 2006 UTC-7... Only these values are
recognized" when building a layout.

- `time` package constants: <https://pkg.go.dev/time#pkg-constants>
- `Time.Format`: <https://pkg.go.dev/time#Time.Format>
