package colon_test

import (
	"strings"
	"testing"

	. "trib/colon"
)

func TestColon(t *testing.T) {
	lb := func(s string) {
		esc := Escape(s)
		got := Unescape(esc)
		if s != got {
			t.Errorf("loopback failed on %q, got %q", s, got)
		}

		if strings.Index(esc, ":") != -1 {
			t.Errorf("found colon in escaped string %q for %q", esc, s)
		}

		t.Logf("esc(%q) = %q", s, esc)
	}

	lb(`|`)
	lb(`||`)
	lb(`|||`)
	lb(`a|:a`)
	lb(`a::a`)
	lb(`a:|a`)
	lb(`    `)
	lb(`::||::||;;||;;||;:`)
}
