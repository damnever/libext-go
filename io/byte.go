package io

import (
	"io"
)

// ByteReadWriter implements io.ByteReader and io.ByteWriter.
type ByteReadWriter struct {
	ByteReader
	ByteWriter
}

// ByteReader implements io.ByteReader.
type ByteReader struct {
	io.Reader

	p [1]byte // We can use slice to avoid copy, but..
}

// NewByteReader creates new ByteReader.
func NewByteReader(r io.Reader) *ByteReader {
	return &ByteReader{Reader: r}
}

func (r *ByteReader) ReadByte() (b byte, err error) {
	p := r.p[:]
	if _, err = r.Reader.Read(p); err != nil {
		return
	}
	b = p[0]
	return
}

// ByteWriter implements io.ByteWriter.
type ByteWriter struct {
	io.Writer
}

// NewByteWriter creates new ByteWriter.
func NewByteWriter(w io.Writer) *ByteWriter {
	return &ByteWriter{Writer: w}
}

func (w *ByteWriter) WriteByte(b byte) error {
	_, err := w.Writer.Write([]byte{b})
	return err
}
