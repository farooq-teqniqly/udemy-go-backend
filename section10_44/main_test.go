package main

import (
	"strconv"
	"strings"
	"testing"
)

var sizes = []int{10, 100, 1000, 10000}

// growSize is an upper bound on the output length for 1..n: every number has
// at most len(Itoa(n)) digits, plus one comma each.
func growSize(n int) int {
	return n * (len(strconv.Itoa(n)) + 1)
}

// stringConcat is the baseline: each += allocates a new string and copies the
// old contents, so the total work grows quadratically with count.
func stringConcat(count int) string {
	s := ""
	for n := range count {
		if n > 0 {
			s += ","
		}
		s += strconv.Itoa(n + 1)
	}
	return s
}

func builderConditional(count int) string {
	b := strings.Builder{}
	for n := range count {
		if n > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.Itoa(n + 1))
	}
	return b.String()
}

func builderTrimSuffix(count int) string {
	b := strings.Builder{}
	for n := range count {
		b.WriteString(strconv.Itoa(n + 1))
		b.WriteString(",")
	}
	return strings.TrimSuffix(b.String(), ",")
}

func stringsJoin(count int) string {
	parts := make([]string, 0, count)
	for n := range count {
		parts = append(parts, strconv.Itoa(n+1))
	}
	return strings.Join(parts, ",")
}

func builderGrow(count int) string {
	b := strings.Builder{}
	b.Grow(growSize(count))
	for n := range count {
		if n > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.Itoa(n + 1))
	}
	return b.String()
}

func builderAppendInt(count int) string {
	b := strings.Builder{}
	b.Grow(growSize(count))
	buf := make([]byte, 0, 8)
	for n := range count {
		if n > 0 {
			b.WriteByte(',')
		}
		buf = strconv.AppendInt(buf[:0], int64(n+1), 10)
		b.Write(buf)
	}
	return b.String()
}

var variants = []struct {
	name string
	fn   func(count int) string
}{
	{"StringConcat", stringConcat},
	{"BuilderConditional", builderConditional},
	{"BuilderTrimSuffix", builderTrimSuffix},
	{"StringsJoin", stringsJoin},
	{"BuilderGrow", builderGrow},
	{"BuilderAppendInt", builderAppendInt},
}

func wantString(count int) string {
	var s string
	for n := 1; n <= count; n++ {
		if n > 1 {
			s += ","
		}
		s += strconv.Itoa(n)
	}
	return s
}

func TestAllVariantsAgree(t *testing.T) {
	for _, size := range sizes {
		want := wantString(size)
		for _, v := range variants {
			if got := v.fn(size); got != want {
				t.Errorf("%s(%d) = %q, want %q", v.name, size, got, want)
			}
		}
	}
}

var sink string

func BenchmarkJoinNumbers(b *testing.B) {
	for _, v := range variants {
		b.Run(v.name, func(b *testing.B) {
			for _, size := range sizes {
				b.Run("n="+strconv.Itoa(size), func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						sink = v.fn(size)
					}
				})
			}
		})
	}
}
