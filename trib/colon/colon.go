package colon

import (
	"bytes"
	"strings"
)

func Escape(s string) string {
	s = strings.Replace(s, `|`, `||`, -1) // double all backslash
	s = strings.Replace(s, `:`, `|;`, -1) // and the replace all colons
	return s
}

func Unescape(s string) string {
	out := new(bytes.Buffer)
	esc := false

	for _, r := range s {
		if !esc {
			if r == '|' {
				esc = true
			} else {
				out.WriteRune(r)
			}
		} else {
			// escaping
			if r == ';' {
				out.WriteRune(':')
			} else if r == '|' {
				out.WriteRune('|')
			} else {
				// should not happen
				out.WriteRune(r)
			}

			esc = false
		}
	}

	return out.String()
}
