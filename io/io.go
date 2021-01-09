package io

import (
	"bufio"
	"io"
)

type (
	// Flusher is the interface that wraps the Flush method.
	Flusher interface {
		Flush() error
	}

	// WithReadWriter combines io.Writer and io.Reader together for easy use.
	WithReadWriter struct {
		io.Writer
		io.Reader
	}
	// WithReadCloser combines io.Reader and io.Closer together for easy use.
	WithReadCloser struct {
		io.Reader
		io.Closer
	}
	// WithWriteCloser combines io.Writer and io.Closer together for easy use.
	WithWriteCloser struct {
		io.Writer
		io.Closer
	}
	// WithCloser combines io.ReadWriter and io.Closer together for easy use.
	WithCloser struct {
		io.ReadWriter
		io.Closer
	}
	// BufferedReadWriterWithCloser combines bufio.ReadWriter and io.Closer together for easy use.
	BufferedReadWriterWithCloser struct {
		*bufio.ReadWriter
		io.Closer
	}
)
