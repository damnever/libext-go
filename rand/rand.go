package rand

import (
	"math/rand"
	"time"

	strconvext "github.com/damnever/libext-go/strconv"
)

// FIXME: Too slow

// Shuffle is the shortcut for (*rand.Rand).Shuffle,
// because the global (package rand).Shuffle will generate same results.
func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}
	if n == 0 {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	r.Shuffle(n, swap)
}

// Bytes is the shortcut for (*rand.Rand).Read, it will panic if the underlying operation failed,
// we should only use it for testing purpose.
func Bytes(p []byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	_, err := r.Read(p)
	if err != nil {
		panic(err)
	}
}

// String returns a string which converted from (*rand.Rand).Read.
// we should only use it for testing purpose.
func String(n int) string {
	p := make([]byte, n, n)
	Bytes(p)
	return strconvext.UnsafeBtoa(p)
}
