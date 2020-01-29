package io

import "io"

var (
	// NopReader does nothing when performing Read operation (XXX: make it panic?).
	NopReader = nopReader{}
	// NopWriter does nothing when performing Write operation (XXX: make it panic?).
	NopWriter = nopWriter{}
	// NopCloser does nothing when performing Close operation (XXX: make it panic?).
	NopCloser = nopCloser{}

	_ io.Reader = NopReader
	_ io.Writer = NopWriter
	_ io.Closer = NopCloser
)

type nopReader struct{}

func (nopReader) Read([]byte) (int, error) { return 0, nil }

type nopWriter struct{}

func (nopWriter) Write([]byte) (int, error) { return 0, nil }

type nopCloser struct{}

func (nopCloser) Close() error { return nil }
