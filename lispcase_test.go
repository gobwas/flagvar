package flagvar

import "testing"

func TestLispCase(t *testing.T) {
	for _, test := range []struct {
		name string
		in   string
		exp  string
	}{
		{
			in:  "LispCase",
			exp: "lisp-case",
		},
		{
			in:  "HTTPServer",
			exp: "http-server",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			act := lispCase(test.in)
			if exp := test.exp; act != exp {
				t.Fatalf(
					"unexpected listCase(%q) = %q; want %q",
					test.in, act, exp,
				)
			}
		})
	}
}
