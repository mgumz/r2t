package main

import (
	"encoding/ascii85"
	"encoding/base64"
	_ "encoding/hex"
	"io"
	"os"
	"strconv"
)

func main() {

	var out io.Writer = os.Stdout

	if wrap := os.Getenv("WRAP"); wrap != "" {
		if wrapColumn, err := strconv.Atoi(wrap); err == nil {
			out = &WrapWriter{W: out, Column: wrapColumn}
		}
	}

	switch os.Getenv("ENC") {
	case "85", "b85", "a85":
		out85 := ascii85.NewEncoder(out)
		defer out85.Close() // flush remaining parts
		out = out85
	case "64", "b64":
		out64 := base64.NewEncoder(base64.StdEncoding, out)
		defer out64.Close() // flush remaining parts
		out = out64
	case "16", "hex":
		// TODO
	case "2":
		// TODO
	}

	io.Copy(out, os.Stdin)
}

type WrapWriter struct {
	W      io.Writer
	Column int

	i int // bytes written in the current block
}

func (w *WrapWriter) Write(p []byte) (int, error) {

	// plan:
	// 1. assume we start in the middle of a block
	//    read m bytes to fill up the previous
	//    block
	// 2. go over each block, print the wrapper
	//    after each block
	// 3. when leaving, set i bytes written for the current
	//    block
	// eg:
	//
	// 111...|
	// -i:3-------------------------
	//    222|222222|22....|
	// --------------i:2------------
	//                 3333|333333|
	//

	const wrapDelim = "\n" // TODO: Windows -> "\n\r"

	total := 0
	block := p

	for {

		// no more work, processed everything from p without errors
		if len(block) == 0 {
			return total, nil
		}

		// process current block
		m := w.Column - w.i
		r := m
		if r > len(block) {
			r = len(block)
		}
		n, err := w.W.Write(block[:r])
		total += n

		if err != nil {
			return total, err
		}

		// progress to the next block
		block = block[r:]

		// wrap if needed
		if n < m {
			w.i = int(n)
			continue
		}
		w.i = 0

		_, err = io.WriteString(w.W, wrapDelim)
		if err != nil {
			return total, err
		}
	}

	return total, nil
}
