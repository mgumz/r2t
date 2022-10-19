package main

import (
	"bytes"
	"testing"
)

func Test_Wrap(t *testing.T) {

	type F struct {
		cols     int
		err      error
		in       []string
		expected string
	}

	fixtures := []F{
		{3, nil, []string{"a"}, "a"},
		{3, nil, []string{"abcd"}, "abc\nd"},
		{3, nil, []string{"ab", "cd"}, "abc\nd"},
		{3, nil, []string{"abcdefghijk"}, "abc\ndef\nghi\njk"},
	}

	b := bytes.NewBuffer(nil)

	for i, f := range fixtures {
		t.Logf("%d -- %#v --", i, f)
		b.Reset()
		w := wrapNewWrap(b, f.cols)
		nb := 0
		for _, in := range f.in {
			n, err := w.Write([]byte(in))
			if err != f.err {
				t.Fatalf("error: expected %s, got %s: %#v", f.err, err, f)
			}
			nb += n
		}
		if b.String() != f.expected {
			t.Fatalf("error, expected\n\n%q\n\ngot\n\n%q\n\n%#v",
				f.expected, b.String(), f)
		}
	}
}
