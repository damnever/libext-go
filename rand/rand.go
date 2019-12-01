package rand

import (
	crand "crypto/rand"
)

// FIXME: Too slow

func Bytes(p []byte) {
	_, _ = crand.Read(p)
}

func String(n int) string {
	p := make([]byte, n, n)
	Bytes(p)
	return string(p)
}
