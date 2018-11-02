package main

import (
	"fmt"
	"io"
)

type PaddedBuffer struct {
	i   int
	buf []byte
}

func NewPaddedBuffer(size int) *PaddedBuffer {
	return &PaddedBuffer{i: 0, buf: make([]byte, size)}
}

func (b *PaddedBuffer) Read(p []byte) (n int, err error) {
	n = copy(p, b.buf[b.i:])
	b.i += n
	return n, nil
}

func (b *PaddedBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	m, err := r.Read(b.buf)
	n = int64(m)
	b.i = 0
	return n, err
}

func (b *PaddedBuffer) Write(p []byte) (n int, err error) {
	n = copy(b.buf[b.i:], p)
	b.i += n

	if n < len(p) {
		err = fmt.Errorf("Wanted to write %v bytes but only %v available", len(p), n)
	} else {
		err = nil
	}
	return n, err
}

func (b *PaddedBuffer) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write(b.buf)
	n = int64(m)
	b.i = 0
	return n, err
}
