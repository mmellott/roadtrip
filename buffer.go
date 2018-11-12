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

func (b *PaddedBuffer) Write(p []byte) (n int, err error) {
	n = copy(b.buf[b.i:], p)
	b.i += n

	if n < len(p) {
		return n, fmt.Errorf("Wanted to write %v bytes but only %v available", len(p), n)
	} else {
		return n, nil
	}
}

func (b *PaddedBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.buf)
	b.i -= n
	if b.i < 0 {
		b.i = 0
	}
	return int64(n), err
}
