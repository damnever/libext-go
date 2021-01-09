package io

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewindableReader(t *testing.T) {
	data := []byte("hello world")
	r := NewRewindableReader(bytes.NewBuffer(data))
	{
		buf := make([]byte, len(data)-5)
		_, err := r.Read(buf)
		require.Nil(t, err)
		require.Equal(t, data[:len(data)-5], buf)
	}

	for i := len(data) - 1; i >= 0; i-- {
		r.Rewind()
		buf := make([]byte, len(data)-i)
		_, err := r.Read(buf)
		require.Nil(t, err)
		require.Equal(t, data[:len(data)-i], buf)
	}
	for i := 0; i < len(data); i++ {
		r.Rewind()
		buf := make([]byte, len(data)-i)
		_, err := r.Read(buf)
		require.Nil(t, err)
		require.Equal(t, data[:len(data)-i], buf)
	}
}
