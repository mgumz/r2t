package main

import (
	"encoding/ascii85"
	"encoding/base64"
	_ "encoding/hex"
	"io"
	"os"
	"strconv"
)

const self = `
H4sIAD8WzFcCA5VU30/bMBB+Tv6KIw9TspVAC0NTtk5CgPY2pjGpD4CmkFxaa4kd2U4DGvzvu7NT0s
L2MEvU4X59d9+dr82LX/kSocmFDEPRtEpbiMMgQlmoUsjlQW4KIT68j7Zld7nBk2MS/YRRuMJ7NhKK
f5XhX2N1oeQ6CpMwrDpZOJg4gd9hGKxzDaqzIFS60MKihjkok17ZkqSkFxX0Om8hc+IvaFGu42jx/f
RblHz0qr05RBEFCzbGZ6ruGjkB1Jr9Bvj01CoRs54cWTWfgxS1cww4hTm8WZDWp/F7kXFeE/DBsq3A
T+RAf0+UnemFLVbbmV18PYu4sqAgciAixiYQ3fkrpysLHdiH95zZQGr6FfsL5g91TLqETEqs0BFD2r
NaGSS6Dg6gqjuzAo1MIJENba6tCTfpO/MNMjWGkekaIE+OGdL3bBtxkBDjF0MPJ/AyC9L/TxYnx5ss
piecBc8EZ0G+Py7PLzfK2a6QCaUxOFPtQ+yo93MgZBKSyj60CGN/uKtdYZnpBbjzPEFh4PsEQroJ4p
vTvnuwaKiPwlKrSAh2hVB0WiPp72pV/GIcN6BxD29HrATcHbdwfctBEogppJsvpf0YU/i2zmXmvqYp
ddZ0DUKPlCexs0FrRFnWCKqCfEBkezoa8xKaIUWroBJ1DV3rnFqNa6E6s7EdHWcpLBWoNdGBOc2h00
zIgStmVx7alhnxnnlld0yd/CiFfkWE1JivXfMNUr4v2KqU3qbrdSq4dLX7+qfTNE0f3fe+yI72/3U2
YWaz2ePMHbrS0XfniGz2N9+dc0Tn8cgdF4I6Q0/fWEfEOdaioQGNbmQEw9RlsBCyVL2B/c+suNER+V
hl85pfy2EYuBL5uyUF08DtZmypoFGaOFbaca4KNAZLQGrHg13xw6i0aqAFWhIrfhtuXozfUzXK2IVO
eA8d+i2k0XaaJoXhJ7yc3KbxcAPAi4ENgoZz69Nh5PfpU5DULb7GQ2n4vA3nkYiHUeYXWjCuzD5d+M
fkDa4zfctGnpd3tDdDH5rN97bW6G4BpN0tYKm5AhpvHiWJ92MVnuW5//9aZ7eDk1vwBCQRSyw9qIRP
9FQcHhVLTjTtseT8uNlWyA6Hgrz6kGP99LXNn9fElaVXsoyp0sk4Hcn/lOUqe92yp/AP/na9104HAA
A=`

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
