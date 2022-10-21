package main

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_B2Enc(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	benc := newBinEncoder(buf)

	for i := 0; i < 256; i++ {
		buf.Reset()

		fixture := []byte{byte(i)}
		benc.Write(fixture)

		expected := fmt.Sprintf("%.8b", i)
		actual := buf.String()

		t.Logf("%d (%q) => %q", fixture[0], expected, actual)

		if expected != actual {
			t.Fatalf("error: expected %q, got %q | %#v",
				expected, actual, fixture)
		}
	}
}
