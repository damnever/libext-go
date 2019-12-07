package io

import (
	"bytes"
	"io"
)

type RewindableReader struct {
	r       io.Reader
	buf     *bytes.Buffer
	bufread int
	rewind  bool
}

func NewRewindableReader(r io.Reader) *RewindableReader {
	buf := &bytes.Buffer{}
	return &RewindableReader{
		r:       io.TeeReader(r, buf),
		buf:     buf,
		bufread: 0,
		rewind:  false,
	}
}
func (rr *RewindableReader) Reset(r io.Reader) {
	rr.rewind = false
	rr.buf.Reset()
	rr.bufread = 0
	rr.r = r
}

func (rr *RewindableReader) Rewind() {
	rr.rewind = true
	rr.bufread = 0
}

func (rr *RewindableReader) Read(p []byte) (int, error) {
	if rr.rewind && rr.buf.Len() > rr.bufread {
		n := copy(p, rr.buf.Bytes()[rr.bufread:])
		rr.bufread += n
		return n, nil
	}
	return rr.r.Read(p)
}
