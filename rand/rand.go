package rand

import (
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
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	r.Shuffle(n, swap)
}

func Bytes(p []byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	for i, n := 0, len(p); i < n; {
		nr, err := r.Read(p[i:])
		if err != nil {
			panic(err)
		}
		i += nr
	}
}

func String(n int) string {
	p := make([]byte, n, n)
	Bytes(p)
	return string(p)
}
