package main

import (
	"io"
)

type wrapWriter struct {
	w   io.Writer
	col int
	i  int
}

func wrapNewWrap(w io.Writer, col int) (ww *wrapWriter) {
	ww = new(wrapWriter)
	ww.w = w
	ww.col = col
	return ww
}

func (ww *wrapWriter) Write(p []byte) (nw int, err error) {

	const wrapDelim = "\n"

	n := 0
	for l := 0; len(p) > 0; p = p[l:] {

		// pick up to next column or
		// end of current p
		if l = ww.col - (ww.i % ww.col); len(p) < l {
			l = len(p)
		}

		// write payload up to ww.col
		n, err = ww.w.Write(p[:l])
		nw += n
		ww.i += n
		if err != nil {
			break
		}

		// write col-delimiter
		if (ww.i % ww.col) == 0 {
			n, err = ww.w.Write([]byte(wrapDelim))
			nw += n
			if err != nil {
				break
			}
		}
	}

	return nw, err
}
