package main

import (
	"encoding/ascii85"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-h" {
		os.Stdout.WriteString(usage)
		return
	}

	var out io.Writer = os.Stdout

	if wrap := os.Getenv("R2T_WRAP"); wrap != "" {
		if wrapColumn, err := strconv.Atoi(wrap); err == nil && wrapColumn > 0 {
			out = wrapNewWrap(out, wrapColumn)
		}
	}

	switch os.Getenv("R2T_ENC") {
	case "85", "b85", "a85":
		out85 := ascii85.NewEncoder(out)
		defer out85.Close() // flush remaining parts
		out = out85
	case "64", "b64":
		out64 := base64.NewEncoder(base64.StdEncoding, out)
		defer out64.Close() // flush remaining parts
		out = out64
	case "16", "b16", "hex":
		out16 := hex.NewEncoder(out)
		out = out16
	case "2", "b2":
		out2 := newBinEncoder(out)
		out = out2
	}

	io.Copy(out, os.Stdin)
}
