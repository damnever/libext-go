package bytes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	p := NewPoolWith(12, 512)
	for _, v := range []struct {
		size   int
		expect int
	}{
		{100, 512},
		{1023, 1024},
		{511, 512},
		{10000, 16384},
		{20000, 32768},
		{32769, 65536},
		{2000, 2048},
		{3000, 4096},
		{5000, 8192},
		{100000, 131072},
		{200000, 262144},
		{524287, 524288},
		{2097151, 2097152},
	} {
		b := p.Get(v.size)
		require.Equal(t, v.expect, cap(b), fmt.Sprintf("%+v", v))
		p.Put(b)
	}
}
