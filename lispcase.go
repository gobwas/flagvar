package flagvar

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func lispCase(s string) string {
	r0, n := utf8.DecodeRuneInString(s)
	if r0 == utf8.RuneError {
		return s
	}
	var (
		rs int
		sb strings.Builder
	)
	for _, r1 := range s[n:] {
		if rs > 0 && unicode.IsUpper(r0) && !unicode.IsUpper(r1) {
			sb.WriteByte('-')
		}
		sb.WriteRune(unicode.ToLower(r0))
		r0 = r1
		rs++
	}
	sb.WriteRune(unicode.ToLower(r0))
	return sb.String()
}
