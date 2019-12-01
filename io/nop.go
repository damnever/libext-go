package io

var (
	// NopReader does nothing when performing Read operation (XXX: make it panic?).
	NopReader = nopReader{}
	// NopWriter does nothing when performing Write operation (XXX: make it panic?).
	NopWriter = nopWriter{}
)

type nopReader struct{}

func (r nopReader) Read([]byte) (int, error) { return 0, nil }

type nopWriter struct{}

func (w nopWriter) Write([]byte) (int, error) { return 0, nil }
