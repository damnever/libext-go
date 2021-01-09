package io

import (
	"bytes"
	"io"
)

// RewindableReader is a reader which can be rewind and reads from begin again.
// It is not thread-safe.
type RewindableReader struct {
	r       io.Reader
	buf     *bytes.Buffer
	bufread int
	rewind  bool
}

// NewRewindableReader creates a new RewindableReader.
func NewRewindableReader(r io.Reader) *RewindableReader {
	rr := &RewindableReader{
		buf: &bytes.Buffer{},
	}
	rr.Reset(r)
	return rr
}

// Reset resets the internal state by given reader.
func (rr *RewindableReader) Reset(r io.Reader) {
	rr.rewind = false
	rr.buf.Reset()
	rr.bufread = 0
	rr.r = io.TeeReader(r, rr.buf)
}

// Rewind resets the read offset to the begin.
func (rr *RewindableReader) Rewind() {
	rr.rewind = true
	rr.bufread = 0
}

// Read implements io.Reader. It always returns len(p) and a nil error, otherwise the error is not nil.
func (rr *RewindableReader) Read(p []byte) (int, error) {
	if rr.rewind && rr.buf.Len() > rr.bufread {
		n := copy(p, rr.buf.Bytes()[rr.bufread:])
		rr.bufread += n
		if n == len(p) {
			return n, nil
		}
		// Short.
		p = p[n:]
	}
	return rr.r.Read(p)
}
