package io

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewindableReader(t *testing.T) {
	data := []byte("hello world")
	r := NewRewindableReader(bytes.NewBuffer(data))
	buf := make([]byte, len(data)-5)
	_, err := r.Read(buf)
	require.Nil(t, err)
	require.Equal(t, data[:len(data)-5], buf)

	for i := 0; i < 5; i++ {
		r.Rewind()
		buf = make([]byte, len(data))
		_, err = io.ReadFull(r, buf)
		require.Nil(t, err)
		require.Equal(t, data, buf)
	}
}
