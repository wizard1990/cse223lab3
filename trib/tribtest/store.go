package tribtest

import (
	"runtime/debug"
	"sort"
	"testing"

	"trib"
)

func CheckStorage(t *testing.T, s trib.Storage) {
	var v string
	var b bool
	var l = new(trib.List)
	var n int

	ne := func(e error) {
		if e != nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			debug.PrintStack()
			t.Fatal("assertion failed")
		}
	}

	kv := func(k, v string) *trib.KeyValue {
		return &trib.KeyValue{k, v}
	}

	pat := func(pre, suf string) *trib.Pattern {
		return &trib.Pattern{pre, suf}
	}

	v = "_"
	ne(s.Get("", &v))
	as(v == "")

	v = "_"
	ne(s.Get("hello", &v))
	as(v == "")

	ne(s.Set(kv("h8liu", "run"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "run")

	ne(s.Set(kv("h8liu", "Run"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "Run")

	ne(s.Set(kv("h8liu", ""), &b))
	as(b)
	v = "_"
	ne(s.Get("h8liu", &v))
	as(v == "")

	ne(s.Set(kv("h8liu", "k"), &b))
	as(b)
	v = "_"
	ne(s.Get("h8liu", &v))
	as(v == "k")

	ne(s.Set(kv("h8he", "something"), &b))
	as(b)
	v = "_"
	ne(s.Get("h8he", &v))
	as(v == "something")

	ne(s.Keys(pat("h8", ""), l))
	sort.Strings(l.L)
	as(len(l.L) == 2)
	as(l.L[0] == "h8he")
	as(l.L[1] == "h8liu")

	ne(s.ListGet("lst", l))
	as(l.L != nil && len(l.L) == 0)

	ne(s.ListAppend(kv("lst", "a"), &b))
	as(b)

	ne(s.ListGet("lst", l))
	as(len(l.L) == 1)
	as(l.L[0] == "a")

	ne(s.ListAppend(kv("lst", "a"), &b))
	as(b)

	ne(s.ListGet("lst", l))
	as(len(l.L) == 2)
	as(l.L[0] == "a")
	as(l.L[1] == "a")

	ne(s.ListRemove(kv("lst", "a"), &n))
	as(n == 2)

	ne(s.ListGet("lst", l))
	as(l.L != nil && len(l.L) == 0)

	ne(s.ListAppend(kv("lst", "h8liu"), &b))
	as(b)
	ne(s.ListAppend(kv("lst", "h7liu"), &b))
	as(b)

	ne(s.ListGet("lst", l))
	as(len(l.L) == 2)
	as(l.L[0] == "h8liu")
	as(l.L[1] == "h7liu")

	ne(s.ListKeys(pat("ls", "st"), l))
	as(len(l.L) == 1)
	as(l.L[0] == "lst")

	ne(s.ListKeys(pat("z", ""), l))
	as(l.L != nil && len(l.L) == 0)

	ne(s.ListKeys(pat("", ""), l))
	as(len(l.L) == 1)
	as(l.L[0] == "lst")
}
