package rand

import (
	crand "crypto/rand"
	"math/rand"
	"time"
)

// FIXME: Too slow

// Shuffle is the shortcut for (new *rand.Rand).Shuffle,
// because the global (package rand).Shuffle will generate same results.
func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}
	if n == 0 {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // seed + rand.Int63n(1000000)?
	r.Shuffle(n, swap)
}

func Bytes(p []byte) {
	_, _ = crand.Read(p)
}

func String(n int) string {
	p := make([]byte, n, n)
	Bytes(p)
	return string(p)
}
