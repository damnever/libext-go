package io

import (
	"bufio"
	"io"
)

type (
	Flusher interface {
		Flush() error
	}

	WithReadWriter struct {
		io.Writer
		io.Reader
	}
	WithReadCloser struct {
		io.Reader
		io.Closer
	}
	WithWriteCloser struct {
		io.Writer
		io.Closer
	}
	WithCloser struct {
		io.ReadWriter
		io.Closer
	}
	BufferedReadWriterWithCloser struct {
		*bufio.ReadWriter
		io.Closer
	}
)
