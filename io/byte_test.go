package io

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteReaderCopy(t *testing.T) {
	buf := &bytes.Buffer{}
	r := NewByteReader(buf)
	r.p[0] = 1
	require.Equal(t, fmt.Sprintf("%p %v", &(r.p[0]), r.p), getrpAddr(r))
}

func getrpAddr(r *ByteReader) string {
	return fmt.Sprintf("%p %v", &(r.p[0]), r.p)
}
