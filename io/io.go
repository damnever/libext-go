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
	WithCloser struct {
		io.ReadWriter
		io.Closer
	}
	BufferedReadWriterWithCloser struct {
		*bufio.ReadWriter
		io.Closer
	}
)
