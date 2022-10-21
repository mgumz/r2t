package main

import (
	"io"
)

type binEncoder struct {
	w io.Writer
}

func newBinEncoder(w io.Writer) *binEncoder {
	be := new(binEncoder)
	be.w = w
	return be
}

func (be *binEncoder) Write(p []byte) (n int, err error) {

	const b2digits = "01"

	bits := [8]byte{}
	for _, b := range p {

		bits[0] = b2digits[(b&0x80)>>7]
		bits[1] = b2digits[(b&0x40)>>6]
		bits[2] = b2digits[(b&0x20)>>5]
		bits[3] = b2digits[(b&0x10)>>4]
		bits[4] = b2digits[(b&0x08)>>3]
		bits[5] = b2digits[(b&0x04)>>2]
		bits[6] = b2digits[(b&0x02)>>1]
		bits[7] = b2digits[(b&0x01)>>0]

		n2, err2 := be.w.Write(bits[:])
		n, err = n+n2, err2
		if err != nil {
			return
		}
	}

	return
}
